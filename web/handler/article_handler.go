package handler

import (
	"encoding/json"
	"golang-api/domain"
	"net/http"
)

// ArticleHandler is struct for handling http request about articles.
type ArticleHandler struct {
	Usecase domain.ArticleUsecase
}

// getAllArticles returns all articles in system.
func (h ArticleHandler) getAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Usecase.GetAll()

	if err != nil {
		// TODO
		return
	}

	json.NewEncoder(w).Encode(articles)
}

// HandleArticleRequest set up routes for requests.
func HandleArticleRequest(h ArticleHandler) {

	http.HandleFunc("/articles", h.getAllArticles)
}
