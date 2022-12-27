package presenters

import (
	"strings"

	httpErr "github.com/valerii-smirnov/petli-test-task/pkg/errors/http"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"
	"github.com/valerii-smirnov/petli-test-task/pkg/token"
	"github.com/valerii-smirnov/petli-test-task/pkg/utils/gin/resp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	AuthorizationHeaderName = "Authorization"
	bearerPrefix            = "Bearer "

	contextIdentityKey = "user-id"
)

type AuthMiddleware struct {
	tokenParser TokenParser
}

func NewAuthMiddleware(tokenParser TokenParser) *AuthMiddleware {
	return &AuthMiddleware{tokenParser: tokenParser}
}

func (m AuthMiddleware) Auth(c *gin.Context) {
	bearer := c.GetHeader(AuthorizationHeaderName)
	if bearer == "" {
		resp.AbortWithError(c, ierr.New(ierr.Unauthenticated, "missing authorization token in header"))
		return
	}

	t, err := m.tokenParser.Parse(strings.TrimPrefix(bearer, bearerPrefix))
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.Unauthenticated, err, "parsing token error"))
		return
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		resp.AbortWithError(c, ierr.New(ierr.Unauthenticated, "getting claims from token error"))
		return
	}

	suid, ok := claims[token.UserIDClaimName]
	if !ok {
		resp.AbortWithError(c, ierr.New(ierr.Unauthenticated, "getting user id from token claims error"))
		return
	}

	uid, err := uuid.Parse(suid.(string))
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.Internal, err, "parsing uuid from string error"))
		return
	}

	c.Set(contextIdentityKey, uid)

	c.Next()
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		httpErr.IErrToHttpBody(c, httpErr.ToIErr(c.Errors[0].Err))
	}
}
