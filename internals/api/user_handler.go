package api

import (
	"log"
	"net/http"
	"time"

	"github.com/erfan-goodarzi/booking-event-api/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/store"
	"github.com/erfan-goodarzi/booking-event-api/pkg/validation"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user     store.UserStore
	logger   *log.Logger
	response *APIResponse
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger, response *APIResponse) *UserHandler {
	return &UserHandler{
		userStore,
		logger,
		response,
	}
}

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body store.User true "User signup payload"
// @Success 201 {object} store.User
// @Failure 409 {object} api.ErrorConflict
// @Failure 422 {object} api.ErrorValidation
// @Failure 500 {object} api.ErrorInternalServer
// @Router /auth/signup [post]
func (handler *UserHandler) Signup(c *gin.Context) {
	var user store.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = validation.Validate.Struct(user)

	if err != nil {
		handler.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	err = handler.user.Create(&user)

	if err != nil {
		switch err.Error() {
		case "EMAIL_ALREADY_EXISTS":
			handler.response.RespondError(c, http.StatusConflict, "EMAIL_ALREADY_EXISTS")
		case "USERNAME_ALREADY_EXISTS":
			handler.response.RespondError(c, http.StatusConflict, "USERNAME_ALREADY_EXISTS")
		default:
			handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		}
		return
	}

	handler.response.RespondSuccess(c, http.StatusCreated, messages.Signup, user)
}

// Login godoc
// @Summary Authenticate user
// @Description Login with email and password, returns access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body store.User true "Login credentials (email and password)"
// @Success 200 {object} api.LoginResponse
// @Failure 401 {object} api.ErrorUnauthorized
// @Failure 500 {object} api.ErrorInternalServer
// @Failure 422 {object} api.ErrorValidation
// @Router /auth/login [post]
func (handler *UserHandler) Login(c *gin.Context) {
	var user store.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = validation.Validate.Struct(user)

	if err != nil {
		handler.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	err = handler.user.ValidateCredential(&user)

	if err != nil {
		if err.Error() == "INVALID_CREDENTIAL" {
			handler.response.RespondError(c, http.StatusNonAuthoritativeInfo, "UNAUTHORIZED")
			return
		}
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	tokens, err := apiUtils.GenerateToken(user.Email, user.ID)

	if err != nil {
		handler.response.RespondError(c, http.StatusNonAuthoritativeInfo, "UNAUTHORIZED")
		return
	}

	err = handler.user.SaveRefreshToken(user.ID, tokens.RefreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, 7*24*60*60, "/", "", false, true)

	handler.response.RespondLogin(c, http.StatusOK, messages.Login, tokens.AccessToken)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Exchange refresh cookie for a new access token
// @Tags Auth
// @Produce json
// @Success 200 {object} api.LoginResponse
// @Failure 401 {object} api.ErrorUnauthorized
// @Failure 500 {object} api.ErrorInternalServer
// @Router /auth/refresh [post]
func (handler *UserHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		handler.response.RespondError(c, http.StatusUnauthorized, "MISSING_REFRESH_TOKEN")
		return
	}

	user, err := handler.user.GetUserByRefreshToken(refreshToken)
	if err != nil {
		handler.response.RespondError(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN")
		return
	}

	err = handler.user.DeleteRefreshToken(refreshToken)
	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	tokens, err := apiUtils.GenerateToken(user.Email, user.ID)
	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	err = handler.user.SaveRefreshToken(user.ID, tokens.RefreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, 7*24*60*60, "/", "", false, true)
	handler.response.RespondLogin(c, http.StatusOK, messages.Refresh, tokens.AccessToken)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token and clear cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} api.LogoutResponse
// @Failure 401 {object} api.ErrorUnauthorized
// @Failure 500 {object} api.ErrorInternalServer
// @Router /auth/logout [post]
func (handler *UserHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		handler.response.RespondError(c, http.StatusUnauthorized, "MISSING_REFRESH_TOKEN")
		return
	}

	err = handler.user.DeleteRefreshToken(refreshToken)
	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	handler.response.RespondSuccess(c, http.StatusOK, messages.Logout)
}
