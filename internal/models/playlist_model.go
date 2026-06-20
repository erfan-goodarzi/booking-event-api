package models

import "time"

type Playlist struct {
	ID        string    `json:"id" example:"e2f1c3a8-7d4b-11ec-90d6-0242ac120003"`
	UserId    string    `json:"userId" db:"user_id" example:"u12345"`
	Name      string    `json:"name" example:"listen later"`
	Color     string    `json:"color" example:"red"`
	CreatedAt time.Time `json:"created_at" db:"created_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" swaggertype:"string" format:"date-time" example:"2026-06-07T15:04:05Z"`
	Events    []Event   `json:"events"`
}

type CreatePlaylistRequest struct {
	Name  string `json:"name" validate:"required" example:"listen later"`
	Color string `json:"color" validate:"required" example:"red"`
}

type PatchPlaylistRequest struct {
	Name  *string `json:"name" validate:"omitempty" example:"listen later"`
	Color *string `json:"color" validate:"omitempty" example:"red"`
}

type PlaylistDeleteSuccess struct {
	Message string `json:"message" example:"Playlist deleted successfully"`
}

type PlaylistResponse struct {
	Data    Playlist  `json:"data"`
	Message string `json:"message"`
}

type PlaylistListResponse struct {
	Data    []Playlist `json:"data"`
	Message string  `json:"message"`
}
