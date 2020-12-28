package main

import (
	"fmt"
	"golang-api/repository"
	"golang-api/usecase"
	"golang-api/web/handler"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	au := usecase.NewArticleUsecase(repository.NewDummyArticleRepository())
	ah := handler.ArticleHandler{Usecase: au}
	handler.HandleArticleRequest(ah)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

// @title golang-sample-api
// @version 1.0
// @description This is a api sample project using  golang.
func main() {
	handleRequests()
}
