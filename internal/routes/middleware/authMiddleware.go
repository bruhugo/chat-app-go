package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/services"
)

func AuthMiddleware() gin.HandlerFunc {
	jwtHandler := services.NewJwtHandler()

	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("X-Auth-Header")
		if err != nil || cookie == nil {
			ctx.Error(exceptions.
				NewHttpErrorWithMessage(err, http.StatusUnauthorized, "Authorization header not found (X-Auth-Header)")).
				SetType(gin.ErrorTypePublic)
			ctx.Abort()
			return
		}

		claims, err := jwtHandler.DecryptJwt(cookie.Value)
		if err != nil {
			ctx.Error(exceptions.NewHttpError(err, http.StatusUnauthorized)).SetType(gin.ErrorTypePublic)
			ctx.Abort()
			return
		}

		log.Printf("User %v authenticated", claims.Id)

		ctx.Set("X-User-Id", claims.Id)
		ctx.Set("X-User-Email", claims.Id)

		ctx.Next()
	}
}
