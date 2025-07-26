package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/thecguygithub/lighthouse/collectors"
)

func main() {
	var agent_port string = os.Getenv("AGENT_PORT")
	var agent_hostname string = os.Getenv("AGENT_HOSTNAME")

	if agent_port == "" {
		agent_port = "8040"
	}
	if agent_hostname == "" {
		agent_hostname = "localhost"
	}

	fmt.Println("Starting Lighthouse Agent...")

	http.HandleFunc("/api/v1/minuteData", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		stats, err := CollectAll()
		if err != nil {
			http.Error(w, "Failed to collect stats: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jsonBytes, err := json.Marshal(stats)
		if err != nil {
			http.Error(w, "Failed to marshal stats: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
	})

	fmt.Println("Running WebServer at", agent_hostname+":"+agent_port)
	http.ListenAndServe(agent_hostname+":"+agent_port, nil)

}

var allCollectors = []collectors.Collector{
	&collectors.CPUStatsCollector{},
	&collectors.MemoryStatsCollector{},
	&collectors.DiskStatsCollector{},
	&collectors.NetworkStatsCollector{},
}

func CollectAll() (collectors.MetricData, error) {
	allData := collectors.MetricData{}

	for _, c := range allCollectors {
		data, err := c.Collect()
		if err != nil {
			return nil, err
		}
		for k, v := range data {
			allData[k] = v
		}
	}

	return allData, nil
}
