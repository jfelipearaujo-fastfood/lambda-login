package interfaces

import "github.com/jfelipearaujo-org/lambda-login/internal/entities"

type Database interface {
	GetUserByCPF(cpf string) (entities.User, error)
}
