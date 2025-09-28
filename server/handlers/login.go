package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prashantkhandelwal/decay/db"
	"github.com/prashantkhandelwal/decay/login"
)

var authConfig = login.DefaultAuthConfig()

func LoginHandler() gin.HandlerFunc {
	fn := func(g *gin.Context) {

		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := g.ShouldBindJSON(&req); err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
			return
		}
		u, ok := login.Users[req.Username]
		if !ok || u.Password != req.Password {
			g.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		access, aExp, err := makeAccessToken(req.Username, u.Role)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "cannot issue access token"})
			return
		}
		refresh, rExp, jti, issuedAt, err := makeRefreshToken(req.Username, u.Role)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "cannot issue refresh token"})
			return
		}
		if err := db.SaveRefresh(g, jti, req.Username, issuedAt, rExp); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "persist refresh token failed"})
			return
		}

		setRefreshCookie(g, refresh, rExp)

		g.JSON(http.StatusOK, gin.H{
			"access_token":       access,
			"expires_in":         int(time.Until(aExp).Seconds()),
			"token_type":         "Bearer",
			"refresh_expires_in": int(time.Until(rExp).Seconds()),
		})

	}
	return gin.HandlerFunc(fn)
}

func RefreshHandler() gin.HandlerFunc {
	fn := func(g *gin.Context) {

		// Prefer cookie, fallback to JSON body { "refresh_token": "..." } for tool callers.
		refreshToken, err := g.Cookie(authConfig.RefreshCookieName)
		if err != nil || refreshToken == "" {
			var req struct {
				Refresh string `json:"refresh_token"`
			}
			if bindErr := g.ShouldBindJSON(&req); bindErr != nil || req.Refresh == "" {
				g.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
				return
			}
			refreshToken = req.Refresh
		}

		claims, token, err := parseToken(refreshToken)
		if err != nil || !token.Valid || claims.TokenType != "refresh" {
			g.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			g.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
			return
		}

		// Enforce per-user logout cutoff on refresh tokens too
		if va, err := db.GetValidAfter(g, claims.Username); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "session check failed"})
			return
		} else if !va.IsZero() && claims.IssuedAt != nil && claims.IssuedAt.Time.Before(va) {
			g.JSON(http.StatusUnauthorized, gin.H{"error": "session invalidated"})
			return
		}

		// Validate against DB (not in-memory)
		if ok, err := db.IsRefreshValid(g, claims.ID, claims.Username, time.Now()); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "refresh lookup failed"})
			return
		} else if !ok {
			g.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not recognized or revoked"})
			return
		}

		// Rotate: revoke old, issue new pair
		_ = db.RevokeRefresh(g, claims.ID)

		// Re-fetch role (authoritative source); here from the demo map
		u := login.Users[claims.Username]
		access, _, err := makeAccessToken(claims.Username, u.Role)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "cannot issue access token"})
			return
		}
		newRefresh, rExp, newJti, issuedAt, err := makeRefreshToken(claims.Username, u.Role)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "cannot issue refresh token"})
			return
		}
		if err := db.SaveRefresh(g, newJti, claims.Username, issuedAt, rExp); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "persist refresh token failed"})
			return
		}
		setRefreshCookie(g, newRefresh, rExp)

		g.JSON(http.StatusOK, gin.H{
			"access_token": access,
			"token_type":   "Bearer",
			"expires_in":   int(authConfig.AccessTTL.Seconds()),
		})
	}
	return gin.HandlerFunc(fn)
}

func LogoutHandler() gin.HandlerFunc {
	fn := func(g *gin.Context) {
		var username string

		log.Println("LogoutHandler called")
		// Revoke presented refresh token if present
		if rt, err := g.Cookie(authConfig.RefreshCookieName); err == nil && rt != "" {
			log.Println("Found refresh token in cookie")
			if cl, tok, err := parseToken(rt); err == nil && tok.Valid && cl.TokenType == "refresh" {
				username = cl.Username
				if cl.ID != "" {
					err := db.RevokeRefresh(g, cl.ID)
					if err != nil {
						log.Printf("ERROR:Database: Error in revoking refresh token. %s", err)
					}
				}
			}
		}

		// Fallback: get username from access token header (if client sends it)
		if username == "" {
			log.Println("No refresh token; checking Authorization header for access token")
			if auth := g.GetHeader("Authorization"); auth != "" {
				parts := strings.SplitN(auth, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					if cl, tok, err := parseToken(parts[1]); err == nil && tok.Valid && cl.TokenType == "access" {
						username = cl.Username
					}
				}
			}
		}

		// Invalidate all tokens for this user: set valid_after and revoke all refresh tokens
		if username != "" {
			log.Println("Logging out user:", username)
			if err := db.SetValidAfterNow(g, username); err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed (valid_after)"})
				return
			}
			if err := db.RevokeAllRefreshForUser(g, username); err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed (revoke refresh)"})
				return
			}
		}

		// Clear cookie (Path MUST match setRefreshCookie)
		http.SetCookie(g.Writer, &http.Cookie{
			Name:     authConfig.RefreshCookieName,
			Value:    "",
			Path:     "/token", // keep consistent with setRefreshCookie
			Domain:   authConfig.CookieDomain,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   authConfig.CookieSecure,
			SameSite: authConfig.CookieSameSiteMode,
		})
		g.Status(http.StatusNoContent)
	}
	return gin.HandlerFunc(fn)
}

func makeAccessToken(username, role string) (string, time.Time, error) {
	exp := time.Now().Add(authConfig.AccessTTL)
	claims := &login.Claims{
		Username:  username,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    authConfig.TokenIssuer,
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// ID:     uuid.NewString(), // add if you also want a denylist
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(authConfig.JWTSecret)
	return signed, exp, err
}

func makeRefreshToken(username, role string) (token string, exp time.Time, jti string, issuedAt time.Time, err error) {
	exp = time.Now().Add(authConfig.RefreshTTL)
	issuedAt = time.Now()
	jti = uuid.NewString()
	claims := &login.Claims{
		Username:  username,
		Role:      role, // carry role if you want to propagate; we re-fetch on refresh anyway
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    authConfig.TokenIssuer,
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        jti,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tok.SignedString(authConfig.JWTSecret)
	return
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

func setRefreshCookie(g *gin.Context, token string, exp time.Time) {
	http.SetCookie(g.Writer, &http.Cookie{
		Name:     authConfig.RefreshCookieName,
		Value:    token,
		Path:     "/token", // keep consistent with logout clear
		Domain:   authConfig.CookieDomain,
		Expires:  exp,
		HttpOnly: true,
		Secure:   authConfig.CookieSecure,
		SameSite: authConfig.CookieSameSiteMode,
	})
}
