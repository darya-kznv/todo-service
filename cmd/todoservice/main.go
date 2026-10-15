package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_clock "github.com/darya-kznv/todo-service/internal/core/clock"
	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_pgx_pool "github.com/darya-kznv/todo-service/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/darya-kznv/todo-service/internal/core/transport/http/middleware"
	core_http_server "github.com/darya-kznv/todo-service/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/darya-kznv/todo-service/internal/features/tasks/repository/postgres"
	tasks_service "github.com/darya-kznv/todo-service/internal/features/tasks/service"
	tasks_transport_http "github.com/darya-kznv/todo-service/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/darya-kznv/todo-service/internal/features/users/repository/postgres"
	users_service "github.com/darya-kznv/todo-service/internal/features/users/service"
	users_transport_http "github.com/darya-kznv/todo-service/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("aplication debug time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	clock := core_clock.NewUTCClock()
	usersTransportHTTP := initUsers(pool, logger, clock)
	tasksTransportHTTP := initTasks(pool, logger, clock)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}

func initUsers(
	pool *core_pgx_pool.Pool,
	logger *core_logger.Logger,
	clock core_clock.Clock,
) *users_transport_http.UsersHTTPHandler {
	logger.Debug("initializing feature", zap.String("feature", "users"))

	repo := users_postgres_repository.NewUsersRepository(pool)
	service := users_service.NewUsersService(repo, clock)

	return users_transport_http.NewUsersHTTPHandler(service)
}

func initTasks(
	pool *core_pgx_pool.Pool,
	logger *core_logger.Logger,
	clock core_clock.Clock,
) *tasks_transport_http.TasksHTTPHandler {
	logger.Debug("initializing feature", zap.String("feature", "tasks"))

	repo := tasks_postgres_repository.NewTasksRepository(pool)
	service := tasks_service.NewTasksService(repo, clock)

	return tasks_transport_http.NewTasksHTTPHandler(service)
}
