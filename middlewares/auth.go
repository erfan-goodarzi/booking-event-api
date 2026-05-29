package middlewares

import (
	"net/http"
	"strings"

	"example.com/booking-event/models"
	"example.com/booking-event/utils"
	"github.com/gin-gonic/gin"
)

var response models.APIResponse

func Authenticate(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")

	if token == "" {
		response.RespondAuthError(c, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	id, err := utils.VerifyToken(token)

	if err != nil {
		response.RespondAuthError(c, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	c.Set("userId", id)
	c.Next()
}
