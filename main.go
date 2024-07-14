package main

import (
    "context"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "log"
    "math/rand"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

const (
    initialLength = 6
    maxRetries    = 8
)

type URL struct {
    ID          uint   `gorm:"primaryKey"`
    ShortURL    string `gorm:"unique;not null"`
    OriginalURL string `gorm:"not null"`
}

type PostgreSQLURLShortener struct {
    db    *gorm.DB
    cache *redis.Client
    rng   *rand.Rand
}

func NewPostgreSQLURLShortener(dsn string, redisOptions *redis.Options) (*PostgreSQLURLShortener, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    db.AutoMigrate(&URL{})

    cache := redis.NewClient(redisOptions)
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    return &PostgreSQLURLShortener{db: db, cache: cache, rng: rng}, nil
}

func (s *PostgreSQLURLShortener) generateShortURL(url string, length int) string {
    // hash := sha256.Sum256([]byte(url + time.Now().String()))
    hash := sha256.Sum256([]byte(url + string(rune(s.rng.Int63()))))
    return base64.URLEncoding.EncodeToString(hash[:length])
}

func (s *PostgreSQLURLShortener) Shorten(url string) (string, error) {

    length := initialLength

    for i := 0; i <= maxRetries; i++ {
        shortURL := s.generateShortURL(url, length)
        newURL := URL{ShortURL: shortURL, OriginalURL: url}
        result := s.db.Create(&newURL)
        if result.Error == nil {
            s.cache.Set(context.Background(), shortURL, url, 0)
            return shortURL, nil
        }
        if result.Error != nil && !errors.Is(result.Error, gorm.ErrDuplicatedKey) {
            return "", result.Error
        }
        if i == maxRetries {
            length++
        }
    }
    return "", errors.New("failed to generate a unique short URL after multiple attempts")
}

func (s *PostgreSQLURLShortener) Resolve(shortURL string) (string, error) {
    ctx := context.Background()

    // Check cache first
    originalURL, err := s.cache.Get(ctx, shortURL).Result()
    if err == nil {
        return originalURL, nil
    }

    // Fallback to database if not found in cache
    var url URL
    result := s.db.First(&url, "short_url = ?", shortURL)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return "", errors.New("URL not found")
    } else if result.Error != nil {
        return "", result.Error
    }

    // Update cache
    s.cache.Set(ctx, shortURL, url.OriginalURL, 0)
    return url.OriginalURL, nil
}

func main() {
    dsn := "user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    redisOptions := &redis.Options{
        Addr: "localhost:6379",
    }
    urlShortener, err := NewPostgreSQLURLShortener(dsn, redisOptions)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    app := fiber.New()
    app.Use(logger.New())

    app.Post("/shorten", func(c *fiber.Ctx) error {
        type request struct {
            URL string `json:"url"`
        }
        var req request
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
        }
        shortURL, err := urlShortener.Shorten(req.URL)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"short_url": shortURL})
    })

    app.Get("/:shortURL", func(c *fiber.Ctx) error {
        shortURL := c.Params("shortURL")
        originalURL, err := urlShortener.Resolve(shortURL)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Redirect(originalURL)
    })

    log.Fatal(app.Listen(":3000"))
}
