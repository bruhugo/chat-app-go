package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
	"github.com/grongoglongo/chatter-go/internal/services"
)

func GetUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Param("id")
		id, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			c.Error(exceptions.NewHttpError(err, http.StatusBadRequest)).SetType(gin.ErrorTypePublic)
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

func PostUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createUserDto dto.CreateUserDto

		if err := ctx.ShouldBindBodyWithJSON(&createUserDto); err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusBadRequest)).SetType(gin.ErrorTypePublic)
			return
		}

		userDto, err := userService.CreateUser(createUserDto)
		if err != nil {
			ctx.Error(err)
			return
		}

		jwtHandler := services.NewJwtHandler()
		jwt, err := jwtHandler.CreateJwt(userDto)
		if err != nil {
			ctx.Error(exceptions.InternalServerError)
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
