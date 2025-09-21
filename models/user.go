package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Hidden in JSON
	Role      string    `json:"role" db:"role"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// UserCreateRequest represents a request to create a new user
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
}

// UserUpdateRequest represents a request to update a user
type UserUpdateRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UserChangePasswordRequest represents a request to change password
type UserChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// JWTPayload represents JWT token payload
type JWTPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

// UserListResponse represents a paginated user list response
type UserListResponse struct {
	Users      []User `json:"users"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	AdminUsers    int64 `json:"admin_users"`
	RegularUsers  int64 `json:"regular_users"`
	NewUsersToday int64 `json:"new_users_today"`
}

// HEPRecord represents a HEP record for analytics
type HEPRecord struct {
	ID            uint64    `json:"id"`
	CallID        string    `json:"call_id"`
	SourceIP      string    `json:"source_ip"`
	DestinationIP string    `json:"destination_ip"`
	Protocol      string    `json:"protocol"`
	Method        string    `json:"method"`
	StatusCode    uint16    `json:"status_code"`
	Timestamp     time.Time `json:"timestamp"`
	RawData       string    `json:"raw_data"`
	CreatedAt     time.Time `json:"created_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
