package models

import "time"

type Event struct {
	ID          string    `json:"id" example:"e2f1c3a8-7d4b-11ec-90d6-0242ac120003"`
	Title       string    `json:"title" example:"Board Meeting"`
	Description string    `json:"description" example:"Quarterly planning meeting"`
	Location    string    `json:"location" example:"Conference Room A"`
	DateTime    time.Time `json:"dateTime" db:"date_time" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	UserId      string    `json:"userId" db:"user_id" example:"u12345"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	Tickets     []Ticket  `json:"tickets,omitempty"`
	TicketCount int       `json:"ticketCount,omitempty"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=50" example:"Board Meeting"`
	Description string    `json:"description" example:"Quarterly planning meeting"`
	Location    string    `json:"location" validate:"required" example:"Conference Room A"`
	DateTime    time.Time `json:"dateTime" validate:"required" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
}

type PatchEventRequest struct {
	Title       *string    `json:"title" validate:"omitempty,min=3,max=50" example:"Board Meeting"`
	Description *string    `json:"description" validate:"omitempty" example:"Updated description"`
	Location    *string    `json:"location" validate:"omitempty" example:"Room B"`
	DateTime    *time.Time `json:"dateTime" validate:"omitempty" swaggertype:"string" format:"date-time" example:"2026-07-01T09:00:00Z"`
}
