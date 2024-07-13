package handlers

import (
    "github.com/gofiber/fiber/v2"

    "github.com/satyambaran/oauth/server/clients/config"
    "github.com/satyambaran/oauth/server/clients/helper"
    "github.com/satyambaran/oauth/server/clients/structs"
    "golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    type request struct {
        Name     string
        Email    string
        Password string
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
    client := &structs.Client{
        Name:     body.Name,
        Email:    body.Email,
        Password: string(hash),
        Salt:     salt,
    }
    client, err = helper.CreateClient(config.DB, client)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create client"})
    }
    if err := db.Where("email = ?", body.Email).First(&client).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create client"})
    }
    return c.JSON(client)
}
