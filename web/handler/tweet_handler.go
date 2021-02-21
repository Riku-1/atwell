package handler

import (
	"atwell/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"

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
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param from query string true "tweets search between 'from' value and 'to' value"
// @Param to query string true "tweets search between 'from' value and 'to' value"
// @Success 200 {object} []domain.Tweet
// @Router /tweets [get]
func (h TweetHandler) Get(c echo.Context) error {
	fromString := c.QueryParam("from")
	toString := c.QueryParam("to")
	from, _ := time.ParseInLocation("2006-01-02", fromString, time.Local)
	to, _ := time.ParseInLocation("2006-01-02", toString, time.Local)

	// set by twelve o'clock midnight of the next day
	to.AddDate(0, 0, 1)

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	tweets, err := h.Usecase.Get(email, from, to)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusBadRequest)
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
	g.Use(middleware.JWT([]byte("secret")))

	g.GET("", h.Get)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}
