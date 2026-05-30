package routes

import (
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/messages"
	"github.com/erfan-goodarzi/booking-event-api/models"
	"github.com/erfan-goodarzi/booking-event-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "SIGNUP")
		return
	}

	user.Create()

	response.RespondSuccess(c, http.StatusCreated, messages.Signup, user)
}

func login(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "LOGIN")
		return
	}

	err = user.ValidateCredential()

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
