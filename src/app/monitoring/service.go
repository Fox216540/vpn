package monitoring

import "vpn/src/domain/monitoring"

type service struct {
	r monitoring.Repository
}

func NewService(r monitoring.Repository) UseCase {
	return &service{r: r}
}

func (s *service) GetTotalMbps() (float64, error) {
	trafficDown, trafficUp, err := s.r.GetTrafficUsage()
	if err != nil {
		//TODO: log error
		return 0, err
	}
	return trafficDown + trafficUp, err
}

func (s *service) GetSplitMbps() (down, up float64, err error) {
	trafficDown, trafficUp, err := s.r.GetTrafficUsage()
	if err != nil {
		//TODO: log error
		return 0, 0, err
	}
	return trafficDown, trafficUp, err
}

func (s *service) GetConnections() (int, error) {
	return s.r.GetActiveConnections()
}
func (s *service) GetCPUPercent() (float64, error) {
	return s.r.GetCPUPercentUsage()
}
func (s *service) GetMemoryPercent() (float64, error) {
	return s.r.GetMemoryPercentUsage()
}
