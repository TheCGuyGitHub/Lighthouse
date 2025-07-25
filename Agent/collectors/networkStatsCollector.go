package collectors

import (
	"sync"

	"github.com/shirou/gopsutil/v4/net"
)

type NetworkStatsCollector struct {
	prevDownload uint64
	prevUpload   uint64
	mu           sync.Mutex
}

func (n *NetworkStatsCollector) Name() string {
	return "network"
}

var filteredInterfaces = []string{"lo", "docker0"}

func (n *NetworkStatsCollector) Collect() (MetricData, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return nil, err
	}

	filtered := make([]net.IOCountersStat, 0)
	for _, counter := range counters {
		if !contains(filteredInterfaces, counter.Name) {
			filtered = append(filtered, counter)
		}
	}
	counters = filtered

	if len(counters) == 0 {
		return nil, nil
	}

	curDownload := counters[0].BytesRecv
	curUpload := counters[0].BytesSent

	// Calculate delta since last Collect()
	deltaDownload := curDownload - n.prevDownload
	deltaUpload := curUpload - n.prevUpload

	// Store current values for next time
	n.prevDownload = curDownload
	n.prevUpload = curUpload

	// Return the bandwidth usage since last call
	return MetricData{
		"bandwidth": map[string]float64{
			"download": float64(deltaDownload) / (1024 * 1024),
			"upload":   float64(deltaUpload) / (1024 * 1024),
		},
	}, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
