package handler

import (
	"atwell/domain"
	"atwell/infrastructure"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthHandler is struct for handling http request about auth.
type AuthHandler struct {
	Usecase domain.UserUsecase
}

func HandleAuthRequest(h AuthHandler, e *echo.Echo) {
	e.POST("/sign-in", h.SignIn)
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
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

// Login creates session for user.
//func (h AuthHandler) Login(c echo.Context) error {
//	// get mail from form
//
//	// get user by email
//
//	// create session
//
//	// return response of jwt
//}
