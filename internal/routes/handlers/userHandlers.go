package handlers

import (
	"errors"
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

		if user == nil {
			c.Error(exceptions.NotFoundError).SetType(gin.ErrorTypePublic)
			return
		}

		c.JSON(http.StatusOK, dto.UserDto{Username: user.Username, ID: user.ID, Email: user.Email})
	}
}

func PostUserHandler(userRepo *repositories.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createUserDto dto.CreateUserDto
		jwtHandler := services.NewJwtHandler()

		if err := ctx.ShouldBindBodyWithJSON(&createUserDto); err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusBadRequest)).SetType(gin.ErrorTypePublic)
			return
		}

		hashService := services.NewHashService()
		hashedPassword := hashService.Hash(createUserDto.Password)

		user := &models.User{
			Username: createUserDto.Username,
			Email:    createUserDto.Email,
			Password: hashedPassword,
		}

		err := userRepo.Create(user)
		if err != nil {
			if err == exceptions.ConflictSqlError {
				ctx.Error(exceptions.NewHttpError(errors.New("Conflict creating user."), http.StatusConflict)).SetType(gin.ErrorTypePublic)
				return
			}
			ctx.Error(exceptions.NewHttpError(errors.New("Error creating user. Try again later."), http.StatusInternalServerError)).SetType(gin.ErrorTypePublic)
			return
		}

		userDto := &dto.UserDto{
			Username: createUserDto.Username,
			Email:    createUserDto.Email,
			ID:       user.ID,
		}

		jwt, err := jwtHandler.CreateJwt(userDto)
		if err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusInternalServerError)).SetType(gin.ErrorTypePublic)
			return
		}

		addAuthCookieToRequest(jwt, ctx.Writer)

		ctx.JSON(http.StatusOK, userDto)
	}

}

func addAuthCookieToRequest(jwt string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "X-Auth-Header",
		Value:    jwt,
		MaxAge:   60 * 60 * 24 * 60, // 60 days
		Secure:   false,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}
