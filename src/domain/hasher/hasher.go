package hasher

type Hasher interface {
	Verify(password string, hash string) error
}
