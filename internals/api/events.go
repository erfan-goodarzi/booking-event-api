package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/models"
	"github.com/erfan-goodarzi/booking-event-api/utils"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventStore models.EventStore
	logger     *log.Logger
}

var response models.APIResponse

func NewEventHandler(eventStore models.EventStore, logger *log.Logger) *EventHandler {
	return &EventHandler{
		eventStore,
		logger,
	}
}

func (handler *EventHandler) GetEvents(c *gin.Context) {
	events, err := handler.eventStore.GetAllEvents()

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondRetrievedSuccess(c, http.StatusOK, events)
}

func (handler *EventHandler) GetEvent(c *gin.Context) {
	id, err := utils.ParseID(c)

	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "ID_NOT_FOUND")
		return
	}

	event, err := handler.eventStore.GetEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusNotFound, err.Error())
		return
	}

	response.RespondRetrievedSuccess(c, http.StatusOK, event)
}

func (handler *EventHandler) CreateEvents(c *gin.Context) {
	var event models.Event
	err := c.ShouldBindJSON(&event)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	id := c.GetInt64("userId")
	event.UserID = id

	createdEvent, err := handler.eventStore.CreateEvent(&event)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	response.RespondSuccess(c, http.StatusCreated, messages.CreateEventSuccess, createdEvent)
}

func (handler *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := utils.ParseID(c)
	currentUserId := c.GetInt64("userId")

	if err != nil {
		response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	existingEvent, err := handler.eventStore.GetEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	if existingEvent == nil {
		response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	var partialEvent models.PatchEventRequest

	err = c.ShouldBindJSON(&partialEvent)

	if err != nil {
		response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	eventOwner, err := handler.eventStore.GetEventOwner(id)

	if errors.Is(err, sql.ErrNoRows) {
		response.RespondError(c, http.StatusUnprocessableEntity, "EVENT_NOT_EXIST")
		return
	}

	if eventOwner != int(currentUserId) {
		response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	handler.eventStore.ApplyPatch(existingEvent, partialEvent)

	updatedEvent, err := handler.eventStore.UpdateEvent(existingEvent)

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	response.RespondSuccess(c, http.StatusOK, messages.UpdateEventSuccess, updatedEvent)
}

func (handler *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := utils.ParseID(c)
	currentUserId := c.GetInt64("userId")

	if err != nil {
		response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	eventOwner, err := handler.eventStore.GetEventOwner(id)

	if errors.Is(err, sql.ErrNoRows) {
		response.RespondError(c, http.StatusUnprocessableEntity, "EVENT_NOT_EXIST")
		return
	}

	if eventOwner != int(currentUserId) {
		response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	err = handler.eventStore.DeleteEvent(id)

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	response.RespondSuccess(c, http.StatusOK, messages.DeletesEventSuccess)
}
