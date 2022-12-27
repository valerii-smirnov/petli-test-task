package presenters

import (
	"context"
	"github.com/gin-gonic/gin"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=./mock_test.go -package=presenters -source=./contracts.go

type AuthUsecase interface {
	SignUp(ctx context.Context, su domain.SignUp) error
	SignIn(ctx context.Context, si domain.SingIn) (domain.Token, error)
}

type DogUsecase interface {
	List(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) (domain.DogList, error)
	Get(ctx context.Context, dogID uuid.UUID) (domain.Dog, error)
	Matches(ctx context.Context, userID, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error)
	Create(ctx context.Context, dog domain.Dog) (domain.Dog, error)
	Update(ctx context.Context, dogID uuid.UUID, dog domain.Dog) (domain.Dog, error)
	Delete(ctx context.Context, dogID uuid.UUID, userID uuid.UUID) error
	AddReaction(ctx context.Context, userID uuid.UUID, reaction domain.Reaction) error
}

type TokenParser interface {
	Parse(token string) (*jwt.Token, error)
}

type IdentityExtractor interface {
	ExtractFromContext(c *gin.Context) (uuid.UUID, error)
}

type Paginator interface {
	GetPagination(c *gin.Context) (domain.Pagination, error)
}
