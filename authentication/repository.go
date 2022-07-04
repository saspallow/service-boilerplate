package authentication

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
)

type Repository interface {
	Authentication(ctx context.Context, username string) (string, string, error)
}

type repository struct {
	contextTimeout time.Duration
}

func NewRepository(
	timeout time.Duration,
) Repository {
	return &repository{
		contextTimeout: timeout,
	}
}

func (r *repository) Authentication(ctx context.Context, username string) (string, string, error) {
	db := ctx.Value("db").(*pg.DB)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	var uuid, password string

	_, err := db.QueryOneContext(ctx, pg.Scan(
		&uuid,
		&password,
	), `
		SELECT
			uuid, password
		FROM public.users
		WHERE username = ? AND deleted_at IS NULL
	 `, username)

	if err == pg.ErrNoRows {
		return "", "", ErrUsernameOrPasswordIncorrect
	}

	if err != nil {
		return "", "", err
	}

	return uuid, password, nil
}
