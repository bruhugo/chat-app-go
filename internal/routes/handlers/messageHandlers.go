package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/services"
	"github.com/grongoglongo/chatter-go/internal/utils"
)

func GetMessagesByChatIdHandler(messageService *services.MessageService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageRequest, err := dto.GetPageRequest(*ctx.Request)
		if err != nil {
			ctx.Error(err)
			return
		}

		chatId, err := strconv.ParseInt(ctx.Param("chatId"), 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Id must be a string.", http.StatusBadRequest))
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		messagesPage, err := messageService.GetMessages(chatId, userId, pageRequest)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, messagesPage)
	}
}

func CreateMessageHandler(messageService *services.MessageService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createMessageDto dto.CreateMessageDto
		err := ctx.ShouldBindBodyWithJSON(&createMessageDto)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			log.Print("Failed to convert user id to int64")
			ctx.Error(exceptions.InternalServerError)
			return
		}

		messageDto, err := messageService.CreateMessage(createMessageDto, userId)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, messageDto)
	}
}
