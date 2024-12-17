package tokenizers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	_ "github.com/go-sql-driver/mysql"
)

type JWTConfig struct {
	Key      string
	Duration time.Duration
}

type JWT struct {
	config JWTConfig
}

func NewTokenizer(config JWTConfig) JWT {
	return JWT{
		config: config,
	}
}

func (tokenizer JWT) GenerateToken(username string, userID int64, userType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":        username,
		"user_id":         userID,
		"tipo":            userType, // Incluye el tipo de usuario
		"expiration_date": time.Now().UTC().Add(tokenizer.config.Duration),
	})

	value, err := token.SignedString([]byte(tokenizer.config.Key))
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %w", err)
	}

	return value, nil
}
