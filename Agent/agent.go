package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/thecguygithub/lighthouse/collectors"
)

func main() {
	var backend_ip string = os.Getenv("BACKEND_IP")

	fmt.Println("Starting Lighthouse Agent...")
	fmt.Println("Trying to connect to backend at:", backend_ip)

	for {
		stats, err := CollectAll()
		if err != nil {
			panic(err)
		}

		jsonBytes, _ := json.Marshal(stats)
		fmt.Println(string(jsonBytes))
		time.Sleep(5 * time.Second)
	}
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
