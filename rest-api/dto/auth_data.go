package dto

// type AuthRequest struct {
// 	Email string `json:"email"`
// 	Password string `json:"password"`
// }

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Password    string `json:"password" validate:"required, min:6"`
}

type RegisterResponse struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Token       string `json:"token"`
}
