package main

import (
	"log"

	"drm-app/app/drm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var engine *drm.Engine

func main() {
	engine = drm.NewEngine()

	app := fiber.New(fiber.Config{
		AppName: "DRM Core v1.0.0",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/request", handleRequest)

	log.Println("Starting DRM (Declarative-Relation Mapping) Core server on :8080")
	log.Fatal(app.Listen(":8080"))
}

type RequestBody struct {
	Query string `json:"query"`
	Token string `json:"token"`
}

func handleRequest(c *fiber.Ctx) error {
	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	result, err := engine.ProcessRequest(c.Context(), req.Query, req.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
		"status": "success",
	})
}
