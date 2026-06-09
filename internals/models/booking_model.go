package models

import "time"

type bookingStatus string

const (
	BookingStatusPending   bookingStatus = "pending"
	BookingStatusConfirmed bookingStatus = "confirmed"
	BookingStatusCancelled bookingStatus = "cancelled"
)

type Booking struct {
	ID        string        `json:"id" example:"e2f1c3a8-7d4b-11ec-90d6-0242ac120003"`
	UserId    string        `json:"userId" db:"user_id" example:"u12345"`
	TicketId  string        `json:"ticketId" db:"ticket_id" example:"t12345"`
	Status    bookingStatus `json:"status" example:"pending"`
	CreatedAt time.Time     `json:"created_at" db:"created_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
}

type CreateBookingRequest struct {
	UserId   string        `json:"userId" validate:"required" example:"u12345"`
	TicketId string        `json:"ticketId" validate:"required" example:"t12345"`
	Status   bookingStatus `json:"status" validate:"required oneof=pending confirmed cancelled" example:"pending"`
}

type PatchBookingRequest struct {
	Status *bookingStatus `json:"status" validate:"omitempty oneof=pending confirmed cancelled" example:"pending"`
}
