package users_service

import (
	"context"

	core_domain "github.com/compositemannn/GoMarket/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	input CreateUserInput,
) (core_domain.User, error) {
	if err := input.Validate(); err != nil {
		return core_domain.User{}, err
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return core_domain.User{}, err
	}

	userWithHash := core_domain.User{
		Role:         input.Role,
		Email:        input.Email,
		PasswordHash: passwordHash,
		PhoneNumber:  input.PhoneNumber,
		Name:         input.Name,
		Status:       "active",
	}

	user, err := s.usersRepository.CreateUser(ctx, userWithHash)
	if err != nil {
		return core_domain.User{}, err
	}

	return user, nil
}
