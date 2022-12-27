package adapters

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUser_Exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")
	email := "test@email.com"

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      bool
		wantErr   bool
	}{
		{
			name: "query execution error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:   context.TODO(),
				email: email,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").WithArgs(email).WillReturnError(testingError)
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "user not exists",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:   context.TODO(),
				email: email,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(false)
				mock.ExpectQuery("select").WithArgs(email).WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "user not exists",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:   context.TODO(),
				email: email,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(true)
				mock.ExpectQuery("select").WithArgs(email).WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.mocksInit()

				u := NewUser(tt.fields.db)
				got, err := u.Exists(tt.args.ctx, tt.args.email)
				assert.Equal(t, tt.wantErr, err != nil)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}

func TestUser_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")

	in := domain.SignUp{
		Email:    "test@email.com",
		Password: "qwerty hashed",
	}

	type fields struct {
		db *sqlx.DB
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
			name: "insertion query error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				in:  in,
			},
			mocksInit: func() {
				mock.ExpectExec("insert").WithArgs(in.Email, in.Password).WillReturnError(testingError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				in:  in,
			},
			mocksInit: func() {
				mock.ExpectExec("insert").WithArgs(in.Email, in.Password).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			u := NewUser(tt.fields.db)
			err := u.Create(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUser_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")

	in := domain.SingIn{
		Email:    "test@email.com",
		Password: "qwerty hashed",
	}

	expectedUser := domain.User{
		ID:           uuid.New(),
		Email:        in.Email,
		RegisteredAt: time.Now(),
	}

	type fields struct {
		db *sqlx.DB
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
		want      domain.User
		wantErr   bool
	}{
		{
			name: "no rows error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				in:  in,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").WithArgs(in.Email, in.Password).WillReturnError(sql.ErrNoRows)
			},
			want:    domain.User{},
			wantErr: true,
		},
		{
			name: "execution select query error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				in:  in,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").WithArgs(in.Email, in.Password).WillReturnError(testingError)
			},
			want:    domain.User{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				in:  in,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "registered_at"}).
					AddRow(expectedUser.ID, expectedUser.Email, in.Password, expectedUser.RegisteredAt)

				mock.ExpectQuery("select").
					WithArgs(in.Email, in.Password).
					WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			u := NewUser(tt.fields.db)
			got, err := u.Get(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
