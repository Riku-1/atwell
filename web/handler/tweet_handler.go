package handler

import (
	"atwell/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// TweetHandler is struct for handling http request about tweets.
type TweetHandler struct {
	Usecase domain.TweetUsecase
}

// Get returns tweets from database.
// @Description Get tweets from database
// @ID get-tweets
// @Accept  json
// @Produce  json
// @Params from query string true "tweets search between 'from' value and 'to' value"
// @Params to query string true "tweets search between 'from' value and 'to' value"
// @Success 200 {object} []domain.Tweet
// @Router /tweets [get]
func (h TweetHandler) Get(c echo.Context) error {
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

// Create creates new tweet.
// @Description create new tweet.
// @ID post-tweets
// @Accept  json
// @Produce  json
// @Params comment formData string true "comment is tweet content"
// @Success 200 {object} domain.Tweet
// @Router /tweets [post]
func (h TweetHandler) Create(c echo.Context) error {
	comment := c.FormValue("comment")

	tweet, err := h.Usecase.Create(comment)

	if err != nil {
		// TODO
		panic(err)
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, tweet)
}

// Delete deletes new tweet.
// @Description delete new tweet.
// @ID delete-tweets-id
// @Accept  json
// @Produce  json
// @Success 200 "OK"
// @Router /tweets/{id} [delete]
func (h TweetHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
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

	g.GET("", h.Get)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}
