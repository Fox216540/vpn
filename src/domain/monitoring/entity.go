package monitoring

type Monitoring struct {
	MbpsDown           float64
	MbpsUp             float64
	CPUPercentUsage    float64
	MemoryPercentUsage float64
	ActiveConnections  int
}
