package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	db_interface "github.com/jfelipearaujo-org/lambda-login/internal/database/interfaces"
	db_interface_mock "github.com/jfelipearaujo-org/lambda-login/internal/database/interfaces/mocks"
	"github.com/jfelipearaujo-org/lambda-login/internal/entities"
	hash_interface "github.com/jfelipearaujo-org/lambda-login/internal/hashs/interfaces"
	hash_interface_mock "github.com/jfelipearaujo-org/lambda-login/internal/hashs/interfaces/mocks"
	token_interface "github.com/jfelipearaujo-org/lambda-login/internal/token/interfaces"
	token_interface_mock "github.com/jfelipearaujo-org/lambda-login/internal/token/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		db     db_interface.Database
		hasher hash_interface.Hasher
		jwt    token_interface.Token
	}
	tests := []struct {
		name string
		args args
		want Handler
	}{
		{
			name: "Should return a new instance correctly",
			args: args{
				db:     db_interface_mock.NewMockDatabase(t),
				hasher: hash_interface_mock.NewMockHasher(t),
				jwt:    token_interface_mock.NewMockToken(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange

			// Act
			got := NewHandler(tt.args.db, tt.args.hasher, tt.args.jwt)

			// Assert
			assert.IsType(t, tt.want, got)
		})
	}
}

func TestHandler_AuthenticateUser(t *testing.T) {
	t.Run("Should return a success response", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("GetUserByCPF", "218.486.310-65").
			Return(entities.User{
				Id:         "1",
				DocumentId: "218.486.310-65",
				Password:   "12345678",
			}, nil).
			Once()

		hasher_mock.On("CheckPassword", "12345678", "12345678").
			Return(nil).
			Once()

		jwt_mock.On("CreateJwtToken", mock.AnythingOfType("entities.User")).
			Return("token", nil).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		resp, err := h.AuthenticateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return error if CPF is invalid", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"1234","pass":"12345678"}`,
		}

		// Act
		resp, err := h.AuthenticateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return error if password is invalid", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"abc"}`,
		}

		// Act
		resp, err := h.AuthenticateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return error while trying to fetch the user by CPF", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("GetUserByCPF", "218.486.310-65").
			Return(entities.User{}, errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		resp, err := h.AuthenticateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return error while trying to check the user's password", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("GetUserByCPF", "218.486.310-65").
			Return(entities.User{
				Id:         "1",
				DocumentId: "218.486.310-65",
				Password:   "12345678",
			}, nil).
			Once()

		hasher_mock.On("CheckPassword", "12345678", "12345678").
			Return(errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		resp, err := h.AuthenticateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return error while trying to create the JWT token", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("GetUserByCPF", "218.486.310-65").
			Return(entities.User{
				Id:         "1",
				DocumentId: "218.486.310-65",
				Password:   "12345678",
			}, nil).
			Once()

		hasher_mock.On("CheckPassword", "12345678", "12345678").
			Return(nil).
			Once()

		jwt_mock.On("CreateJwtToken", mock.AnythingOfType("entities.User")).
			Return("", errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		resp, err := h.AuthenticateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})
}
