package collectors

import (
	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryStatsCollector struct{}

func (m *MemoryStatsCollector) Name() string {
	return "memory"
}

func (m *MemoryStatsCollector) Collect() (MetricData, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	const mb = 1024 * 1024

	return MetricData{
		"memory_total":  vmStat.Total / mb,
		"memory_used":   vmStat.Used / mb,
		"memory_cached": vmStat.Cached / mb,
		"memory_free":   vmStat.Free / mb,
	}, nil
}
