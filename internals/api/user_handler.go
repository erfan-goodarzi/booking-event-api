package api

import (
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/store"
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

func (handler *UserHandler) Signup(c *gin.Context) {
	var user store.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
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

func (handler *UserHandler) Login(c *gin.Context) {
	var user store.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
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

	token, err := apiUtils.GenerateToken(user.Email, user.ID)

	if err != nil {
		handler.response.RespondError(c, http.StatusNonAuthoritativeInfo, "UNAUTHORIZED")
		return
	}

	handler.response.RespondLogin(c, http.StatusOK, messages.Login, token)
}
