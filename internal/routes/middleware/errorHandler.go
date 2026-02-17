package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

func NewErrorResponse(errorMessage string, status int, validationErrors map[string]string) dto.ErrorResponse {
	return dto.ErrorResponse{
		ErrorMessage:     errorMessage,
		StatusCode:       status,
		IssuedAt:         time.Now(),
		ValidationErrors: validationErrors,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, err := range ctx.Errors {
			log.Printf("Error: %s", err.Err.Error())

			if httpError, ok := err.Err.(*exceptions.HttpError); ok {
				message := httpError.Error()
				if httpError.Message != "" {
					message = httpError.Message
				}

				ctx.JSON(httpError.Status, NewErrorResponse(message, httpError.Status, make(map[string]string)))
				return
			}
		}
	}
}
