package presenters

import (
	_ "github.com/valerii-smirnov/petli-test-task/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Petly App API
//	@version		1.0
//	@description	Petly app endpoints description.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Valerii Smirnov
//	@contact.email	smirnov.valeriy90@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					As value you have to use string Bearer + 'received token after sign-in action'

func InitRoutes(engine *gin.Engine, injectors ...RoutesInjector) *gin.Engine {
	apiGroup := engine.Group("/api", ErrorHandler)
	for _, injector := range injectors {
		injector.Inject(apiGroup)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

type RoutesInjector interface {
	Inject(r gin.IRouter)
}
