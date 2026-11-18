package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/compositemannn/GoMarket/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New() // в глобальной области видимости для кэшей и оптимизации

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"%w: decode json: %v",
			core_errors.ErrInvalidArgument,
			err,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"%w: validate request: %v",
			core_errors.ErrInvalidArgument,
			err,
		)
	}

	return nil
}
