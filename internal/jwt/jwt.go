package jwt

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey     []byte
	secretKeyOnce sync.Once
)

func GetSecretKey() []byte {
	secretKeyOnce.Do(func() {
		secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	})
	return secretKey
}

type CustomClaims struct {
	CID          uint64 `json:"cid"`
	Email        string `json:"email"`
	Roles        string `json:"roles"`
	CountryName  string `json:"country_name"`
	DivisionName string `json:"division_name"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT
func GenerateToken(cid uint64, email, roles, countryName, divisionName string) (string, error) {
	claims := CustomClaims{
		CID:          cid,
		Email:        email,
		Roles:        roles,
		CountryName:  countryName,
		DivisionName: divisionName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// VerifyToken парсит и проверяет JWT
func VerifyToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Защита от алгоритмической атаки
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
