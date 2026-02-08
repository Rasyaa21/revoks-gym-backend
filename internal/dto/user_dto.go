package dto

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Role      string `json:"role" validate:"required,oneof=SUPERADMIN ADMIN CASHIER MEMBER PT"`
	FullName  string `json:"full_name" validate:"required,min=2,max=255"`
	Email     string `json:"email" validate:"omitempty,email"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	Password  string `json:"password" validate:"required,min=6"`
	PhotoURL  string `json:"photo_url" validate:"omitempty,url"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=255"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
	Status   string `json:"status" validate:"omitempty,oneof=active blocked"`
}

// UserResponse represents the response body for a user
type UserResponse struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PhotoURL  string `json:"photo_url"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
