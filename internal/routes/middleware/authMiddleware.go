package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/auth"
	"github.com/grongoglongo/chatter-go/internal/config"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
)

func AuthMiddleware() gin.HandlerFunc {
	jwtHandler := auth.NewJwtHandler(config.EnvConfig.JwtSecret)

	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}

		cookie, err := ctx.Request.Cookie("X-Auth-Header")
		if err != nil || cookie == nil {
			ctx.Error(exceptions.
				NewHttpError("Authorization header not found (X-Auth-Header)", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		claims, err := jwtHandler.DecryptJwt(cookie.Value)
		if err != nil {
			ctx.Error(exceptions.NewHttpError("Invalid token", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		log.Printf("User %v authenticated", claims.Id)

		ctx.Set("userId", claims.Id)
		ctx.Set("userEmail", claims.Id)

		ctx.Next()
	}
}
