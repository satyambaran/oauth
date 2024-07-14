package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/clients/helper"
)

func Middleware() fiber.Handler { //imp  small first letter in a package makes it private
    return func(c *fiber.Ctx) error {
        authToken := c.Get("Authorization")
        if authToken == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
        }
        claims, msg := helper.ValidateToken(authToken)
        if msg != "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
        }
        c.Locals("id", claims.ID)
        c.Locals("client_id", claims.ClientID)
        return c.Next()
    }
}
func UserMiddleware() fiber.Handler { //imp  small first letter in a package makes it private
    return func(c *fiber.Ctx) error {
        authToken := c.Get("Authorization")
        if authToken == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
        }
        claims, msg := helper.ValidateUserToken(authToken)
        if msg != "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
        }
        c.Locals("id", claims.ID)
        // c.Locals("name", claims.Name)
        // c.Locals("role", claims.Role)
        // c.Locals("email", claims.Email)
        return c.Next()
    }
}
