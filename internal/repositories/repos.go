package repositories

import "database/sql"

type Repositories struct {
	UserRepository    *UserRepository
	MessageRepository *MessageRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:    NewUserRepository(db),
		MessageRepository: NewMessageRepository(db),
	}
}
