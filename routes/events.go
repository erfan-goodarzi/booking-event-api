package routes

import (
	"net/http"
	"strconv"

	"example.com/booking-event/messages"
	"example.com/booking-event/models"
	"github.com/gin-gonic/gin"
)

var response models.APIResponse

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondRetrievedSuccess(c, http.StatusOK, events)
}

func getEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "ID_NOT_FOUND")
		return
	}

	event, err := models.GetEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "EVENT_NOT_FOUND")
		return
	}

	response.RespondRetrievedSuccess(c, http.StatusOK, event)
}

func createEvents(c *gin.Context) {
	var event models.Event
	err := c.ShouldBindJSON(&event)

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "EVENTS")
		return
	}

	event.Create()

	response.RespondSuccess(c, http.StatusOK, messages.CreateEventSuccess, event)
}
