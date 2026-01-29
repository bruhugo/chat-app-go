package repositories

import (
	"database/sql"
	"log"

	"github.com/grongoglongo/chatter-go/internal/models"
)

type MessageRepository struct {
	DB *sql.DB
}

func (repo *MessageRepository) Create(m *models.Message) error {
	result, err := repo.DB.Exec("INSERT INTO messages (content, sender_id, receiver_id) VALUES (?, ?, ?)",
		m.Content,
		m.Sender.ID,
		m.Receiver.ID,
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

func (repo *MessageRepository) FindById(id int64) (*models.Message, error) {
	row := repo.DB.QueryRow(
		"SELECT m.id, m.content, s.id, s.username, s.email, r.id, r.username, r.email "+
			"FROM messages m "+
			"JOIN users s ON m.sender_id = s.id "+
			"JOIN users r ON m.receiver_id = m.id "+
			"WHERE id = ?",
		id)

	u, err := scanMessage(row)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *MessageRepository) FindBySenderAndReceiver(sId int64, rId int64) (*models.Message, error) {
	row := repo.DB.QueryRow(
		"SELECT m.id, m.content, s.id, s.username, s.email, r.id, r.username, r.email "+
			"FROM messages m "+
			"JOIN users s ON m.sender_id = s.id "+
			"JOIN users r ON m.receiver_id = m.id "+
			"WHERE s.id = ? AND r.id = ?",
		sId, rId)

	u, err := scanMessage(row)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *MessageRepository) DeleteById(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return err
	}

	log.Printf("Message with id %d deleted", id)

	return nil
}

func (repo *MessageRepository) PatchContent(id int64, content string) error {
	_, err := repo.DB.Exec("UPDATE messages SET content = ? WHERE id = id", content, id)
	if err != nil {
		return err
	}

	log.Printf("Message with id %d updated.", id)

	return nil
}

func scanMessage(row *sql.Row) (*models.Message, error) {
	var m models.Message
	err := row.Scan(
		&m.ID,
		&m.Content,
		&m.Sender.ID,
		&m.Sender.Username,
		&m.Sender.Email,
		&m.Receiver.ID,
		&m.Receiver.Username,
		&m.Receiver.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}
