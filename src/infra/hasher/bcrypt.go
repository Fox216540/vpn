package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Verify(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
