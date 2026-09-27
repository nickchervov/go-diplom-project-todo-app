package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nickchervov/go-diplom-project/internal/domain"
	"github.com/nickchervov/go-diplom-project/internal/dto"
)

func (s *SchedulerService) SignIn(ctx context.Context, input dto.SignInInput) (dto.SignInOutput, error) {
	claims := domain.Claims{
		PasswordHash: sha256.Sum256([]byte(s.password)),
	}

	if input.Password == s.password {
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.password))
		if err != nil {
			return dto.SignInOutput{}, fmt.Errorf("signing token by secret key: %w", err)
		}

		output := dto.SignInOutput{Token: token}
		return output, nil
	}

	return dto.SignInOutput{}, fmt.Errorf("%w", domain.ErrIncorrectPassword)
}
