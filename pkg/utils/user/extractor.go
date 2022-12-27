package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"
)

const contextIdentityKey = "user-id"

type IdentityExtractor struct{}

func NewIdentityExtractor() *IdentityExtractor {
	return &IdentityExtractor{}
}

func (e IdentityExtractor) ExtractFromContext(c *gin.Context) (uuid.UUID, error) {
	sUid, ok := c.Get(contextIdentityKey)
	if !ok || sUid == nil {
		return uuid.Nil, ierr.New(ierr.Internal, "getting user unique id from context error")
	}

	uid, ok := sUid.(uuid.UUID)
	if !ok {
		return uuid.Nil, ierr.New(ierr.Internal, "casting user unique id to uuid error")
	}

	return uid, nil
}
