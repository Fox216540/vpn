package monitoring

type UseCase interface {
	GetTotalMbps() (float64, error)
	GetSplitMbps() (down, up float64, err error)
	GetConnections() (int, error)
	GetCPUPercent() (float64, error)
	GetMemoryPercent() (float64, error)
}
