package handler

import (
	"atwell/authentication/usecase"
	"atwell/domain"
	"atwell/infrastructure"
	"errors"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
)

// AuthHandler is struct for handling http request about auth.
type AuthHandler struct {
	Usecase domain.AuthenticationUsecase
}

func HandleAuthRequest(h AuthHandler, e *echo.Echo) {
	g := e.Group("/yahoo-japan")

	g.POST("/sign-up", h.SignIn)
	g.POST("/login", h.Login)
}

// SignIn creates account for a user.
// TODO: verify email address is valid
func (h AuthHandler) SignIn(c echo.Context) error {
	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, "code param should not be empty")
	}

	err := h.Usecase.SignUp(&usecase.YahooJapanAuthenticationInformation{Token: code})

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
	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, "code param should not be empty")
	}

	token, err := h.Usecase.Login(&usecase.YahooJapanAuthenticationInformation{Token: code})
	if errors.Is(err, infrastructure.NotFoundError{}) {
		return c.JSON(http.StatusBadRequest, "user is not registered")
	}

	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, token)
}
