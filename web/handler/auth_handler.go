package handler

import (
	"atwell/config"
	"atwell/domain"
	"atwell/infrastructure/db"
	"atwell/web"
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

	g.POST("/sign-up", h.SignUp, middleware.JWT([]byte(c.Secret)))
	g.POST("/login", h.Login, middleware.JWT([]byte(c.Secret)))
	g.POST("/before-login", h.BeforeLogin)
}

// SignUp creates account for a user.
// @Description create account for user by using yahoo japan authorization.
// @ID post-yahoo-japan-sign-up
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param code formData string true "authorization code"
// @Success 200
// @Failure 400 {object} web.ErrorResponse
// @Router /yahoo-japan/sign-up [post]
func (h AuthHandler) SignUp(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	nonce := claims["yahoo_japan_nonce"].(string)

	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "code param should not be empty",
			Code:    web.NotEnoughParameters,
		})
	}

	err := h.Usecase.SignUp(code, nonce)

	if _, ok := err.(db.DuplicateError); ok {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "user is already registered.",
			Code:    web.UserIsAlreadyRegistered,
		})
	}

	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "internal error",
			Code:    web.OtherError,
		})
	}

	return c.NoContent(http.StatusOK)
}

// Login creates session for user.
// @Description login by using yahoo japan authorization.
// @ID post-yahoo-japan-login
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param code formData string true "yahoo japan authorization code"
// @Success 200
// @Failure 400 {object} web.ErrorResponse
// @Router /yahoo-japan/login [post]
func (h AuthHandler) Login(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	nonce := claims["yahoo_japan_nonce"].(string)

	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "code param should not be empty",
			Code:    web.NotEnoughParameters,
		})
	}

	token, err := h.Usecase.Login(code, nonce)
	if errors.Is(err, db.NotFoundError{}) {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "user is not registered",
			Code:    web.UserIsNotRegistered,
		})
	}

	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "internal error",
			Code:    web.OtherError,
		})
	}

	return c.JSON(http.StatusOK, token)
}

// BeforeLogin creates temporary session for saving nonce.
// @Description creates temporary session which contains nonce for auth.
// @ID post-yahoo-japan-before-login
// @Accept  json
// @Produce  json
// @Param nonce formData string true "nonce for yahoo japan authorization"
// @Success 200
// @Failure 400 {object} web.ErrorResponse
// @Router /yahoo-japan/before-login [post]
func (h AuthHandler) BeforeLogin(c echo.Context) error {
	nonce := c.FormValue("nonce")
	if nonce == "" {
		return c.JSON(http.StatusBadRequest, "nonce should not be empty")
	}

	token, err := h.Usecase.BeforeLogin(nonce)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Message: "internal error",
			Code:    web.OtherError,
		})
	}

	return c.JSON(http.StatusOK, token)
}
