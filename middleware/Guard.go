package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
)

func API(c fiber.Ctx) error {
	apiKey := c.GetHeaders()["X-Token"][0]
	token := os.Getenv("FRONTEND_API_KEY")
	fmt.Printf("%s '=' %s", apiKey, token)
	if len(token) <= 0 {

		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.serverError"})
		// log.Fatal("Error while loading environment variable")
	}
	if apiKey != token {
		fmt.Printf("issue")

		c.Status(http.StatusUnauthorized)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.unAuthorized"})
	}
	return c.Next()
}
