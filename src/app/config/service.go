package config

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
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
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return nil, err
	}

	if err := os.Remove(path); err != nil {
		log.Fatalf("Ошибка при удалении файла: %v", err)
	}

	return data, nil

}

func (s *service) DeleteConfig(configID uuid.UUID) error {
	if err := s.r.Delete(configID); err != nil {
		return err
	}
	return nil
}

func (s *service) StartServerConfig() error {
	if err := s.r.CreateServer(); err != nil {
		return err
	}
	return nil
}
