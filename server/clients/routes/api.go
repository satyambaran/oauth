package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/clients/config"
    clientHandlers "github.com/satyambaran/oauth/server/clients/handlers/clientHandlers"
    "github.com/satyambaran/oauth/server/clients/handlers/handlers"
    // server/clients/handlers/oauthHandlers/oauthHandlers.go
    oauthHandlers "github.com/satyambaran/oauth/server/clients/handlers/oauthHandlers"
    middleware "github.com/satyambaran/oauth/server/clients/middlewares"
)

func SetupRoutes(app *fiber.App, config *config.Config) {
    app.Get("/", handlers.Welcome)
    api := app.Group("/api")
    SetupV1Routes(api, config)
    SetupV2Routes(api, config)
}
func SetupV1Routes(api fiber.Router, config *config.Config) {
    // API: /api/v1
    v := api.Group("/v1")
    // User
    ClientRoutersV1(v, config)
    OAuthRoutesV1(v, config)
}
func SetupV2Routes(api fiber.Router, config *config.Config) {
    // API: /api/v2
    v := api.Group("/v1")
    // User
    ClientRoutersV2(v, config)
    OAuthRoutesV2(v, config)
}

func ClientRoutersV1(v fiber.Router, config *config.Config) {
    v.Post("/register", func(c *fiber.Ctx) error { return clientHandlers.Register(c, config) })
    v.Post("/login", func(c *fiber.Ctx) error { return clientHandlers.Login(c, config) })
    v.Post("/reset", func(c *fiber.Ctx) error { return clientHandlers.ResetPassword(c, config) })
    v.Post("/refresh", func(c *fiber.Ctx) error { return clientHandlers.RefreshToken(c, config) })
}
func ClientRoutersV2(v fiber.Router, config *config.Config) {
    v.Post("/register", func(c *fiber.Ctx) error { return clientHandlers.Register(c, config) })
    v.Post("/login", func(c *fiber.Ctx) error { return clientHandlers.Login(c, config) })
    v.Post("/reset", func(c *fiber.Ctx) error { return clientHandlers.ResetPassword(c, config) })
    v.Post("/refresh", func(c *fiber.Ctx) error { return clientHandlers.RefreshToken(c, config) })
}
func OAuthRoutesV1(v fiber.Router, config *config.Config) {
    v.Post("/consent",middleware.UserMiddleware(), func(c *fiber.Ctx) error { return oauthHandlers.Consent(c, config) })
}
func OAuthRoutesV2(v fiber.Router, config *config.Config) {
}
