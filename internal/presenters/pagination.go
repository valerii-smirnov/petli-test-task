package presenters

import (
	"strconv"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"

	"github.com/gin-gonic/gin"
)

const (
	pageRequestQueryParamName    = "page"
	perPageRequestQueryParamName = "per-page"

	defaultPage    = "1"
	defaultPerPage = "10"
)

type UrlPagination struct{}

func NewUrlPagination() *UrlPagination {
	return &UrlPagination{}
}

// GetPagination helper function parses pagination data from request and retuning domain.Pagination object.
func (p UrlPagination) GetPagination(c *gin.Context) (domain.Pagination, error) {
	page, err := strconv.Atoi(c.DefaultQuery(pageRequestQueryParamName, defaultPage))
	if err != nil {
		return domain.Pagination{}, ierr.WrapCode(ierr.InvalidArgument, err, "wrong page param")
	}

	perPage, err := strconv.Atoi(c.DefaultQuery(perPageRequestQueryParamName, defaultPerPage))
	if err != nil {
		return domain.Pagination{}, ierr.WrapCode(ierr.InvalidArgument, err, "wrong per-page param")
	}

	return domain.Pagination{
		Page:    page,
		PerPage: perPage,
	}, nil
}
