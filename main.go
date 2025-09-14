package main

import (
	"fmt"

	users "github.com/arshiaas1973/iransima-backend/api/v1"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
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
		AllowOrigins:     []string{"https://localhost:3000"},
		AllowHeaders:     []string{"X-XSRF-Token"},
		AllowCredentials: true,
	}))
	app.Use(logger.New(logger.ConfigDefault))
	users.Init(app)
	app.Listen(":4000",fiber.ListenConfig{
		CertFile: "./certs/cert.pem",
		CertKeyFile: "./certs/key.pem",
		EnablePrefork: true,
		EnablePrintRoutes: true,
	})
}
