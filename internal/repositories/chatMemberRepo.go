package repositories

import (
	"database/sql"

	"github.com/grongoglongo/chatter-go/internal/models"
)

type ChatMemberRepository interface {
	Create(member *models.ChatMember) error
	FindById(id int64) (*models.ChatMember, error)
	FindByUserIdAndChatId(userId int64, chatId int64) (*models.ChatMember, error)
	Update(id int64, member *models.ChatMember) error
	Delete(id int64) error
}

type MySQLChatMemberRepository struct {
	DB *sql.DB
}

func NewMySQLChatMemberRepository(db *sql.DB) *MySQLChatMemberRepository {
	return &MySQLChatMemberRepository{
		DB: db,
	}
}

func (repo *MySQLChatMemberRepository) Create(member *models.ChatMember) error {
	result, err := repo.DB.Exec(
		"INSERT INTO chat_members (user_id, chat_id, role) VALUES (?, ?, ?)",
		member.User.ID,
		member.Chat.ID,
		member.Role,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	member.ID = id

	createdMember, err := repo.FindById(id)
	if err != nil {
		return err
	}
	if createdMember != nil {
		*member = *createdMember
	}

	return nil
}

func (repo *MySQLChatMemberRepository) FindById(id int64) (*models.ChatMember, error) {
	row := repo.DB.QueryRow(
		"SELECT cm.id, cm.role, u.id, u.username, u.email, "+
			"c.id, c.name, c.description, c.created_at, "+
			"cr.id, cr.username, cr.email "+
			"FROM chat_members cm "+
			"JOIN users u ON cm.user_id = u.id "+
			"JOIN chats c ON cm.chat_id = c.id "+
			"JOIN users cr ON c.creator_id = cr.id "+
			"WHERE cm.id = ?",
		id,
	)

	member, err := scanChatMemberRow(row)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (repo *MySQLChatMemberRepository) FindByUserIdAndChatId(userId int64, chatId int64) (*models.ChatMember, error) {
	row := repo.DB.QueryRow(
		"SELECT cm.id, cm.role, u.id, u.username, u.email, "+
			"c.id, c.name, c.description, c.created_at, "+
			"cr.id, cr.username, cr.email "+
			"FROM chat_members cm "+
			"JOIN users u ON cm.user_id = u.id "+
			"JOIN chats c ON cm.chat_id = c.id "+
			"JOIN users cr ON c.creator_id = cr.id "+
			"WHERE u.id = ? AND c.id = ?",
		userId,
		chatId,
	)

	member, err := scanChatMemberRow(row)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (repo *MySQLChatMemberRepository) Update(id int64, member *models.ChatMember) error {
	_, err := repo.DB.Exec(
		"UPDATE chat_members SET user_id = ?, chat_id = ?, role = ? WHERE id = ?",
		member.User.ID,
		member.Chat.ID,
		member.Role,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MySQLChatMemberRepository) Delete(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM chat_members WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func scanChatMemberRow(row *sql.Row) (*models.ChatMember, error) {
	cm := models.ChatMember{
		Chat: models.Chat{
			Creator: &models.User{},
		},
	}
	err := row.Scan(
		&cm.ID,
		&cm.Role,
		&cm.User.ID,
		&cm.User.Username,
		&cm.User.Email,
		&cm.Chat.ID,
		&cm.Chat.Name,
		&cm.Chat.Description,
		&cm.Chat.CreatedAt,
		&cm.Chat.Creator.ID,
		&cm.Chat.Creator.Username,
		&cm.Chat.Creator.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cm, nil
}
