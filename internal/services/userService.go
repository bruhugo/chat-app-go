package services

import (
	"fmt"
	"net/http"

	"github.com/grongoglongo/chatter-go/internal/auth"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
)

type UserService struct {
	repository  repositories.UserRepository
	hashService auth.HashService
}

func NewUserService(
	repository repositories.UserRepository,
	hashService auth.HashService,
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

func (us *UserService) LoginUser(loginDto *dto.LoginUserDto) (*dto.UserDto, error) {
	user, err := us.repository.FindByEmail(loginDto.Email)
	if err != nil {
		return nil, exceptions.InternalServerError
	}

	if user == nil {
		return nil, exceptions.NotFoundError
	}

	newHashedPassword := us.hashService.Hash(loginDto.Password)

	if newHashedPassword != user.Password {
		return nil, exceptions.NewHttpError("Invalid credentials", http.StatusUnauthorized)
	}

	return user.ToDto(), nil
}

func (us *UserService) SearchByUsername(username string) ([]models.User, error) {

	users, err := us.repository.SearchByUsername(username)
	if err != nil {
		fmt.Print(err.Error())
		return nil, exceptions.NewHttpError("Error fetching users", http.StatusInternalServerError)
	}

	return users, nil
}
