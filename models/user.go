package models

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type RegisterPayload struct {
	Username string `json:"username" db:"username" validate:"required,min=4,max=50"`
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=6,max=50"`
}

type LoginPayload struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=6,max=50"`
}
