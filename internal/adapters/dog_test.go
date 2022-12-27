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

func TestDog_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")
	userID := uuid.New()
	pag := domain.Pagination{
		Page:    1,
		PerPage: 2,
	}

	dogsTime := time.Now()
	dog1ID := uuid.New()
	dog2ID := uuid.New()
	expectedList := domain.DogList{
		{
			ID:        dog1ID,
			UserID:    userID,
			Name:      "dog1",
			Sex:       "male",
			Age:       2,
			Breed:     "test_breed_1",
			Image:     "http://dog-images.com/test.jpg",
			CreatedAt: dogsTime,
			UpdatedAt: dogsTime,
		},
		{
			ID:        dog2ID,
			UserID:    userID,
			Name:      "dog2",
			Sex:       "female",
			Age:       3,
			Breed:     "test_breed_1",
			Image:     "http://dog-images.com/test.jpg",
			CreatedAt: dogsTime,
			UpdatedAt: dogsTime,
		},
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx        context.Context
		userID     uuid.UUID
		pagination domain.Pagination
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      domain.DogList
		wantErr   bool
	}{
		{
			name: "query execution error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				pagination: pag,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").
					WithArgs(userID, pag.PerPage, pag.PerPage*(pag.Page-1)).
					WillReturnError(testingError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "struct scanning error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				pagination: pag,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dog1ID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime).
					AddRow(dog2ID, userID, "dog2", "male", "wrong-age-type", "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime)

				mock.ExpectQuery("select").
					WithArgs(userID, pag.PerPage, pag.PerPage*(pag.Page-1)).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				pagination: pag,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dog1ID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime).
					AddRow(dog2ID, userID, "dog2", "female", 3, "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime)

				mock.ExpectQuery("select").
					WithArgs(userID, pag.PerPage, pag.PerPage*(pag.Page-1)).
					WillReturnRows(rows)
			},
			want:    expectedList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			got, err := d.List(tt.args.ctx, tt.args.userID, tt.args.pagination)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")
	dogID := uuid.New()
	userID := uuid.New()
	dogTime := time.Now()

	expectedDog := domain.Dog{
		ID:        dogID,
		UserID:    userID,
		Name:      "dog1",
		Sex:       "male",
		Age:       2,
		Breed:     "test_breed_1",
		Image:     "http://dog-images.com/test.jpg",
		CreatedAt: dogTime,
		UpdatedAt: dogTime,
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		uid uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      domain.Dog
		wantErr   bool
	}{
		{
			name: "no rows error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").
					WithArgs(dogID).
					WillReturnError(sql.ErrNoRows)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "execution query error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").
					WithArgs(dogID).
					WillReturnError(testingError)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dogID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogTime, dogTime)

				mock.ExpectQuery("select").
					WithArgs(dogID).
					WillReturnRows(rows)
			},
			want:    expectedDog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			got, err := d.Get(tt.args.ctx, tt.args.uid)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Matches(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")
	userID := uuid.New()
	dID := uuid.New()

	pag := domain.Pagination{
		Page:    1,
		PerPage: 2,
	}

	dogsTime := time.Now()
	dog1ID := uuid.New()
	dog2ID := uuid.New()
	expectedList := domain.DogList{
		{
			ID:        dog1ID,
			UserID:    userID,
			Name:      "dog1",
			Sex:       "male",
			Age:       2,
			Breed:     "test_breed_1",
			Image:     "http://dog-images.com/test.jpg",
			CreatedAt: dogsTime,
			UpdatedAt: dogsTime,
		},
		{
			ID:        dog2ID,
			UserID:    userID,
			Name:      "dog2",
			Sex:       "female",
			Age:       3,
			Breed:     "test_breed_1",
			Image:     "http://dog-images.com/test.jpg",
			CreatedAt: dogsTime,
			UpdatedAt: dogsTime,
		},
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx        context.Context
		dogID      uuid.UUID
		pagination domain.Pagination
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      domain.DogList
		wantErr   bool
	}{
		{
			name: "select query execution error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:        context.TODO(),
				dogID:      dID,
				pagination: pag,
			},
			mocksInit: func() {
				mock.ExpectQuery("select").
					WithArgs(dID, domain.Like, pag.PerPage, pag.PerPage*(pag.Page-1)).
					WillReturnError(testingError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "result scanning error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:        context.TODO(),
				dogID:      dID,
				pagination: pag,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dog1ID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime).
					AddRow(dog2ID, userID, "dog2", "female", "wrong-age", "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime)

				mock.ExpectQuery("select").
					WithArgs(dID, domain.Like, pag.PerPage, pag.PerPage*(pag.Page-1)).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:        context.TODO(),
				dogID:      dID,
				pagination: pag,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dog1ID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime).
					AddRow(dog2ID, userID, "dog2", "female", 3, "test_breed_1", "http://dog-images.com/test.jpg", dogsTime, dogsTime)

				mock.ExpectQuery("select").
					WithArgs(dID, domain.Like, pag.PerPage, pag.PerPage*(pag.Page-1)).
					WillReturnRows(rows)
			},
			want:    expectedList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			got, err := d.Matches(tt.args.ctx, tt.args.dogID, tt.args.pagination)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")

	dogID := uuid.New()
	userID := uuid.New()
	dogTime := time.Now()

	dogIn := domain.Dog{
		UserID:    userID,
		Name:      "dog1",
		Sex:       "male",
		Age:       2,
		Breed:     "test_breed_1",
		Image:     "http://dog-images.com/test.jpg",
		CreatedAt: dogTime,
		UpdatedAt: dogTime,
	}

	dogOut := domain.Dog{
		ID:        dogID,
		UserID:    userID,
		Name:      "dog1",
		Sex:       "male",
		Age:       2,
		Breed:     "test_breed_1",
		Image:     "http://dog-images.com/test.jpg",
		CreatedAt: dogTime,
		UpdatedAt: dogTime,
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		dog domain.Dog
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      domain.Dog
		wantErr   bool
	}{
		{
			name: "execution insert query error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				dog: dogIn,
			},
			mocksInit: func() {
				mock.ExpectQuery("insert").
					WithArgs(dogIn.UserID, dogIn.Name, dogIn.Sex, dogIn.Age, dogIn.Breed, dogIn.Image).
					WillReturnError(testingError)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				dog: dogIn,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dogOut.ID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogTime, dogTime)

				mock.ExpectQuery("insert").
					WithArgs(dogIn.UserID, dogIn.Name, dogIn.Sex, dogIn.Age, dogIn.Breed, dogIn.Image).
					WillReturnRows(rows)
			},
			want:    dogOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			got, err := d.Create(tt.args.ctx, tt.args.dog)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")

	dogID := uuid.New()
	userID := uuid.New()
	dogTime := time.Now()

	dogIn := domain.Dog{
		UserID:    userID,
		Name:      "dog1",
		Sex:       "male",
		Age:       2,
		Breed:     "test_breed_1",
		Image:     "http://dog-images.com/test.jpg",
		CreatedAt: dogTime,
		UpdatedAt: dogTime,
	}

	dogOut := domain.Dog{
		ID:        dogID,
		UserID:    userID,
		Name:      "dog1",
		Sex:       "male",
		Age:       2,
		Breed:     "test_breed_1",
		Image:     "http://dog-images.com/test.jpg",
		CreatedAt: dogTime,
		UpdatedAt: dogTime,
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		uid uuid.UUID
		dog domain.Dog
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		want      domain.Dog
		wantErr   bool
	}{
		{
			name: "execution update query error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
				dog: dogIn,
			},
			mocksInit: func() {
				mock.ExpectQuery("update").
					WithArgs(dogIn.Name, dogIn.Sex, dogIn.Age, dogIn.Breed, dogIn.Image, dogID).
					WillReturnError(testingError)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
				dog: dogIn,
			},
			mocksInit: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "name", "sex", "age", "breed", "image", "created_at", "updated_at"}).
					AddRow(dogID, userID, "dog1", "male", 2, "test_breed_1", "http://dog-images.com/test.jpg", dogTime, dogTime)

				mock.ExpectQuery("update").
					WithArgs(dogIn.Name, dogIn.Sex, dogIn.Age, dogIn.Breed, dogIn.Image, dogID).
					WillReturnRows(rows)
			},
			want:    dogOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			got, err := d.Update(tt.args.ctx, tt.args.uid, tt.args.dog)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")
	dogID := uuid.New()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		uid uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		wantErr   bool
	}{
		{
			name: "execution delete query error",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
			},
			mocksInit: func() {
				mock.ExpectExec("delete").WithArgs(dogID).WillReturnError(testingError)
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
				uid: dogID,
			},
			mocksInit: func() {
				mock.ExpectExec("delete").WithArgs(dogID).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			err := d.Delete(tt.args.ctx, tt.args.uid)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDog_AddReaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testingError := errors.New("testing-error")

	likerID := uuid.New()
	likedID := uuid.New()

	inReaction := domain.Reaction{
		Liker:  likerID,
		Liked:  likedID,
		Action: domain.Like,
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx      context.Context
		reaction domain.Reaction
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
				ctx:      context.TODO(),
				reaction: inReaction,
			},
			mocksInit: func() {
				mock.ExpectExec("insert").
					WithArgs(inReaction.Liker, inReaction.Liked, inReaction.Action).
					WillReturnError(testingError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:      context.TODO(),
				reaction: inReaction,
			},
			mocksInit: func() {
				mock.ExpectExec("insert").
					WithArgs(inReaction.Liker, inReaction.Liked, inReaction.Action).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.db)
			err := d.AddReaction(tt.args.ctx, tt.args.reaction)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
