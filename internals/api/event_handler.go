package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/models"
	"github.com/erfan-goodarzi/booking-event-api/internals/store"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/pkg/validation"
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

// GetEvents godoc
// @Summary List events
// @Description Get all events
// @Tags Events
// @Produce json
// @Success 200 {object} api.EventListResponse
// @Failure 500 {object} api.ErrorInternalServer
// @Router /events [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	events, err := h.eventStore.GetAllEvents()

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	h.response.RespondRetrievedSuccess(c, http.StatusOK, events)
}

// GetEvent godoc
// @Summary Get event by ID
// @Description Get an event by its ID
// @Tags Events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} api.EventResponse
// @Failure 400 {object} api.ErrorBadRequest
// @Failure 404 {object} api.ErrorNotFound
// @Failure 500 {object} api.ErrorInternalServer
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	event, err := h.eventStore.GetEvent(id)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	h.response.RespondRetrievedSuccess(c, http.StatusOK, event)
}

// CreateEvent godoc
// @Summary Create an event
// @Description Create a new event (authenticated)
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body models.CreateEventRequest true "Event payload"
// @Success 201 {object} models.Event
// @Failure 401 {object} api.ErrorUnauthorized
// @Failure 422 {object} api.ErrorValidation
// @Failure 500 {object} api.ErrorInternalServer
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var payload models.CreateEventRequest
	err := c.ShouldBindJSON(&payload)

	if err != nil {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = validation.Validate.Struct(payload)

	if err != nil {
		h.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	event := models.Event{
		Title:       payload.Title,
		Description: payload.Description,
		Location:    payload.Location,
		DateTime:    payload.DateTime,
		UserId:      c.GetString("userId"),
		Duration:    *payload.Duration,
	}

	createdEvent, err := h.eventStore.CreateEvent(&event)
	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.response.RespondSuccess(c, http.StatusCreated, messages.CreateEventSuccess, createdEvent)
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an existing event (authenticated, owner only)
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param event body models.PatchEventRequest true "Patch payload"
// @Success 200 {object} models.Event
// @Failure 401 {object} api.ErrorUnauthorized
// @Failure 403 {object} api.ErrorForbidden
// @Failure 404 {object} api.ErrorNotFound
// @Failure 422 {object} api.ErrorValidation
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	existingEvent, err := h.eventStore.GetEvent(id)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	if existingEvent == nil {
		h.response.RespondError(c, http.StatusNotFound, "EVENT_NOT_FOUND")
		return
	}

	var partialEvent models.PatchEventRequest

	err = c.ShouldBindJSON(&partialEvent)

	if err != nil {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = validation.Validate.Struct(partialEvent)

	if err != nil {
		h.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	eventOwner, err := h.eventStore.GetEventOwner(id)

	if errors.Is(err, sql.ErrNoRows) {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "EVENT_NOT_EXIST")
		return
	}

	if eventOwner != currentUserId {
		h.response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	store.ApplyEventPatch(existingEvent, partialEvent)

	updatedEvent, err := h.eventStore.UpdateEvent(existingEvent)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	h.response.RespondSuccess(c, http.StatusOK, messages.UpdateEventSuccess, updatedEvent)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an event by ID (authenticated, owner only)
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 200 {object} api.EventDeleteSuccess
// @Failure 401 {object} api.ErrorUnauthorized
// @Failure 403 {object} api.ErrorForbidden
// @Failure 422 {object} api.ErrorValidation
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	eventOwner, err := h.eventStore.GetEventOwner(id)

	if errors.Is(err, sql.ErrNoRows) {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "EVENT_NOT_EXIST")
		return
	}

	if eventOwner != currentUserId {
		h.response.RespondError(c, http.StatusForbidden, "ACCESS_DENIED")
		return
	}

	err = h.eventStore.DeleteEvent(id)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "UNKNOWN_ERROR")
		return
	}

	h.response.RespondSuccess(c, http.StatusOK, messages.DeletesEventSuccess)
}
