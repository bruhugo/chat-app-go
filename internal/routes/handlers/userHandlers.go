package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/repositories"
	"github.com/grongoglongo/chatter-go/internal/services"
)

func GetUserHandler(userRepo *repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Param("id")
		id, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			c.Error(exceptions.NewHttpError(err, http.StatusBadRequest)).SetType(gin.ErrorTypePublic)
			return
		}

		user, err := userRepo.FindById(id)
		if err != nil {
			c.Error(exceptions.NewHttpError(err, http.StatusNotFound)).SetType(gin.ErrorTypePublic)
			return
		}

		c.JSON(http.StatusOK, dto.UserDto{Username: user.Username, ID: user.ID, Email: user.Email})
	}
}

func PostUserHandler(userRepo *repositories.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createUserDto dto.CreateUserDto
		if err := ctx.ShouldBindBodyWithJSON(&createUserDto); err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusBadRequest)).SetType(gin.ErrorTypePublic)
			return
		}

		conflictUser, err := userRepo.FindByEmail(createUserDto.Email)
		if err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusInternalServerError)).SetType(gin.ErrorTypePublic)
			return
		}
		if conflictUser != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusConflict)).SetType(gin.ErrorTypePublic)
			return
		}

		hashedPassword := services.Hash(createUserDto.Password)

		user := &models.User{
			Username: createUserDto.Username,
			Email:    createUserDto.Email,
			Password: hashedPassword,
		}

		userRepo.Create(user)

		userDto := &dto.UserDto{
			Username: createUserDto.Username,
			Email:    createUserDto.Email,
			ID:       user.ID,
		}

		jwt, err := services.CreateJwt(userDto)
		if err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusInternalServerError)).SetType(gin.ErrorTypePublic)
			return
		}

		ctx.JSON(http.StatusOK, dto.AuthResponse{User: *userDto, Jwt: jwt})
	}
}
