package users_postgres_repository

import "time"

type UserModel struct {
	ID           int64
	Role         string
	PasswordHash string
	Email        string
	PhoneNumber  *string
	Status       string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
