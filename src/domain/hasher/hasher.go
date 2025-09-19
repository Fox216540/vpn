package hasher

type Hasher interface {
	Verify(hash string) error
}
