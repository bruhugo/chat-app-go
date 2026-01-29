package repositories

import (
	"database/sql"
	"log"

	"github.com/grongoglongo/chatter-go/internal/models"
)

type FriendshipRepository struct {
	DB *sql.DB
}

func (repo *FriendshipRepository) Create(friendship *models.Friendship) error {
	result, err := repo.DB.Exec(
		"INSERT INTO friendships (user1_id, user2_id) VALUES (?, ?)",
		friendship.User1.ID,
		friendship.User2.ID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	friendship.ID = id

	log.Printf("Friendship created with id %d", id)

	return nil
}

func (repo *FriendshipRepository) DeleteById(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM friendships WHERE id = ?", id)
	if err != nil {
		return err
	}

	log.Printf("Friendship with id %d deleted", id)

	return nil
}

func (repo *FriendshipRepository) FindByUser(userId int64) ([]models.Friendship, error) {
	rows, err := repo.DB.Query(
		"SELECT f.id, u1.id, u1.username, u1.email, u2.id, u2.username, u2.email, f.created_at, f.updated_at "+
			"FROM friendships f "+
			"JOIN users u1 ON f.user1_id = u1.id "+
			"JOIN users u2 ON f.user2_id = u2.id "+
			"WHERE f.user1_id = ? OR f.user2_id = ?",
		userId, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var friendships []models.Friendship

	for rows.Next() {
		friendship, err := scanFriendshipRow(rows)
		if err != nil {
			return nil, err
		}
		friendships = append(friendships, *friendship)
	}

	return friendships, nil
}

func scanFriendshipRow(rows *sql.Rows) (*models.Friendship, error) {
	var f models.Friendship
	err := rows.Scan(
		&f.ID,
		&f.User1.ID,
		&f.User1.Username,
		&f.User1.Email,
		&f.User2.ID,
		&f.User2.Username,
		&f.User2.Email,
		&f.CreatedAt,
		&f.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &f, nil
}
