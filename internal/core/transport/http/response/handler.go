package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/compositemannn/GoMarket/internal/core/errors"
	core_logger "github.com/compositemannn/GoMarket/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(log *core_logger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(body any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(body); err != nil {
		h.log.Error("Write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		http.StatusInternalServerError,
		"internal server error",
	)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, message string) {
	var (
		statusCode    int
		publicMessage string
		logFunc       func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		publicMessage = "conflict"
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		publicMessage = "invalid argument"
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		publicMessage = "not found"
		logFunc = h.log.Debug

	default:
		statusCode = http.StatusInternalServerError
		publicMessage = "internal server error"
		logFunc = h.log.Error
	}

	logFunc(message, zap.Error(err))

	h.errorResponse(statusCode, publicMessage)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	msg string,
) {
	response := map[string]string{
		"message": msg,
	}

	h.JSONResponse(response, statusCode)
}
