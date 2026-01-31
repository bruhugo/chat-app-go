package utils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

func CreateUserDto() *dto.UserDto {
	return &dto.UserDto{Username: "username", ID: 12, Email: "email"}
}

func GenerateKey() string {
	b := make([]byte, 32)

	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
