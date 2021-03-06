package integration_test

import (
	"atwell/config"
	"atwell/domain"
	"atwell/tweet/repository"
	"atwell/tweet/usecase"
	handler2 "atwell/web/handler"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4/middleware"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func TestTweetHandler_Get(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user and tweets
	email := "tweet_handler_get_test@email.com"
	user := domain.User{
		Email: email,
		Tweets: []domain.Tweet{
			{Comment: "tweet1"},
			{Comment: "tweet2"},
		},
	}
	tx.Create(&user)

	// get token
	uu := usecase.NewAuthUsecase(repository.NewMysqlUserRepository(tx))
	token, err := uu.Login(email)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// set request
	e := echo.New()
	query := make(url.Values)
	query.Set("from", "2000-01-01")
	query.Add("to", "3000-01-01")
	req := httptest.NewRequest(http.MethodGet, "/?"+query.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(
		echo.HeaderAuthorization,
		"Bearer "+token,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Get)(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	var resTWList []domain.Tweet
	err = json.Unmarshal(rec.Body.Bytes(), &resTWList)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "tweet1", resTWList[0].Comment)
	assert.Equal(t, "tweet2", resTWList[1].Comment)
}

func TestTweetHandler_Get_NoLogin(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user and tweets
	email := "tweet_handler_get_test@email.com"
	user := domain.User{
		Email: email,
		Tweets: []domain.Tweet{
			{Comment: "tweet1"},
			{Comment: "tweet2"},
		},
	}
	tx.Create(&user)

	// set request with no token
	e := echo.New()
	query := make(url.Values)
	query.Set("from", "2000-01-01")
	query.Add("to", "3000-01-01")
	req := httptest.NewRequest(http.MethodGet, "/?"+query.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Get)(c)

	assert.Error(t, err)
}

func TestTweetHandler_Create(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user
	email := "tweet_handler_get_test@email.com"
	user := domain.User{
		Email:  email,
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user)

	// get token
	uu := usecase.NewAuthUsecase(repository.NewMysqlUserRepository(tx))
	token, err := uu.Login(email)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// set request
	e := echo.New()
	query := make(url.Values)
	testComment := "tweet_api_create_test"
	query.Set("comment", testComment)
	req := httptest.NewRequest(http.MethodPost, "/?"+query.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(
		echo.HeaderAuthorization,
		"Bearer "+token,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Create)(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	var resTweet domain.Tweet
	err = json.Unmarshal(rec.Body.Bytes(), &resTweet)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, testComment, resTweet.Comment)
}

func TestTweetHandler_Create_WithEmptyComment(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user
	email := "tweet_handler_get_test@email.com"
	user := domain.User{
		Email:  email,
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user)

	// get token
	uu := usecase.NewAuthUsecase(repository.NewMysqlUserRepository(tx))
	token, err := uu.Login(email)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// set request
	e := echo.New()
	query := make(url.Values)
	query.Set("comment", "") // empty comment
	req := httptest.NewRequest(http.MethodPost, "/?"+query.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(
		echo.HeaderAuthorization,
		"Bearer "+token,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Create)(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTweetHandler_Create_WithNoLogin(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// set request with no login
	e := echo.New()
	query := make(url.Values)
	testComment := "tweet_api_create_test"
	query.Set("comment", testComment)
	req := httptest.NewRequest(http.MethodPost, "/?"+query.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Create)(c)

	assert.NotNil(t, err)
}

func TestTweetHandler_Delete(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user and tweet
	email := "tweet_handler_get_test@email.com"
	user := domain.User{
		Email:  email,
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user)

	tweet := domain.Tweet{Comment: "test_delete", UserID: user.ID}
	tx.Create(&tweet)

	// get token
	uu := usecase.NewAuthUsecase(repository.NewMysqlUserRepository(tx))
	token, err := uu.Login(email)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// set request
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(
		echo.HeaderAuthorization,
		"Bearer "+token,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(tweet.ID)))

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Delete)(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusOK, rec.Code)

	var deletedTweet domain.Tweet
	tx.Find(&deletedTweet, tweet.ID)
	assert.NotNil(t, deletedTweet.DeletedAt)
}

func TestTweetHandler_Delete_WithNoLogin(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user and tweet
	email := "tweet_handler_get_test@email.com"
	user := domain.User{
		Email:  email,
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user)

	tweet := domain.Tweet{Comment: "test_delete", UserID: user.ID}
	tx.Create(&tweet)

	// set request with no login
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(tweet.ID)))

	// do request
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Delete)(c)

	// assertion
	assert.NotNil(t, err)
}

func TestTweetHandler_Delete_ByNotOwner(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	u := usecase.NewTweetUsecase(
		repository.NewMysqlTweetRepository(tx),
		repository.NewMysqlUserRepository(tx),
	)
	handler := handler2.TweetHandler{Usecase: u}

	// create user1 and tweet
	user1 := domain.User{
		Email:  "tweet_handler_get_test@email.com",
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user1)

	tweet := domain.Tweet{Comment: "test_delete", UserID: user1.ID}
	tx.Create(&tweet)

	noOwnerEmail := "no_owner@email.com"
	user2 := domain.User{
		Email: noOwnerEmail,
	}
	tx.Create(&user2)

	// get token by user2
	uu := usecase.NewAuthUsecase(repository.NewMysqlUserRepository(tx))
	token, err := uu.Login(noOwnerEmail)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// set request
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(
		echo.HeaderAuthorization,
		"Bearer "+token,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(tweet.ID)))

	// do request for user1's tweet by user2
	err = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})(handler.Delete)(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
