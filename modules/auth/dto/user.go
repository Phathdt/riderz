package dto

import "github.com/phathdt/service-context/core"

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type SignupRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type UserCreate struct {
	core.SQLModel
	Email    string `json:"email"`
	Password string `json:"password"`
}
