package usecases

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"testing"
	"time"
)

func TestDog_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	testErr := errors.New("testing error")

	dog := domain.Dog{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Name:      "test name",
		Sex:       "test sex",
		Age:       3,
		Breed:     "test breed",
		Image:     "http://test-image/image.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	listDog := domain.DogList{dog, dog}

	pag := domain.Pagination{
		Page:    1,
		PerPage: 2,
	}

	userID := uuid.New()

	type fields struct {
		dogAdapter DogAdapter
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
			name: "getting_list_error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				pagination: pag,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().List(gomock.Any(), gomock.Eq(userID), gomock.Eq(pag)).Return(nil, testErr)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				pagination: pag,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().List(gomock.Any(), gomock.Eq(userID), gomock.Eq(pag)).Return(listDog, nil)
			},
			want:    listDog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			got, err := d.List(tt.args.ctx, tt.args.userID, tt.args.pagination)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Get(t *testing.T) {
	type fields struct {
		dogAdapter DogAdapter
	}
	type args struct {
		ctx context.Context
		uid uuid.UUID
	}

	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	testErr := errors.New("testing error")

	dogUuid := uuid.New()

	dog := domain.Dog{
		ID:        dogUuid,
		UserID:    uuid.New(),
		Name:      "test name",
		Sex:       "test sex",
		Age:       3,
		Breed:     "test breed",
		Image:     "http://test-image/image.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
			name: "getting_dog_error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				uid: dogUuid,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogUuid)).Return(domain.Dog{}, testErr)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				uid: dogUuid,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogUuid)).Return(dog, nil)
			},
			want:    dog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			got, err := d.Get(tt.args.ctx, tt.args.uid)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Matches(t *testing.T) {
	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	dogID := uuid.New()
	userID := uuid.New()

	testErr := errors.New("testing error")

	wrongDog := domain.Dog{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	goodDog := domain.Dog{
		ID:     dogID,
		UserID: userID,
	}

	pag := domain.Pagination{
		Page:    1,
		PerPage: 5,
	}

	goodDogs := domain.DogList{goodDog, goodDog}

	//listDog := domain.DogList{dog, dog}

	type fields struct {
		dogAdapter DogAdapter
	}
	type args struct {
		ctx        context.Context
		userID     uuid.UUID
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
			name: "getting dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				dogID:      dogID,
				pagination: pag,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(domain.Dog{}, testErr)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "getting matches for not your dog",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				dogID:      dogID,
				pagination: pag,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(wrongDog, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:        context.TODO(),
				userID:     userID,
				dogID:      dogID,
				pagination: pag,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(goodDog, nil)
				dogAdapterMock.EXPECT().Matches(gomock.Any(), gomock.Eq(dogID), gomock.Eq(pag)).Return(goodDogs, nil)
			},
			want:    goodDogs,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			got, err := d.Matches(tt.args.ctx, tt.args.userID, tt.args.dogID, tt.args.pagination)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	testError := errors.New("testing-error")
	dogID := uuid.New()
	userID := uuid.New()
	dogName := "test-name"

	dog := domain.Dog{
		UserID: userID,
		Name:   dogName,
	}

	createdDog := domain.Dog{
		ID:     dogID,
		UserID: userID,
		Name:   dogName,
	}

	type fields struct {
		dogAdapter DogAdapter
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
			name: "creation dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				dog: dog,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Create(gomock.Any(), gomock.Eq(dog)).Return(domain.Dog{}, testError)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				dog: dog,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Create(gomock.Any(), gomock.Eq(dog)).Return(createdDog, nil)
			},
			want:    createdDog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			got, err := d.Create(tt.args.ctx, tt.args.dog)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	testError := errors.New("testing-error")
	userID := uuid.New()
	dogID := uuid.New()

	dogIn := domain.Dog{
		UserID: userID,
		Name:   "test-name",
	}

	dogOut := domain.Dog{
		ID:     dogID,
		UserID: userID,
		Name:   "test-name",
	}

	wrongFoundDog := domain.Dog{
		ID:     dogID,
		UserID: uuid.New(),
		Name:   "wrong found dog",
	}

	type fields struct {
		dogAdapter DogAdapter
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
			name: "getting dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
				dog: dogIn,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(domain.Dog{}, testError)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "updating not your dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
				dog: dogIn,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(wrongFoundDog, nil)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "updating dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
				dog: dogIn,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(dogOut, nil)
				dogAdapterMock.EXPECT().Update(gomock.Any(), gomock.Eq(dogID), dogIn).Return(domain.Dog{}, testError)
			},
			want:    domain.Dog{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx: context.TODO(),
				uid: dogID,
				dog: dogIn,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(dogOut, nil)
				dogAdapterMock.EXPECT().Update(gomock.Any(), gomock.Eq(dogID), dogIn).Return(dogOut, nil)
			},
			want:    dogOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			got, err := d.Update(tt.args.ctx, tt.args.uid, tt.args.dog)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDog_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	testError := errors.New("testing-error")
	dogID := uuid.New()
	userID := uuid.New()

	wrongDogOut := domain.Dog{
		ID:     dogID,
		UserID: uuid.New(),
	}

	dogOut := domain.Dog{
		ID:     dogID,
		UserID: userID,
	}

	type fields struct {
		dogAdapter DogAdapter
	}
	type args struct {
		ctx     context.Context
		dogUid  uuid.UUID
		userUid uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mocksInit func()
		wantErr   bool
	}{
		{
			name: "getting dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:     context.TODO(),
				dogUid:  dogID,
				userUid: userID,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(domain.Dog{}, testError)
			},
			wantErr: true,
		},
		{
			name: "deleting not your dog",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:     context.TODO(),
				dogUid:  dogID,
				userUid: userID,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(wrongDogOut, nil)
			},
			wantErr: true,
		},
		{
			name: "deleting error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:     context.TODO(),
				dogUid:  dogID,
				userUid: userID,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(dogOut, nil)
				dogAdapterMock.EXPECT().Delete(gomock.Any(), gomock.Eq(dogID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:     context.TODO(),
				dogUid:  dogID,
				userUid: userID,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(dogID)).Return(dogOut, nil)
				dogAdapterMock.EXPECT().Delete(gomock.Any(), gomock.Eq(dogID)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			err := d.Delete(tt.args.ctx, tt.args.dogUid, tt.args.userUid)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDog_AddReaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	dogAdapterMock := NewMockDogAdapter(ctrl)

	testError := errors.New("testing-error")

	userID := uuid.New()

	likerID := uuid.New()
	likedID := uuid.New()

	reactionToItself := domain.Reaction{
		Liker:  likerID,
		Liked:  likerID,
		Action: domain.Like,
	}

	reactionNotYourDog := domain.Reaction{
		Liker:  uuid.New(),
		Liked:  likedID,
		Action: domain.Like,
	}

	correctReaction := domain.Reaction{
		Liker:  likerID,
		Liked:  likedID,
		Action: domain.Like,
	}

	wrongDog := domain.Dog{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	correctDog := domain.Dog{
		ID:     uuid.New(),
		UserID: userID,
	}

	type fields struct {
		dogAdapter DogAdapter
	}
	type args struct {
		ctx      context.Context
		uid      uuid.UUID
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
			name: "like to itself",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:      context.TODO(),
				uid:      userID,
				reaction: reactionToItself,
			},
			mocksInit: func() {},
			wantErr:   true,
		},
		{
			name: "getting dog error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:      context.TODO(),
				uid:      userID,
				reaction: reactionNotYourDog,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(reactionNotYourDog.Liker)).Return(domain.Dog{}, testError)
			},
			wantErr: true,
		},
		{
			name: "owner error",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:      context.TODO(),
				uid:      userID,
				reaction: reactionNotYourDog,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(reactionNotYourDog.Liker)).Return(wrongDog, nil)
			},
			wantErr: true,
		},
		{
			name: "error adding reaction",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:      context.TODO(),
				uid:      userID,
				reaction: correctReaction,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(correctReaction.Liker)).Return(correctDog, nil)
				dogAdapterMock.EXPECT().AddReaction(gomock.Any(), gomock.Eq(correctReaction)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				dogAdapter: dogAdapterMock,
			},
			args: args{
				ctx:      context.TODO(),
				uid:      userID,
				reaction: correctReaction,
			},
			mocksInit: func() {
				dogAdapterMock.EXPECT().Get(gomock.Any(), gomock.Eq(correctReaction.Liker)).Return(correctDog, nil)
				dogAdapterMock.EXPECT().AddReaction(gomock.Any(), gomock.Eq(correctReaction)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksInit()

			d := NewDog(tt.fields.dogAdapter)
			err := d.AddReaction(tt.args.ctx, tt.args.uid, tt.args.reaction)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
