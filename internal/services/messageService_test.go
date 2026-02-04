package services

import (
	"testing"

	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/stretchr/testify/require"
)

type MockMessageRepository struct {
	messagePage *dto.Page[dto.MessageDto]
}

func (mr *MockMessageRepository) Create(m *models.Message) error             { return nil }
func (mr *MockMessageRepository) FindById(id int64) (*models.Message, error) { return nil, nil }
func (mr *MockMessageRepository) FindBySenderAndReceiver(sId int64, rId int64, pageRequest dto.PageRequest) (*dto.Page[dto.MessageDto], error) {
	if mr.messagePage != nil {
		return mr.messagePage, nil
	}
	return &dto.Page[dto.MessageDto]{}, nil
}
func (mr *MockMessageRepository) PatchContent(id int64, content string) error { return nil }
func (mr *MockMessageRepository) Delete(id int64) error                       { return nil }

func TestMessageService_GetMessages(t *testing.T) {
	messageRepo := &MockMessageRepository{
		messagePage: &dto.Page[dto.MessageDto]{
			Content:  []dto.MessageDto{*getMessageDto(), *getMessageDto(), *getMessageDto()},
			Number:   3,
			Page:     0,
			PageSize: 5,
		},
	}

	messageService := NewMessageService(messageRepo, &MockUserRepo{})

	_, err := messageService.GetMessages(dto.GetMessagesDto{
		SenderId:    1,
		ReceiverId:  2,
		PageRequest: dto.PageRequest{Page: 0, PageSize: 5},
	})

	require.NoError(t, err)
}

func getMessageDto() *dto.MessageDto {
	return &dto.MessageDto{
		ID:       12,
		Content:  "huge content",
		Sender:   &dto.UserDto{},
		Receiver: &dto.UserDto{},
	}
}
