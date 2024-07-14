package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/users/config"
    "github.com/satyambaran/oauth/server/users/handlers/handlers"
    "github.com/satyambaran/oauth/server/users/handlers/resourceHandlers"
    "github.com/satyambaran/oauth/server/users/handlers/userHandlers"
    middleware "github.com/satyambaran/oauth/server/users/middlewares"
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
    UserRoutersV1(v, config)

    // Resources
    ResourceRoutersV1(v, config)
}
func SetupV2Routes(api fiber.Router, config *config.Config) {
    // API: /api/v1
    v := api.Group("/v2")

    // User
    UserRoutersV2(v, config)

    // Resources
    ResourceRoutersV2(v, config)
}
func UserRoutersV1(v fiber.Router, config *config.Config) {
    v.Post("/register", func(c *fiber.Ctx) error { return userHandlers.Register(c, config) })
    v.Post("/login", func(c *fiber.Ctx) error { return userHandlers.Login(c, config) })
    v.Post("/reset", func(c *fiber.Ctx) error { return userHandlers.ResetPassword(c, config) })
    v.Post("/refresh", func(c *fiber.Ctx) error { return userHandlers.RefreshToken(c, config) })
}
func UserRoutersV2(v fiber.Router, config *config.Config) {
    v.Post("/register", func(c *fiber.Ctx) error { return userHandlers.Register(c, config) })
    v.Post("/login", func(c *fiber.Ctx) error { return userHandlers.Login(c, config) })
    v.Post("/reset", func(c *fiber.Ctx) error { return userHandlers.ResetPassword(c, config) })
    v.Post("/refresh", func(c *fiber.Ctx) error { return userHandlers.RefreshToken(c, config) })
}
func ResourceRoutersV1(v fiber.Router, config *config.Config) {
    v.Get("/resources", middleware.Middleware(), func(c *fiber.Ctx) error { return resourceHandlers.Resources(c, config) })
    v.Post("/resource", middleware.Middleware(), func(c *fiber.Ctx) error { return resourceHandlers.Resource(c, config) })
    v.Delete("/resource/:id", middleware.Middleware(), func(c *fiber.Ctx) error { return resourceHandlers.DeleteResource(c, config) })
}

func ResourceRoutersV2(v fiber.Router, config *config.Config) {
    v.Get("/resources", middleware.Middleware(), func(c *fiber.Ctx) error { return resourceHandlers.Resources(c, config) })
    v.Post("/resource", middleware.Middleware(), func(c *fiber.Ctx) error { return resourceHandlers.Resource(c, config) })
    v.Delete("/resource/:id", middleware.Middleware(), func(c *fiber.Ctx) error { return resourceHandlers.DeleteResource(c, config) })
}
