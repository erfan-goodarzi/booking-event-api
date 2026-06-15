package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internal/messages"
	"github.com/erfan-goodarzi/booking-event-api/internal/models"
	"github.com/erfan-goodarzi/booking-event-api/internal/store"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/pkg/validation"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketStore store.TicketStore
	eventStore  store.EventStore
	logger      *log.Logger
	response    *APIResponse
}

func NewTicketHandler(ticketStore store.TicketStore, eventStore store.EventStore, logger *log.Logger, response *APIResponse) *TicketHandler {
	return &TicketHandler{
		ticketStore,
		eventStore,
		logger,
		response,
	}
}

// CreateTicket godoc
// @Summary Create a ticket
// @Description Create a new ticket (authenticated)
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param ticket body models.CreateTicketRequest true "Ticket payload"
// @Success 201 {object} models.TicketResponse
// @Failure 422 {object} models.ErrorBadRequest
// @Failure 404 {object} models.ErrorNotFound
// @Failure 500 {object} models.ErrorInternalServer
// @Router /events/{id}/tickets [post]
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var payload models.CreateTicketRequest
	id, err := apiUtils.ParseID(c)
	currentUserId := c.GetString("userId")

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	err = c.ShouldBindJSON(&payload)

	if err != nil {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
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

	ticket := &models.Ticket{
		EventId:  id,
		Quantity: payload.Quantity,
		Price:    payload.Price,
		Type:     payload.Type,
		UserId:   c.GetString("userId"),
	}

	err = validation.Validate.Struct(payload)

	if err != nil {
		h.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	ticket, err = h.ticketStore.CreateTicket(id, ticket)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "FAILED_TO_REGISTER")
		return
	}

	h.response.RespondSuccess(c, http.StatusCreated, messages.CreateTicketSuccess, ticket)
}
