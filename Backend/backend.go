package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/thecguygithub/lighthouse/backend/servers"
)

func main() {
	app := pocketbase.New()

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		initCollections(app)
		registerCronJobs(app)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func registerCronJobs(app core.App) {
	app.Cron().MustAdd("Pull Server Stats", "*/1 * * * *", func() {
		log.Println("Pulling server stats...")

		servers, err := app.FindAllRecords("servers")
		if err != nil {
			log.Println("Error fetching servers:", err)
			return
		}

		// Fetch the target collection once instead of inside the loop
		serverStatsCollection, _ := app.FindCollectionByNameOrId("server_stats")
		if serverStatsCollection == nil {
			log.Println("Error: server_stats collection not found")
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}

		for _, server := range servers {
			hostname := server.GetString("hostname")
			port := server.GetString("port")
			name := server.GetString("name")

			if hostname == "" || port == "" {
				log.Printf("Skipping server %s: missing hostname or port", name)
				continue
			}

			url := fmt.Sprintf("http://%s:%s/api/v1/minuteData", hostname, port)

			resp, err := client.Get(url)
			if err != nil {
				log.Printf("Error fetching stats for server %s: %v", name, err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("Failed to fetch stats for server %s: %s", name, resp.Status)
				resp.Body.Close()
				continue
			}

			var stats map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
				log.Printf("Error decoding stats for server %s: %v", name, err)
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			serverStats := core.NewRecord(serverStatsCollection)
			serverStats.Set("server", server.Id)
			serverStats.Set("stats", stats)
			serverStats.Set("type", "1m")

			if err := app.Save(serverStats); err != nil {
				log.Printf("Error saving stats for server %s: %v", name, err)
			} else {
				log.Printf("Successfully saved stats for server %s", name)
			}
		}
	})

	app.Cron().MustAdd("Summarize Records", "*/1 * * * *", func() {
		now := time.Now().UTC()
		log.Printf("⏳ CRON START at %s", now.Format(time.RFC3339Nano))

		summarizeRecords(app)

		log.Printf("✅ CRON END at %s", time.Now().UTC().Format(time.RFC3339Nano))
	})
}

func initCollections(app core.App) {
	usersCollection, err := app.FindCollectionByNameOrId("users")
	if err != nil {

		return
	}

	serversCollection := core.NewBaseCollection("servers")
	serverStatsCollection := core.NewBaseCollection("server_stats")

	serversCollection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})
	serversCollection.Fields.Add(&core.TextField{
		Name:     "hostname",
		Required: true,
	})
	serversCollection.Fields.Add(&core.TextField{
		Name:     "port",
		Required: true,
	})
	serversCollection.Fields.Add(&core.SelectField{
		Name:     "state",
		Required: true,
		Values:   []string{"online", "offline", "maintenance", "unknown"},
	})
	serversCollection.Fields.Add(&core.JSONField{
		Name:     "quickdata",
		Required: true,
	})
	serversCollection.Fields.Add(&core.RelationField{
		Name:          "user",
		Required:      true,
		CollectionId:  usersCollection.Id,
		CascadeDelete: true,
	})
	serversCollection.Fields.Add(&core.AutodateField{
		Name:     "created",
		OnCreate: true,
	})
	serversCollection.Fields.Add(&core.AutodateField{
		Name:     "updated",
		OnCreate: true,
		OnUpdate: true,
	})

	// ServerStats collection
	serverStatsCollection.Fields.Add(&core.RelationField{
		Name:          "server",
		CollectionId:  serversCollection.Id,
		CascadeDelete: true,
	})
	serverStatsCollection.Fields.Add(&core.JSONField{
		Name:     "stats",
		Required: true,
	})
	serverStatsCollection.Fields.Add(&core.SelectField{
		Name:     "type",
		Required: true,
		Values:   []string{"1m", "5m", "15m", "30m", "60m", "120m"},
	})
	serverStatsCollection.Fields.Add(&core.AutodateField{
		Name:     "created",
		OnCreate: true,
	})
	serverStatsCollection.Fields.Add(&core.AutodateField{
		Name:     "updated",
		OnCreate: true,
		OnUpdate: true,
	})

	app.Save(serversCollection)
	app.Save(serverStatsCollection)
}

func summarizeRecords(app core.App) {
	summarySteps := []struct {
		fromType      string
		toType        string
		bucket        time.Duration
		neededRecords int
	}{
		{"1m", "5m", -5 * time.Minute, 5},
		{"5m", "15m", -15 * time.Minute, 3},
		{"15m", "30m", -30 * time.Minute, 2},
		{"30m", "60m", -time.Hour, 2},
		{"60m", "120m", -2 * time.Hour, 2},
	}

	app.RunInTransaction(func(txApp core.App) error {
		systemsCollection, err := txApp.FindCollectionByNameOrId("servers")
		if err != nil {
			log.Fatal("Error whilst trying to get 'servers' collection:", err)
			return err
		}
		serverStatsCollection, err := txApp.FindCollectionByNameOrId("server_stats")
		if err != nil {
			log.Fatal("Error whilst trying to get 'server_stats' collection:", err)
			return err
		}

		systems, err := txApp.FindRecordsByFilter(systemsCollection, "state ?= 'online'", "", 0, 0)
		if err != nil {
			log.Fatal("Error whilst trying to get records from 'servers' collection ", err)
			return err
		}

		for _, system := range systems {
			log.Println("Summarizing for system named:", system.GetString("name"))
			for _, step := range summarySteps {
				log.Println("Summarizing " + step.fromType + " to " + step.toType)

				created := time.Now().UTC().Add(step.bucket)

				count, err := txApp.CountRecords(serverStatsCollection.Id, dbx.NewExp(
					"server={:server} && type={:type} && created > {:created}",
					dbx.Params{"server": system.Id, "type": step.fromType, "created": created},
				))

				if err != nil || count > 0 {
					continue
				}

				records, err := txApp.FindRecordsByFilter(serverStatsCollection, "server={:server} && type={:type} && created > {:created}", "-created", 0, 0, dbx.Params{"server": system.Id, "type": step.fromType, "created": created})
				if err != nil {
					log.Fatal("Error whilst trying to get records from 'server_stats' collection", err)
					return err
				}

				if len(records) == 0 {
					log.Println("Could not find any records!                       Line 225")
					return nil
				}

				log.Println("Found Records! ", len(records))

				if len(records) < step.neededRecords {
					log.Println("not enought records")
					continue
				}

				summarizedRecord := core.NewRecord(serverStatsCollection)
				summarizedRecord.Set("server", system.Id)
				summarizedRecord.Set("type", step.toType)
				summarizedRecord.Set("stats", GetSummarizedSystemStats(app, records))

				txApp.Save(summarizedRecord)
			}
		}

		return nil
	})
}

func GetSummarizedSystemStats(app core.App, records []*core.Record) servers.SystemStats {
	var totalBandwidthDownload float64
	var totalBandwidthUpload float64
	var totalCPUUsage float64
	var totalDiskRead int
	var totalDiskWrite int
	var totalDiskTotal int
	var totalDiskUsed int
	var totalMemoryCached int
	var totalMemoryFree int
	var totalMemoryTotal int
	var totalMemoryUsed int

	for _, record := range records {
		statsRaw := record.Get("stats")

		statsBytes, err := json.Marshal(statsRaw)
		if err != nil {
			log.Println("Malformed Json!")
			continue // skip malformed entries
		}

		var stats servers.SystemStats
		if err := json.Unmarshal(statsBytes, &stats); err != nil {
			log.Println("Malformed Json! (2)")
			continue // skip malformed entries
		}

		totalBandwidthDownload += stats.Bandwidth.Download
		totalBandwidthUpload += stats.Bandwidth.Upload
		totalCPUUsage += stats.CPUUsage
		totalDiskRead += stats.DiskIO.Read
		totalDiskWrite += stats.DiskIO.Write
		totalDiskTotal += stats.DiskTotal
		totalDiskUsed += stats.DiskUsed
		totalMemoryCached += stats.MemoryCached
		totalMemoryFree += stats.MemoryFree
		totalMemoryTotal += stats.MemoryTotal
		totalMemoryUsed += stats.MemoryUsed
	}

	count := float64(len(records))
	if count == 0 {
		return servers.SystemStats{} // fallback for no records
	}

	return servers.SystemStats{
		Bandwidth: servers.Bandwidth{
			Download: totalBandwidthDownload / count,
			Upload:   totalBandwidthUpload / count,
		},
		CPUUsage: totalCPUUsage / count,
		DiskIO: servers.DiskIO{
			Read:  totalDiskRead / int(count),
			Write: totalDiskWrite / int(count),
		},
		DiskTotal:    totalDiskTotal / int(count),
		DiskUsed:     totalDiskUsed / int(count),
		MemoryCached: totalMemoryCached / int(count),
		MemoryFree:   totalMemoryFree / int(count),
		MemoryTotal:  totalMemoryTotal / int(count),
		MemoryUsed:   totalMemoryUsed / int(count),
	}
}
