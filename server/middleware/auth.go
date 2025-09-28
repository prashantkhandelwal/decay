package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prashantkhandelwal/decay/db"
	"github.com/prashantkhandelwal/decay/login"
)

var authConfig = login.DefaultAuthConfig()

func AuthMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
		return
	}
	claims, token, err := parseToken(parts[1])
	if err != nil || token == nil || !token.Valid || claims.TokenType != "access" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}
	if claims.Issuer != authConfig.TokenIssuer {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad issuer"})
		return
	}
	// Per-user logout cutoff (DB)
	if va, err := db.GetValidAfter(c, claims.Username); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "session check failed"})
		return
	} else if !va.IsZero() && claims.IssuedAt != nil && claims.IssuedAt.Time.Before(va) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session invalidated"})
		return
	}

	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
	c.Next()
}

func parseToken(tokenStr string) (*login.Claims, *jwt.Token, error) {
	claims := &login.Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return authConfig.JWTSecret, nil
	})
	return claims, tok, err
}
