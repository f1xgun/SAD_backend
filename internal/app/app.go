package app

import (
	"log"
	"sad/internal/config"
	middleware "sad/internal/middlewares/auth"
	routes "sad/internal/routes/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type App struct {
	serviceProvider *serviceProvider

	router *fiber.App
}

func NewApp() (*App, error) {
	app := &App{}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("Failed to load environment variables: %s", err.Error())
	}

	app.initDeps(config)

	return app, nil
}

func (a *App) initDeps(config config.Config) {
	a.serviceProvider = newServiceProvider()
	a.router = a.setupRouter(config)
}

func (a *App) setupRouter(config config.Config) *fiber.App {
	r := fiber.New()

	r.Use(recover.New())
	r.Use(requestid.New())
	r.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	authHandler := a.serviceProvider.AuthHandler()

	routes.AuthRoutes(r, authHandler)

	api := r.Group("/api").Use(middleware.AuthMiddleware(config))

	api.Post("/example", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to SAD",
		})
	})

	return r
}

func (a *App) Run() error {
	err := a.router.Listen(":8080")
	return err
}
