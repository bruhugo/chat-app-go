package repositories

import (
	"database/sql"
	"log"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type MessageRepository interface {
	Create(m *models.Message) error
	FindById(id int64) (*models.Message, error)
	FindByChat(chatId int64, pageRequest dto.PageRequest) (*dto.Page[dto.MessageDto], error)
	PatchContent(id int64, content string) error
	Delete(id int64) error
}

type MySQLMessageRepository struct {
	DB *sql.DB
}

func (repo *MySQLMessageRepository) Create(m *models.Message) error {
	result, err := repo.DB.Exec("INSERT INTO messages (content, chat_id, user_id) VALUES (?, ?, ?)",
		m.Content,
		m.Chat.ID,
		m.User.ID,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = id

	log.Printf("Message created with id %d", id)

	return nil
}

func (repo *MySQLMessageRepository) FindById(id int64) (*models.Message, error) {
	row := repo.DB.QueryRow(
		"SELECT m.id, m.content, m.created_at, u.id, u.username, u.email "+
			"FROM messages m "+
			"JOIN users u ON m.user_id = u.id "+
			"WHERE m.id = ?",
		id)

	u, err := scanMessage(row)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *MySQLMessageRepository) FindByChat(chatId int64, pageRequest dto.PageRequest) (*dto.Page[dto.MessageDto], error) {
	page := &dto.Page[dto.MessageDto]{}
	rows, err := repo.DB.Query(
		"SELECT m.id, m.content, m.created_at, u.id, u.username, u.email "+
			"FROM messages m "+
			"JOIN users u ON m.user_id = u.id "+
			"WHERE m.chat_id = ? "+
			"ORDER BY m.created_at DESC "+
			"LIMIT ? OFFSET ? ",
		chatId, pageRequest.PageSize, pageRequest.Page*page.PageSize)

	if err != nil {
		log.Print(err.Error())
		return nil, exceptions.InternalServerError
	}

	messages, err := scanMessages(rows)
	if err != nil {
		log.Print(err.Error())
		return nil, exceptions.InternalServerError
	}

	dtoContent := []dto.MessageDto{}
	for _, m := range messages {
		dtoContent = append(dtoContent, *m.ToDto())
	}

	page.Content = dtoContent
	page.Number = len(messages)
	page.Page = pageRequest.Page
	page.PageSize = pageRequest.PageSize

	return page, nil
}

func (repo *MySQLMessageRepository) DeleteById(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return err
	}

	log.Printf("Message with id %d deleted", id)

	return nil
}

func (repo *MySQLMessageRepository) PatchContent(id int64, content string) error {
	_, err := repo.DB.Exec("UPDATE messages SET content = ? WHERE id = id", content, id)
	if err != nil {
		return err
	}

	log.Printf("Message with id %d updated.", id)

	return nil
}

func (repo *MySQLMessageRepository) Delete(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM message WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func scanMessage(row *sql.Row) (*models.Message, error) {
	var m models.Message
	err := row.Scan(
		&m.ID,
		&m.Content,
		&m.CreatedAt,
		&m.User.ID,
		&m.User.Username,
		&m.User.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func scanMessages(rows *sql.Rows) (messages []models.Message, _ error) {
	for rows.Next() {
		m := models.Message{
			User: &models.User{},
			Chat: &models.Chat{
				Creator: &models.User{},
			},
		}
		err := rows.Scan(
			&m.ID,
			&m.Content,
			&m.CreatedAt,
			&m.User.ID,
			&m.User.Username,
			&m.User.Email,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}
