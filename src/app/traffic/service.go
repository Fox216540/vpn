package traffic

import "vpn/src/domain/traffic"

type service struct {
	r traffic.Repository
}

func NewService(r traffic.Repository) UseCase {
	return &service{r: r}
}

func (s *service) GetTotalMbps() (float64, error) {
	trafficDown, trafficUp, err := s.r.GetTraffic()
	if err != nil {
		//TODO: log error
		return 0, err
	}
	return trafficDown + trafficUp, err
}

func (s *service) GetSplitMbps() (down, up float64, err error) {
	trafficDown, trafficUp, err := s.r.GetTraffic()
	if err != nil {
		//TODO: log error
		return 0, 0, err
	}
	return trafficDown, trafficUp, err
}
