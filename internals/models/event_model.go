package models

import "time"

type Event struct {
	ID          string    `json:"id" example:"e2f1c3a8-7d4b-11ec-90d6-0242ac120003"`
	Title       string    `json:"title" example:"Board Meeting"`
	Description string    `json:"description" example:"Quarterly planning meeting"`
	Location    string    `json:"location" example:"Conference Room A"`
	DateTime    time.Time `json:"dateTime" db:"date_time" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	UserId      string    `json:"userId" db:"user_id" example:"u12345"`
	Tickets     []Ticket  `json:"tickets"`
	Host        Host      `json:"host"`
	TicketCount int       `json:"ticketCount,omitempty"`
	Duration    int       `json:"duration" db:"duration" example:"60"` // Duration in minutes
	CreatedAt   time.Time `json:"created_at" db:"created_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	Version     int       `json:"version" db:"version" example:"1"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=50" example:"Board Meeting"`
	Description string    `json:"description" example:"Quarterly planning meeting"`
	Location    string    `json:"location" validate:"required" example:"Conference Room A"`
	DateTime    time.Time `json:"dateTime" validate:"required" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	Duration    *int      `json:"duration,omitempty" validate:"omitempty" example:"60"`
}

type PatchEventRequest struct {
	Title       *string    `json:"title" validate:"omitempty,min=3,max=50" example:"Board Meeting"`
	Description *string    `json:"description" validate:"omitempty" example:"Updated description"`
	Location    *string    `json:"location" validate:"omitempty" example:"Room B"`
	DateTime    *time.Time `json:"dateTime" validate:"omitempty" swaggertype:"string" format:"date-time" example:"2026-07-01T09:00:00Z"`
	Duration    *int       `json:"duration,omitempty" validate:"omitempty" example:"60"`
}

type Host struct {
	ID       string `json:"id" example:"u12345"`
	Username string `json:"username" example:"john"`
	Email    string `json:"email" example:"john@example.com"`
}

type EventFilter struct {
	Search   string
	Location string
	From     time.Time
	To       time.Time
}
