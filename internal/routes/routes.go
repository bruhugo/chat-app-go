package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/repositories"
	"github.com/grongoglongo/chatter-go/internal/routes/handlers"
	"github.com/grongoglongo/chatter-go/internal/routes/middleware"
	"github.com/grongoglongo/chatter-go/internal/services"
)

func ApplyRoutes(router *gin.Engine, db *sql.DB) {

	repos := repositories.NewRepositories(db)
	userService := services.NewUserService(repos.UserRepository, services.NewShaH256Service())

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(middleware.ErrorHandler())
	v1.Use(gin.Logger())

	{
		users := v1.Group("/users")
		users.POST("/", handlers.PostUserHandler(userService))

		//	PROTECTED ROUTES
		users.Use(middleware.AuthMiddleware())
		users.GET("/:id", handlers.GetUserHandler(userService))
	}
}
