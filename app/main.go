package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "DRM Core v1.0.0",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/request", handleRequest)

	log.Println("Starting DRM (Declarative-Relation Mapping) Core server on :8080")
	log.Fatal(app.Listen(":8080"))
}

func handleRequest(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "DRM request handler",
		"status":  "ready",
	})
}