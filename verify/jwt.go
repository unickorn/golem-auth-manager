package verify

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math"
	"time"
)

// refreshTokens holds a map of emails to refresh tokens.
var refreshTokens = make(map[string]string)

// secret is the HMAC key used by the authentication server to sign JWTs.
// unfortunately RSA is not fully supported yet on tinygo, so we're down to using symmetric encryption with HMAC.
var secret []byte

func Initialize(secretString string) {
	// unfortunately env is not supported yet
	secret = []byte(secretString)
}

// CreateTokens creates an access token with 1 hour expiry and a refresh token with 3 day expiry.
func CreateTokens(email string) (string, error) {
	now := time.Now().Unix()
	accessTokenDurationSeconds := int64(math.Floor(time.Hour.Seconds()))
	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "golem-poll",
			"sub": email,
			"iat": now,
			"exp": now + accessTokenDurationSeconds,
		})
	accessToken, err := accessTokenJWT.SignedString(secret)
	if err != nil {
		return "", err
	}
	refreshTokenDurationSeconds := int64(math.Floor((time.Hour * 24 * 3).Seconds()))
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "golem-poll",
			"sub": email,
			"exp": now + refreshTokenDurationSeconds,
		})
	refreshToken, err := refreshTokenJWT.SignedString(secret)
	if err != nil {
		return "", err
	}
	// old one gets overwritten if it exists
	refreshTokens[email] = refreshToken
	result := map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires":       accessTokenDurationSeconds,
	}
	jsonResultString, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(jsonResultString), nil
}

// RefreshToken creates a new access token and refresh token for the user.
func RefreshToken(email string, refreshToken string) (string, error) {
	if rt, ok := refreshTokens[email]; !ok || rt != refreshToken {
		return "", jwt.ErrInvalidKey
	}
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if s, err := token.Claims.GetSubject(); err != nil || s != email {
			return nil, fmt.Errorf("subject does not match email")
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrTokenExpired
	}
	return CreateTokens(email)
}

// InvalidateToken invalidates the refresh token for the user.
func InvalidateToken(accessToken string) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return jwt.ErrTokenExpired
	}
	email, err := token.Claims.GetSubject()
	if _, ok := refreshTokens[email]; !ok {
		return jwt.ErrInvalidKey
	}
	delete(refreshTokens, email)
	return nil
}

// CheckToken verifies the user's access token.
func CheckToken(accessToken string) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return jwt.ErrTokenExpired
	}
	return nil
}
