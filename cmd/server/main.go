package main

import (
	"context"
	"errors"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nnemirovsky/pacgen/internal/handler"
	"github.com/nnemirovsky/pacgen/internal/repository"
	"github.com/nnemirovsky/pacgen/internal/router"
	"github.com/nnemirovsky/pacgen/internal/service"
	"github.com/nnemirovsky/pacgen/pkg/logutil"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	opts           *options
	logger         zerolog.Logger
	db             *sqlx.DB
	server         *http.Server
	ruleRepo       *repository.RuleRepository
	profileRepo    *repository.ProxyProfileRepository
	ruleService    *service.RuleService
	profileService *service.ProxyProfileService
	pacService     *service.PACService
	ruleHandler    *handler.RuleHandler
	profileHandler *handler.ProxyProfileHandler
	mux            http.Handler
)

type options struct {
	LogLevel string `short:"l" long:"loglevel" env:"APP_LOG_LEVEL" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" description:"Log level" default:"info"`
	Port     int    `short:"p" long:"port" env:"APP_PORT" description:"Http port to listen on" default:"8080"`
	User     string `short:"U" long:"user" env:"APP_USER" description:"User for http basic auth" default:"admin"`
	Password string `short:"P" long:"password" env:"APP_PASSWORD" description:"Password for http basic auth" default:"admin"`
}

func main() {
	initOpts()
	initLogger()
	initDB()
	initRepositories()
	initServices()
	initHandlers()
	initRouter()
	initServer()

	logger.Info().Str("addr", server.Addr).Msg("Application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-quit

	logger.Info().Msg("Application is shutting down...")

	shutdownServer()
	shutdownDB()
}

func initRouter() {
	mux = router.New(ruleHandler, profileHandler, logger, map[string]string{opts.User: opts.Password})
}

func initHandlers() {
	ruleHandler = handler.NewRuleHandler(ruleService, logutil.WithLayer[handler.RuleHandler](logger))
	profileHandler = handler.NewProxyProfileHandler(profileService, logutil.WithLayer[handler.ProxyProfileHandler](logger))
}

func initServices() {
	pacService = service.NewPACService(ruleRepo, logutil.WithLayer[service.PACService](logger))
	ruleService = service.NewRuleService(ruleRepo, pacService, logutil.WithLayer[service.RuleService](logger))
	profileService = service.NewProxyProfileService(profileRepo, pacService, logutil.WithLayer[service.ProxyProfileService](logger))
}

func initRepositories() {
	ruleRepo = repository.NewRuleRepository(db, logutil.WithLayer[repository.RuleRepository](logger))
	profileRepo = repository.NewProxyProfileRepository(db, logutil.WithLayer[repository.ProxyProfileRepository](logger))
}

func initOpts() {
	opts = &options{}
	if _, err := flags.Parse(opts); err != nil {
		os.Exit(1)
	}
}

func initLogger() {
	logger = logutil.Logger
	lvl, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse log level")
	}
	logger = logger.Level(lvl)
}

func initDB() {
	var err error
	if db, err = sqlx.Connect("sqlite3", "./data/data.db"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to open db connection")
	}
	db.MustExec("PRAGMA foreign_keys = ON")
}

func shutdownDB() {
	if err := db.Close(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to close db connection")
	}
}

func initServer() {
	server = &http.Server{
		Addr:    ":" + strconv.Itoa(opts.Port),
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Error occurred while running http server")
		}
	}()
}

func shutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Error occurred while shutting down server")
	}
}
