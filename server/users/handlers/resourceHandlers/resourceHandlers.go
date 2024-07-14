package resourceHandlers

import (
    "fmt"

    "github.com/gofiber/fiber/v2"
    "github.com/satyambaran/oauth/server/users/config"
    "github.com/satyambaran/oauth/server/users/structs"
)

func Resource(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    type request struct {
        Type int
        Name string
        URI  string
    }
    var body request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }
    ID, err := c.Locals("id").(int)
    if !err {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "internal server error"})
    }
    resource := structs.Resource{
        UserId: ID,
        Name:   body.Name,
        Type:   body.Type,
        URI:    body.URI,
    }
    if err := db.Create(&resource).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
    }

    if err := db.Where(&structs.Resource{Name: body.Name, URI: body.URI}).First(&resource).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
    }

    return c.JSON(resource)

}
func Resources(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    ID, err := c.Locals("id").(int)
    if !err {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "internal server error"})
    }
    var resources []structs.Resource
    if err := db.Where("user_id = ?", ID).Find(&resources).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
    }

    return c.JSON(resources)
}
func DeleteResource(c *fiber.Ctx, config *config.Config) error {
    db := config.DB
    UserID, err := c.Locals("id").(int)
    ResourceID := c.Params("id")
    if !err {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "internal server error"})
    }
    res := db.Where("user_id = ?", UserID).Where("id = ?", ResourceID).Delete(&structs.Resource{})
    if err := res.Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
    }

    return c.JSON(fiber.Map{"message": fmt.Sprintf("deleted %v column", res.RowsAffected)})
}
