package service

import (
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nickchervov/go-diplom-project/internal/domain"
)

func (s *SchedulerService) IsAuthEnabled() bool {
	return s.password != ""
}

func (s *SchedulerService) ValidateToken(tokenString string) error {
	var claims domain.Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method")
		}
		return []byte(s.password), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("authorization issue")
	}

	if claims.PasswordHash != sha256.Sum256([]byte(s.password)) {
		return fmt.Errorf("authorization require")
	}

	return nil
}
