package main

import (
	"fmt"
	"log"
	"os"

	users "github.com/arshiaas1973/iransima-backend/api/v1"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Backend Running!!!")
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "Iransima",
	})
	app.Use(csrf.New(csrf.Config{
		CookieName: "xsrf_protection",
		Extractor: csrf.Chain(
			csrf.FromForm("_token"),
			csrf.FromHeader("X-XSRF-Token"),
		),
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	}))
	app.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowOrigins:     []string{"https://localhost:3000", "https://localhost:3000"},
		AllowHeaders:     []string{"X-XSRF-Token", "X-Token"},
		AllowCredentials: true,
	}))
	app.Use(logger.New(logger.ConfigDefault))
	users.Init(app)
	domain := os.Getenv("APP_DOMAIN_WITHOUT_PROTOCOL")
	if domain == "" {
		domain = "localhost:4000"
	}
	app.Listen(domain, fiber.ListenConfig{
		CertFile:          "./certs/cert.pem",
		CertKeyFile:       "./certs/key.pem",
		EnablePrefork:     true,
		EnablePrintRoutes: true,
	})
}
