package main

import (
	usecase2 "atwell/authentication/usecase"
	"atwell/config"
	"atwell/infrastructure"
	infrastructure2 "atwell/infrastructure/api"
	"atwell/repository"
	"atwell/usecase"
	"atwell/web/handler"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func handleRequests(db *gorm.DB, e *echo.Echo) {
	tu := usecase.NewTweetUsecase(repository.NewMysqlTweetRepository(db), repository.NewMysqlUserRepository(db))
	th := handler.TweetHandler{Usecase: tu}
	handler.HandleTweetRequest(th, e)

	c, err := config.GetYahooAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	yEmailUsecase := usecase2.NewYahooJapanGetUserEmailUsecase(
		&infrastructure2.YahooJapanAuthAPI{Conf: c},
	)
	userRepo := repository.NewMysqlUserRepository(db)
	authUsecase := usecase2.NewAuthenticationUsecase(yEmailUsecase, userRepo)
	uh := handler.AuthHandler{Usecase: authUsecase}
	handler.HandleAuthRequest(uh, e)

	log.Fatal(e.Start(":10000"))
}

// @title atwell
// @version 0.1.0
// @description Atwell is a Twitter for one person.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	db, err := infrastructure.GetPrdGormDB()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	handleRequests(db, e)
}
