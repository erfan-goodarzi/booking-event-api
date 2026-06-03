package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (res APIResponse) RespondError(c *gin.Context, status int, err string) {
	c.JSON(status, gin.H{
		"message": http.StatusText(status),
		"error":   err,
	})
}

func (res APIResponse) RespondAuthError(c *gin.Context, status int, err string) {
	c.AbortWithStatusJSON(status, gin.H{
		"message": http.StatusText(status),
		"error":   err,
	})
}

func (res APIResponse) RespondSuccess(c *gin.Context, status int, message string, data ...any) {
	response := gin.H{
		"message": message,
	}

	if len(data) > 0 {
		response["data"] = data[0]
	}

	c.JSON(status, response)
}

func (res APIResponse) RespondRetrievedSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"data":    data,
		"message": "data retrieved successfully",
	})
}

func (res APIResponse) RespondLogin(c *gin.Context, status int, message string, token string) {
	c.JSON(status, gin.H{
		"token":   token,
		"message": message,
	})
}
