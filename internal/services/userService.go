package services

import (
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
)

type UserService struct {
	repository  repositories.UserRepository
	hashService HashService
}

func NewUserService(
	repository repositories.UserRepository,
	hashService HashService,
) *UserService {
	return &UserService{
		repository:  repository,
		hashService: hashService,
	}
}

func (us *UserService) CreateUser(createUserDto dto.CreateUserDto) (*dto.UserDto, error) {
	hashedPassword := us.hashService.Hash(createUserDto.Password)

	user := &models.User{
		Username: createUserDto.Username,
		Email:    createUserDto.Email,
		Password: hashedPassword,
	}

	err := us.repository.Create(user)
	if err != nil {
		if err == exceptions.ConflictSqlError {
			return nil, exceptions.ConflictSqlError
		}
		return nil, exceptions.InternalServerError
	}

	return user.ToDto(), nil
}

func (us *UserService) FindUserById(id int64) (*dto.UserDto, error) {
	user, err := us.repository.FindById(id)
	if err != nil {
		return nil, exceptions.InternalServerError
	}

	if user == nil {
		return nil, exceptions.NotFoundError
	}

	return user.ToDto(), nil
}
