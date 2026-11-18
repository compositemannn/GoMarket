package users_service

import (
	"fmt"

	core_errors "github.com/compositemannn/GoMarket/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var inputValidator = validator.New()

type CreateUserInput struct {
	Role        string  `validate:"required,oneof=buyer seller"`
	Email       string  `validate:"required,email"`
	Password    string  `validate:"required,min=8"`
	PhoneNumber *string `validate:"omitempty,e164"`
	Name        string  `validate:"required,min=3,max=100"`
}

func (i *CreateUserInput) Validate() error {
	if err := inputValidator.Struct(i); err != nil {
		return fmt.Errorf(
			"%w: validate create user input: %v",
			core_errors.ErrInvalidArgument,
			err,
		)
	}
	return nil
}
