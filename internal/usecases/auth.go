package usecases

import (
	"context"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"
)

type Auth struct {
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
	userAdapter    UserAdapter
}

func NewAuth(
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
	userAdapter UserAdapter,
) *Auth {
	return &Auth{
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
		userAdapter:    userAdapter,
	}
}

func (a Auth) SignUp(ctx context.Context, in domain.SignUp) error {
	exists, err := a.userAdapter.Exists(ctx, in.Email)
	if err != nil {
		return err
	}

	if exists {
		return ierr.New(ierr.AlreadyExists, "user with provided email already exists")
	}

	in.Password = a.passwordHasher.StringHash(in.Password)

	if err := a.userAdapter.Create(ctx, in); err != nil {
		return ierr.WrapCode(ierr.Internal, err, "creating user error")
	}

	return nil
}

func (a Auth) SignIn(ctx context.Context, in domain.SingIn) (domain.Token, error) {
	in.Password = a.passwordHasher.StringHash(in.Password)

	user, err := a.userAdapter.Get(ctx, in)
	if err != nil {
		return "", err
	}

	token, err := a.tokenGenerator.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return domain.Token(token), nil
}
