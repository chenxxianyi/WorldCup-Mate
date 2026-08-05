package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

const (
	accessTokenTTL = 30 * time.Minute // SEC-04: short-lived access token
	tokenIssuer    = "worldcup-mate"
)

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
	return GenerateTokenAt(userID, role, time.Now())
}

// GenerateTokenAt issues a token with a caller-controlled issued-at time
// (test hook; production uses GenerateToken).
func GenerateTokenAt(userID uint, role string, issuedAt time.Time) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			Subject:   strconv.FormatUint(uint64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithIssuer(tokenIssuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
