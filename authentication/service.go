package authentication

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidArgument             = errors.New("invalid argument")
	ErrUsernameOrPasswordIncorrect = errors.New("username or password incorrect")
)

type Service interface {
	SignIn(ctx context.Context, username string, password string) (*SignInResponseModel, error)
}

type service struct {
	repo           Repository
	contextTimeout time.Duration
}

func NewService(
	repo Repository,
	timeout time.Duration,
) Service {
	return &service{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (s *service) SignIn(ctx context.Context, username string, password string) (*SignInResponseModel, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if len(username) < 1 || len(password) < 1 {
		return nil, ErrInvalidArgument
	}

	username = strings.ToLower(username)

	uuid, hashedPassword, err := s.repo.Authentication(ctx, username)
	if err != nil {
		return nil, err
	}

	err = VerifyPassword(hashedPassword, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrUsernameOrPasswordIncorrect
	}

	accessToken, err := CreateToken(uuid, nil, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	return &SignInResponseModel{
		accessToken,
		"bearer",
		int64(accessTokenDuration.Seconds()),
		uuid,
	}, nil
}
