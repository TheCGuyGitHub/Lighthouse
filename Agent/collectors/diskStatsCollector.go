package collectors

import (
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStatsCollector struct{}

func (d *DiskStatsCollector) Name() string {
	return "disk"
}

func (d *DiskStatsCollector) Collect() (MetricData, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	ioCounters, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}
	// Convert bytes to megabytes (MB)
	const mb = 1024 * 1024
	usage.Total = usage.Total / mb
	usage.Used = usage.Used / mb
	var readBytes, writeBytes uint64
	for _, counter := range ioCounters {
		readBytes += counter.ReadBytes / mb
		writeBytes += counter.WriteBytes / mb
	}

	return MetricData{
		"disk_total": usage.Total,
		"disk_used":  usage.Used,
		"disk_io": map[string]uint64{
			"read":  readBytes,
			"write": writeBytes,
		},
	}, nil
}
