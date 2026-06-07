package api

import (
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internals/store"
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type LoginResponse struct {
	Token   string `json:"token" example:"eyJhbGci..."`
	Message string `json:"message" example:"user logged in successfully"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"user logged out successfully"`
}

type ErrorBadRequest struct {
	Message string `json:"message" example:"Bad Request"`
	Error   string `json:"error" example:"PAYLOAD_NOT_VALID"`
}

type ErrorNotFound struct {
	Message string `json:"message" example:"Not Found"`
	Error   string `json:"error" example:"EVENT_NOT_FOUND"`
}

type ErrorInternalServer struct {
	Message string `json:"message" example:"Internal Server Error"`
	Error   string `json:"error" example:"UNKNOWN_ERROR"`
}

type ErrorConflict struct {
	Message string `json:"message" example:"Conflict"`
	Error   string `json:"error" example:"EMAIL_ALREADY_EXISTS"`
}

type ErrorValidation struct {
	Message string            `json:"message" example:"Validation Failed"`
	Error   string            `json:"error" example:"VALIDATION_FAILED"`
	Fields  map[string]string `json:"fields" example:"{\"email\": \"invalid email format\"}"`
}

type ErrorUnauthorized struct {
	Message string `json:"message" example:"Unauthorized"`
	Error   string `json:"error" example:"INVALID_CREDENTIALS"`
}

type ErrorForbidden struct {
	Message string `json:"message" example:"Forbidden"`
	Error   string `json:"error" example:"ACCESS_DENIED"`
}

type EventDeleteSuccess struct {
	Message string `json:"message" example:"Event deleted successfully"`
}

type HealthCheckResponse struct {
	Status string `json:"status" example:"ok"`
}

type HealthCheckErrorResponse struct {
	Status string `json:"status" example:"unavailable"`
	Error  string `json:"error" example:"DB_UNAVAILABLE"`
}

type EventResponse struct {
	Data    store.Event `json:"data"`
	Message string      `json:"message"`
}

type EventListResponse struct {
	Data    []store.Event `json:"data"`
	Message string        `json:"message"`
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

func (res APIResponse) ValidationError(c *gin.Context, status int, message string, fields map[string]string) {
	c.JSON(status, gin.H{
		"message": message,
		"fields":  fields,
	})
}
