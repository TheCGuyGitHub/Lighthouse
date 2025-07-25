package collectors

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUStatsCollector struct{}

func (c *CPUStatsCollector) Name() string {
	return "cpu"
}

// Collect gathers CPU usage statistics.
func (c *CPUStatsCollector) Collect() (MetricData, error) {
	percentages, err := cpu.Percent(0*time.Second, false)
	if err != nil {
		return nil, err
	}
	if len(percentages) == 0 {
		return nil, fmt.Errorf("no CPU usage data returned")
	}
	return MetricData{
		"cpu_usage": percentages[0],
	}, nil
}
