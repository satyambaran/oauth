package main

import (
    "fmt"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/clients/config"
    "github.com/satyambaran/oauth/server/clients/routes"
)

func main() {
    app := fiber.New()
    config := &config.Config{
        DB:  config.InitDatabase(),
        RDB: config.InitRedis(),
    }
    routes.SetupRoutes(app, config)
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    app.Listen(fmt.Sprintf(":%s", port))
}
