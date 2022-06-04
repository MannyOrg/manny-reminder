package models

import "github.com/google/uuid"

type User struct {
	Id    *uuid.UUID `json:"id"`
	Email *string    `json:"email"`
	Token *string    `json:"-"`
}
