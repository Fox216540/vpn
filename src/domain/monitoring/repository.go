package monitoring

type Repository interface {
	GetTrafficUsage() (down, up float64, err error)
	GetCPUPercentUsage() (float64, error)
	GetMemoryPercentUsage() (float64, error)
	GetActiveConnections() (int, error)
}
