package auth

import (
	"database/sql"
	"log"
	"manny-reminder/internal/models"
)

type AuthRepository interface {
	GetUsers() ([]models.User, error)
	AddUser(authCode string, token string) error
	GetUser(id string) (*models.User, error)
}

type RepositoryImpl struct {
	l  *log.Logger
	db *sql.DB
}

func NewRepository(l *log.Logger, db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{l, db}
}

func (r RepositoryImpl) GetUsers() ([]models.User, error) {
	var res models.User
	var users []models.User
	rows, err := r.db.Query("SELECT id, email, token FROM users")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			r.l.Fatal(err)
		}
	}()
	for rows.Next() {
		err := rows.Scan(&res.Id, &res.Email, &res.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, res)
	}
	return users, nil
}

func (r RepositoryImpl) GetUser(userId string) (*models.User, error) {
	var user models.User
	row := r.db.QueryRow("SELECT id, email, token FROM users WHERE id = $1 LIMIT 1", userId)
	err := row.Scan(&user.Id, &user.Email, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r RepositoryImpl) AddUser(id, token string) error {
	_, err := r.db.Exec("INSERT INTO users (id, token) VALUES ($1, $2)", id, token)
	if err != nil {
		return err
	}

	return nil
}
