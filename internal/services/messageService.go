package services

import (
	"log"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
)

type MessageService struct {
	messageRepo repositories.MessageRepository
	userRepo    repositories.UserRepository
}

func NewMessageService(messageRepo repositories.MessageRepository, userRepo repositories.UserRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (ms *MessageService) CreateMessage(createMessageDto dto.CreateMessageDto, userId int64) (*dto.MessageDto, error) {
	foundTarget, err := ms.userRepo.FindById(createMessageDto.TargetUserId)
	if err != nil {
		log.Print(err.Error())
		return nil, exceptions.InternalServerError
	}
	if foundTarget == nil {
		return nil, exceptions.NotFoundError
	}

	message := &models.Message{
		Content:  createMessageDto.Content,
		Sender:   &models.User{ID: userId},
		Receiver: &models.User{ID: createMessageDto.TargetUserId},
	}

	err = ms.messageRepo.Create(message)
	if err != nil {
		log.Print(err.Error())
		return nil, exceptions.InternalServerError
	}

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
	if message.Sender.ID != userId {
		return exceptions.ForbiddenError
	}

	err = ms.messageRepo.Delete(messageId)
	if err != nil {
		log.Print(err.Error())
		return exceptions.InternalServerError
	}

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
	if message.Sender.ID != updateMessageDto.SenderId {
		return exceptions.ForbiddenError
	}

	err = ms.messageRepo.PatchContent(updateMessageDto.MessageId, updateMessageDto.NewContent)
	if err != nil {
		log.Print(err.Error())
		return exceptions.InternalServerError
	}

	return nil
}

func (ms *MessageService) GetMessages(getMessagesDto dto.GetMessagesDto) (*dto.Page[dto.MessageDto], error) {
	page, err := ms.messageRepo.FindBySenderAndReceiver(getMessagesDto.SenderId, getMessagesDto.ReceiverId, getMessagesDto.PageRequest)
	if err != nil {
		return nil, err
	}

	return page, nil
}
