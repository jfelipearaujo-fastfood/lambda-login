package interfaces

import "github.com/jfelipearaujo-org/lambda-login/src/entities"

type Database interface {
	GetUserByCPF(cpf string) (entities.User, error)
}
