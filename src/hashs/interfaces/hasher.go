package interfaces

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, passwordHashed string) error
}
