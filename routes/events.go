package routes

import (
	"net/http"

	"example.com/booking-event/messages"
	"example.com/booking-event/models"
	"example.com/booking-event/utils"
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
	id, err := utils.ParseID(c)

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "ID_NOT_FOUND")
		return
	}

	event, err := models.GetEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	response.RespondRetrievedSuccess(c, http.StatusOK, event)
}

func createEvents(c *gin.Context) {

	var event models.Event
	err := c.ShouldBindJSON(&event)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	id := c.GetInt64("userId")
	event.UserID = id
	event.Create()

	response.RespondSuccess(c, http.StatusCreated, messages.CreateEventSuccess, event)
}

func updateEvent(c *gin.Context) {
	id, err := utils.ParseID(c)

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "ID_NOT_FOUND")
		return
	}

	_, err = models.GetEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	var updatedEvent models.Event
	err = c.ShouldBindJSON(&updatedEvent)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	updatedEvent.ID = id
	err = updatedEvent.Update()

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "EVENTS")
		return
	}

	response.RespondSuccess(c, http.StatusOK, messages.UpdateEventSuccess, updatedEvent)
}

func deleteEvent(c *gin.Context) {
	id, err := utils.ParseID(c)

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "ID_NOT_FOUND")
		return
	}

	event, err := models.GetEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	err = event.Delete()

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "EVENTS")
		return
	}

	response.RespondSuccess(c, http.StatusOK, messages.DeletesEventSuccess)
}
