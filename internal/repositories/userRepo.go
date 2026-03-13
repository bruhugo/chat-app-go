package repositories

import (
	"database/sql"
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindById(id int64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	SearchByUsername(username string) ([]models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(id int64, user *models.User) error
}

type MySQLUserRepository struct {
	DB *sql.DB
}

func (repo *MySQLUserRepository) Create(user *models.User) error {
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

	createdUser, err := repo.FindById(id)
	if err != nil {
		return err
	}
	if createdUser != nil {
		*user = *createdUser
	}

	log.Printf("Created user %s with id %d", user.Username, id)

	return nil
}

func (repo *MySQLUserRepository) FindById(id int64) (*models.User, error) {
	row := repo.DB.QueryRow("SELECT username, password_hash, id, email, created_at, updated_at FROM users WHERE id = ?", id)
	user, err := repo.scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *MySQLUserRepository) FindByEmail(email string) (*models.User, error) {
	row := repo.DB.QueryRow("SELECT username, password_hash, id, email, created_at, updated_at FROM users WHERE email = ?", email)
	user, err := repo.scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *MySQLUserRepository) FindByUsername(username string) (*models.User, error) {
	row := repo.DB.QueryRow("SELECT username, password_hash, id, email, created_at, updated_at FROM users WHERE username = ?", username)
	user, err := repo.scanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *MySQLUserRepository) SearchByUsername(username string) ([]models.User, error) {
	rows, err := repo.DB.Query("SELECT username, password_hash, id, email, created_at, updated_at FROM users WHERE username LIKE CONCAT(?, '%')", username)
	if err != nil {
		return nil, err
	}

	users, err := repo.scanUserRows(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *MySQLUserRepository) Update(id int64, user *models.User) error {
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

func (repo *MySQLUserRepository) scanUserRows(rows *sql.Rows) ([]models.User, error) {
	users := []models.User{}
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Username, &u.Password, &u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (repo *MySQLUserRepository) scanUserRow(row *sql.Row) (*models.User, error) {
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
