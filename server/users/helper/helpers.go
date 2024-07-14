package helper

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/golang-jwt/jwt/v4"

    "github.com/satyambaran/oauth/server/users/config"
    "github.com/satyambaran/oauth/server/users/structs"
)

var ctx = context.Background()

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
func SaveAllTokens(user *structs.User, token string, refreshToken string, atExp int, rtExp int, config *config.Config) error {
    rdb := config.RDB
    err := SaveAccessToken(rdb, user, token, atExp)
    if err != nil {
        return err
    }

    err = SaveRefreshToken(rdb, user, refreshToken, rtExp)
    if err != nil {
        return err
    }
    return nil
}
func SaveRefreshToken(rdb *redis.Client, user *structs.User, refreshToken string, rtExp int) error {
    err := rdb.Set(ctx, strconv.Itoa(user.ID)+":refresh_token", refreshToken, time.Duration(rtExp)).Err()
    return err
}
func SaveAccessToken(rdb *redis.Client, user *structs.User, token string, atExp int) error {
    err := rdb.Set(ctx, strconv.Itoa(user.ID)+":access_token", token, time.Duration(atExp)).Err()
    return err
}
func CreateAllTokens(user *structs.User, secret string, config *config.Config) (string, string, error) {
    atExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTE"))
    if err != nil {
        atExp = 30
    }
    rtExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"))
    if err != nil {
        rtExp = 72
    }
    refreshToken, err := CreateRefreshToken(user, secret, rtExp)
    if err != nil {
        return "", "", fmt.Errorf("could not generate refresh token")
    }
    token, err := CreateAccessToken(user, secret, atExp)
    if err != nil {
        return "", "", fmt.Errorf("could not generate token")
    }
    err = SaveAllTokens(user, token, refreshToken, atExp, rtExp, config)
    return refreshToken, token, err
}
func CreateAccessToken(user *structs.User, secret string, atExp int) (token string, err error) {
    claims := &structs.JwtCustomClaims{
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
        ID:    user.ID,
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
func CreateRefreshToken(user *structs.User, secret string, rtExp int) (token string, err error) {
    claimsRefresh := &structs.JwtCustomRefreshClaims{
        ID: user.ID,
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
func GetRefreshToken(rdb *redis.Client, user *structs.User) (token string, err error) {
    refreshToken, err := rdb.Get(ctx, strconv.Itoa(user.ID)+":refresh_token").Result()
    if err == redis.Nil {
        return "", fmt.Errorf("key does not exist")
    } else if err != nil {
        return "", err
    } else {
        return refreshToken, nil
    }
}
