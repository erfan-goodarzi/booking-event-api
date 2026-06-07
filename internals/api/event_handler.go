package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/store"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventStore store.EventStore
	logger     *log.Logger
	response   *APIResponse
}

func NewEventHandler(eventStore store.EventStore, logger *log.Logger, response *APIResponse) *EventHandler {
	return &EventHandler{
		eventStore,
		logger,
		response,
	}
}

func (handler *EventHandler) GetEvents(c *gin.Context) {
	events, err := handler.eventStore.GetAllEvents()

	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "EVENT_NOT_FOUND")
		return
	}

	handler.response.RespondRetrievedSuccess(c, http.StatusOK, events)
}

func (handler *EventHandler) GetEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)

	if err != nil {
		handler.response.RespondError(c, http.StatusBadRequest, "ID_NOT_FOUND")
		return
	}

	event, err := handler.eventStore.GetEvent(*id)

	if err != nil {
		handler.response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	handler.response.RespondRetrievedSuccess(c, http.StatusOK, event)
}

func (handler *EventHandler) CreateEvents(c *gin.Context) {
	var event store.Event
	err := c.ShouldBindJSON(&event)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	id := c.GetString("userId")
	event.UserID = id

	createdEvent, err := handler.eventStore.CreateEvent(&event)
	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.response.RespondSuccess(c, http.StatusCreated, messages.CreateEventSuccess, createdEvent)
}

func (handler *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		handler.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	existingEvent, err := handler.eventStore.GetEvent(*id)

	if err != nil {
		handler.response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	if existingEvent == nil {
		handler.response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	var partialEvent store.PatchEventRequest

	err = c.ShouldBindJSON(&partialEvent)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	eventOwner, err := handler.eventStore.GetEventOwner(*id)

	if errors.Is(err, sql.ErrNoRows) {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "EVENT_NOT_EXIST")
		return
	}

	if *eventOwner != currentUserId {
		handler.response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	handler.eventStore.ApplyPatch(existingEvent, partialEvent)

	updatedEvent, err := handler.eventStore.UpdateEvent(existingEvent)

	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	handler.response.RespondSuccess(c, http.StatusOK, messages.UpdateEventSuccess, updatedEvent)
}

func (handler *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		handler.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	eventOwner, err := handler.eventStore.GetEventOwner(*id)

	if errors.Is(err, sql.ErrNoRows) {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "EVENT_NOT_EXIST")
		return
	}

	if *eventOwner != currentUserId {
		handler.response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	err = handler.eventStore.DeleteEvent(*id)

	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	handler.response.RespondSuccess(c, http.StatusOK, messages.DeletesEventSuccess)
}
