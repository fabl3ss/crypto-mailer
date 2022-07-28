package main

import (
	"genesis_test_case/src/config"
	"genesis_test_case/src/pkg/delivery/http/middleware"
	"genesis_test_case/src/pkg/routes"
	"genesis_test_case/src/pkg/utils"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Define Fiber config.
	config := config.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	if err := routes.InitRoutes(app); err != nil {
		panic("Unable to initialize routes")
	}

	// Start server (with graceful shutdown).
	utils.StartServerWithGracefulShutdown(app)
}
