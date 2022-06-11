package events

import (
	"database/sql"
	"log"
)

type Repository interface {
}

type RepositoryImpl struct {
	l  *log.Logger
	db *sql.DB
}

func NewRepository(l *log.Logger, db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{l, db}
}
