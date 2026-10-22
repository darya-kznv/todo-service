package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/darya-kznv/todo-service/docs"
	core_config "github.com/darya-kznv/todo-service/internal/core/config"
	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_pgx_pool "github.com/darya-kznv/todo-service/internal/core/repository/postgres/pool/pgx"
	core_redis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool"
	core_goredis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool/goredis"
	core_http_middleware "github.com/darya-kznv/todo-service/internal/core/transport/http/middleware"
	core_http_server "github.com/darya-kznv/todo-service/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/darya-kznv/todo-service/internal/features/statistics/repository/postgres"
	statistics_service "github.com/darya-kznv/todo-service/internal/features/statistics/service"
	statistics_transport_http "github.com/darya-kznv/todo-service/internal/features/statistics/transport/http"
	tasks_cached_repository "github.com/darya-kznv/todo-service/internal/features/tasks/repository/cached"
	tasks_postgres_repository "github.com/darya-kznv/todo-service/internal/features/tasks/repository/postgres"
	tasks_service "github.com/darya-kznv/todo-service/internal/features/tasks/service"
	tasks_transport_http "github.com/darya-kznv/todo-service/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/darya-kznv/todo-service/internal/features/users/repository/postgres"
	users_service "github.com/darya-kznv/todo-service/internal/features/users/service"
	users_transport_http "github.com/darya-kznv/todo-service/internal/features/users/transport/http"
	web_fs_repository "github.com/darya-kznv/todo-service/internal/features/web/repository/fs"
	web_service "github.com/darya-kznv/todo-service/internal/features/web/service"
	web_transport_http "github.com/darya-kznv/todo-service/internal/features/web/transport/http"
	"go.uber.org/zap"
)

// @title 		Golang Todo API
// @version 	1.0
// @description Todo Application REST API
// @host 		127.0.0.1:5050
// @BasePath 	/api/v1
func main() {
	config := core_config.NewConfigMust()
	time.Local = config.TimeZone

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

	logger.Debug("aplication time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	postgresPool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer postgresPool.Close()

	logger.Debug("initializing Redis connection pool")
	redisPool, err := core_goredis_pool.NewPool(ctx, core_goredis_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init Redis connection pool", zap.Error(err))
	}
	defer func() {
		if err := redisPool.Close(); err != nil {
			logger.Error("failed to close Redis connection pool", zap.Error(err))
		}
	}()

	usersTransportHTTP := initUsers(postgresPool, logger)
	tasksTransportHTTP := initTasks(postgresPool, redisPool, logger)
	statisticsTransportHTTP := initStatistics(postgresPool, logger)
	webTransportHTTP := initWeb(logger)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}

func initWeb(logger *core_logger.Logger) *web_transport_http.WebHTTPHandler {
	logger.Debug("initializing feature", zap.String("feature", "web"))

	repo := web_fs_repository.NewWebRepository()
	service := web_service.NewWebService(repo)

	return web_transport_http.NewWebHTTPHandler(service)
}

func initUsers(
	pool *core_pgx_pool.Pool,
	logger *core_logger.Logger,
) *users_transport_http.UsersHTTPHandler {
	logger.Debug("initializing feature", zap.String("feature", "users"))

	repo := users_postgres_repository.NewUsersRepository(pool)
	service := users_service.NewUsersService(repo)

	return users_transport_http.NewUsersHTTPHandler(service)
}

func initTasks(
	postgresPool *core_pgx_pool.Pool,
	redisPool core_redis_pool.Pool,
	logger *core_logger.Logger,
) *tasks_transport_http.TasksHTTPHandler {
	logger.Debug("initializing feature", zap.String("feature", "tasks"))

	postgresRepo := tasks_postgres_repository.NewTasksRepository(postgresPool)
	cachedRepo := tasks_cached_repository.NewCachedRepository(redisPool, postgresRepo, logger)
	service := tasks_service.NewTasksService(cachedRepo)

	return tasks_transport_http.NewTasksHTTPHandler(service)
}

func initStatistics(
	pool *core_pgx_pool.Pool,
	logger *core_logger.Logger,
) *statistics_transport_http.StatisticsHTTPHandler {
	logger.Debug("initializing feature", zap.String("feature", "statistics"))

	repo := statistics_postgres_repository.NewStatisticsRepository(pool)
	service := statistics_service.NewStatisticsService(repo)

	return statistics_transport_http.NewStatisticsHTTPHandler(service)
}
