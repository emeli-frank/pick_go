package mysql

import (
	"database/sql"
	"github.com/emeli-frank/pick_go/pkg/domain/user"
	errors2 "github.com/emeli-frank/pick_go/pkg/errors"
	"github.com/go-sql-driver/mysql"
)

type userStorage struct {
	DB *sql.DB
}

func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{db}
}

func (r *userStorage) Tx() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *userStorage) SaveUser(user *user.User, hashedPassword string) (int, error) {
	const op = "userStorage.SaveUser"

	query := "INSERT INTO users (names, email, password) VALUE (?, ?, ?)"
	result, err := r.DB.Exec(query, user.Names, user.Email, hashedPassword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				err := &errors2.Conflict{Err:err, Item:"email"}
				return 0, errors2.Wrap(err, op, "executing insert query")
			}
		}
		return 0, errors2.Wrap(err, op, "executing insert query")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors2.Wrap(err, op, "getting last insert id")
	}

	return int(id), nil
}

func (r *userStorage) GetUserToAuthenticate(email string) (*user.User, string, error) {
	const op = "userStorage.GetUserToAuthenticate"
	query := `SELECT id, names, password FROM users WHERE email = ?`
	row := r.DB.QueryRow(query, email)

	user := user.User{}
	var hashedPassword string

	err := row.Scan(&user.ID, &user.Names, &hashedPassword)
	if err == sql.ErrNoRows {
		return nil, "", errors2.Wrap(&errors2.NotFound{Err: err}, op,"user with provided email does not exist")
	} else if err != nil {
		return nil, "", errors2.Wrap(err, op, "scanning rows into user struct")
	}

	user.Email = email

	return &user, hashedPassword, nil
}

func (r *userStorage) GetUser(userId int) (*user.User, error) {
	const op = "userStorage.Get"
	query := `SELECT id, names, email FROM users WHERE id = ?`
	row := r.DB.QueryRow(query, userId)

	user := user.User{}

	err := row.Scan(&user.ID, &user.Names, &user.Email)
	if err == sql.ErrNoRows {
		return nil, errors2.Wrap(&errors2.NotFound{Err:err}, op,"user not found")
	} else if err != nil {
		return nil, errors2.Wrap(err, op, "scanning rows into user struct")
	}

	return &user, nil
}
