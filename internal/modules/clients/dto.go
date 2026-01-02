package clients

// UpdateClientProfileRequest DTO for updating client profile
type UpdateClientProfileRequest struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

// ClientProfileResponse DTO for client profile response
type ClientProfileResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"isActive"`
}

// ClientsListResponse DTO for paginated clients list (admin)
type ClientsListResponse struct {
	Data       []ClientProfileResponse `json:"data"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	Total      int64                   `json:"total"`
	TotalPages int                     `json:"totalPages"`
}
