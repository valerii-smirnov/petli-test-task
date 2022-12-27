package models

import (
	"time"

	"github.com/google/uuid"
)

type Dog struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Sex       string    `db:"sex"`
	Age       uint      `db:"age"`
	Breed     string    `db:"breed"`
	Image     string    `db:"image"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
