package main

import (
	"log"
	"os"

	aiHandler "github.com/haramurti/-KWU--food-web-be/internal/app/gemini/handler"
	aiService "github.com/haramurti/-KWU--food-web-be/internal/app/gemini/service"

	"github.com/haramurti/-KWU--food-web-be/pkg/config"
	"github.com/haramurti/-KWU--food-web-be/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading env vars directly")
	}

	cfg := config.New()

	// ── Dependency Injection ──────────────────────────────────────────
	geminiSvc := aiService.NewGeminiService(cfg.GeminiAPIKey)
	aiH := aiHandler.NewAIHandler(geminiSvc)

	// ── App ───────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName: "NutriBar API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.InternalError(c, err.Error())
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST",
	}))

	// ── Routes ────────────────────────────────────────────────────────
	api := app.Group("/api")

	ai := api.Group("/ai")
	ai.Post("/chat", aiH.Chat)

	app.Get("/health", func(c *fiber.Ctx) error {
		return response.OK(c, "Server is running", nil)
	})

	// ── Start ─────────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
