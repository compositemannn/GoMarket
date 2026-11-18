package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/compositemannn/GoMarket/internal/core/logger"
	core_postgres_pool "github.com/compositemannn/GoMarket/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/compositemannn/GoMarket/internal/core/transport/http/middleware"
	core_http_server "github.com/compositemannn/GoMarket/internal/core/transport/http/server"
	users_postgres_repository "github.com/compositemannn/GoMarket/internal/features/users/repository/postgres"
	users_service "github.com/compositemannn/GoMarket/internal/features/users/service"
	users_transport_http "github.com/compositemannn/GoMarket/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("Failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Error("failed to init postgres connection pool", zap.Error(err))
		return
	}
	defer pool.Close()

	usersRepository := users_postgres_repository.NewUsersRepository(pool)

	usersService := users_service.NewUsersService(usersRepository)

	logger.Debug("Starting GoMarket")

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing HTTP server")
	server := core_http_server.NewHTTPServer(
		core_http_server.NewHTTPConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	routerAPIv1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	routerAPIv1.RegisterRoutes(usersTransportHTTP.Routes()...)
	server.RegisterAPIRouters(&routerAPIv1)

	if err := server.Run(ctx); err != nil {
		logger.Error("HTTP-server run error", zap.Error(err))
	}
}
