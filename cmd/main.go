package main

import (
	"fiber-gorm-app/internal/config"
	"fiber-gorm-app/internal/routes"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	config.InitDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	routes.SetupRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Log startup info
	log.Printf("Revok's Gym API starting...")
	log.Printf("Server running at: http://localhost:%s", port)
	log.Printf("API Documentation: http://localhost:%s/api/v1", port)
	log.Printf("Environment: %s", os.Getenv("APP_ENV"))
	log.Printf("Database: %s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Show registered routes
	log.Println("\n Registered Routes:")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	routeCount := 0
	for _, stack := range app.Stack() {
		for _, route := range stack {
			if route.Path != "" {
				routeCount++
				log.Printf(" %-7s\t%s",
					route.Method,
					route.Path,
				)
			}
		}
	}

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Total routes: %d\n", routeCount)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}

