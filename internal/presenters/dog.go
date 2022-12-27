package presenters

import (
	"net/http"

	"github.com/valerii-smirnov/petli-test-task/internal/domain"
	"github.com/valerii-smirnov/petli-test-task/internal/presenters/messages"
	"github.com/valerii-smirnov/petli-test-task/pkg/errors/ierr"
	"github.com/valerii-smirnov/petli-test-task/pkg/utils/gin/resp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Dog presenter.
type Dog struct {
	dogUsecase        DogUsecase
	identityExtractor IdentityExtractor
	paginator         Paginator

	middlewares []gin.HandlerFunc
}

// NewDog constructor.
func NewDog(
	dogUsecase DogUsecase,
	identityExtractor IdentityExtractor,
	paginator Paginator,
	middlewares ...gin.HandlerFunc,
) *Dog {
	return &Dog{
		dogUsecase:        dogUsecase,
		identityExtractor: identityExtractor,
		paginator:         paginator,
		middlewares:       middlewares,
	}
}

// Inject Injector implementation.
func (d Dog) Inject(r gin.IRouter) {
	dogsGroup := r.Group("/dog")
	if len(d.middlewares) > 0 {
		dogsGroup.Use(d.middlewares...)
	}

	dogsGroup.GET("", d.List)
	dogsGroup.GET("/:id", d.Get)
	dogsGroup.GET("/:id/matches", d.Matches)
	dogsGroup.POST("", d.Create)
	dogsGroup.PUT("/:id", d.Update)
	dogsGroup.DELETE("/:id", d.Delete)
	dogsGroup.POST("/reaction", d.Reaction)
}

// List http handler func to retrieve list of dogs.
// @Summary      Dogs list
// @Description  Getting dogs list
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 page query string false "pagination page number"
// @Param 		 per-page query string false "pagination per page items number"
// @Success      200 {object} messages.DogListResponseBody
// @Failure      400  {object}  messages.BadRequestError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog [get]
func (d Dog) List(c *gin.Context) {
	uid, err := d.identityExtractor.ExtractFromContext(c)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	pag, err := d.paginator.GetPagination(c)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	list, err := d.dogUsecase.List(c, uid, pag)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, d.domainDogListToMessageList(list))
}

// Get http handler func to get dog by ID.
// @Summary      Dogs list
// @Description  Getting dogs list
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 id path string true "dog ID"
// @Success      200 {object} messages.DogResponseBody
// @Failure      400  {object}  messages.BadRequestError
// @Failure      404  {object}  messages.NotFoundError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog/{id} [get]
func (d Dog) Get(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.InvalidArgument, err, "wrong dog id param"))
		return
	}

	dog, err := d.dogUsecase.Get(c, uid)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, d.domainDogToMessage(dog))
}

// Matches http handler func to get all matches for provided dog.
// @Summary      Dog matches
// @Description  Getting dog matches with another dogs
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 id path string true "dog ID"
// @Param 		 page query string false "pagination page number"
// @Param 		 per-page query string false "pagination per page items number"
// @Success      200 {object} messages.DogListResponseBody
// @Failure      400  {object}  messages.BadRequestError
// @Failure      404  {object}  messages.NotFoundError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog/{id}/matches [get]
func (d Dog) Matches(c *gin.Context) {
	pag, err := d.paginator.GetPagination(c)
	if err != nil {
		resp.AbortWithError(c, err)
	}

	dogUid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.Internal, err, "error parsing dog id param"))
		return
	}

	userUid, err := d.identityExtractor.ExtractFromContext(c)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	list, err := d.dogUsecase.Matches(c, userUid, dogUid, pag)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, d.domainDogListToMessageList(list))
}

// Create http handler func to create new dog.
// @Summary      Create dog
// @Description  Creates new dog
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 input body messages.CreateOrUpdateDogRequestBody true "dog object body"
// @Success      200 {object} messages.DogResponseBody
// @Failure      400  {object}  messages.BadRequestError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog [post]
func (d Dog) Create(c *gin.Context) {
	var req messages.CreateOrUpdateDogRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	uid, err := d.identityExtractor.ExtractFromContext(c)
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.Internal, err, "getting user id error"))
		return
	}

	newDog := domain.Dog{
		UserID: uid,
		Name:   req.Name,
		Sex:    domain.DogSex(req.Sex),
		Age:    req.Age,
		Breed:  req.Breed,
		Image:  req.Image,
	}

	dog, err := d.dogUsecase.Create(c, newDog)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, d.domainDogToMessage(dog))
}

// Update http handler func to update dog.
// @Summary      Dog update
// @Description  Updates existing dog
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 id path string true "dog ID"
// @Param 		 input body messages.CreateOrUpdateDogRequestBody true "dog object body"
// @Success      200 {object} messages.DogResponseBody
// @Failure      400  {object}  messages.BadRequestError
// @Failure      404  {object}  messages.NotFoundError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog/{id} [put]
func (d Dog) Update(c *gin.Context) {
	dogUid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.InvalidArgument, err, "wrong dog id"))
		return
	}

	var req messages.CreateOrUpdateDogRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	uid, err := d.identityExtractor.ExtractFromContext(c)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	newDog := domain.Dog{
		UserID: uid,
		Name:   req.Name,
		Sex:    domain.DogSex(req.Sex),
		Age:    req.Age,
		Breed:  req.Breed,
		Image:  req.Image,
	}

	dog, err := d.dogUsecase.Update(c, dogUid, newDog)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, d.domainDogToMessage(dog))
}

// Delete http handler func to delete tog.
// @Summary      Dog update
// @Description  Updates existing dog
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 id path string true "dog ID"
// @Success      204
// @Failure      400  {object}  messages.BadRequestError
// @Failure      404  {object}  messages.NotFoundError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog/{id} [delete]
func (d Dog) Delete(c *gin.Context) {
	dogUid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		resp.AbortWithError(c, ierr.WrapCode(ierr.InvalidArgument, err, "wrong dog id"))
		return
	}

	uid, err := d.identityExtractor.ExtractFromContext(c)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	if err := d.dogUsecase.Delete(c, dogUid, uid); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// Reaction http handler func to save reaction of one dog to another.
// @Summary      Reaction
// @Description  React to another dog
// @Tags         dogs
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param 		 input body messages.ReactionRequestBody true "reaction body"
// @Success      204
// @Failure      400  {object}  messages.BadRequestError
// @Failure      404  {object}  messages.NotFoundError
// @Failure      500  {object}  messages.InternalServerError
// @Router       /dog/reaction [post]
func (d Dog) Reaction(c *gin.Context) {
	var req messages.ReactionRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	reaction := domain.Reaction{
		Liker:  uuid.MustParse(req.Liker),
		Liked:  uuid.MustParse(req.Liked),
		Action: domain.Action(req.Action),
	}

	uid, err := d.identityExtractor.ExtractFromContext(c)
	if err != nil {
		resp.AbortWithError(c, err)
		return
	}

	if err := d.dogUsecase.AddReaction(c, uid, reaction); err != nil {
		resp.AbortWithError(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (d Dog) domainDogToMessage(dog domain.Dog) messages.DogResponseBody {
	return messages.DogResponseBody{
		ID:    dog.ID.String(),
		Name:  dog.Name,
		Sex:   dog.Sex.String(),
		Age:   dog.Age,
		Breed: dog.Breed,
		Image: dog.Image,
	}
}

func (d Dog) domainDogListToMessageList(dogs []domain.Dog) messages.DogListResponseBody {
	list := make(messages.DogListResponseBody, 0, len(dogs))
	for _, dog := range dogs {
		list = append(list, d.domainDogToMessage(dog))
	}

	return list
}
