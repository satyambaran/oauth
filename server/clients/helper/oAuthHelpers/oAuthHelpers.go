package helper

import (
    "context"
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/golang-jwt/jwt/v4"
    "github.com/satyambaran/oauth/server/clients/config"
    "github.com/satyambaran/oauth/server/clients/structs"
)

var ctx = context.Background()

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
func SaveAllTokens(userOAuth *structs.UserOAuth, token string, refreshToken string, atExp int, rtExp int, config *config.Config) error {
    rdb := config.RDB
    err := SaveAccessToken(rdb, userOAuth, token, atExp)
    if err != nil {
        return err
    }

    err = SaveRefreshToken(rdb, userOAuth, refreshToken, rtExp)
    if err != nil {
        return err
    }
    return nil
}
func SaveRefreshToken(rdb *redis.Client, userOAuth *structs.UserOAuth, refreshToken string, rtExp int) error {
    err := rdb.Set(ctx, userOAuth.ClientID+":refresh_token", refreshToken, time.Duration(rtExp)).Err()
    return err
}
func SaveAccessToken(rdb *redis.Client, userOAuth *structs.UserOAuth, token string, atExp int) error {
    err := rdb.Set(ctx, userOAuth.ClientID+":access_token", token, time.Duration(atExp)).Err()
    return err
}
func CreateAllTokens(userOAuth *structs.UserOAuth, secret string, config *config.Config) (string, string, error) {
    atExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTE"))
    if err != nil {
        atExp = 30
    }
    rtExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"))
    if err != nil {
        rtExp = 72
    }
    refreshToken, err := CreateRefreshToken(userOAuth, secret, rtExp)
    if err != nil {
        return "", "", fmt.Errorf("could not generate refresh token")
    }
    token, err := CreateAccessToken(userOAuth, secret, atExp)
    if err != nil {
        return "", "", fmt.Errorf("could not generate token")
    }
    err = SaveAllTokens(userOAuth, token, refreshToken, atExp, rtExp, config)
    return refreshToken, token, err
}
func CreateAccessToken(userOAuth *structs.UserOAuth, secret string, atExp int) (token string, err error) {
    claims := &structs.JwtCustomClaims{
        Name:     userOAuth.Name,
        Email:    userOAuth.Email,
        ClientID: userOAuth.ClientID,
        ID:       userOAuth.ID,
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
func CreateRefreshToken(userOAuth *structs.UserOAuth, secret string, rtExp int) (token string, err error) {
    claimsRefresh := &structs.JwtCustomRefreshClaims{
        ID:       userOAuth.ID,
        ClientID: userOAuth.ClientID,
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
