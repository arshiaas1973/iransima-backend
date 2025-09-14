package v1

import (
	"net/http"

	"github.com/arshiaas1973/iransima-backend/middleware"
	"github.com/arshiaas1973/iransima-backend/orm"
	"github.com/gofiber/fiber/v3"
)

func Init(engine *fiber.App) {
	v1 := engine.Group("/api/v1", middleware.API)
	v1.Get("/users", middleware.Member, GetUser)
}

func GetUser(c fiber.Ctx) error {
	user := c.Locals("AuthUser").(orm.User)
	type Map map[string]interface{}
	c.Status(http.StatusOK)
	return c.JSON(Map{"status": "success", "result": user})
}
