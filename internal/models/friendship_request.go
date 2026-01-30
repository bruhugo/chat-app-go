package models

import "time"

type FriendshipRequest struct {
	ID        int64 `db:"id"`
	Sender    User
	Receiver  User
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}
