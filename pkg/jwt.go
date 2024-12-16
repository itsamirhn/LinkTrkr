package pkg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService[T any] interface {
	Encode(data T, expiry *time.Duration) (string, error)
	Decode(tokenString string) (*T, error)
}

type jwtService[T any] struct {
	secretKey []byte
}

func NewJWTService[T any](secretKey string) JWTService[T] {
	return &jwtService[T]{secretKey: []byte(secretKey)}
}

func (j *jwtService[T]) Encode(data T, expiry *time.Duration) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	claims := jwt.MapClaims{
		"data": string(dataBytes),
	}

	if expiry != nil {
		claims["exp"] = time.Now().Add(*expiry).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService[T]) Decode(tokenString string) (*T, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token has expired")
			}
		}

		dataJSON, ok := claims["data"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid data format in token")
		}

		var data T
		if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		return &data, nil
	}
	return nil, fmt.Errorf("invalid token or claims")
}
