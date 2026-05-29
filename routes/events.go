package routes

import (
	"net/http"

	"example.com/booking-event/messages"
	"example.com/booking-event/models"
	"github.com/gin-gonic/gin"
)

var response models.APIResponse

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "EVENTS")
		return
	}

	response.RespondRetrievedSuccess(c, http.StatusOK, events)
}

func createEvents(c *gin.Context) {
	var event models.Event
	err := c.ShouldBindJSON(&event)
	event.Create()

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "EVENTS")
		return
	}

	response.RespondSuccess(c, http.StatusOK, messages.CreateEventSuccess, event)
}
