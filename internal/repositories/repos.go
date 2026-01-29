package repositories

import "database/sql"

type Repositories struct {
	UserRepository        *UserRepository
	MessageRepository     *MessageRepository
	FriendshipRequestRepo *FriendshipRequestRepo
	FriendshipRepository  *FriendshipRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:        &UserRepository{DB: db},
		MessageRepository:     &MessageRepository{DB: db},
		FriendshipRequestRepo: &FriendshipRequestRepo{DB: db},
		FriendshipRepository:  &FriendshipRepository{DB: db},
	}
}
