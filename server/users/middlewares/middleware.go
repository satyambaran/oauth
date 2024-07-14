package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/users/helper"
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
        return c.Next()
    }
}
func ClientMiddleware() fiber.Handler { //imp  small first letter in a package makes it private
    return func(c *fiber.Ctx) error {
        return c.Next()
    }
}
func UserOAuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        return c.Next()
    }
}
