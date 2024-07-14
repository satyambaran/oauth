package userHandlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/satyambaran/oauth/server/users/config"
	"github.com/satyambaran/oauth/server/users/helper"
	"github.com/satyambaran/oauth/server/users/structs"
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
    var user structs.User
    if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
    }
    salt, err := helper.GenerateSalt()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate salt"})
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword+salt), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
    }

    user.Password = string(hash)
    user.Salt = salt
    db.Save(&user)

    return c.JSON(user)
}
func ResetPasswordV2(c *fiber.Ctx, config *config.Config) error {
    return nil
}
