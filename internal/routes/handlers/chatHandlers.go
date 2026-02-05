package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/services"
	"github.com/grongoglongo/chatter-go/internal/utils"
)

func CreateChatHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createChatDto dto.CreateChatDto
		err := ctx.ShouldBindBodyWithJSON(&createChatDto)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		createChatDto.CreatorId = userId

		chat, err := chatService.CreateChat(createChatDto)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, chat)
	}
}
