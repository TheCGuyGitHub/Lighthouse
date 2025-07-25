package collectors

type MetricData map[string]interface{}

type Collector interface {
	Name() string                 // e.g. "cpu", "memory"
	Collect() (MetricData, error) // returns the collected metrics
}
