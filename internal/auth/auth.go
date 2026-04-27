package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// key
var accessSecretKey = []byte(config.GetEnv("ACCESS_SECRET_KEY", "access-secret-key"))
var refreshSecretKey = []byte(config.GetEnv("REFRESH_SECRET_KEY", "refresh-secret-key"))

// var cookieName string =

type tokenType string

const (
	TokenAccess  tokenType = "access_token"
	TokenRefresh tokenType = "refresh_token"
)

var (
	expTimeTokenAccess  time.Time = time.Now().Add(time.Hour * 24)
	expTimeTokenRefresh time.Time = time.Now().Add(time.Hour * 24)
)

// var expTime time.Time =

// var secretKey = []byte("secret-key")

type Claims struct {
	UserID   uint16    `json:"user_id"`
	Username string    `json:"username,omitempty"`
	Type     tokenType `json:"type"`
	jwt.RegisteredClaims
}

// type RefreshClaims struct {
// 	UserID uint16 `json:"user_id"`
// 	Type   string `json:"type"`
// 	jwt.RegisteredClaims
// }

// HashPassword hashes password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID uint16, username string, cookieName tokenType) (string, error) {
	var claims Claims
	now := time.Now()

	switch cookieName {
	case TokenAccess:
		claims = Claims{
			UserID:   userID,
			Username: username,
			Type:     TokenAccess,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "your-app",
				Subject:   strconv.Itoa(int(userID)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(accessSecretKey)

	case TokenRefresh:
		claims = Claims{
			UserID: userID,
			Type:   TokenRefresh,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "your-app",
				Subject:   strconv.Itoa(int(userID)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(refreshSecretKey)

	default:
		return "", fmt.Errorf("invalid token type")
	}
}

func parseToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString, string(accessSecretKey))
	if err != nil {
		return nil, err
	}

	if claims.Type != TokenAccess {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString, string(refreshSecretKey))
	if err != nil {
		return nil, err
	}

	if claims.Type != TokenRefresh {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// func ValidateToken(tokenString string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
// 		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
// 			return nil, fmt.Errorf("unexpected signing method")
// 		}
// 		return secretKey, nil
// 	})

// 	if err != nil {
// 		if errors.Is(err, jwt.ErrTokenExpired) {
// 			return nil, fmt.Errorf("token expired")
// 		}
// 		return nil, fmt.Errorf("invalid token: %w", err)
// 	}

// 	claims, ok := token.Claims.(*Claims)
// 	if ok && token.Valid {
// 		return claims, nil
// 	}

// 	if claims.ExpiresAt.Time.Before(time.Now()) {
// 		return nil, fmt.Errorf("token expired")
// 	}

// 	return nil, fmt.Errorf("invalid token")
// }

func SetCookie(cookieName tokenType, token string, w http.ResponseWriter) {
	now := time.Now()

	var maxAge int
	var expires time.Time

	switch cookieName {
	case TokenAccess:
		expires = now.Add(15 * time.Minute)
		maxAge = 15 * 60

	case TokenRefresh:
		expires = now.Add(7 * 24 * time.Hour)
		maxAge = 7 * 24 * 60 * 60

	default:
		return
	}

	cookie := &http.Cookie{
		Name:     string(cookieName),
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // set false di dev kalau pakai HTTP
		SameSite: http.SameSiteNoneMode, // ubah ke Lax jika same domain
		Expires:  expires,
		MaxAge:   maxAge,
	}

	http.SetCookie(w, cookie)
}

func DeleteCookie(cookieName string, w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{}

	c.Name = cookieName
	c.Value = ""
	c.Expires = time.Unix(0, 0)
	c.MaxAge = -1

	http.SetCookie(w, c)

	// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// func SetCookies() {

// }
