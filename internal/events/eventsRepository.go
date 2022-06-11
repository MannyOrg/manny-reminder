package events

import (
	"database/sql"
	"log"
)

type EventsRepository interface {
}

type RepositoryImpl struct {
	l  *log.Logger
	db *sql.DB
}

func NewRepository(l *log.Logger, db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{l, db}
}
