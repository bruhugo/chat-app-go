package repositories

import (
	"database/sql"
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo *UserRepository) Create(user *models.User) error {
	result, err := repo.DB.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		user.Username, user.Email, user.Password,
	)
	if err != nil {
		var myErr *mysql.MySQLError
		if errors.As(err, &myErr) {
			if myErr.Number == 1062 {
				return exceptions.ConflictSqlError
			}
		}

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
	row := repo.DB.QueryRow("SELECT username, password_hash, id, email, created_at, updated_at FROM users WHERE id = ?", id)
	user, err := scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	row := repo.DB.QueryRow("SELECT username, password_hash, id, email, created_at, updated_at FROM users WHERE email = ?", email)
	user, err := scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) Update(id int64, user *models.User) error {
	_, err := repo.DB.Exec("UPDATE users SET username = ?, password_hash = ?, email = ?, id = ? WHERE id = ?",
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
	err := row.Scan(&u.Username, &u.Password, &u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}
