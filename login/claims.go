package login

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"typ"` // "access" or "refresh"
	jwt.RegisteredClaims
}
