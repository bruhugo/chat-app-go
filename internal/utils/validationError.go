package utils

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

func GetValidationError(err error) *dto.ErrorResponse {

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {

		errorsMap := make(map[string]string)

		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()

			switch tag {
			case "required":
				errorsMap[field] = "This field is required"
			case "email":
				errorsMap[field] = "Invalid email format"
			case "min":
				errorsMap[field] = "Minimum length is " + fe.Param()
			default:
				errorsMap[field] = "Invalid value"
			}
		}

		return &dto.ErrorResponse{
			IssuedAt:         time.Now(),
			ErrorMessage:     "Bad request provided",
			StatusCode:       400,
			ValidationErrors: errorsMap,
		}
	}

	return &dto.ErrorResponse{
		IssuedAt:         time.Now(),
		ErrorMessage:     "Bad request provided",
		StatusCode:       400,
		ValidationErrors: make(map[string]string),
	}
}
