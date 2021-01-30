package main

import (
	"atwell/config"
	adb "atwell/db"
	"atwell/repository"
	"atwell/usecase"
	"atwell/web/handler"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func handleRequests(db *gorm.DB, e *echo.Echo) {
	tu := usecase.NewTweetUsecase(repository.NewMysqlTweetRepository(db))
	th := handler.TweetHandler{Usecase: tu}
	handler.HandleTweetRequest(th, e)
	log.Fatal(e.Start(":10000"))
}

// @title atwell
// @version 0.1.0
// @description Atwell is a Twitter for one person.
func main() {
	dc, err := config.GetPrdDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := adb.CreateGormDB(&dc)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	handleRequests(db, e)
}
