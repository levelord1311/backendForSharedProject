package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/config"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/handlers"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot/db"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/metric"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/mysql"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("initializing config...")
	cfg := config.GetConfig()

	logger.Println("initializing router..")
	router := httprouter.New()

	logger.Println("initializing helpers..")
	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	logger.Println("initializing database..")
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.MysqlDB.Username,
		cfg.MysqlDB.Password,
		cfg.MysqlDB.Host,
		cfg.MysqlDB.Port,
		cfg.MysqlDB.DBName)
	mysqlClient, err := mysql.NewClient(logger, dbConnString)
	if err != nil {
		logger.Fatalln(err)
	}

	lotStorage := db.NewStorage(mysqlClient, logger)
	lotService, err := lot.NewService(lotStorage, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("initializing handlers..")
	lotsHandler := handlers.Handler{
		Logger:     logger,
		LotService: lotService,
	}
	lotsHandler.Register(router)

	logger.Println("starting application...")
	start(router, logger, cfg)

}

func start(router http.Handler, logger logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Infof("socket path: %s", socketPath)

		logger.Info("creating and starting to listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("bind application to host: %s and port: %s\n", cfg.Listen.BindIP, cfg.Listen.Port)

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
			logger.Warn("server is shutting down")
		default:
			logger.Fatal(err)
		}
	}
}
