package login

import (
	"net/http"
	"time"
)

type AuthConfig struct {
	JWTSecret          []byte        `json:"jwt_secret"`
	AccessTTL          time.Duration `json:"access_ttl"`
	RefreshTTL         time.Duration `json:"refresh_ttl"`
	RefreshCookieName  string        `json:"refresh_cookie_name"`
	CookieDomain       string        `json:"cookie_domain"`
	CookieSecure       bool          `json:"cookie_secure"`
	CookieSameSiteMode http.SameSite `json:"cookie_same_site_mode"`
	TokenIssuer        string        `json:"token_issuer"`
	DBPath             string        `json:"db_path"`
}

func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:          []byte("dev-secret-change-me"),
		AccessTTL:          15 * time.Minute,
		RefreshTTL:         7 * 24 * time.Hour,
		RefreshCookieName:  "refresh_token",
		CookieDomain:       "", // empty OK for localhost
		CookieSecure:       false,
		CookieSameSiteMode: http.SameSiteNoneMode,
		TokenIssuer:        "decay",
		DBPath:             "data\\decay.db",
	}
}

// var (
// 	jwtSecret          = []byte(getEnv("JWT_SECRET", "dev-secret-change-me"))
// 	accessTTL          = getDurationEnv("ACCESS_TTL", 15*time.Minute)
// 	refreshTTL         = getDurationEnv("REFRESH_TTL", 7*24*time.Hour)
// 	refreshCookieName  = getEnv("REFRESH_COOKIE", "refresh_token")
// 	cookieDomain       = os.Getenv("COOKIE_DOMAIN")                            // empty OK for localhost
// 	cookieSecure       = strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true") // TRUE in prod (HTTPS)
// 	cookieSameSiteMode = http.SameSiteLaxMode                                  // Lax is fine for refresh
// 	tokenIssuer        = getEnv("JWT_ISS", "myapp")
// 	dbPath             = getEnv("DB_PATH", "auth.db")
// )

// Util functions to get env variables with defaults
// func getEnv(k, def string) string {
// 	if v := os.Getenv(k); v != "" {
// 		return v
// 	}
// 	return def
// }
// func getDurationEnv(k string, def time.Duration) time.Duration {
// 	if v := os.Getenv(k); v != "" {
// 		if d, err := time.ParseDuration(v); err == nil {
// 			return d
// 		}
// 	}
// 	return def
// }
