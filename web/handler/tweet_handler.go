package handler

import (
	"golang-api/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
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
func (h TweetHandler) get(c echo.Context) error {
	// TODO: when from and to is empty
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	_from, _ := time.ParseInLocation("2006-01-02", from, time.Local)
	_to, _ := time.ParseInLocation("2006-01-02", to, time.Local)
	// set by twelve o'clock midnight of the next day.
	tweets, err := h.Usecase.Get(_from, _to.AddDate(0, 0, 1))

	if err != nil {
		// TODO
		panic(err)
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, tweets)
}

// create creates new tweet.
// @Description create new tweet.
// @ID post-tweets
// @Accept  json
// @Produce  json
// @Params comment formData string true "comment is tweet content"
// @Success 200 {object} domain.Tweet
// @Router /tweets [post]
func (h TweetHandler) create(c echo.Context) error {
	comment := c.FormValue("comment")

	tweet, err := h.Usecase.Create(comment)

	if err != nil {
		// TODO
		panic(err)
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, tweet)
}

// delete deletes new tweet.
// @Description delete new tweet.
// @ID delete-tweets-id
// @Accept  json
// @Produce  json
// @Success 200 "OK"
// @Router /tweets/{id} [delete]
func (h TweetHandler) delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	err := h.Usecase.Delete(id)
	if err != nil {
		// TODO: error response
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, nil)
}

// HandleTweetRequest set up routes for requests.
func HandleTweetRequest(h TweetHandler, e *echo.Echo) {
	g := e.Group("/tweets")

	g.GET("", h.get)
	g.POST("", h.create)
	g.DELETE("/:id", h.delete)
}
