package main

import (
	"golang-api/repository"
	"golang-api/usecase"
	"golang-api/web/handler"
	"log"
	"net/http"

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
	dsn := "root:example@tcp(127.0.0.1:3306)/sample_project?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	handleRequests(db)
}
