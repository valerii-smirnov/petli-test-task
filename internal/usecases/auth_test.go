package usecases

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"testing"
)

func TestAuth_SignUp(t *testing.T) {
	controller := gomock.NewController(t)
	passwordHasherMock := NewMockPasswordHasher(controller)
	userAdapterMock := NewMockUserAdapter(controller)

	email := "test@test.com"
	password := "aaaa"
	passwordHashed := "aaaa bbbb"
	testingError := errors.New("testing-error")

	req := domain.SignUp{
		Email:    email,
		Password: password,
	}

	su := domain.SignUp{
		Email:    email,
		Password: passwordHashed,
	}

	type fields struct {
		passwordHasher PasswordHasher
		tokenGenerator TokenGenerator
		userAdapter    UserAdapter
	}
	type args struct {
		ctx context.Context
		in  domain.SignUp
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		wantErr   bool
	}{
		{
			name: "call exists error",
			fields: fields{
				passwordHasher: passwordHasherMock,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				userAdapterMock.EXPECT().Exists(gomock.Any(), gomock.Eq(email)).Return(false, testingError)
			},
			wantErr: true,
		},
		{
			name: "user already exists error",
			fields: fields{
				passwordHasher: passwordHasherMock,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				userAdapterMock.EXPECT().Exists(gomock.Any(), gomock.Eq(email)).Return(true, nil)
			},
			wantErr: true,
		},
		{
			name: "creation user error",
			fields: fields{
				passwordHasher: passwordHasherMock,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				userAdapterMock.EXPECT().Exists(gomock.Any(), gomock.Eq(email)).Return(false, nil)
				passwordHasherMock.EXPECT().StringHash(gomock.Eq(password)).Return(passwordHashed)
				userAdapterMock.EXPECT().Create(gomock.Any(), gomock.Eq(su)).Return(testingError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				passwordHasher: passwordHasherMock,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				userAdapterMock.EXPECT().Exists(gomock.Any(), gomock.Eq(email)).Return(false, nil)
				passwordHasherMock.EXPECT().StringHash(gomock.Eq(password)).Return(passwordHashed)
				userAdapterMock.EXPECT().Create(gomock.Any(), gomock.Eq(su)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			a := NewAuth(tt.fields.passwordHasher, tt.fields.tokenGenerator, tt.fields.userAdapter)
			err := a.SignUp(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAuth_SignIn(t *testing.T) {
	controller := gomock.NewController(t)
	passwordHasherMock := NewMockPasswordHasher(controller)
	userAdapterMock := NewMockUserAdapter(controller)
	tokenGenerator := NewMockTokenGenerator(controller)

	password := "aaaa"
	passwordHashed := "aaaa bbbb"
	testingError := errors.New("testing-error")

	userID := uuid.New()

	foundUser := domain.User{
		ID: userID,
	}

	req := domain.SingIn{
		Email:    "test@test.com",
		Password: password,
	}

	reqHashed := domain.SingIn{
		Email:    "test@test.com",
		Password: passwordHashed,
	}

	token := "tokentoken"

	type fields struct {
		passwordHasher PasswordHasher
		tokenGenerator TokenGenerator
		userAdapter    UserAdapter
	}
	type args struct {
		ctx context.Context
		in  domain.SingIn
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      domain.Token
		wantErr   bool
	}{
		{
			name: "getting user error",
			fields: fields{
				passwordHasher: passwordHasherMock,
				tokenGenerator: tokenGenerator,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				passwordHasherMock.EXPECT().StringHash(gomock.Eq(password)).Return(passwordHashed)
				userAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(reqHashed)).Return(domain.User{}, testingError)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "token generation error",
			fields: fields{
				passwordHasher: passwordHasherMock,
				tokenGenerator: tokenGenerator,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				passwordHasherMock.EXPECT().StringHash(gomock.Eq(password)).Return(passwordHashed)
				userAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(reqHashed)).Return(foundUser, nil)
				tokenGenerator.EXPECT().Generate(gomock.Eq(userID)).Return("", testingError)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				passwordHasher: passwordHasherMock,
				tokenGenerator: tokenGenerator,
				userAdapter:    userAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				in:  req,
			},
			mocksInit: func() {
				passwordHasherMock.EXPECT().StringHash(gomock.Eq(password)).Return(passwordHashed)
				userAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(reqHashed)).Return(foundUser, nil)
				tokenGenerator.EXPECT().Generate(gomock.Eq(userID)).Return(token, nil)
			},
			want:    domain.Token(token),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			a := NewAuth(tt.fields.passwordHasher, tt.fields.tokenGenerator, tt.fields.userAdapter)
			got, err := a.SignIn(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
