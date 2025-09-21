package di

import (
	"vpn/src/app/traffic"
	r "vpn/src/infra/traffic"
)

func GetTrafficService() traffic.UseCase {
	repo := r.NewRepository()
	return traffic.NewService(repo)
}
