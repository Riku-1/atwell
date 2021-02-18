package handler

import (
	"atwell/domain"
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

// SignIn is create acount for a user.
// TODO: verify email address is valid
func (h AuthHandler) SignIn(c echo.Context) error {
	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, "email param should not be empty")
	}

	_, err := h.Usecase.SignIn(email)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}
