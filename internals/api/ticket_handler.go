package api

import (
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internals/messages"
	"github.com/erfan-goodarzi/booking-event-api/internals/models"
	"github.com/erfan-goodarzi/booking-event-api/internals/store"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketStore store.TicketStore
	logger      *log.Logger
	response    *APIResponse
}

func NewTicketHandler(ticketStore store.TicketStore, logger *log.Logger, response *APIResponse) *TicketHandler {
	return &TicketHandler{
		ticketStore,
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
// @Success 201 {object} api.TicketResponse
// @Failure 400 {object} api.ErrorBadRequest
// @Failure 422 {object} api.ErrorValidation
// @Failure 500 {object} api.ErrorInternalServer
// @Router /events/{id}/tickets [post]
func (handler *TicketHandler) CreateTicket(c *gin.Context) {
	var payload models.CreateTicketRequest
	id, err := apiUtils.ParseID(c)

	if err != nil {
		handler.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	err = c.ShouldBindJSON(&payload)

	if err != nil {
		handler.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	ticket := &models.Ticket{
		EventId:  id,
		Quantity: payload.Quantity,
		Price:    payload.Price,
		Type:     payload.Type,
		UserId:   c.GetString("userId"),
	}

	ticket, err = handler.ticketStore.CreateTicket(id, ticket)

	if err != nil {
		handler.response.RespondError(c, http.StatusInternalServerError, "FAILED_TO_CREATE_TICKET")
		return
	}

	handler.response.RespondSuccess(c, http.StatusCreated, messages.CreateTicketSuccess, ticket)
}
