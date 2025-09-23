package hasher

import (
	"vpn/src/core/settings"
	"vpn/src/domain/hasher"
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Verify(hash string) error {
	if hash != settings.Config.HashPass {
		return hasher.NewWrongPasswordError(nil)
	}
	return nil
}
