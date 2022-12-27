package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/internal/presenters/messages"
	"github.com/valerii-smirnov/petli-test-task/pkg/utils/gin/resp"
)

type Auth struct {
	authUsecase AuthUsecase

	middlewares []gin.HandlerFunc
}

func NewAuth(authUsecase AuthUsecase, middlewares ...gin.HandlerFunc) *Auth {
	return &Auth{
		authUsecase: authUsecase,
		middlewares: middlewares,
	}
}

func (a Auth) Inject(r gin.IRouter) {
	authGroup := r.Group("/auth")
	if len(a.middlewares) > 0 {
		authGroup.Use(a.middlewares...)
	}

	authGroup.POST("sign-in", a.SignIn)
	authGroup.POST("sign-up", a.SignUp)
}

// SignUp godoc
// @Summary      User registration
// @Description  User registration endpoint
// @ID 			 Create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 		 input body messages.SignUpRequestBody true "sign up info"
// @Success      204
// @Failure      400  {object}  messages.BadRequestError
// @Failure      409  {object}  messages.ConflictError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /auth/sign-up [post]
func (a Auth) SignUp(c *gin.Context) {
	var req messages.SignUpRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	domainSignUP := domain.SignUp{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := a.authUsecase.SignUp(c, domainSignUP); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// SignIn godoc
// @Summary      User login
// @Description  User login endpoint
// @ID 			 Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 		 input body messages.SignInRequestBody true "sign in info"
// @Success      200 {object} messages.SignInResponseBody
// @Failure      400  {object}  messages.BadRequestError
// @Failure      404  {object}  messages.NotFoundError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /auth/sign-in [post]
func (a Auth) SignIn(c *gin.Context) {
	var req messages.SignInRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	domainSignIn := domain.SingIn{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := a.authUsecase.SignIn(c, domainSignIn)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	rb := messages.SignInResponseBody{
		Token: string(token),
	}

	c.JSON(http.StatusOK, rb)
}
