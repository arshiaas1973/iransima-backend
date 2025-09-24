package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
)

func API(c fiber.Ctx) error {
	apiKey := c.GetHeaders()["X-Token"]
	if len(apiKey) <= 0 {
		c.Status(http.StatusUnauthorized)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.unAuthorized", "message": "Set the token"})
	}
	key := apiKey[0]
	token := os.Getenv("FRONTEND_API_KEY")
	fmt.Printf("%s '=' %s", key, token)
	if len(token) <= 0 {

		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.serverError"})
		// log.Fatal("Error while loading environment variable")
	}
	if key != token {
		fmt.Printf("issue")

		c.Status(http.StatusUnauthorized)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.unAuthorized", "message": "Token isn't set properly!"})
	}
	return c.Next()
}
