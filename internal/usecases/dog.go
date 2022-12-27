package usecases

import (
	"context"

	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"

	"github.com/google/uuid"
)

type Dog struct {
	dogAdapter DogAdapter
}

func NewDog(dogAdapter DogAdapter) *Dog {
	return &Dog{
		dogAdapter: dogAdapter,
	}
}

func (d Dog) List(ctx context.Context, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	list, err := d.dogAdapter.List(ctx, dogID, pagination)
	if err != nil {
		return nil, ierr.WrapCode(ierr.Internal, err, "getting dogs list error")
	}

	return list, nil
}

func (d Dog) Get(ctx context.Context, uid uuid.UUID) (domain.Dog, error) {
	dog, err := d.dogAdapter.Get(ctx, uid)
	if err != nil {
		return domain.Dog{}, err
	}

	return dog, nil
}

func (d Dog) Matches(ctx context.Context, userID, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	dog, err := d.dogAdapter.Get(ctx, dogID)
	if err != nil {
		return nil, err
	}

	if dog.UserID != userID {
		return nil, ierr.New(ierr.PermissionDenied, "cannot get matches of not your dog")
	}

	return d.dogAdapter.Matches(ctx, dogID, pagination)
}

func (d Dog) Create(ctx context.Context, dog domain.Dog) (domain.Dog, error) {
	dog, err := d.dogAdapter.Create(ctx, dog)
	if err != nil {
		return domain.Dog{}, ierr.WrapCode(ierr.Internal, err, "creation dog error")
	}

	return dog, nil
}

func (d Dog) Update(ctx context.Context, uid uuid.UUID, dog domain.Dog) (domain.Dog, error) {
	dDog, err := d.dogAdapter.Get(ctx, uid)
	if err != nil {
		return domain.Dog{}, err
	}

	if dDog.UserID != dog.UserID {
		return domain.Dog{}, ierr.WrapCode(ierr.PermissionDenied, err, "cannot edit a dog that isn't yours")
	}

	uDog, err := d.dogAdapter.Update(ctx, uid, dog)
	if err != nil {
		return domain.Dog{}, ierr.WrapCode(ierr.Internal, err, "updating dog error")
	}

	return uDog, nil
}

func (d Dog) Delete(ctx context.Context, dogUid, userUid uuid.UUID) error {
	dDog, err := d.dogAdapter.Get(ctx, dogUid)
	if err != nil {
		return err
	}

	if dDog.UserID != userUid {
		return ierr.WrapCode(ierr.PermissionDenied, err, "cannot delete a dog that isn't yours")
	}

	if err := d.dogAdapter.Delete(ctx, dogUid); err != nil {
		return ierr.WrapCode(ierr.Internal, err, "deletion dog error")
	}

	return nil
}

func (d Dog) AddReaction(ctx context.Context, uid uuid.UUID, reaction domain.Reaction) error {
	if reaction.Liker == reaction.Liked {
		return ierr.New(ierr.InvalidArgument, "the dog can't react to itself")
	}

	liker, err := d.dogAdapter.Get(ctx, reaction.Liker)
	if err != nil {
		return ierr.WrapCode(ierr.Internal, err, "getting liker dog error")
	}

	if liker.UserID != uid {
		return ierr.New(ierr.InvalidArgument, "you're not an owner of liker dog")
	}

	if err := d.dogAdapter.AddReaction(ctx, reaction); err != nil {
		return ierr.WrapCode(ierr.Internal, err, "adding reaction error")
	}

	return nil
}
