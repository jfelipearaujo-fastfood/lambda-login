package hashs

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
}

func NewHasher() Hasher {
	return Hasher{}
}

func (h Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h Hasher) CheckPassword(password string, passwordHashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))
}
