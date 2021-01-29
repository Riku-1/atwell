package main

import (
	"fmt"
	"golang-api/config"
	"golang-api/repository"
	"golang-api/usecase"
	"golang-api/web/handler"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
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
	var dc config.DatabaseConfigurations
	err := envconfig.Process("atwell_db", &dc)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	e := echo.New()
	handleRequests(db, e)
}
