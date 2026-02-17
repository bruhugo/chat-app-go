package dto

import "time"

type ErrorResponse struct {
	ErrorMessage     string            `json:"errorMessage"`
	StatusCode       int               `json:"statusCode"`
	IssuedAt         time.Time         `json:"issuedAt"`
	ValidationErrors map[string]string `json:"validationErrors"`
}
