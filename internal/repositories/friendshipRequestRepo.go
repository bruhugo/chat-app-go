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
		request.Receiver.ID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	request.ID = id

	createdRequest, err := repo.findById(id)
	if err != nil {
		return err
	}
	if createdRequest != nil {
		*request = *createdRequest
	}

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

func (repo *MySQLFriendshipRequestRepository) findById(id int64) (*models.FriendshipRequest, error) {
	row := repo.DB.QueryRow(
		"SELECT fr.id, s.id, s.username, s.email, r.id, r.username, r.email, fr.created_at, fr.updated_at "+
			"FROM friendship_requests fr "+
			"JOIN users s ON fr.sender_id = s.id "+
			"JOIN users r ON fr.receiver_id = r.id "+
			"WHERE fr.id = ?",
		id,
	)

	var request models.FriendshipRequest
	err := row.Scan(
		&request.ID,
		&request.Sender.ID,
		&request.Sender.Username,
		&request.Sender.Email,
		&request.Receiver.ID,
		&request.Receiver.Username,
		&request.Receiver.Email,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &request, nil
}
