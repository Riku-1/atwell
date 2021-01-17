package handler

import (
	"encoding/json"
	"golang-api/domain"
	"net/http"
	"time"
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
// @Params from query string true "tweets search between 'from' value and 'to' value"
// @Params to query string true "tweets search between 'from' value and 'to' value"
// @Success 200 {object} []domain.Tweet
// @Router /tweets [get]
func (h TweetHandler) get(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	_from, _ := time.ParseInLocation("2006-01-02", from, time.Local)
	_to, _ := time.ParseInLocation("2006-01-02", to, time.Local)
	// set by twelve o'clock midnight of the next day.
	tweets, err := h.Usecase.Get(_from, _to.AddDate(0, 0, 1))

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
