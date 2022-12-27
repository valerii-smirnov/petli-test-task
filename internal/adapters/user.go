package adapters

import (
	"context"
	"database/sql"
	"github.com/valerii-smirnov/petli-test-task/internal/adapters/models"
	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"

	"github.com/jmoiron/sqlx"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}

func (u User) Exists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := u.db.QueryRowContext(ctx, "select exists(select * FROM users where email=$1)", email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, ierr.WrapCode(ierr.Internal, err, "execution exists query error")
	}

	return exists, nil
}

func (u User) Create(ctx context.Context, in domain.SignUp) error {
	query := "insert into users (email, password_hash) values ($1, $2)"

	if _, err := u.db.ExecContext(ctx, query, in.Email, in.Password); err != nil {
		return ierr.WrapCode(ierr.Internal, err, "execution insert query error")
	}

	return nil
}

func (u User) Get(ctx context.Context, in domain.SingIn) (domain.User, error) {
	query := "select * from users WHERE email=$1 AND password_hash=$2"

	var user models.User

	if err := u.db.GetContext(ctx, &user, query, in.Email, in.Password); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, ierr.WrapCode(ierr.NotFound, err, "user not found")
		}

		return domain.User{}, ierr.WrapCode(ierr.Internal, err, "execution select query error")
	}

	return domain.User{
		ID:           user.ID,
		Email:        user.Email,
		RegisteredAt: user.RegisteredAt,
	}, nil
}
