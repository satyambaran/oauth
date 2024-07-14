package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/clients/config"
    "github.com/satyambaran/oauth/server/clients/helper"
    oAuthHelper "github.com/satyambaran/oauth/server/clients/helper/oAuthHelpers"
    "github.com/satyambaran/oauth/server/clients/structs"
)

var ctx = context.Background()

func Consent(c *fiber.Ctx, config *config.Config) error {
    userID := c.Locals("id").(int)
    email := c.Locals("email").(string)
    role := c.Locals("role").(string)
    name := c.Locals("name").(string)
    grantType := c.Query("grant_type", "all") // default to "all" if not provided, code
    code := c.Query("code", "1")              // default to "1" if not provided
    clientId := c.Query("client_id", "")      // default to "1" if not provided
    if clientId == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot find client code"})
    }
    redirectURI := c.Query("redirect_uri", "google.com")
    authCode, err := helper.GenerateAuthCode(6)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate auth code"})
    }

    claims := &structs.UserOAuth{
        ID:        userID,
        Name:      name,
        Email:     email,
        Role:      role,
        ClientID:  clientId,
        GrantType: grantType,
        Code:      code,
        AuthCode:  authCode,
    }
    data, err := json.Marshal(claims)
    if err != nil {
        return err
    }
    rdb := config.RDB
    err = rdb.Set(ctx, authCode+fmt.Sprintf("_%v_", userID)+clientId, data, time.Minute*30).Err()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
    }
    return c.JSON(fiber.Map{"authCode": authCode, "uri": redirectURI})
}
func GetTokenFromAuthCode(c *fiber.Ctx, config *config.Config) error {
    type AuthCodeBody struct {
        UserID   int    `json:"user_id"`
        AuthCode string `json:"auth_code"`
    }
    clientId := c.Locals("client_id").(string)
    body := &AuthCodeBody{}
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }
    rdb := config.RDB
    claimsFromRedis, err := rdb.Get(ctx, body.AuthCode+fmt.Sprintf("_%v_", body.UserID)+clientId).Result()
    if err != nil {
        if err == redis.Nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key does not exist"})
        }
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }
    claims := &structs.UserOAuth{}
    err = json.Unmarshal([]byte(claimsFromRedis), claims)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "value does not exist"})
    }
    secret_key := os.Getenv("SECRET_KEY")
    if secret_key == "" {
        secret_key = "secret_key"
    }

    refreshToken, token, err := oAuthHelper.CreateAllTokens(claims, secret_key, config)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"token": token, "refreshToken": refreshToken})
}
