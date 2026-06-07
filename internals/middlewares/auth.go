package middlewares

import (
	"net/http"
	"strings"

	"github.com/erfan-goodarzi/booking-event-api/internals/api"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/gin-gonic/gin"
)

var response api.APIResponse

func Authenticate(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")

	if token == "" {
		response.RespondAuthError(c, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	id, err := apiUtils.VerifyToken(token)

	if err != nil {
		response.RespondAuthError(c, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	c.Set("userId", id)
	c.Next()
}
