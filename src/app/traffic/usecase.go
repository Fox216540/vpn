package traffic

type UseCase interface {
	GetTotalMbps() (float64, error)
	GetSplitMbps() (down, up float64, err error)
}
