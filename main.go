package main

import (
	"fmt"
	"golang-api/web/handler"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	handler.HandleArticleRequest()
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
