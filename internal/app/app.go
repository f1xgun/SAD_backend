package app

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sad/internal/config"
	"sad/internal/routes/auth"
	"sad/internal/routes/grades"
	"sad/internal/routes/groups"
	"sad/internal/routes/subjects"
	users "sad/internal/routes/user"

	authMiddlewares "sad/internal/middlewares/auth"
	"sad/internal/middlewares/logger"

	usersMiddlewares "sad/internal/middlewares/users"
	usersModels "sad/internal/models/users"
	slogpretty "sad/internal/utils/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type App struct {
	serviceProvider *serviceProvider

	config config.Config

	logger *slog.Logger

	routerChi *chi.Mux
}

func NewApp() (*App, error) {
	app := &App{}

	loadedConfig, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("Failed to load environment variables: %s", err.Error())
	}

	app.config = loadedConfig

	app.logger = setupLogger(app.config.Environment)

	app.initDeps()

	return app, nil
}

func (a *App) CloseDBConnection() {
	a.serviceProvider.db.Close()
}

func (a *App) RunServer() error {
	a.logger.Info("starting server", slog.String("address", a.config.Address))
	srv := &http.Server{
		Addr:         a.config.Address,
		Handler:      a.routerChi,
		ReadTimeout:  a.config.Timeout,
		WriteTimeout: a.config.Timeout,
		IdleTimeout:  a.config.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		a.logger.Error("failed to start server")
		return err
	}

	a.logger.Error("server stopped")
	return nil
}

func (a *App) initDeps() {
	serviceProvider, err := newServiceProvider(a.config, a.logger)

	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
		return
	}

	a.serviceProvider = serviceProvider
	a.routerChi = a.setupRouterChi()
}

type ValidationError struct {
	Message string `json:"message"`
}

func (a *App) setupRouterChi() *chi.Mux {
	router := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: a.config.AllowedOrigins,
		Debug:          a.config.Environment == config.Local,
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	router.Use(c.Handler)

	router.Use(middleware.RequestID)
	router.Use(logger.New(a.logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		err := json.NewEncoder(w).Encode(ValidationError{
			Message: "Not found",
		})

		if err != nil {
			return
		}
	})

	setupApiRoutes(router, a.serviceProvider, a.config, a.logger)

	return router
}

func setupApiRoutes(router *chi.Mux, serviceProvider *serviceProvider, config config.Config, logger *slog.Logger) {
	apiRouter := chi.NewRouter()

	authHandler := serviceProvider.NewAuthHandler()

	auth.Routes(apiRouter, authHandler)

	usersHandler := serviceProvider.NewUserHandler()

	authMiddleware := authMiddlewares.NewAuthMiddleware(config, logger)

	adminMiddleware := usersMiddlewares.AllowedRoleMiddleware(
		serviceProvider.userService,
		[]usersModels.UserRole{usersModels.Admin},
		logger,
	)

	teacherAndAdminMiddleware := usersMiddlewares.AllowedRoleMiddleware(
		serviceProvider.userService,
		[]usersModels.UserRole{usersModels.Admin, usersModels.Teacher},
		logger,
	)

	users.Routes(apiRouter, usersHandler, authMiddleware, adminMiddleware)

	gradesHandler := serviceProvider.NewGradesHandler()

	grades.Routes(apiRouter, gradesHandler, authMiddleware, teacherAndAdminMiddleware)

	groupsHandler := serviceProvider.NewGroupsHandler()

	groups.Routes(apiRouter, groupsHandler, authMiddleware, adminMiddleware)

	subjectsHandler := serviceProvider.NewSubjectsHandler()

	subjects.Routes(apiRouter, subjectsHandler, authMiddleware, adminMiddleware)

	router.Mount("/api", apiRouter)
}

func setupLogger(env config.Env) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.Local:
		log = setupPrettySlog()
	case config.Dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.Prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
