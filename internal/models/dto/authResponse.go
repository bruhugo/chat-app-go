package dto

type AuthResponse struct {
	User UserDto
	Jwt  string
}
