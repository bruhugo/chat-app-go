package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/messenger"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
	"github.com/grongoglongo/chatter-go/internal/services"
	"github.com/grongoglongo/chatter-go/internal/utils"
)

// @Summary Create chat
// @Description Creates a new chat and returns it.
// @Tags chats
// @Accept json
// @Produce json
// @Param body body dto.CreateChatDto true "Chat payload"
// @Success 201 {object} dto.ChatDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Router /chats/ [post]
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

// @Summary Delete chat
// @Description Deletes a chat by ID.
// @Tags chats
// @Param chatId path int true "Chat ID"
// @Success 204
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId} [delete]
func DeleteChatHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatId, err := strconv.ParseInt(ctx.Param("chatId"), 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Id must be a string.", http.StatusBadRequest))
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		if err := chatService.DeleteChat(userId, chatId); err != nil {
			ctx.Error(err)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

// @Summary Update chat
// @Description Updates a chat by ID.
// @Tags chats
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param body body dto.UpdateChatDto true "Chat update payload"
// @Success 200 {object} dto.ChatDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId} [put]
func UpdateChatHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatId, err := strconv.ParseInt(ctx.Param("chatId"), 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Id must be a string.", http.StatusBadRequest))
			return
		}

		var updateChatDto dto.UpdateChatDto
		if err := ctx.ShouldBindBodyWithJSON(&updateChatDto); err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		chatDto, err := chatService.Update(&updateChatDto, userId, chatId)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, chatDto)
	}
}

// @Summary Add chat member
// @Description Adds a member to a chat.
// @Tags chats
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param body body dto.AddChatMemberDto true "Member payload"
// @Success 201 {object} dto.ChatMemberDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId}/members [post]
func AddChatMemberHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var addChatMemberDto dto.AddChatMemberDto
		if err := ctx.ShouldBindBodyWithJSON(&addChatMemberDto); err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		rawChatId := ctx.Param("chatId")
		chatId, err := strconv.ParseInt(rawChatId, 10, 64)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Invalid chat id provided.", http.StatusBadRequest))
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		chatMemberDto, err := chatService.AddMember(userId, chatId, addChatMemberDto)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, chatMemberDto)
	}
}

// @Summary Update chat member role
// @Description Updates a member role in a chat.
// @Tags chats
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param body body dto.ChangeRoleDto true "Role update payload"
// @Success 200 {object} dto.ChatMemberDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId}/members [put]
func UpdateChatMemberHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var changeRoleDto dto.ChangeRoleDto
		err := ctx.ShouldBindJSON(&changeRoleDto)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatIdString, ok := ctx.Params.Get("chatId")
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatId, err := strconv.ParseInt(chatIdString, 10, 64)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatMemberDto, err := chatService.UpdateMemberRole(userId, chatId, changeRoleDto)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, chatMemberDto)
	}
}

// @Summary Delete chat member
// @Description Deletes a member from a chat.
// @Tags chats
// @Accept json
// @Param chatId path int true "Chat ID"
// @Param body body dto.DeleteChatMemberDto true "Member delete payload"
// @Success 204
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId}/members [delete]
func DeleteChatMemberHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var deleteChatMemberDto dto.DeleteChatMemberDto
		err := ctx.ShouldBindJSON(&deleteChatMemberDto)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatIdString, ok := ctx.Params.Get("chatId")
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatId, err := strconv.ParseInt(chatIdString, 10, 64)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		err = chatService.DeleteMember(userId, chatId, deleteChatMemberDto)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

// @Summary Finds by user id
// @Tags chats
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ChatDto
// @Failure 401 {object} exceptions.HttpError
// @Failure 403 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats [get]
func GetChatsByUserIdHandler(chatService *services.ChatService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.NewHttpError("Error converting user id.", http.StatusInternalServerError))
			return
		}

		chats, err := chatService.FindByUserId(userId)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, chats)
	}
}

// @Summary Post typing status
// @Description Emits a typing (writing) event for a chat to connected members.
// @Tags chats
// @Accept json
// @Produce json
// @Param chatId path int true "Chat ID"
// @Param body body dto.WritingEventDto true "Typing event payload"
// @Success 204
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /chats/{chatId}/typing [put]
func TypingHandler(
	userRepo repositories.UserRepository,
	chatRepo repositories.ChatRepository,
	messenger *messenger.EventBus) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var writingEventDto dto.WritingEventDto
		err := ctx.BindJSON(&writingEventDto)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatIdString, ok := ctx.Params.Get("chatId")
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chatId, err := strconv.ParseInt(chatIdString, 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}

		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		chat, err := chatRepo.FindById(chatId)
		if err != nil {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		user, err := userRepo.FindById(userId)
		if err != nil {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		messenger.PostTypingEvent(*chat, *user, writingEventDto.Typing)
	}
}
