package hasher

import (
	"fmt"
	"vpn/src/core/settings"
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Verify(hash string) error {
	if hash != settings.Config.HashPass {
		return fmt.Errorf("invalid hash")
	}
	return nil
}
