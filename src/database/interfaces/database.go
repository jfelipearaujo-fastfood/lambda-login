package interfaces

import "github.com/jsfelipearaujo/lambda-login/src/entities"

type Database interface {
	GetUserByCPF(cpf string) (entities.User, error)
}
