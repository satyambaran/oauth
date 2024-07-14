package helper

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/golang-jwt/jwt/v4"
    "github.com/jackc/pgconn"
    "github.com/satyambaran/oauth/server/clients/config"
    "github.com/satyambaran/oauth/server/clients/structs"
    "gorm.io/gorm"
)

var ctx = context.Background()

func GenerateAuthCode(length int) (string, error) {
    authCode := GenerateRandomString(length)
    return authCode, nil
}
func CreateClient(db *gorm.DB, client *structs.Client) (*structs.Client, error) {
    var retClient *structs.Client
    err := db.Transaction(func(tx *gorm.DB) error {
        for {
            client.ClientID = GenerateRandomString(7)
            if err := tx.Create(client).Error; err != nil {
                if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
                    continue
                }
                return err
            }
            retClient = client
            return nil
        }
    })
    return retClient, err
}
func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, length)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }
    for i := range b {
        b[i] = charset[int(b[i])%len(charset)]
    }
    return string(b)
}
func ValidateUserToken(authToken string) (*structs.JwtUserCustomClaims, string) {
    token, err := jwt.ParseWithClaims(authToken, &structs.JwtUserCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("your_secret_key"), nil
    })
    if err != nil {
        return nil, err.Error()
    }

    claims, ok := token.Claims.(*structs.JwtUserCustomClaims)
    if !ok || !token.Valid {
        return nil, "invalid token claims"
    }
    if claims.ExpiresAt < time.Now().Local().Unix() {
        return nil, "token expired"
    }
    return claims, ""
}
func ValidateToken(authToken string) (*structs.JwtCustomClaims, string) {
    token, err := jwt.ParseWithClaims(authToken, &structs.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("your_secret_key"), nil
    })
    if err != nil {
        return nil, err.Error()
    }

    claims, ok := token.Claims.(*structs.JwtCustomClaims)
    if !ok || !token.Valid {
        return nil, "invalid token claims"
    }
    if claims.ExpiresAt < time.Now().Local().Unix() {
        return nil, "token expired"
    }
    return claims, ""
}
func ValidateRefreshToken(authToken string) (*structs.JwtCustomRefreshClaims, string) {
    token, err := jwt.ParseWithClaims(authToken, &structs.JwtCustomRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("your_secret_key"), nil
    })
    if err != nil {
        return nil, err.Error()
    }

    claims, ok := token.Claims.(*structs.JwtCustomRefreshClaims)
    if !ok || !token.Valid {
        return nil, "invalid token claims"
    }
    if claims.ExpiresAt < time.Now().Local().Unix() {
        return nil, "refresh token expired"
    }
    return claims, ""
}
func SaveAllTokens(client *structs.Client, token string, refreshToken string, atExp int, rtExp int, config *config.Config) error {
    rdb := config.RDB
    err := SaveAccessToken(rdb, client, token, atExp)
    if err != nil {
        return err
    }

    err = SaveRefreshToken(rdb, client, refreshToken, rtExp)
    if err != nil {
        return err
    }
    return nil
}
func SaveRefreshToken(rdb *redis.Client, client *structs.Client, refreshToken string, rtExp int) error {
    err := rdb.Set(ctx, client.ClientID+":refresh_token", refreshToken, time.Duration(rtExp)).Err()
    return err
}
func SaveAccessToken(rdb *redis.Client, client *structs.Client, token string, atExp int) error {
    err := rdb.Set(ctx, client.ClientID+":access_token", token, time.Duration(atExp)).Err()
    return err
}
func CreateAllTokens(client *structs.Client, secret string, config *config.Config) (string, string, error) {
    atExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTE"))
    if err != nil {
        atExp = 30
    }
    rtExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"))
    if err != nil {
        rtExp = 72
    }
    refreshToken, err := CreateRefreshToken(client, secret, rtExp)
    if err != nil {
        return "", "", fmt.Errorf("could not generate refresh token")
    }
    token, err := CreateAccessToken(client, secret, atExp)
    if err != nil {
        return "", "", fmt.Errorf("could not generate token")
    }
    err = SaveAllTokens(client, token, refreshToken, atExp, rtExp, config)
    return refreshToken, token, err
}
func CreateAccessToken(client *structs.Client, secret string, atExp int) (token string, err error) {
    claims := &structs.JwtCustomClaims{
        Name:     client.Name,
        Email:    client.Email,
        ClientID: client.ClientID,
        ID:       client.ID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().UTC().Local().Add(time.Minute * time.Duration(atExp)).Unix(),
        },
    }

    token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
    if err != nil {
        return "", err
    }
    return token, err
}
func CreateRefreshToken(client *structs.Client, secret string, rtExp int) (token string, err error) {
    claimsRefresh := &structs.JwtCustomRefreshClaims{
        ID:       client.ID,
        ClientID: client.ClientID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().UTC().Local().Add(time.Hour * time.Duration(rtExp)).Unix(),
        },
    }
    token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh).SignedString([]byte(secret))
    if err != nil {
        return "", err
    }
    return token, err
}
func GenerateSalt() (string, error) {
    salt := make([]byte, 16)
    _, err := rand.Read(salt)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(salt), nil
}
func CreateOAuthToken(client *structs.Client, secret string, atExp int) (token string, err error) {
    claims := &structs.OAuthJwtCustomClaims{
        ID:       client.ID,
        Name:     client.Name,
        Email:    client.Email,
        ClientID: client.ClientID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().UTC().Local().Add(time.Minute * time.Duration(atExp)).Unix(),
        },
    }

    token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
    if err != nil {
        return "", err
    }
    return token, err
}
func CreateOAuthRefreshToken(client *structs.Client, secret string, rtExp int) (token string, err error) {
    claimsRefresh := &structs.JwtCustomRefreshClaims{
        ID:       client.ID,
        ClientID: client.ClientID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().UTC().Local().Add(time.Hour * time.Duration(rtExp)).Unix(),
        },
    }
    token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh).SignedString([]byte(secret))
    if err != nil {
        return "", err
    }
    return token, err
}
func GetRefreshToken(rdb *redis.Client, client *structs.Client) (token string, err error) {
    refreshToken, err := rdb.Get(ctx, client.ClientID+":refresh_token").Result()
    if err == redis.Nil {
        return "", fmt.Errorf("key does not exist")
    } else if err != nil {
        return "", err
    } else {
        return refreshToken, nil
    }
}
