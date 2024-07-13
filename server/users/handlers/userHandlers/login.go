package userHandlers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/satyambaran/oauth/server/users/config"
	"github.com/satyambaran/oauth/server/users/helper"
	"github.com/satyambaran/oauth/server/users/structs"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    type request struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    var body request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    var user structs.User
    if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password+user.Salt), []byte(body.Password)); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect password"})
    }
    // Generate JWT token
    secret_key := os.Getenv("SECRET_KEY")
    if secret_key == "" {
        secret_key = "secret_key"
    }
    refreshToken, token, err := helper.CreateAllTokens(&user, secret_key, config)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"token": token, "refreshToken": refreshToken})
}
