package structs

import jwt "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
    ID int `json:"id"`
    jwt.StandardClaims
}
