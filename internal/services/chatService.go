package services

import (
	"net/http"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/messenger"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
)

type ChatService struct {
	chatRepo       repositories.ChatRepository
	chatMemberRepo repositories.ChatMemberRepository
	userRepo       repositories.UserRepository
	eventBus       *messenger.EventBus
}

func NewChatService(
	chatRepo repositories.ChatRepository,
	chatMemberRepo repositories.ChatMemberRepository,
	userRepo repositories.UserRepository,
	eventBus *messenger.EventBus,
) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		chatMemberRepo: chatMemberRepo,
		userRepo:       userRepo,
		eventBus:       eventBus,
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

	err = s.chatMemberRepo.Create(&models.ChatMember{
		User: models.User{ID: createChatDto.CreatorId},
		Role: dto.ADMIN,
		Chat: models.Chat{ID: chat.ID},
	})

	s.eventBus.PostEnterChatEvent(*chat, *chat.Creator, *chat.Creator)

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
	if chatMember == nil || chatMember.Role != dto.ADMIN {
		return exceptions.ForbiddenError
	}

	err = s.chatRepo.Delete(userId)
	return err
}

func (s *ChatService) Update(updateChatDto *dto.UpdateChatDto, userId, chatId int64) (*dto.ChatDto, error) {
	chatMember, err := s.chatMemberRepo.FindByUserIdAndChatId(userId, chatId)
	if err != nil {
		return nil, exceptions.InternalServerError
	}
	if chatMember == nil || chatMember.Role != dto.ADMIN {
		return nil, exceptions.ForbiddenError
	}

	s.chatRepo.Update(chatId, updateChatDto)

	chat, err := s.chatRepo.FindById(chatId)
	if err != nil {
		return nil, err
	}

	return chat.ToDto(), nil
}

func (s *ChatService) AddMember(userId, chatId int64, addChatMemberDto dto.AddChatMemberDto) (*dto.ChatMemberDto, error) {
	chatMember, err := s.chatMemberRepo.FindByUserIdAndChatId(userId, chatId)
	if err != nil {
		return nil, exceptions.InternalServerError
	}
	if chatMember == nil || chatMember.Role != dto.ADMIN {
		return nil, exceptions.ForbiddenError
	}

	user, err := s.userRepo.FindById(addChatMemberDto.TargetId)
	if err != nil {
		return nil, exceptions.InternalServerError
	}
	if user == nil {
		return nil, exceptions.NotFoundError
	}

	conflictChatMember, err := s.chatMemberRepo.FindByUserIdAndChatId(addChatMemberDto.TargetId, chatId)
	if err != nil {
		return nil, exceptions.InternalServerError
	}
	if conflictChatMember != nil {
		return nil, exceptions.NewHttpError("That user is already in the chat", http.StatusConflict)
	}

	newChatMember := &models.ChatMember{
		Chat: chatMember.Chat,
		User: models.User{ID: addChatMemberDto.TargetId, Username: user.Username, Email: user.Email},
		Role: addChatMemberDto.Role,
	}

	s.chatMemberRepo.Create(newChatMember)

	s.eventBus.PostEnterChatEvent(chatMember.Chat, chatMember.User, *user)

	return newChatMember.ToDto(), nil
}

func (s *ChatService) FindByUserId(userId int64) ([]*dto.ChatResponseDto, error) {
	chats, err := s.chatRepo.FindByUser(userId)
	if err != nil {
		return nil, exceptions.InternalServerError
	}

	return chats, nil
}
