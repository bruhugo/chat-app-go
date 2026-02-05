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
	messageService := services.NewMessageService(repos.MessageRepository, repos.ChatRepository)
	chatService := services.NewChatService(repos.ChatRepository, repos.ChatMemberRepository)

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(middleware.ErrorHandler())
	v1.Use(gin.Logger())

	// USERS
	{
		users := v1.Group("/users")
		users.POST("/", handlers.PostUserHandler(userService))
		users.POST(("/login"), handlers.LoginUserHandler(userService))
		users.GET("/logout", handlers.LogoutUserHandler())

		users.Use(middleware.AuthMiddleware())
		users.GET("/:id", handlers.GetUserHandler(userService))
	}

	// CHATS
	{
		chats := v1.Group("/chats")
		chats.Use(middleware.AuthMiddleware())
		chats.GET("/:chatId/messages", handlers.GetMessagesByChatIdHandler(messageService))
		chats.POST("/", handlers.CreateChatHandler(chatService))
	}

	// MESSAGES
	{
		messages := v1.Group("/messages")
		messages.Use(middleware.AuthMiddleware())
		messages.POST("/", handlers.CreateMessageHandler(messageService))
	}
}
