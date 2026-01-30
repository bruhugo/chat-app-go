package services

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grongoglongo/chatter-go/internal/config"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

func CreateJwt(user *dto.UserDto) (string, error) {
	if config.EnvConfig.JwtSecret == "" {
		return "", errors.New("Must provide a key as env variable.")
	}

	claims := jwt.StandardClaims{
		Id:        strconv.FormatInt(user.ID, 10),
		Subject:   user.Email,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "Go-Chatter-Backend",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60 days
	}

	realKey, err := base64.StdEncoding.DecodeString(config.EnvConfig.JwtSecret)
	if err != nil {
		return "", err
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rawJwt, err := jwt.SignedString(realKey)
	if err != nil {
		return "", err
	}

	log.Printf("JWT created for user %d", user.ID)

	return rawJwt, nil
}

func DecryptJwt(token string) (*jwt.StandardClaims, error) {
	var claims jwt.StandardClaims

	jwt, err := jwt.ParseWithClaims(token, &claims, KeyFunc)
	if err != nil {
		return nil, err
	}

	if !jwt.Valid {
		return nil, errors.New("Invalid token provided.")
	}

	log.Printf("JWT decoded from user %v", claims.Id)

	return &claims, nil
}

func KeyFunc(t *jwt.Token) (interface{}, error) {

	b, err := base64.StdEncoding.DecodeString(config.EnvConfig.JwtSecret)
	if err != nil {
		return nil, err
	}

	return b, nil
}
