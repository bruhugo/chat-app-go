package services

import (
	"testing"

	"github.com/grongoglongo/chatter-go/internal/messenger"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/stretchr/testify/require"
)

var eventBus = messenger.NewEventBus(messenger.NewInMemoryMessenger(), messenger.NewConnectionHub())

type MockMessageRepository struct {
	messagePage *dto.Page[dto.MessageDto]
}

func (mr *MockMessageRepository) Create(m *models.Message) error             { return nil }
func (mr *MockMessageRepository) FindById(id int64) (*models.Message, error) { return nil, nil }
func (mr *MockMessageRepository) FindByChat(chatId int64, pageRequest dto.PageRequest) (*dto.Page[dto.MessageDto], error) {
	if mr.messagePage != nil {
		return mr.messagePage, nil
	}
	return &dto.Page[dto.MessageDto]{}, nil
}
func (mr *MockMessageRepository) PatchContent(id int64, content string) error { return nil }
func (mr *MockMessageRepository) Delete(id int64) error                       { return nil }

type MockChatRepository struct {
}

func (MockChatRepository) Create(chat *models.Chat) error                          { return nil }
func (MockChatRepository) Delete(id int64) error                                   { return nil }
func (MockChatRepository) FindByUser(userId int64) ([]*dto.ChatResponseDto, error) { return nil, nil }
func (MockChatRepository) Update(id int64, newChat *dto.UpdateChatDto) error       { return nil }
func (MockChatRepository) FindById(id int64) (*models.Chat, error)                 { return &models.Chat{}, nil }
func (MockChatRepository) IsUserMember(chatId, userId int64) (bool, error)         { return true, nil }

func TestMessageService_GetMessages(t *testing.T) {
	messageRepo := &MockMessageRepository{
		messagePage: &dto.Page[dto.MessageDto]{
			Content:  []dto.MessageDto{*getMessageDto(), *getMessageDto(), *getMessageDto()},
			Number:   3,
			Page:     0,
			PageSize: 5,
		},
	}

	messageService := NewMessageService(messageRepo, &MockChatRepository{}, eventBus)

	page, err := messageService.GetMessages(1, 1, &dto.PageRequest{Page: 0, PageSize: 5})

	require.NoError(t, err)
	require.Equal(t, 5, page.PageSize)
	require.Equal(t, 0, page.Page)
	require.Equal(t, 3, page.Number)
	require.Equal(t, 3, len(page.Content))
}

func TestMessageService_CreateMessage(t *testing.T) {
	messageRepo := &MockMessageRepository{}
	messageService := NewMessageService(messageRepo, &MockChatRepository{}, eventBus)
	createMessageDto := dto.CreateMessageDto{ChatId: 1, Content: "content"}
	messageDto, err := messageService.CreateMessage(createMessageDto, 1)

	require.NoError(t, err)
	require.Equal(t, messageDto.Content, createMessageDto.Content)
	require.Equal(t, messageDto.Chat.ID, createMessageDto.ChatId)
}

func getMessageDto() *dto.MessageDto {
	return &dto.MessageDto{
		ID:      12,
		Content: "huge content",
		User:    &dto.UserDto{},
		Chat:    &dto.ChatDto{},
	}
}
