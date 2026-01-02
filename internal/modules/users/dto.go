package users

// CreateUserRequest DTO for creating a new user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullName" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=ADMIN STAFF"` // Only ADMIN or STAFF
}

// UpdateUserRequest DTO for updating a user
type UpdateUserRequest struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"omitempty,oneof=ADMIN STAFF"`
	IsActive *bool  `json:"isActive"`
}

// UserResponse DTO for user response
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

// UsersListResponse DTO for paginated users list
type UsersListResponse struct {
	Data       []UserResponse `json:"data"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"totalPages"`
}
