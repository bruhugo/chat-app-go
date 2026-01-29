package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/services"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("X-Auth-Header")
		if err != nil {
			ctx.Error(exceptions.NewHttpErrorWithMessage(err, http.StatusUnauthorized, "Authorization header not found (X-Auth-Header)"))
			return
		}

		token := cookie.Value[7:]

		claims, err := services.DecryptJwt(token)
		if err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusUnauthorized))
			return
		}

		log.Printf("User %v authenticated", claims.Id)

		ctx.Request.Header.Set("X-User-Id", claims.Id)
		ctx.Request.Header.Set("X-User-Email", claims.Id)

		ctx.Next()
	}
}
