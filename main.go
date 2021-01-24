package main

import (
	"fmt"
	"golang-api/config"
	"golang-api/repository"
	"golang-api/usecase"
	"golang-api/web/handler"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
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
// @description Atwell is a Twitter api for one person.
func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var configuration config.Configurations
	err = viper.Unmarshal(&configuration)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configuration.Database.User,
		configuration.Database.Password,
		configuration.Database.Host,
		configuration.Database.Port,
		configuration.Database.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	e := echo.New()
	handleRequests(db, e)
}
