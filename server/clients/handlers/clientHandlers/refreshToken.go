package handlers

import (
    "os"
    "strconv"

    "github.com/gofiber/fiber/v2"

    "github.com/satyambaran/oauth/server/clients/config"
    "github.com/satyambaran/oauth/server/clients/helper"
    "github.com/satyambaran/oauth/server/clients/structs"
)

func RefreshToken(c *fiber.Ctx, config *config.Config) error {
    type request struct {
        ID int `json:"id"`
    }
    var body request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }

    var client *structs.Client
    if err := config.DB.Where("id = ?", body.ID).First(&client).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "client not found"})
    }
    rdb := config.RDB
    refreshToken, err := helper.GetRefreshToken(rdb, client)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    _, msg := helper.ValidateRefreshToken(refreshToken)
    if msg != "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
    }
    secret_key := os.Getenv("SECRET_KEY")
    if secret_key == "" {
        secret_key = "secret_key"
    }
    atExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTE"))
    if err != nil {
        atExp = 30
    }
    token, err := helper.CreateAccessToken(client, secret_key, atExp)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
    }
    err = helper.SaveAccessToken(rdb, client, token, atExp)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
