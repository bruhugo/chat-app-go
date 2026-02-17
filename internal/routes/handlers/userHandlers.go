package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	auth "github.com/grongoglongo/chatter-go/internal/auth"
	"github.com/grongoglongo/chatter-go/internal/config"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/services"
	"github.com/grongoglongo/chatter-go/internal/utils"
)

const COOKIE_NAME = "X-Auth-Header"

// @Summary Get user by ID
// @Description Returns a single user by ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 404 {object} exceptions.HttpError
// @Router /users/{id} [get]
func GetUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Param("id")
		id, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			c.Error(exceptions.NewHttpError("Id must be a string", http.StatusBadRequest)).SetType(gin.ErrorTypePublic)
			return
		}

		userDto, err := userService.FindUserById(id)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, *userDto)
	}
}

// @Summary Create user
// @Description Creates a new user and returns the created user.
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.CreateUserDto true "User payload"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 409 {object} exceptions.HttpError
// @Router /users/ [post]
func PostUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createUserDto dto.CreateUserDto

		if err := ctx.ShouldBindBodyWithJSON(&createUserDto); err != nil {
			ctx.JSON(http.StatusBadGateway, utils.GetValidationError(err))
			return
		}

		userDto, err := userService.CreateUser(createUserDto)
		if err != nil {
			ctx.Error(err)
			return
		}

		jwtHandler := auth.NewJwtHandler(config.EnvConfig.JwtSecret)
		jwt, err := jwtHandler.CreateJwt(userDto)
		if err != nil {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		cookie := buildCookie(jwt)
		http.SetCookie(ctx.Writer, cookie)

		ctx.JSON(http.StatusOK, userDto)
	}
}

// @Summary Login user
// @Description Authenticates a user and returns the user data.
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.LoginUserDto true "Login payload"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} exceptions.HttpError
// @Failure 401 {object} exceptions.HttpError
// @Router /users/login [post]
func LoginUserHandler(userService *services.UserService) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var loginUserDto dto.LoginUserDto

		err := ctx.ShouldBindBodyWithJSON(&loginUserDto)
		if err != nil {
			ctx.Error(exceptions.BadRequestError)
			return
		}

		userDto, err := userService.LoginUser(&loginUserDto)

		if err != nil {
			ctx.Error(err)
			return
		}

		jwtHandler := auth.NewJwtHandler(config.EnvConfig.JwtSecret)
		jwt, err := jwtHandler.CreateJwt(userDto)
		if err != nil {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		cookie := buildCookie(jwt)
		http.SetCookie(ctx.Writer, cookie)

		ctx.JSON(http.StatusOK, userDto)
	}
}

// @Summary Logout user
// @Description Clears the auth cookie.
// @Tags users
// @Success 200 {string} string "ok"
// @Router /users/logout [get]
func LogoutUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie := buildCookie("")
		cookie.MaxAge = 0
		http.SetCookie(ctx.Writer, cookie)
	}
}

// @Summary Get user
// @Description Get current authenticated user
// @Tags users
// @Success 200 {string} string "ok"
// @Router /users/me [get]
func GetMeHandler(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.ConvertAnyToInt64(ctx.Value("userId"))
		if !ok {
			ctx.Error(exceptions.InternalServerError)
			return
		}

		userDto, err := userService.FindUserById(userId)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, userDto)
	}
}

func buildCookie(jwt string) *http.Cookie {
	return &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    jwt,
		MaxAge:   60 * 60 * 24 * 60, // 60 days
		Secure:   false,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}
