package api

import (
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/models"
	"github.com/erfan-goodarzi/booking-event-api/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user   models.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore models.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore,
		logger,
	}
}

func (handler *UserHandler) Signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = handler.user.Create(&user)

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusCreated, messages.Signup, user)
}

func (handler *UserHandler) Login(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = handler.user.ValidateCredential(&user)

	if err != nil {
		response.RespondError(c, http.StatusNonAuthoritativeInfo, "UNAUTHORIZED")
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		response.RespondError(c, http.StatusNonAuthoritativeInfo, "UNAUTHORIZED")
		return
	}

	response.RespondLogin(c, http.StatusOK, messages.Login, token)
}
