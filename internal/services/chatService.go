package services

import (
	"net/http"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
)

type ChatService struct {
	chatRepo       repositories.ChatRepository
	chatMemberRepo repositories.ChatMemberRepository
}

func NewChatService(chatRepo repositories.ChatRepository, chatMemberRepo repositories.ChatMemberRepository) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		chatMemberRepo: chatMemberRepo,
	}
}

func (s *ChatService) CreateChat(createChatDto dto.CreateChatDto) (*dto.ChatDto, error) {
	chat := &models.Chat{
		Name:        createChatDto.Name,
		Description: createChatDto.Description,
		Creator: &models.User{
			ID: createChatDto.CreatorId,
		},
	}

	err := s.chatRepo.Create(chat)
	if err != nil {
		return nil, exceptions.NewHttpError("Error creating chat. Please try again later.", http.StatusInternalServerError)
	}

	return chat.ToDto(), nil
}

func (s *ChatService) DeleteChat(userId, chatId int64) error {
	chat, err := s.chatRepo.FindById(chatId)
	if err != nil {
		return err
	}
	if chat == nil {
		return nil
	}

	chatMember, err := s.chatMemberRepo.FindByUserIdAndChatId(userId, chatId)
	if err != nil {
		return nil
	}
	if chatMember == nil || chatMember.Role != models.ADMIN {
		return exceptions.ForbiddenError
	}

	err = s.chatRepo.Delete(userId)
	return err
}
