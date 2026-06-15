package api

import (
	"log"
	"net/http"

	"github.com/erfan-goodarzi/booking-event-api/internal/messages"
	"github.com/erfan-goodarzi/booking-event-api/internal/models"
	"github.com/erfan-goodarzi/booking-event-api/internal/store"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/erfan-goodarzi/booking-event-api/pkg/validation"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingStore store.BookingStore
	eventStore   store.EventStore
	logger       *log.Logger
	response     *APIResponse
}

func NewBookingHandler(bookingStore store.BookingStore, eventStore store.EventStore, logger *log.Logger, response *APIResponse) *BookingHandler {
	return &BookingHandler{
		bookingStore,
		eventStore,
		logger,
		response,
	}
}

// RegisterEvent godoc
// @Summary Book an event
// @Description Book a new ticket
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 201 {object} models.BookingResponse
// @Failure 404 {object} models.ErrorNotFound
// @Failure 500 {object} models.ErrorInternalServer
// @Router /events/tickets/{id}/register [post]
func (h *BookingHandler) RegisterEvent(c *gin.Context) {
	id, err := apiUtils.ParseID(c)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	booking := &models.Booking{
		UserId:   c.GetString("userId"),
		TicketId: id,
		Status:   models.BookingStatusPending,
	}

	booking, err = h.bookingStore.CreateBooking(id, booking)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "FAILED_TO_REGISTER")
		return
	}

	h.response.RespondSuccess(c, http.StatusCreated, messages.CreateTicketSuccess, booking)
}

// UpdateRegistrationStatus godoc
// @Summary Update registration status
// @Description update the status of registered ticket
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "registration ID"
// @Param status body models.PatchBookingRequest true "Patch payload"
// @Success 201 {object} models.BookingResponse
// @Failure 404 {object} models.ErrorNotFound
// @Failure 500 {object} models.ErrorInternalServer
// @Router /events/tickets/register/{id}/status [put]
func (h *BookingHandler) UpdateRegistrationStatus(c *gin.Context) {
	id, err := apiUtils.ParseID(c)

	if err != nil {
		h.response.RespondError(c, http.StatusNotFound, "ID_NOT_FOUND")
		return
	}

	var partialBooking models.PatchBookingRequest

	err = c.ShouldBindJSON(&partialBooking)

	if err != nil {
		h.response.RespondError(c, http.StatusUnprocessableEntity, "PAYLOAD_NOT_VALID")
		return
	}

	err = validation.Validate.Struct(partialBooking)

	if err != nil {
		h.response.ValidationError(c, http.StatusUnprocessableEntity, "VALIDATION_FAILED", validation.FormatValidationErrors(err))
		return
	}

	booking := &models.Booking{}

	store.ApplyRegistrationPatch(booking, partialBooking)

	booking, err = h.bookingStore.UpdateBookingStatus(id, booking)

	if err != nil {
		h.response.RespondError(c, http.StatusInternalServerError, "FAILED_TO_UPDATE_STATUS")
		return
	}

	h.response.RespondSuccess(c, http.StatusCreated, messages.RegisterStatusSuccess, booking)
}
