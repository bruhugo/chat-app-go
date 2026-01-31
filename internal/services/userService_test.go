package services

import (
	"testing"

	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/stretchr/testify/require"
)

type MockUserRepo struct {
	CreateFunc   func(user *models.User) error
	FindByIdFunc func(id int64) (*models.User, error)
}

func (m *MockUserRepo) Create(user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}

	user.ID = 1
	return nil
}

func (m *MockUserRepo) FindById(id int64) (*models.User, error) {
	if m.FindByIdFunc != nil {
		return m.FindByIdFunc(id)
	}

	return &models.User{ID: id, Username: "username", Password: "password", Email: "email"}, nil
}

func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) { return nil, nil }
func (m *MockUserRepo) Update(id int64, user *models.User) error       { return nil }

func TestCreateUser_Success(t *testing.T) {
	tests := []struct {
		name          string
		createUserDto dto.CreateUserDto
	}{
		{"Case_1", dto.CreateUserDto{Username: "username", Password: "password", Email: "email@com"}},
		{"Case_2", dto.CreateUserDto{Username: "123", Password: "1221", Email: "email@com"}},
		{"Case_3", dto.CreateUserDto{Username: "33", Password: "111111111111", Email: "dwa@com"}},
		{"Case_4", dto.CreateUserDto{Username: "dwwwa", Password: "passwdwaword", Email: "email@dwaww"}},
	}

	userService := NewUserService(&MockUserRepo{}, NewShaH256Service())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userDto, err := userService.CreateUser(tt.createUserDto)

			require.NoError(t, err)
			require.Equal(t, tt.createUserDto.Email, userDto.Email)
			require.Equal(t, tt.createUserDto.Username, userDto.Username)
			require.NotEmpty(t, userDto.ID)
		})
	}
}

func TestCreateUser_Conflict(t *testing.T) {
	repo := &MockUserRepo{
		CreateFunc: func(user *models.User) error { return exceptions.ConflictSqlError },
	}

	userService := NewUserService(repo, NewShaH256Service())
	createUserDto := dto.CreateUserDto{Username: "username", Password: "password", Email: "Email"}

	_, err := userService.CreateUser(createUserDto)

	require.Error(t, err, "Conflict error should have been thrown.")
	require.Equal(t, err, exceptions.ConflictSqlError)

}

func TestCreateUser_HashFunction(t *testing.T) {
	var savedUser *models.User

	repo := &MockUserRepo{
		CreateFunc: func(user *models.User) error {
			savedUser = user
			return nil
		},
	}

	hashService := NewShaH256Service()
	userService := NewUserService(repo, hashService)
	createUserDto := &dto.CreateUserDto{Username: "username", Password: "password", Email: "Email"}

	_, err := userService.CreateUser(*createUserDto)

	require.NoError(t, err)
	require.NotEqual(t, savedUser.Password, createUserDto.Password)

	verifyHashedPassword := hashService.Hash(createUserDto.Password)

	require.Equal(t, verifyHashedPassword, savedUser.Password)
}

func TestFindUser_Success(t *testing.T) {
	userService := NewUserService(&MockUserRepo{}, NewShaH256Service())

	var id int64 = 12
	userDto, err := userService.FindUserById(id)

	require.NoError(t, err)
	require.NotEmpty(t, userDto)
	require.Equal(t, id, userDto.ID)
}

func TestFindUser_NoFound(t *testing.T) {
	repo := &MockUserRepo{
		FindByIdFunc: func(id int64) (*models.User, error) {
			return nil, nil
		},
	}

	userService := NewUserService(repo, NewShaH256Service())
	userDto, err := userService.FindUserById(12)

	require.Empty(t, userDto)
	require.Error(t, err)
	require.Equal(t, err, exceptions.NotFoundError)
}
