package core_http_server

import (
	"context"
	"fmt"
	"net/http"

	core_logger "github.com/compositemannn/GoMarket/internal/core/logger"
	core_http_middleware "github.com/compositemannn/GoMarket/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct { // хранит необходимые зависимости для http.Server
	mux        *http.ServeMux // регистрирует разные роутеры (версии API)
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion) // например: /api/v1

		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("Start HTTP server", zap.String("addr", h.config.Addr))
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			ch <- err
		}
	}()

	select { // ожидает либо отмены контекста, либо ошибки из канала
	case err := <-ch:
		if err != nil { // если канал закрыт, то прочитанное значение - nil
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("Shutdown HTTP server")

		shutDownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel() // отменяем сами, если Shutdown завершился до истечения таймаута

		if err := server.Shutdown(shutDownCtx); err != nil { // ждет либо отмену контекста (истек таймаут), либо завершение работы обработчиков
			_ = server.Close() // принудительное заврешение, если истек таймаут
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}
	return nil
}
