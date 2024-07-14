package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/satyambaran/oauth/server/clients/config"
	"github.com/satyambaran/oauth/server/clients/helper"
	"github.com/satyambaran/oauth/server/clients/structs"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    type request struct {
        OTP         string `json:"otp"`
        Email       string `json:"email"`
        NewPassword string `json:"new_password"`
    }

    var body request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    // verify OTP
    if body.OTP == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "incorrect found"})
    }
    // update password
    var client structs.Client
    if err := db.Where("email = ?", body.Email).First(&client).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "client not found"})
    }
    salt, err := helper.GenerateSalt()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate salt"})
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword+salt), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
    }

    client.Password = string(hash)
    client.Salt = salt
    db.Save(&client)

    return c.JSON(client)
}
func ResetPasswordV2(c *fiber.Ctx, config *config.Config) error {
    return nil
}
