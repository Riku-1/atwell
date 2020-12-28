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

func main() {
	handleRequests()
}
