package handler

import (
	"encoding/json"
	"golang-api/domain"
	"golang-api/usecase"
	"net/http"
)

// articleHandler is struct for handling http request about articles.
type articleHandler struct {
	usecase domain.ArticleUsecase
}

// getAllArticles returns all articles in system.
func (h articleHandler) getAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.usecase.GetAll()

	if err != nil {
		// TODO
		return
	}

	json.NewEncoder(w).Encode(articles)
}

// HandleArticleRequest set up routes for requests.
func HandleArticleRequest() {
	h := articleHandler{
		usecase.NewArticleUsecase(),
	}

	http.HandleFunc("/articles", h.getAllArticles)
}
