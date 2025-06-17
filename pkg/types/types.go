package types

import "auth/db/models"

// Common types

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type OKResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PublicUser struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Method/Handler types

type SignUpStruct struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateThought struct {
	Thought    string            `json:"thought"`
	Visibility models.Visibility `json:"visibility"`
}

type UpdateThought struct {
	Visibility models.Visibility `json:"visibility"`
}
