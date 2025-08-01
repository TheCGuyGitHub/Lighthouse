package servers

type Bandwidth struct {
	Download float64 `json:"download"`
	Upload   float64 `json:"upload"`
}

type DiskIO struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type SystemStats struct {
	Bandwidth    Bandwidth `json:"bandwidth"`
	CPUUsage     float64   `json:"cpu_usage"`
	DiskIO       DiskIO    `json:"disk_io"`
	DiskTotal    int       `json:"disk_total"`
	DiskUsed     int       `json:"disk_used"`
	MemoryCached int       `json:"memory_cached"`
	MemoryFree   int       `json:"memory_free"`
	MemoryTotal  int       `json:"memory_total"`
	MemoryUsed   int       `json:"memory_used"`
}
