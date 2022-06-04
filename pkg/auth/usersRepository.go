package auth

import (
	"database/sql"
	"log"
	"manny-reminder/pkg/models"
)

type IRepository interface {
	GetUsers() ([]models.User, error)
	AddUser(authCode string) error
}

type Repository struct {
	l  *log.Logger
	db *sql.DB
}

func NewRepository(l *log.Logger, db *sql.DB) *Repository {
	return &Repository{l, db}
}

func (r Repository) GetUsers() ([]models.User, error) {
	var res models.User
	var users []models.User
	rows, err := r.db.Query("SELECT id, email FROM users")
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
		err := rows.Scan(&res.Id, &res.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, res)
	}
	return users, nil
}

func (r Repository) AddUser(id string) error {
	_, err := r.db.Exec("INSERT INTO users (id) VALUES ($1)", id)
	if err != nil {
		return err
	}

	return nil
}
