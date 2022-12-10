package service

import (
	"api/internal/adapters"
	"api/internal/application"
	"api/internal/application/command"
	"api/internal/application/query"
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	App                string
	Address            string
	DBConnectionString string
	LogLevel           string
}

type Service struct {
	config *config
	logger *zerolog.Logger
	db     *sql.DB
	app    *application.Application
	server *http.Server
}

func (s *Service) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s.loadConfig()
	s.initLogger()

	s.logInfo("application is starting up...")

	if err := s.initDB(ctx); err != nil {
		s.logFatal(err)
	}

	s.initApplication()
	s.initHTTPServer(ctx)

	if err := s.serve(ctx); err != nil {
		s.logFatal(err)
	}

	s.logInfo("application is shut down")
}

func (s *Service) loadConfig() {
	s.config = &config{
		App:                env("APP_NAME", "api"),
		Address:            env("APP_ADDRESS", ":8080"),
		DBConnectionString: env("APP_DB_CONNECTION_STRING", "./items.db"),
		LogLevel:           env("APP_LOG_LEVEL", "warn"),
	}
}

func (s *Service) initLogger() {
	level, err := zerolog.ParseLevel(s.config.LogLevel)
	if err != nil {
		level = zerolog.WarnLevel
	}

	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Str("app", s.config.App).
		Logger()

	s.logger = &logger
}

func (s *Service) initDB(ctx context.Context) error {
	if _, err := os.Stat(s.config.DBConnectionString); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", s.config.DBConnectionString)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Service) initApplication() {
	r := adapters.NewSQLiteItemRepository(s.db)

	s.app = &application.Application{
		Commands: application.Commands{
			CreateItem: command.NewCreateItemHandler(r),
		},
		Queries: application.Queries{
			GetItems: query.NewGetItemsHandler(r),
		},
	}
}

func (s *Service) initHTTPServer(ctx context.Context) {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	handlers := adapters.NewHTTPHandlers(s.app, s.logger)

	router.MethodFunc(http.MethodGet, "/items", handlers.GetItems)
	router.MethodFunc(http.MethodPost, "/items", handlers.CreateItem)

	router.NotFound(adapters.NotFound)

	s.server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         s.config.Address,
		Handler:      router,
	}
}

func (s *Service) serve(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		<-ctx.Done()

		s.logInfo("application is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = s.server.Shutdown(ctx)
		_ = s.db.Close()

		errChan <- nil
	}()

	s.logInfo("application is running at %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return <-errChan
}
