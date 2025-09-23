package config

import (
	"errors"
	"github.com/google/uuid"
	"os"
	"vpn/src/core/exception"
	"vpn/src/domain/config"
	"vpn/src/domain/hasher"
)

type service struct {
	r config.Repository
	h hasher.Hasher
}

func NewService(r config.Repository, h hasher.Hasher) UseCase {
	return &service{r: r, h: h}
}

func (s *service) CreateConfig(configID uuid.UUID) ([]byte, error) {
	path, err := s.r.Create(configID)
	if err != nil {
		var serverError *exception.ServerError
		if errors.As(err, &serverError) {
			return nil, serverError
		}
		return nil, NewInvalidCreateConfig(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, NewInvalidReadFile(err)
	}

	if err = os.Remove(path); err != nil {
		return nil, NewInvalidRemoveFile(err)
	}

	return data, nil

}

func (s *service) DeleteConfig(configID uuid.UUID) error {
	if err := s.r.Delete(configID); err != nil {
		var serverError *exception.ServerError
		if errors.As(err, &serverError) {
			return serverError
		}
		return NewInvalidDeleteConfig(err)
	}
	return nil
}

func (s *service) DeleteConfigs(configIDs []uuid.UUID) error {
	if err := s.r.DeleteIDs(configIDs); err != nil {
		return err
	}
	return nil
}

func (s *service) StartServerConfig() error {
	if err := s.r.CreateServer(); err != nil {
		var serverError *exception.ServerError
		if errors.As(err, &serverError) {
			return serverError
		}
		return NewInvalidStartServer(err)
	}
	return nil
}
