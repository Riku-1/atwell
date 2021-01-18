package handler

import (
	"encoding/json"
	"golang-api/domain"
	"net/http"
	"strconv"
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

// create creates new tweet.
// @Description create new tweet.
// @ID post-tweets
// @Accept  json
// @Produce  json
// @Params comment formData string true "comment is tweet content"
// @Success 200 {object} domain.Tweet
// @Router /tweets [post]
func (h TweetHandler) create(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")

	tweet, err := h.Usecase.Create(comment)

	if err != nil {
		// TODO
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(tweet)
}

func (h TweetHandler) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		// TODO: error response
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := h.Usecase.Delete(id)
	if err != nil {
		// TODO: error response
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(nil)
}

// tweetsRouting set up "/tweets" routes for requests.
func (h TweetHandler) tweetsRouting(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.get(w, r)
	}

	if r.Method == "POST" {
		h.create(w, r)
	}

	// TODO: error response
}

// HandleTweetRequest set up routes for requests.
func HandleTweetRequest(h TweetHandler) {
	http.HandleFunc("/tweets", h.tweetsRouting)
	http.HandleFunc("/tweets/{id}", h.tweetsRouting)
}
