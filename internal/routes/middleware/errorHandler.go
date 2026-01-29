package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
)

type errorResponse struct {
	ErrorMessage string    `json:"errorMessage"`
	StatusCode   int       `json:"statusCode"`
	IssuedAt     time.Time `json:"issuedAt"`
}

func NewErrorResponse(errorMessage string, status int) errorResponse {
	return errorResponse{
		ErrorMessage: errorMessage,
		StatusCode:   status,
		IssuedAt:     time.Now(),
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, err := range ctx.Errors {
			log.Printf("Error: %s", err.Err.Error())

			if httpError, ok := err.Err.(*exceptions.HttpError); err.Type != gin.ErrorTypePrivate && ok {
				message := httpError.Error()
				if httpError.Message != "" {
					message = httpError.Message
				}

				ctx.JSON(httpError.Status, NewErrorResponse(message, httpError.Status))
				return
			}
		}
	}
}
