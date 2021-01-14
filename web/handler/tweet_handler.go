package handler

import (
	"encoding/json"
	"golang-api/domain"
	"net/http"
)

// TweetHandler is struct for handling http request about tweets.
type TweetHandler struct {
	Usecase domain.TweetUsecase
}

// get returns tweets from database.
// @Description get tweets from database
// @ID get-tweets
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Tweet
// @Router /tweets [get]
func (h TweetHandler) get(w http.ResponseWriter, r *http.Request) {
	tweets, err := h.Usecase.Get()

	if err != nil {
		// TODO
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(tweets)
}

// HandleArticleRequest set up routes for requests.
func HandleTweetRequest(h TweetHandler) {
	http.HandleFunc("/tweets", h.get)
}
