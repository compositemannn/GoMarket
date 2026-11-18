package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/compositemannn/GoMarket/internal/core/domain"
	core_errors "github.com/compositemannn/GoMarket/internal/core/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	ctxSQLQuery, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO gomarket.users (role, name, status, email, password_hash, phone)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, role, password_hash, email, phone,
        status, name, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctxSQLQuery,
		query,
		user.Role,
		user.Name,
		user.Status,
		user.Email,
		user.PasswordHash,
		user.PhoneNumber,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Role,
		&userModel.PasswordHash,
		&userModel.Email,
		&userModel.PhoneNumber,
		&userModel.Status,
		&userModel.Name,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 - cущность в бд уже существует
			return core_domain.User{}, fmt.Errorf(
				"create user: %w: %w",
				core_errors.ErrConflict,
				err,
			)
		}

		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}

	userDomain := core_domain.NewUser(
		userModel.ID,
		userModel.Role,
		userModel.PasswordHash,
		userModel.Email,
		userModel.PhoneNumber,
		userModel.Status,
		userModel.Name,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)

	return userDomain, nil
}
