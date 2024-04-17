package interfaces

import (
	"github.com/jsfelipearaujo/lambda-login/src/entities"
)

type Token interface {
	CreateJwtToken(user entities.User) (string, error)
}
