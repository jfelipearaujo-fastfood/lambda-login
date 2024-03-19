package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jsfelipearaujo/fast-food-lambda-auth/src/cpf"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	engine = "postgres"
)

var (
	signingKey = []byte(os.Getenv("SIGN_KEY"))

	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")

	connectionStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
)

type Request struct {
	CPF      string `json:"cpf"`
	Password string `json:"pass"`
}

type Response struct {
	Status      int    `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token,omitempty"`
}

type User struct {
	Id         string `json:"Id"`
	DocumentId string `json:"DocumentId"`
	Password   string `json:"Password"`
}

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	log := slog.New(handler)

	slog.SetDefault(log)
}

func buildResponse(status int, message string, token string) events.APIGatewayProxyResponse {
	response := Response{
		Status:      status,
		Message:     message,
		AccessToken: token,
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Error("error while trying to marshal the response", "error", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func getUserByCpf(cpf string) (*User, error) {
	conn, err := sql.Open(engine, connectionStr)
	if err != nil {
		slog.Error("error while trying to connect to the database", "error", err)
		return nil, err
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		slog.Error("error while trying to ping the database", "error", err)
		return nil, err
	}

	statement, err := conn.Query(`SELECT c."Id", c."DocumentId", c."Password" FROM clients c WHERE c."DocumentId" = $1`, cpf)
	if err != nil {
		slog.Error("error while trying to execute the query", "error", err)
		return nil, err
	}

	var user User
	for statement.Next() {
		if err := statement.Scan(&user.Id, &user.DocumentId, &user.Password); err != nil {
			slog.Error("error while trying to scan the result", "error", err)
			return nil, err
		}
	}

	return &user, nil
}

func maskCpf(cpf string) string {
	return strings.ReplaceAll(cpf, cpf[3:(len(cpf)-2)], strings.Repeat("*", len(cpf)-5))
}

func clearSpecialChars(cpf string) string {
	charsToRemove := []string{"/", ".", "-", " "}
	for _, char := range charsToRemove {
		cpf = strings.ReplaceAll(cpf, char, "")
	}
	return cpf
}

func checkPassword(password string, passwordHashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))
}

func createJwtToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(signingKey)
}

func validateJwtToken(token string, key []byte) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}

func handleAuth(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request Request

	slog.Info("checking user credentials")

	slog.Debug("unmarshalling the request")
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		slog.Error("error while trying to unmarshal the request", "error", err)
		return buildResponse(http.StatusUnauthorized, "error to parse the request body", ""), nil
	}

	cpf := cpf.NewCPF(clearSpecialChars(request.CPF))

	slog.Debug("validating the cpf")
	if !cpf.IsValid() {
		slog.Error("invalid cpf", "cpf", request.CPF)
		return buildResponse(http.StatusUnauthorized, "invalid CPF or password", ""), nil
	}

	slog.Debug("validating the password")
	if len(request.Password) < 8 {
		slog.Error("invalid password", "password_length", len(request.Password))
		return buildResponse(http.StatusUnauthorized, "invalid CPF or password", ""), nil
	}

	slog.Debug("getting the user by cpf")
	user, err := getUserByCpf(request.CPF)
	if err != nil {
		slog.Error("error while trying to get the user by cpf", "error", err)
		return buildResponse(http.StatusUnauthorized, "invalid CPF or password", ""), nil
	}

	if user == nil {
		slog.Error("user not found", "cpf", maskCpf(request.CPF))
		return buildResponse(http.StatusUnauthorized, "invalid CPF or password", ""), nil
	}

	slog.Debug("checking the password hash")
	if err := checkPassword(request.Password, user.Password); err != nil {
		slog.Error("invalid password, hash not match", "cpf", maskCpf(request.CPF))
		return buildResponse(http.StatusUnauthorized, "invalid CPF or password", ""), nil
	}

	slog.Debug("creating the jwt token")
	token, err := createJwtToken(user)
	if err != nil {
		slog.Error("error while trying to create the jwt token", "error", err)
		return buildResponse(http.StatusInternalServerError, "internal server error", token), nil
	}

	slog.Info("user authenticated successfully", "cpf", maskCpf(request.CPF))

	return buildResponse(http.StatusCreated, "success", token), nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("received a request", "path", req.Path, "method", req.HTTPMethod)

	if req.Path == "/login" && req.HTTPMethod == "POST" {
		return handleAuth(req)
	}

	return buildResponse(http.StatusMethodNotAllowed, "method not allowed", ""), nil
}

func main() {
	lambda.Start(router)
}
