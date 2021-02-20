package handler

import (
	"atwell/domain"
	"atwell/infrastructure"
	"atwell/repository"
	"atwell/usecase"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getAuthHandler(db *gorm.DB) AuthHandler {
	r := repository.NewMysqlUserRepository(db)
	u := usecase.NewAuthUsecase(r)
	return AuthHandler{Usecase: u}
}

func TestAuthHandler_SignIn(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	h := getAuthHandler(db)

	// request
	e := echo.New()
	f := make(url.Values)
	email := "test_handler_sign_in@email.com"
	f.Set("email", email)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sign-in")
	err = h.SignIn(c)
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
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
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
	err = h.SignIn(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_SignIn_DuplicateUser(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	h := getAuthHandler(db)

	// request
	e := echo.New()
	f := make(url.Values)
	email := "test_handler_sign_in_dulicate@email.com"
	f.Set("email", email)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sign-in")
	err = h.SignIn(c)
	if err != nil {
		t.Fatal(err)
	}

	// one more request by same address
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/sign-in")
	err = h.SignIn(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "user is already registered")

	// delete record to avoid duplicate email
	var user domain.User
	db.First(&user, "email = ?", email)
	db.Unscoped().Delete(&user)
}
