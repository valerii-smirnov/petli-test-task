package usecases

import (
	"context"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=./mock_test.go -package=usecases -source=./contracts.go

type DogAdapter interface {
	List(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) (domain.DogList, error)
	Get(ctx context.Context, dogID uuid.UUID) (domain.Dog, error)
	Matches(ctx context.Context, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error)
	Create(ctx context.Context, dog domain.Dog) (domain.Dog, error)
	Update(ctx context.Context, dogID uuid.UUID, dog domain.Dog) (domain.Dog, error)
	Delete(ctx context.Context, dogID uuid.UUID) error
	AddReaction(ctx context.Context, reaction domain.Reaction) error
}

type UserAdapter interface {
	Create(ctx context.Context, su domain.SignUp) error
	Get(ctx context.Context, si domain.SingIn) (domain.User, error)
	Exists(ctx context.Context, email string) (bool, error)
}

type PasswordHasher interface {
	StringHash(in string) string
}

type TokenGenerator interface {
	Generate(userID uuid.UUID) (string, error)
}
