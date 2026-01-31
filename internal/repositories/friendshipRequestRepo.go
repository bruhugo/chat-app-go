package repositories

import (
	"database/sql"

	"github.com/grongoglongo/chatter-go/internal/models"
)

type FriendshipRequestRepository interface {
	Create(request *models.FriendshipRequest) error
	DeleteById(id int64) error
	FindByReceiverId(id int64) ([]models.FriendshipRequest, error)
}

type MySQLFriendshipRequestRepository struct {
	DB *sql.DB
}

func (repo *MySQLFriendshipRequestRepository) Create(request *models.FriendshipRequest) error {
	result, err := repo.DB.Exec("INSERT INTO friendship_requests (sender_id, receiver_id) VALUES (?, ?)",
		request.Sender.ID,
		request.Sender.ID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	request.ID = id

	return nil
}

func (repo *MySQLFriendshipRequestRepository) DeleteById(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM friendship_requests WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MySQLFriendshipRequestRepository) FindByReceiverId(id int64) ([]models.FriendshipRequest, error) {
	rows, err := repo.DB.Query("SELECT (fr.id, s.id, s.username, s.email, r.id, r.username, r.email) FROM friendship_requests fr "+
		"JOIN users s ON fr.sender_id = s.id "+
		"JOIN users r ON fr.sender_id = r.id "+
		"WHERE r.id = ?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var requests []models.FriendshipRequest

	for rows.Next() {
		var request models.FriendshipRequest
		rows.Scan(&request.ID,
			&request.Sender.ID, &request.Sender.ID, &request.Sender.ID,
			&request.Receiver.ID, &request.Receiver.Username, &request.Receiver.Email,
		)
		requests = append(requests, request)
	}

	return requests, nil
}
