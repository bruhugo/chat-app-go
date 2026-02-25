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

// MessagePageResponse is the paged response for message lists.
type MessagePageResponse struct {
	Content  []dto.MessageDto
	Page     int
	PageSize int
	Number   int
}

// @Summary List messages by chat ID
// @Description Returns paged messages for a chat.
// @Tags messages
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} MessagePageResponse
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId}/messages [get]
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

// @Summary Create message
// @Description Creates a new message.
// @Tags messages
// @Accept json
// @Produce json
// @Param body body dto.CreateMessageDto true "Message payload"
// @Success 201 {object} dto.MessageDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /messages/ [post]
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

// @Summary Delete message
// @Description Deletes a message by ID.
// @Tags messages
// @Param messageId path int true "Message ID"
// @Success 204
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /messages/{messageId} [delete]
func DeleteMessageHandler(messageService *services.MessageService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		messageId, err := strconv.ParseInt(ctx.Param("messageId"), 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Id must be a string.", http.StatusBadRequest))
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			log.Print("Failed to convert user id to int64")
			ctx.Error(exceptions.InternalServerError)
			return
		}

		err = messageService.DeleteMessage(messageId, userId)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
