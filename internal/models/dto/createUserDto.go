package dto

type CreateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email" binding:"email"`
}
