package database

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/lambda-login/internal/entities"
	"github.com/jfelipearaujo-org/lambda-login/internal/providers/interfaces/mocks"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabase(t *testing.T) {
	// Arrange
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeProviderMock := mocks.NewMockTimeProvider(t)

	// Act
	database := NewDatabase(db, timeProviderMock)

	// Assert
	assert.NotNil(t, database)
}

func TestNewDatabaseFromConnStr(t *testing.T) {
	// Arrange
	timeProviderMock := mocks.NewMockTimeProvider(t)

	// Act
	database := NewDatabaseFromConnStr(timeProviderMock)

	// Assert
	assert.NotNil(t, database)
}

func TestDatabase_GetUserByCPF(t *testing.T) {
	t.Run("Should return a user when CPF is found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		timeProviderMock := mocks.NewMockTimeProvider(t)

		database := NewDatabase(db, timeProviderMock)

		rows := sqlmock.
			NewRows([]string{"id", "document_id", "password"}).
			AddRows([]driver.Value{"1", "123", "123456"})

		mock.ExpectQuery("SELECT (.+) FROM customers c").
			WithArgs("123").
			WillReturnRows(rows)

		expected := entities.User{
			Id:         "1",
			DocumentId: "123",
			Password:   "123456",
		}

		// Act
		user, err := database.GetUserByCPF("123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expected, user)
	})

	t.Run("Should return an error when user is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		timeProviderMock := mocks.NewMockTimeProvider(t)

		database := NewDatabase(db, timeProviderMock)

		rows := sqlmock.NewRows([]string{})

		mock.ExpectQuery("SELECT (.+) FROM customers c").
			WithArgs("123").
			WillReturnRows(rows)

		expected := entities.User{}

		// Act
		user, err := database.GetUserByCPF("123")

		// Assert
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.NotNil(t, user)
		assert.Equal(t, expected, user)
	})
}
