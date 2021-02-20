package handler

import (
	"atwell/domain"
	"atwell/infrastructure"
	"errors"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
)

// AuthHandler is struct for handling http request about auth.
type AuthHandler struct {
	Usecase domain.UserUsecase
}

func HandleAuthRequest(h AuthHandler, e *echo.Echo) {
	e.POST("/sign-in", h.SignIn)
	e.POST("/login", h.Login)
}

// SignIn creates account for a user.
// TODO: verify email address is valid
func (h AuthHandler) SignIn(c echo.Context) error {
	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, "email param should not be empty")
	}

	_, err := h.Usecase.SignIn(email)

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
	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, "email param should not be empty")
	}

	token, err := h.Usecase.Login(email)
	if errors.Is(err, infrastructure.NotFoundError{}) {
		return c.JSON(http.StatusBadRequest, "user is not registered")
	}

	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, token)
}
