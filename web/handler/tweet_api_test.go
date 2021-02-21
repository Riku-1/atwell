package handler

import (
	"atwell/domain"
	"atwell/infrastructure"
	"atwell/repository"
	"atwell/usecase"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4/middleware"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func setup() {
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func TestTweetHandler_Get(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
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
	handler := TweetHandler{Usecase: u}

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
	db, err := infrastructure.GetDevGormDB()
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
	handler := TweetHandler{Usecase: u}

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
	db, err := infrastructure.GetDevGormDB()
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
	handler := TweetHandler{Usecase: u}

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
	db, err := infrastructure.GetDevGormDB()
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
	handler := TweetHandler{Usecase: u}

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
	db, err := infrastructure.GetDevGormDB()
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
	handler := TweetHandler{Usecase: u}

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

//
//func TestDelete(t *testing.T) {
//	tx := db.Begin()
//
//	tw := domain.Tweet{Comment: "test_delete"}
//	err := db.Create(&tw).Error
//
//	if err != nil {
//		tx.Rollback()
//		t.Fatal(err)
//	}
//	tx.Commit()
//
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodDelete, "/", nil)
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//	c.SetPath("/tweets/:id")
//	c.SetParamNames("id")
//	c.SetParamValues(strconv.Itoa(int(tw.ID)))
//	err = handler.Delete(c)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// assertion
//	var _tw domain.Tweet
//	db.Find(&_tw, tw.ID)
//	assert.Equal(t, http.StatusOK, rec.Code)
//	assert.Equal(t, "", _tw.Comment)
//}
