package handler

import (
	"atwell/domain"
	"atwell/infrastructure"
	"atwell/repository"
	"atwell/usecase"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var db *gorm.DB
var handler TweetHandler

func setup() {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		log.Fatal(err)
	}

	r := repository.NewMysqlTweetRepository(db)
	u := usecase.NewTweetUsecase(r)
	handler = TweetHandler{Usecase: u}
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func TestGet(t *testing.T) {
	// TODO: Check order
	tx := db.Begin()

	comment := "test_get"
	tw := domain.Tweet{Comment: comment}
	err := db.Create(&tw).Error

	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	tx.Commit()

	e := echo.New()
	query := make(url.Values)
	query.Set("from", "2000-01-01")
	query.Add("to", "3000-01-01")
	req := httptest.NewRequest(http.MethodGet, "/?"+query.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")
	err = handler.Get(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	var resTWList []domain.Tweet
	err = json.Unmarshal(rec.Body.Bytes(), &resTWList)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, comment, resTWList[0].Comment)
}

func TestCreate(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	comment := "test_create"
	f.Set("comment", comment)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets")
	err := handler.Create(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	var tw domain.Tweet
	err = json.Unmarshal(rec.Body.Bytes(), &tw)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, comment, tw.Comment)
}

func TestDelete(t *testing.T) {
	tx := db.Begin()

	tw := domain.Tweet{Comment: "test_delete"}
	err := db.Create(&tw).Error

	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	tx.Commit()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tweets/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(tw.ID)))
	err = handler.Delete(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	var _tw domain.Tweet
	db.Find(&_tw, tw.ID)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "", _tw.Comment)
}
