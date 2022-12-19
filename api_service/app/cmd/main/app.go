package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/client/lot_service"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/client/user_service"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/config"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/handlers/auth"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/handlers/lots"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/jwt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/metric"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

// TODO Makefile
// TODO swagger
// TODO отрефакторить из монолита в микросервисы
// TODO сделать ответы микросервисов доступными только для запросов от api service
// TODO JWT refresh токены, хотя бы в кеше (попробовать в redis?). Или сгенерить как JWT?
/*
   "Зачем держать в кеше refresh_token?
   При рестарте в течение часа 100% пользователей разлогинит.
   Его же точно так же можно сгенерить как jwt и не хранить вообще ничего"
*/
// TODO тесты

func main() {

	// TODO "Пионерский код. Логер без закрытия в defer c потерей данных логирования" - проверить, исправить
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("initializing config...")
	cfg := config.GetConfig()

	logger.Println("initializing router...")
	router := httprouter.New()

	logger.Println("initializing helpers...")
	jwtHelper := jwt.NewHelper(logger)

	logger.Println("creating and registering handlers...")

	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	userService := user_service.NewService(cfg.UserService.URL, "/users", logger)
	authHandler := auth.Handler{JWTHelper: jwtHelper, UserService: userService, Logger: logger}
	authHandler.Register(router)

	lotService := lot_service.NewService(cfg.LotService.URL, "/lots", logger)
	lotsHandler := lots.Handler{LotService: lotService, Logger: logger}
	lotsHandler.Register(router)

	logger.Println("starting application...")
	start(router, logger, cfg)

	//if err := main.StartMainHTTP(cfg); err != nil {
	//	logger.Errorln("error starting http server: ", err)
	//	os.Exit(2)
	//}
}

func start(router *httprouter.Router, logger logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Infof("socket path: %s", socketPath)

		logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

		var err error

		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}