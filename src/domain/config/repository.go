package config

import (
	"github.com/google/uuid"
)

type Repository interface {
	Create(configID uuid.UUID) (string, error)
	Delete(configID uuid.UUID) error
	CreateServer() error
}
