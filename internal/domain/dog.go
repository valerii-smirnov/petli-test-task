package domain

import (
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	Like    Action = "like"
	Dislike Action = "dislike"
)

type DogSex string

func (s DogSex) String() string {
	return string(s)
}

type Dog struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Sex       DogSex
	Age       uint
	Breed     string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Pagination struct {
	Page    int
	PerPage int
}

type DogList []Dog

type Reaction struct {
	Liker  uuid.UUID
	Liked  uuid.UUID
	Action Action
}
