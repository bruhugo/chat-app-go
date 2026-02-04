package dto

type CreateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email" binding:"email"`
}

type LoginUserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
