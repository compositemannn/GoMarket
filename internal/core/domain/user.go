package core_domain

import "time"

type User struct {
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

func NewUser(
	id int64,
	role string,
	passwordHash string,
	email string,
	phoneNumber *string,
	status string,
	name string,
	createdAt time.Time,
	updatedAt *time.Time,
) User {
	return User{
		ID:           id,
		Role:         role,
		PasswordHash: passwordHash,
		Email:        email,
		PhoneNumber:  phoneNumber,
		Status:       status,
		Name:         name,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
