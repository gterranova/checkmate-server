package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FormResponse[T any] struct {
	Changes T `json:"changes"`
	Model   T `json:"model"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	Email string `json:"email"`
}

type User struct {
	ID        int    `json:"uid"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	SessionID string `json:"-"`
}
