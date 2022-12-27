package messages

import (
	"net/http"

	"github.com/valerii-smirnov/petli-test-task/internal/presenters/messages"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

const (
	InternalServerErrorDefaultText  = "something went wrong"
	UnauthenticatedErrorDefaultText = "unauthenticated"
	BadRequestErrorDefaultText      = "validation error"
)

func ToIErr(err error) *ierr.Error {
	switch e := err.(type) {
	case *ierr.Error:
		return e
	case validator.ValidationErrors:
		return fromValidationError(e)
	case *jwt.ValidationError:
		return fromJWTTokenValidationError(e)
	default:
		return ierr.WrapCode(ierr.Internal, e, InternalServerErrorDefaultText)
	}
}

func IErrToHttpBody(c *gin.Context, err *ierr.Error) {
	switch err.Code() {
	case ierr.Internal:
		c.JSON(http.StatusInternalServerError, messages.NewInternalServerError(InternalServerErrorDefaultText, err))
	case ierr.NotFound:
		c.JSON(http.StatusNotFound, messages.NewNotFoundError(err.Message()))
	case ierr.InvalidArgument:
		c.JSON(http.StatusBadRequest, messages.NewBadRequestError(err.Message(), err.Props()))
	case ierr.PermissionDenied:
		c.JSON(http.StatusForbidden, messages.NewForbiddenError(err.Message()))
	case ierr.Unauthenticated:
		c.JSON(http.StatusUnauthorized, messages.NewUnauthenticatedError(UnauthenticatedErrorDefaultText))
	case ierr.AlreadyExists:
		c.JSON(http.StatusConflict, messages.NewConflictError(err.Message()))
	default:
		c.JSON(http.StatusInternalServerError, messages.NewInternalServerError(InternalServerErrorDefaultText, err))
	}
}

func fromValidationError(vErrs validator.ValidationErrors) *ierr.Error {
	iErr := ierr.WrapCode(ierr.InvalidArgument, vErrs, BadRequestErrorDefaultText)

	props := make(ierr.KV)
	for _, vErr := range vErrs {
		props[vErr.Field()] = vErr.Error()
	}

	return iErr.SetProps(props)
}

func fromJWTTokenValidationError(validationError *jwt.ValidationError) *ierr.Error {
	if validationError.Is(jwt.ErrTokenExpired) {
		return ierr.WrapCode(ierr.Unauthenticated, validationError, "token is expired")
	}

	return ierr.WrapCode(ierr.Unauthenticated, validationError, UnauthenticatedErrorDefaultText)
}
