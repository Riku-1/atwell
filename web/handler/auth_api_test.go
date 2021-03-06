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
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	h := getAuthHandler(tx)

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
	err = h.SignUp(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusOK, rec.Code)

	var user domain.User
	tx.First(&user, "email = ?", email)
	assert.Equal(t, email, user.Email)
}

func TestAuthHandler_SignIn_WhenEmailIsEmpty(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	h := getAuthHandler(tx)

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
	err = h.SignUp(c)
	if err != nil {
		tx.Rollback()
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
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	h := getAuthHandler(tx)

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
	err = h.SignUp(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// one more request by same address
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/sign-in")
	err = h.SignUp(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "user is already registered")
}

func TestAuthHandler_Login(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	h := getAuthHandler(tx)

	// create user before login
	email := "test_auth_handler_login@email.com"
	user := domain.User{Email: email}
	tx.Create(&user)

	// request
	e := echo.New()
	HandleAuthRequest(h, e)
	f := make(url.Values)
	f.Set("email", email)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")
	err = h.Login(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusOK, rec.Code)

	var token string
	err = json.Unmarshal(rec.Body.Bytes(), &token)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	assert.NotEqual(t, "", token)
}

func TestAuthHandler_Login_WhenEmailIsEmpty(t *testing.T) {
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
	c.SetPath("/login")
	err = h.Login(c)
	if err != nil {
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Login_ByNoRegisteredUser(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	h := getAuthHandler(tx)

	// request
	e := echo.New()
	HandleAuthRequest(h, e)
	f := make(url.Values)
	f.Set("email", "no_registered_user@email.com") // login by no registered user

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")
	err = h.Login(c)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assertion
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var token string
	err = json.Unmarshal(rec.Body.Bytes(), &token)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	assert.Contains(t, rec.Body.String(), "user is not registered")
}
