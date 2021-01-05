package main

import (
	"fmt"
	"golang-api/config"
	"golang-api/repository"
	"golang-api/usecase"
	"golang-api/web/handler"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func handleRequests(db *gorm.DB) {
	au := usecase.NewArticleUsecase(repository.NewMysqlArticleRepository(db))
	ah := handler.ArticleHandler{Usecase: au}
	handler.HandleArticleRequest(ah)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

// @title golang-sample-api
// @version 1.0
// @description This is a api sample project using  golang.
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

	handleRequests(db)
}
