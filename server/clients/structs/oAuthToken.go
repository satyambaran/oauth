package structs

import "github.com/golang-jwt/jwt/v4"

type OAuthJwtCustomClaims struct {
    ID       int    `json:"user_id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    ClientID string `json:"client_id"`
    Scope    int    `json:"scope"`
    jwt.StandardClaims
}

type OAuthJwtCustomRefreshClaims struct {
    ID       int    `json:"user_id"`
    ClientID string `json:"client_id"`
    jwt.StandardClaims
}
