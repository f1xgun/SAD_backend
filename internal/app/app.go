package app

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"sad/internal/config"
	middleware "sad/internal/middlewares/auth"
	"sad/internal/middlewares/users"
	usersModels "sad/internal/models/users"
	"sad/internal/routes/auth"
	"sad/internal/routes/grades"
	"sad/internal/routes/groups"
	"sad/internal/routes/subjects"
	usersRoutes "sad/internal/routes/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type App struct {
	serviceProvider *serviceProvider

	router *fiber.App

	config config.Config
}

func NewApp() (*App, error) {
	app := &App{}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("Failed to load environment variables: %s", err.Error())
	}

	app.config = config

	app.initDeps()

	return app, nil
}

func (a *App) initDeps() {
	serviceProvider, err := newServiceProvider(a.config)

	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
		return
	}

	a.serviceProvider = serviceProvider
	a.router = a.setupRouter()
}

func (a *App) setupRouter() *fiber.App {
	r := fiber.New()

	r.Use(cors.New())
	r.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("Unhandled error occurred: %v", e)
			c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	}))
	r.Use(logger.New())

	authHandler := a.serviceProvider.NewAuthHandler()

	auth.AuthRoutes(r, authHandler)

	userHandler := a.serviceProvider.NewUserHandler()

	usersRoutes.UserRoutes(
		r,
		userHandler,
		middleware.NewAuthMiddleware(a.config),
		users.AllowedRoleMiddleware(a.serviceProvider.userService, []usersModels.UserRole{usersModels.Admin}),
	)

	groupsHandler := a.serviceProvider.NewGroupsHandler()

	groups.GroupsRoutes(
		r,
		groupsHandler,
		middleware.NewAuthMiddleware(a.config),
		users.AllowedRoleMiddleware(a.serviceProvider.userService, []usersModels.UserRole{usersModels.Admin}),
	)

	subjectsHandler := a.serviceProvider.NewSubjectsHandler()

	subjects.SubjectsRoutes(
		r,
		subjectsHandler,
		middleware.NewAuthMiddleware(a.config),
		users.AllowedRoleMiddleware(a.serviceProvider.userService, []usersModels.UserRole{usersModels.Admin}),
	)

	gradesHandler := a.serviceProvider.NewGradesHandler()

	grades.GradesRoutes(
		r,
		gradesHandler,
		middleware.NewAuthMiddleware(a.config),
		users.AllowedRoleMiddleware(a.serviceProvider.userService, []usersModels.UserRole{usersModels.Admin, usersModels.Teacher}),
	)

	return r
}

func (a *App) Run() error {
	err := a.router.Listen(":8080")
	return err
}

func (a *App) CloseDBConnection() {
	a.serviceProvider.db.Close()
}
