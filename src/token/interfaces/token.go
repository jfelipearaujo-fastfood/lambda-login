package interfaces

import (
	"github.com/jfelipearaujo-org/lambda-login/src/entities"
)

type Token interface {
	CreateJwtToken(user entities.User) (string, error)
}
