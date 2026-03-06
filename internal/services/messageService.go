package services

import (
	"log"
	"time"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/messenger"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
)

type MessageService struct {
	messageRepo repositories.MessageRepository
	chatRepo    repositories.ChatRepository
	eventBus    *messenger.EventBus
}

func NewMessageService(
	messageRepo repositories.MessageRepository,
	chatRepository repositories.ChatRepository,
	eventBus *messenger.EventBus,
) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepository,
		eventBus:    eventBus,
	}
}

func (ms *MessageService) CreateMessage(content string, userId, chatId int64) (*dto.MessageDto, error) {
	chat, err := ms.chatRepo.FindById(chatId)
	if err != nil {
		log.Print(err.Error())
		return nil, exceptions.InternalServerError
	}
	if chat == nil {
		return nil, exceptions.NotFoundError
	}

	message := &models.Message{
		Content:   content,
		User:      &models.User{ID: userId},
		Chat:      &models.Chat{ID: chatId, Creator: &models.User{}},
		CreatedAt: time.Now(),
	}

	err = ms.messageRepo.Create(message)
	if err != nil {
		log.Print(err.Error())
		return nil, exceptions.InternalServerError
	}

	ms.eventBus.PostCreateMessageEvent(*message)

	return message.ToDto(), nil
}

func (ms *MessageService) DeleteMessage(messageId int64, userId int64) error {
	message, err := ms.messageRepo.FindById(messageId)
	if err != nil {
		log.Print(err.Error())
		return exceptions.InternalServerError
	}
	if message == nil {
		return exceptions.NewHttpError("Message does not exist.", exceptions.NotFoundError.Status)
	}
	if message.User.ID != userId {
		return exceptions.ForbiddenError
	}

	err = ms.messageRepo.Delete(messageId)
	if err != nil {
		log.Print(err.Error())
		return exceptions.InternalServerError
	}

	ms.eventBus.PostDeleteMessageEvent(messageId, *message.Chat)

	return nil
}

func (ms *MessageService) UpdateMessageContent(updateMessageDto dto.UpdateMessageDto) error {
	message, err := ms.messageRepo.FindById(updateMessageDto.MessageId)
	if err != nil {
		log.Print(err.Error())
		return exceptions.InternalServerError
	}
	if message == nil {
		return exceptions.NewHttpError("Message does not exist.", exceptions.NotFoundError.Status)
	}
	if message.User.ID != updateMessageDto.UserId {
		return exceptions.ForbiddenError
	}

	err = ms.messageRepo.PatchContent(updateMessageDto.MessageId, updateMessageDto.NewContent)
	if err != nil {
		log.Print(err.Error())
		return exceptions.InternalServerError
	}

	return nil
}

func (ms *MessageService) GetMessages(chatId, userId int64, pageRequest *dto.PageRequest) (*dto.Page[dto.MessageDto], error) {
	isMember, err := ms.chatRepo.IsUserMember(chatId, userId)
	if err != nil {
		return nil, exceptions.InternalServerError
	}
	if !isMember {
		return nil, exceptions.ForbiddenError
	}

	page, err := ms.messageRepo.FindByChat(chatId, *pageRequest)

	if err != nil {
		return nil, err
	}

	return page, nil
}
