package handler

import (
	"atwell/config"
	"atwell/domain"
	"atwell/infrastructure"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthHandler is struct for handling http request about auth.
type AuthHandler struct {
	Usecase domain.AuthenticationUsecase
}

func HandleAuthRequest(h AuthHandler, e *echo.Echo) {
	g := e.Group("/yahoo-japan")
	c, _ := config.GetAppConfig() // TODO: constructor injection

	g.POST("/sign-up", h.SignIn, middleware.JWT([]byte(c.Secret)))
	g.POST("/login", h.Login, middleware.JWT([]byte(c.Secret)))
	g.POST("/before-login", h.BeforeLogin)
}

// SignIn creates account for a user.
func (h AuthHandler) SignIn(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	nonce := claims["yahoo_japan_nonce"].(string)

	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, "code param should not be empty")
	}

	err := h.Usecase.SignUp(code, nonce)

	if _, ok := err.(infrastructure.DuplicateError); ok {
		return c.JSON(http.StatusBadRequest, "user is already registered.")
	}

	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

// Login creates session for user.
func (h AuthHandler) Login(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	nonce := claims["yahoo_japan_nonce"].(string)

	if nonce == "" {
		return c.JSON(http.StatusBadRequest, "nonce should not be empty")
	}

	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, "code param should not be empty")
	}

	token, err := h.Usecase.Login(code, nonce)
	if errors.Is(err, infrastructure.NotFoundError{}) {
		return c.JSON(http.StatusBadRequest, "user is not registered")
	}

	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, token)
}

// BeforeLogin creates temporary session for saving nonce.
func (h AuthHandler) BeforeLogin(c echo.Context) error {
	nonce := c.FormValue("nonce")
	if nonce == "" {
		return c.JSON(http.StatusBadRequest, "nonce should not be empty")
	}

	token, err := h.Usecase.BeforeLogin(nonce)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, token)
}
