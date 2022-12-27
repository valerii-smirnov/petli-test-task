package adapters

import (
	"context"
	"database/sql"

	"github.com/valerii-smirnov/petli-test-task/internal/adapters/models"
	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Dog struct {
	db *sqlx.DB
}

func NewDog(db *sqlx.DB) *Dog {
	return &Dog{
		db: db,
	}
}

func (d Dog) List(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	query := "select * from dogs WHERE user_id != $1 order by created_at desc limit $2 offset $3"
	rows, err := d.db.QueryxContext(ctx, query, userID, pagination.PerPage, pagination.PerPage*(pagination.Page-1))
	if err != nil {
		return nil, ierr.WrapCode(ierr.Internal, err, "execution select query error")
	}

	list := make([]models.Dog, 0, 1)
	for rows.Next() {
		var dog models.Dog
		if err := rows.StructScan(&dog); err != nil {
			return nil, ierr.WrapCode(ierr.Internal, err, "struct scanning error")
		}

		list = append(list, dog)
	}

	return d.dogListToDomainDogList(list)
}

func (d Dog) Get(ctx context.Context, uid uuid.UUID) (domain.Dog, error) {
	var dog models.Dog
	query := "select * from dogs where id=$1"
	if err := d.db.GetContext(ctx, &dog, query, uid); err != nil {
		if err == sql.ErrNoRows {
			return domain.Dog{}, ierr.WrapCode(ierr.NotFound, err, "dog not found")
		}

		return domain.Dog{}, ierr.WrapCode(ierr.Internal, err, "getting dog error")
	}

	return d.dogToDomainDog(dog), nil
}

func (d Dog) Matches(ctx context.Context, dogID uuid.UUID, pagination domain.Pagination) (domain.DogList, error) {
	query := `
			select d.* from reactions r0
			inner join reactions r1 on r0.liker_id = r1.liked_id and r1.liker_id = r0.liked_id
			inner join dogs d on d.id = r1.liker_id
			where r0.liker_id = $1 AND r0.action = $2
			order by r1.created_at DESC
			limit $3 offset $4
		`

	rows, err := d.db.QueryxContext(ctx, query, dogID, domain.Like, pagination.PerPage, pagination.PerPage*(pagination.Page-1))
	if err != nil {
		return nil, ierr.WrapCode(ierr.Internal, err, "getting matches error")
	}

	list := make([]models.Dog, 0, 1)
	for rows.Next() {
		var dog models.Dog
		if err := rows.StructScan(&dog); err != nil {
			return nil, ierr.WrapCode(ierr.Internal, err, "struct scanning error")
		}

		list = append(list, dog)
	}

	return d.dogListToDomainDogList(list)
}

func (d Dog) Create(ctx context.Context, dog domain.Dog) (domain.Dog, error) {
	query := `insert into dogs 
    			(user_id, name, sex, age, breed, image) VALUES 
				($1, $2, $3, $4, $5, $6) RETURNING *`

	var mDog models.Dog
	if err := d.db.GetContext(ctx, &mDog, query, dog.UserID, dog.Name, dog.Sex.String(), dog.Age, dog.Breed, dog.Image); err != nil {
		return domain.Dog{}, ierr.WrapCode(ierr.Internal, err, "creating dog error")
	}

	return d.dogToDomainDog(mDog), nil
}

func (d Dog) Update(ctx context.Context, uid uuid.UUID, dog domain.Dog) (domain.Dog, error) {
	query := "update dogs set name=$1, sex=$2, age=$3, breed=$4, image=$5, updated_at=now() WHERE id=$6 returning *"

	var mDog models.Dog
	if err := d.db.GetContext(ctx, &mDog, query, dog.Name, dog.Sex.String(), dog.Age, dog.Breed, dog.Image, uid); err != nil {
		return domain.Dog{}, ierr.WrapCode(ierr.Internal, err, "updating dog error")
	}

	return d.dogToDomainDog(mDog), nil
}

func (d Dog) Delete(ctx context.Context, uid uuid.UUID) error {
	query := "delete from dogs where id=$1"
	if _, err := d.db.ExecContext(ctx, query, uid); err != nil {
		return ierr.WrapCode(ierr.Internal, err, "execution delete query error")
	}

	return nil
}

func (d Dog) AddReaction(ctx context.Context, reaction domain.Reaction) error {
	query := `insert into reactions (liker_id, liked_id, action, created_at) 
				values ($1, $2, $3, now()) 
				on conflict (liker_id, liked_id) do update set action=$3, created_at=now()`

	if _, err := d.db.ExecContext(ctx, query, reaction.Liker, reaction.Liked, reaction.Action); err != nil {
		return ierr.WrapCode(ierr.Internal, err, "adding reaction error")
	}

	return nil
}

func (d Dog) dogToDomainDog(dog models.Dog) domain.Dog {
	return domain.Dog{
		ID:        dog.ID,
		UserID:    dog.UserID,
		Name:      dog.Name,
		Sex:       domain.DogSex(dog.Sex),
		Age:       dog.Age,
		Breed:     dog.Breed,
		Image:     dog.Image,
		CreatedAt: dog.CreatedAt,
		UpdatedAt: dog.UpdatedAt,
	}
}

func (d Dog) dogListToDomainDogList(dogs []models.Dog) (domain.DogList, error) {
	dDogs := make(domain.DogList, 0, len(dogs))
	for _, dog := range dogs {
		dDogs = append(dDogs, d.dogToDomainDog(dog))
	}

	return dDogs, nil
}
