package handler

import (
	"atwell/config"
	atwellDB "atwell/db"
	"atwell/domain"
	"atwell/repository"
	"atwell/usecase"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	dc, err := config.GetTestDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err = atwellDB.CreateGormDB(&dc)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getAuthHandler(db *gorm.DB) AuthHandler {
	r := repository.NewMysqlUserRepository(db)
	u := usecase.NewAuthUsecase(r)
	return AuthHandler{Usecase: u}
}

func TestAuthHandler_SignIn(t *testing.T) {
	db := getDB()
	h := getAuthHandler(db)

	// request
	e := echo.New()
	f := make(url.Values)
	email := "test_handler_sign_in"
	f.Set("email", email)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sign-in")
	err := h.SignIn(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusOK, rec.Code)

	var user domain.User
	db.First(&user, "email = ?", email)

	assert.Equal(t, email, user.Email)
	// delete record to avoid duplicate email
	db.Unscoped().Delete(&user)
}

func TestAuthHandler_SignIn_WhenEmailIsEmpty(t *testing.T) {
	db := getDB()
	h := getAuthHandler(db)

	// request
	e := echo.New()
	f := make(url.Values)
	email := "" // email is empty
	f.Set("email", email)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sign-in")
	err := h.SignIn(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
