package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/services"
)

func GetMessagesHandler(messageService *services.MessageService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageRequest, err := dto.GetPageRequest(*ctx.Request)
		if err != nil {
			ctx.Error(err)
			return
		}

		senderId, err := strconv.ParseInt(ctx.Param("senderId"), 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Id must be a string.", http.StatusBadRequest))
			return
		}

		receiverId, err := strconv.ParseInt(ctx.Param("receiverId"), 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Id must be a string.", http.StatusBadRequest))
			return
		}

		realUserId := ctx.Value("userId")
		if userId, ok := realUserId.(int64); ok && userId != senderId {
			ctx.Error(exceptions.UnauthorizedError)
			return
		}

		dto := dto.GetMessagesDto{
			SenderId:    senderId,
			ReceiverId:  receiverId,
			PageRequest: *pageRequest,
		}

		messagesPage, err := messageService.GetMessages(dto)

		ctx.JSON(http.StatusOK, messagesPage)
	}
}
