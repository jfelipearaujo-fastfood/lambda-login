package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jfelipearaujo-org/lambda-login/internal/entities"
	"github.com/jfelipearaujo-org/lambda-login/internal/providers/interfaces"
	_ "github.com/lib/pq"
)

const (
	engine = "postgres"

	DOCUMENT_TYPE_CPF = 1
)

var (
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")

	ErrUserNotFound = errors.New("user not found")
)

type Database struct {
	conn         *sql.DB
	timeProvider interfaces.TimeProvider
}

func NewDatabase(db *sql.DB, timeProvider interfaces.TimeProvider) *Database {
	return &Database{
		conn:         db,
		timeProvider: timeProvider,
	}
}

func NewDatabaseFromConnStr(timeProvider interfaces.TimeProvider) *Database {
	db, err := sql.Open(engine, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		panic(err)
	}

	return &Database{
		conn:         db,
		timeProvider: timeProvider,
	}
}

func (db *Database) GetUserByCPF(cpf string) (entities.User, error) {
	var user entities.User

	statement, err := db.conn.Query("SELECT c.id, c.document_id, c.password FROM customers c WHERE c.document_id = $1;", cpf)
	if err != nil {
		return user, err
	}

	for statement.Next() {
		if err := statement.Scan(&user.Id, &user.DocumentId, &user.Password); err != nil {
			return user, err
		}
	}

	if user == (entities.User{}) {
		return user, ErrUserNotFound
	}

	return user, nil
}
