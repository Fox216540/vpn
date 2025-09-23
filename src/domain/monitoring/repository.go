package monitoring

type Repository interface {
	GetTrafficUsage() (down, up float64)
	GetCPUPercentUsage() float64
	GetMemoryPercentUsage() float64
	GetActiveConnections() int
}
