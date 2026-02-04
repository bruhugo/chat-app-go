package dto

import (
	"net/http"
	"strconv"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
)

type Page[T any] struct {
	Content  []T
	Page     int
	PageSize int
	Number   int
}

type PageRequest struct {
	Page     int
	PageSize int
}

func GetPageRequest(r http.Request) (*PageRequest, error) {
	pageRequest := &PageRequest{}
	queries := r.URL.Query()

	i, ok := queries["page"]
	if ok {
		page, err := strconv.Atoi(i[0])
		if err != nil {
			return nil, exceptions.BadRequestError
		}
		pageRequest.Page = page
	}

	i, ok = queries["pageSize"]
	if ok {
		pageSize, err := strconv.Atoi(i[0])
		if err != nil {
			return nil, exceptions.BadRequestError
		}
		pageRequest.PageSize = pageSize
	}

	return pageRequest, nil
}
