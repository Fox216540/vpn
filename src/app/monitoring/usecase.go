package monitoring

type UseCase interface {
	GetTotalMbps() float64
	GetSplitMbps() (down, up float64)
	GetConnections() int
	GetCPUPercent() float64
	GetMemoryPercent() float64
}
