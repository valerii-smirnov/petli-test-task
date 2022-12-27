package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valerii-smirnov/petli-test-task/pkg/utils/user"
	"time"

	"github.com/valerii-smirnov/petli-test-task/internal/adapters"
	"github.com/valerii-smirnov/petli-test-task/internal/presenters"
	"github.com/valerii-smirnov/petli-test-task/internal/usecases"
	"github.com/valerii-smirnov/petli-test-task/pkg/db/sqlx"
	"github.com/valerii-smirnov/petli-test-task/pkg/hasher"
	"github.com/valerii-smirnov/petli-test-task/pkg/token"

	"github.com/urfave/cli/v2"
)

type applicationConfig struct {
	Port                   uint
	DBHost                 string
	DBPort                 uint
	DBUser                 string
	DBPass                 string
	DBName                 string
	JWTTokenSecret         string
	JWTTokenExpirationTime time.Duration
	PasswordSalt           string
}

type App struct {
	appConfig applicationConfig
}

func New() *App {
	return &App{}
}

func (a *App) Init() *cli.App {
	app := cli.NewApp()
	app.Commands = cli.Commands{
		&cli.Command{
			Name:        "serve",
			Usage:       "app serve",
			Description: "command runs a web server",
			Action:      a.serveAction,
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:        "port",
					Usage:       "server port {uint}",
					Destination: &a.appConfig.Port,
					Required:    true,
					EnvVars:     []string{"PORT"},
				},
				&cli.UintFlag{
					Name:        "db-port",
					Usage:       "postgres database port {uint}",
					Destination: &a.appConfig.DBPort,
					Required:    true,
					EnvVars:     []string{"DB_PORT"},
				},
				&cli.StringFlag{
					Name:        "db-host",
					Usage:       "postgres database host {string}",
					Destination: &a.appConfig.DBHost,
					Required:    true,
					EnvVars:     []string{"DB_HOST"},
				},
				&cli.StringFlag{
					Name:        "db-user",
					Usage:       "postgres database user {string}",
					Destination: &a.appConfig.DBUser,
					Required:    true,
					EnvVars:     []string{"DB_USER"},
				},
				&cli.StringFlag{
					Name:        "db-pass",
					Usage:       "postgres database password {string}",
					Destination: &a.appConfig.DBPass,
					Required:    true,
					EnvVars:     []string{"DB_PASS"},
				},
				&cli.StringFlag{
					Name:        "db-name",
					Usage:       "postgres database name {string}",
					Destination: &a.appConfig.DBName,
					Required:    true,
					EnvVars:     []string{"DB_NAME"},
				},
				&cli.StringFlag{
					Name:        "jwt-token-secret",
					Usage:       "jwt token secret {string}",
					Destination: &a.appConfig.JWTTokenSecret,
					Required:    false,
					EnvVars:     []string{"JWT_TOKEN_SECRET"},
					DefaultText: "super-secret-secret",
				},
				&cli.DurationFlag{
					Name:        "jwt-token-expiration-time",
					Usage:       "jwt token expiration time {string}",
					Destination: &a.appConfig.JWTTokenExpirationTime,
					Required:    false,
					EnvVars:     []string{"JWT_TOKEN_EXPIRATION_TIME"},
					DefaultText: "24h",
				},
				&cli.StringFlag{
					Name:        "user-password-salt",
					Usage:       "user password salt {string}",
					Destination: &a.appConfig.PasswordSalt,
					Required:    false,
					EnvVars:     []string{"USER_PASSWORD_SALT"},
					DefaultText: "super-secret-user-password-salt",
				},
			},
		},
	}

	return app
}

func (a *App) serveAction(_ *cli.Context) error {

	db, err := sqlx.NewConnection(
		a.appConfig.DBHost,
		a.appConfig.DBUser,
		a.appConfig.DBPass,
		a.appConfig.DBName,
		a.appConfig.DBPort,
		false,
	)

	if err != nil {
		return err
	}

	passwordHasher := hasher.NewMD5(a.appConfig.PasswordSalt)
	tokenProcessor := token.NewJWT(a.appConfig.JWTTokenSecret, a.appConfig.JWTTokenExpirationTime)

	userAdapter := adapters.NewUser(db)
	dogAdapter := adapters.NewDog(db)

	authUsecase := usecases.NewAuth(passwordHasher, tokenProcessor, userAdapter)
	dogUsecase := usecases.NewDog(dogAdapter)

	authMiddleware := presenters.NewAuthMiddleware(tokenProcessor)

	authPresenter := presenters.NewAuth(authUsecase)
	dogPresenter := presenters.NewDog(
		dogUsecase,
		user.NewIdentityExtractor(),
		presenters.NewUrlPagination(),
		authMiddleware.Auth,
	)

	engine := gin.New()
	presenters.InitRoutes(engine, authPresenter, dogPresenter)
	return engine.Run(fmt.Sprintf(":%d", a.appConfig.Port))
}
