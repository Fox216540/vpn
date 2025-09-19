package di

import (
	"vpn/src/app/config"
	r "vpn/src/infra/config"
	hr "vpn/src/infra/hasher"
)

func GetConfigService() config.UseCase {
	repo := r.NewRepository()
	hRepo := hr.NewHasher()
	return config.NewService(repo, hRepo)
}
