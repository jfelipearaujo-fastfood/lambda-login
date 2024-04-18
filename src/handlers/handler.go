package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jfelipearaujo-org/lambda-login/src/cpf"
	db_interface "github.com/jfelipearaujo-org/lambda-login/src/database/interfaces"
	"github.com/jfelipearaujo-org/lambda-login/src/entities"
	hash_interface "github.com/jfelipearaujo-org/lambda-login/src/hashs/interfaces"
	"github.com/jfelipearaujo-org/lambda-login/src/router"
	token_interface "github.com/jfelipearaujo-org/lambda-login/src/token/interfaces"
)

type Handler struct {
	db     db_interface.Database
	hasher hash_interface.Hasher
	jwt    token_interface.Token
}

func NewHandler(
	db db_interface.Database,
	hasher hash_interface.Hasher,
	jwt token_interface.Token,
) Handler {
	return Handler{
		db:     db,
		hasher: hasher,
		jwt:    jwt,
	}
}

func (h Handler) AuthenticateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request entities.Request
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		return router.InvalidRequestBody(), nil
	}

	cpf := cpf.NewCPF(request.CPF)

	if !cpf.IsValid() || !request.IsPasswordValid() {
		return router.InvalidCPFOrPassword(), nil
	}

	user, err := h.db.GetUserByCPF(request.CPF)
	if err != nil {
		return router.InvalidCPFOrPassword(), nil
	}

	if err := h.hasher.CheckPassword(request.Password, user.Password); err != nil {
		return router.InvalidCPFOrPassword(), nil
	}

	token, err := h.jwt.CreateJwtToken(user)
	if err != nil {
		return router.InternalServerError(), nil
	}

	return router.Success(token), nil
}
