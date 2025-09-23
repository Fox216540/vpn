package monitoring

import "vpn/src/domain/monitoring"

type service struct {
	r monitoring.Repository
}

func NewService(r monitoring.Repository) UseCase {
	return &service{r: r}
}

func (s *service) GetTotalMbps() float64 {
	trafficDown, trafficUp := s.r.GetTrafficUsage()
	return trafficDown + trafficUp
}

func (s *service) GetSplitMbps() (down, up float64) {
	trafficDown, trafficUp := s.r.GetTrafficUsage()
	return trafficDown, trafficUp
}

func (s *service) GetConnections() int {
	return s.r.GetActiveConnections()
}
func (s *service) GetCPUPercent() float64 {
	return s.r.GetCPUPercentUsage()
}
func (s *service) GetMemoryPercent() float64 {
	return s.r.GetMemoryPercentUsage()
}
