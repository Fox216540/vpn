package di

import (
	"vpn/src/app/monitoring"
	r "vpn/src/infra/monitoring"
)

func GetTrafficService() monitoring.UseCase {
	repo := r.NewRepository()
	return monitoring.NewService(repo)
}
