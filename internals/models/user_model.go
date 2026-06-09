package models

type User struct {
	ID       string `json:"id" example:"u12345"`
	Username string `binding:"omitempty" json:"username" validate:"omitempty,min=3,max=20" example:"johndoe"`
	Email    string `binding:"required" json:"email" validate:"required,email" example:"user@example.com"`
	Password string `binding:"required" json:"password" validate:"required,min=8" example:"P@ssw0rd!"`
}
