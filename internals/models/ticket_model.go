package models

import "time"

type TicketType string

const (
	TicketTypeVIP     TicketType = "vip"
	TicketTypeGeneral TicketType = "general"
)

type Ticket struct {
	ID        string     `json:"id" example:"e2f1c3a8-7d4b-11ec-90d6-0242ac120003"`
	UserId    string     `json:"userId" db:"user_id" example:"u12345"`
	EventId   string     `json:"eventId" db:"event_id" example:"e12345"`
	Type      TicketType `json:"type" example:"VIP"`
	Price     float64    `json:"price" example:"99.99"`
	Quantity  int        `json:"quantity" example:"2"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	Bookings  []Booking  `json:"bookings,omitempty"`
}

type CreateTicketRequest struct {
	Type     TicketType `json:"type" validate:"required oneof=vip general" example:"vip"`
	Price    float64    `json:"price" validate:"required,gt=0" example:"99.99"`
	Quantity int        `json:"quantity" validate:"required,gt=0" example:"2"`
}

type PatchTicketRequest struct {
	Type     *TicketType `json:"type" validate:"omitempty" example:"vip"`
	Price    *float64    `json:"price" validate:"omitempty,gt=0" example:"99.99"`
	Quantity *int        `json:"quantity" validate:"omitempty,gt=0" example:"2"`
}

type TicketResponse struct {
	Data    Ticket `json:"data"`
	Message string `json:"message"`
}
