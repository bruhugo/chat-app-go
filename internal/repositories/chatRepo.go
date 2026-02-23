package repositories

import (
	"database/sql"
	"errors"

	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type ChatRepository interface {
	Create(chat *models.Chat) error
	Delete(id int64) error
	FindByUser(userId int64) ([]*dto.ChatResponseDto, error)
	Update(id int64, newChat *dto.UpdateChatDto) error
	FindById(id int64) (*models.Chat, error)
	IsUserMember(chatId, userId int64) (bool, error)
}

type MySQLChatRepository struct {
	DB *sql.DB
}

func NewMySQLChatRepository(db *sql.DB) *MySQLChatRepository {
	return &MySQLChatRepository{
		DB: db,
	}
}

func (r *MySQLChatRepository) FindById(id int64) (*models.Chat, error) {
	row := r.DB.QueryRow(
		"SELECT c.id, c.name, c.description, c.created_at, cr.id, cr.email, cr.username "+
			"FROM chats c JOIN users cr ON c.creator_id = cr.id "+
			"WHERE c.id = ?", id,
	)

	chat := &models.Chat{
		Creator: &models.User{},
	}
	err := row.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatedAt, &chat.Creator.ID, &chat.Creator.Email, &chat.Creator.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return chat, nil
}

func (r *MySQLChatRepository) IsUserMember(chatId, userId int64) (bool, error) {
	row := r.DB.QueryRow("SELECT * FROM chat_members WHERE chat_id = ? AND user_id = ?", chatId, userId)
	var any any
	err := row.Scan(&any, &any, &any, &any, &any)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, err
}

func (r *MySQLChatRepository) Create(chat *models.Chat) error {
	result, err := r.DB.Exec("INSERT INTO chats (name, description, creator_id) VALUES (?, ?, ?)", chat.Name, chat.Description, chat.Creator.ID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	chat.ID = id

	createdChat, err := r.FindById(id)
	if err != nil {
		return err
	}
	if createdChat != nil {
		*chat = *createdChat
	}

	return nil
}

func (r *MySQLChatRepository) Delete(id int64) error {
	_, err := r.DB.Exec("DELETE FROM chats WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLChatRepository) FindByUser(userId int64) ([]*dto.ChatResponseDto, error) {
	rows, err := r.DB.Query(
		"SELECT c.id, c.name, c.description, c.created_at, "+
			"cr.id, cr.username, cr.email, "+
			"m.content, ma.username "+
			"FROM chat_members cm "+
			"JOIN chats c ON cm.chat_id = c.id "+
			"JOIN users cr ON c.creator_id = cr.id "+
			"LEFT JOIN ( "+
			"	SELECT * "+
			"	FROM ( "+
			"		SELECT m.*, "+
			"			ROW_NUMBER() OVER ( "+
			"				PARTITION BY chat_id "+
			"				ORDER BY created_at DESC "+
			"			) as rn"+
			"		FROM messages m"+
			"	) ranked WHERE rn = 1"+
			") m ON c.id = m.chat_id "+
			"LEFT JOIN users ma ON m.user_id = ma.id "+
			"WHERE cm.user_id = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}

	chats := []*dto.ChatResponseDto{}
	for rows.Next() {
		c := &dto.ChatDto{
			Creator: &dto.UserDto{},
		}

		var content sql.NullString
		var username sql.NullString

		err = rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.Creator.ID, &c.Creator.Username, &c.Creator.Email, &content, &username)
		if err != nil {
			return nil, err
		}

		chatResponseDto := &dto.ChatResponseDto{
			ChatDto: c,
		}

		if content.Valid {
			m := &dto.MessageDto{
				User: &dto.UserDto{
					Username: username.String,
				},
				Chat:    c,
				Content: content.String,
			}
			chatResponseDto.LastMessage = m
		}

		chats = append(chats, chatResponseDto)
	}

	return chats, nil
}

func (r *MySQLChatRepository) Update(chatId int64, chat *dto.UpdateChatDto) error {
	_, err := r.DB.Exec("UPDATE chats SET name = ?, description = ? WHERE id = ?", chat.Name, chat.Description, chatId)
	if err != nil {
		return err
	}

	return nil
}
