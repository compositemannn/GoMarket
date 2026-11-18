package users_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/compositemannn/GoMarket/internal/core/domain"
	core_http_server "github.com/compositemannn/GoMarket/internal/core/transport/http/server"
	users_service "github.com/compositemannn/GoMarket/internal/features/users/service"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		input users_service.CreateUserInput,
	) (core_domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{usersService: usersService}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.CreateUser,
		},
	}
}
