package traffic

type Repository interface {
	GetTraffic() (down, up float64, err error)
}
