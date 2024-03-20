package app

import (
	routes "sad/internal/routes/auth"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	serviceProvider *serviceProvider

	router *fiber.App
}

func NewApp() (*App, error) {
	app := &App{}

	app.initDeps()

	return app, nil
}

func (a *App) initDeps() {
	a.serviceProvider = newServiceProvider()
	a.router = a.setupRouter()
}

func (a *App) setupRouter() *fiber.App {
	r := fiber.New()

	authHandler := a.serviceProvider.AuthHandler()

	routes.AuthRoutes(r, authHandler)

	return r
}

func (a *App) Run() {
	a.router.Listen(":8080")
}
