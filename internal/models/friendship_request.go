package models

type FriendshipRequest struct {
	ID        int64 `db:"id"`
	Sender    User
	Receiver  User
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}
