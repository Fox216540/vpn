package config

import "github.com/google/uuid"

type UseCase interface {
	CreateConfig(configID uuid.UUID) ([]byte, error)
	DeleteConfig(configID uuid.UUID) error
	DeleteConfigs(configIDs []uuid.UUID) error
	StartServerConfig() error
}
