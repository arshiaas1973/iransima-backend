package middleware

import (
	"fmt"
	"net/http"
	"os"

	orm "github.com/arshiaas1973/iransima-backend/orm/models"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserPayload struct {
	jwt.Payload
	ID        uint
	Email     string
	FirstName string
	LastName  string
}

func Member(c fiber.Ctx) error {
	// return func(c *gin.Context) {
	cookie := c.Cookies("User", "")
	if cookie == "" {
		c.Status(http.StatusUnauthorized)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.unAuthorized"})
	}
	secret := os.Getenv("SECRET_KEY")
	if len(secret) <= 0 {
		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.serverError"})
	}
	var pl UserPayload
	alg := jwt.NewHS256([]byte(secret))
	_, err := jwt.Verify([]byte(secret), alg, &pl)
	if err != nil {
		c.Cookie(&fiber.Cookie{
			Name:     "User",
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			Domain:   os.Getenv("APP_DOMAIN"),
			HTTPOnly: true,
			Secure:   true,
		})
		fmt.Println(err)
		// c.SetCookie("User", "", -1, "/", "localhost", true, true)
		c.Status(http.StatusUnauthorized)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.unAuthorized"})
	}
	dsn := os.Getenv("DB_URL")
	if len(dsn) <= 0 {
		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.serverError"})
		// log.Fatal("Error while loading environment variable")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]interface{}{"status": "failed", "result": "errors.serverError"})
		// log.Fatal("Error while loading connecting to database")
	}
	user := orm.User{
		Model:     gorm.Model{ID: uint(pl.ID)},
		Email:     pl.Email,
		FirstName: pl.FirstName,
		LastName:  pl.LastName,
	}
	result := db.Model(&orm.User{}).Omit("Password", "ResetToken").Find(&user)
	if result.RowsAffected <= 0 {
		dbi, err := db.DB()
		if err != nil {
			fmt.Println(err)
		}
		dbi.Close()
		c.Cookie(&fiber.Cookie{
			Name:     "User",
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			Domain:   os.Getenv("APP_DOMAIN"),
			HTTPOnly: true,
			Secure:   true,
		})
		c.Status(http.StatusUnauthorized)
		type Map map[string]interface{}
		return c.JSON(map[string]interface{}{"status": "failed", "result": Map{"icon": "material-symbols:person-cancel-outline-rounded", "content": "errors.unAuthorized"}})
	}
	c.Locals("AuthUser", user)
	return c.Next()
	// }
}

func Guest(c fiber.Ctx) error {
	cookie := c.Cookies("User", "")
	if cookie == "" {
		return c.Next()
	}
	c.Status(http.StatusBadRequest)
	type Map map[string]interface{}
	return c.JSON(map[string]interface{}{"status": "failed", "result": Map{"icon": "material-symbols:person-check-outline-rounded", "content": "errors.alreadyAuthorized"}})
}
