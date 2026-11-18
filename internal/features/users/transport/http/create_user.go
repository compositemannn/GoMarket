package users_transport_http

import (
	"net/http"

	core_domain "github.com/compositemannn/GoMarket/internal/core/domain"
	core_logger "github.com/compositemannn/GoMarket/internal/core/logger"
	core_http_request "github.com/compositemannn/GoMarket/internal/core/transport/http/request"
	core_http_response "github.com/compositemannn/GoMarket/internal/core/transport/http/response"
	users_service "github.com/compositemannn/GoMarket/internal/features/users/service"
)

type CreateUserRequest struct {
	Role        string  `json:"role" validate:"required,oneof=buyer seller"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	PhoneNumber *string `json:"phone" validate:"omitempty,e164"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
}

type CreateUserResponse struct {
	ID          int64   `json:"id"`
	Role        string  `json:"role"`
	Email       string  `json:"email"`
	Status      string  `json:"status"`
	PhoneNumber *string `json:"phone"`
	Name        string  `json:"name"`
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidate(r, &request); err != nil { // декодируем запрос в структуру и валидируем
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userInput := inputFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userInput)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := dtoFromDomain(userDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)

}

func inputFromDTO(dto CreateUserRequest) users_service.CreateUserInput {
	return users_service.CreateUserInput{
		Role:        dto.Role,
		Email:       dto.Email,
		Password:    dto.Password,
		PhoneNumber: dto.PhoneNumber,
		Name:        dto.Name,
	}
}

func dtoFromDomain(user core_domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:          user.ID,
		Role:        user.Role,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Status:      user.Status,
		Name:        user.Name,
	}
}
