package structs

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    ClientID string `json:"client_id"`
    jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
    ID       int    `json:"id"`
    ClientID string `json:"client_id"`
    jwt.StandardClaims
}

type JwtUserCustomClaims struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.StandardClaims
}
type JwtUserOAuthCustomClaims struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    ClientID  string `json:"client_id"`
    GrantType string `json:"grant_type"`
    Code      string    `json:"code"`
    jwt.StandardClaims
}

type UserOAuth struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    ClientID  string `json:"client_id"`
    GrantType string `json:"grant_type"`
    Code      string    `json:"code"`
    AuthCode  string `json:"auth_code"`
}
