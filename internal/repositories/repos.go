package repositories

import "database/sql"

type Repositories struct {
	UserRepository        UserRepository
	MessageRepository     MessageRepository
	FriendshipRequestRepo FriendshipRequestRepository
	FriendshipRepository  FriendshipRepository
	ChatRepository        ChatRepository
	ChatMemberRepository  ChatMemberRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:        &MySQLUserRepository{DB: db},
		MessageRepository:     &MySQLMessageRepository{DB: db},
		FriendshipRequestRepo: &MySQLFriendshipRequestRepository{DB: db},
		FriendshipRepository:  &MySQLFriendshipRepository{DB: db},
		ChatRepository:        &MySQLChatRepository{DB: db},
		ChatMemberRepository:  &MySQLChatMemberRepository{DB: db},
	}
}
