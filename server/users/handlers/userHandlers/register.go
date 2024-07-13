package userHandlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/satyambaran/oauth/server/users/config"
	"github.com/satyambaran/oauth/server/users/helper"
	"github.com/satyambaran/oauth/server/users/structs"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    type request struct {
        Name     string
        Email    string
        Password string
        Role     string
    }

    var body request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    salt, err := helper.GenerateSalt()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate salt"})
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password+salt), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
    }

    user := structs.User{
        Name:     body.Name,
        Email:    body.Email,
        Role:     body.Role,
        Password: string(hash),
        Salt:     salt,
    }

    if err := db.Create(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
    }
    if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
    }
    return c.JSON(user)
}
