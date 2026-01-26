package repositories

import (
	"database/sql"
	"log"

	"github.com/grongoglongo/chatter-go/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Create(user *models.User) error {
	result, err := repo.db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		user.Username, user.Email, user.Password,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id

	log.Printf("Created user %s with id %d", user.Username, id)

	return nil
}

func (repo *UserRepository) FindById(id int64) (*models.User, error) {
	row := repo.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	user, err := scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	row := repo.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	user, err := scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) Update(id int64, user *models.User) error {
	_, err := repo.db.Exec("UPDATE users SET username = ?, password_hash = ?, email = ?, id = ? WHERE id = ?",
		user.Username,
		user.Password,
		user.Email,
		user.ID,
		id,
	)

	if err != nil {
		return err
	}

	log.Printf("Updated user with id %d", id)

	return nil
}

func scanUserRow(row *sql.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.ID, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
