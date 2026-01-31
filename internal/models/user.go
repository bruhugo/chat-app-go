package models

import (
	"time"

	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToDto() *dto.UserDto {
	return &dto.UserDto{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
