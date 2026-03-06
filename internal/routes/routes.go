package routes

import (
	"database/sql"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	docs "github.com/grongoglongo/chatter-go/docs"
	"github.com/grongoglongo/chatter-go/internal/auth"
	"github.com/grongoglongo/chatter-go/internal/config"
	"github.com/grongoglongo/chatter-go/internal/messenger"
	"github.com/grongoglongo/chatter-go/internal/repositories"
	"github.com/grongoglongo/chatter-go/internal/routes/handlers"
	"github.com/grongoglongo/chatter-go/internal/routes/middleware"
	"github.com/grongoglongo/chatter-go/internal/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ApplyRoutes(router *gin.Engine, db *sql.DB) {

	repos := repositories.NewRepositories(db)

	connectionHub := messenger.NewConnectionHub()
	eventBus := messenger.NewEventBus(messenger.NewInMemoryMessenger(), connectionHub)

	userService := services.NewUserService(repos.UserRepository, auth.NewShaH256Service())
	messageService := services.NewMessageService(repos.MessageRepository, repos.ChatRepository, eventBus)
	chatService := services.NewChatService(repos.ChatRepository, repos.ChatMemberRepository, repos.UserRepository, eventBus)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.EnvConfig.FrontendHost},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Set-Cookie", "Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(middleware.ErrorHandler())
	v1.Use(gin.Logger())

	// USERS
	{
		users := v1.Group("/users")
		users.POST("", handlers.PostUserHandler(userService))
		users.POST(("/login"), handlers.LoginUserHandler(userService))
		users.GET("/logout", handlers.LogoutUserHandler())

		users.Use(middleware.AuthMiddleware())
		users.GET("/:id", handlers.GetUserHandler(userService))
		users.GET("/me", handlers.GetMeHandler(userService))
	}

	// CHATS
	{
		chats := v1.Group("/chats")
		chats.Use(middleware.AuthMiddleware())

		chats.GET("/:chatId/messages", handlers.GetMessagesByChatIdHandler(messageService))
		chats.POST("/:chatId/messages", handlers.CreateMessageHandler(messageService))

		chats.POST("", handlers.CreateChatHandler(chatService))
		chats.DELETE("/:chatId", handlers.DeleteChatHandler(chatService))
		chats.PUT("/:chatId", handlers.UpdateChatHandler(chatService))
		chats.POST("/:chatId/members", handlers.AddChatMemberHandler(chatService))
		chats.PUT("/:chatId/members", handlers.UpdateChatMemberHandler(chatService))
		chats.DELETE("/:chatId/members", handlers.DeleteChatMemberHandler(chatService))
		chats.GET("", handlers.GetChatsByUserIdHandler(chatService))
		chats.PUT("/:chatId/typing", handlers.TypingHandler(repos.UserRepository, repos.ChatRepository, eventBus))
	}

	// MESSAGES
	{
		messages := v1.Group("/messages")
		messages.Use(middleware.AuthMiddleware())
		messages.DELETE("/:messageId", handlers.DeleteMessageHandler(messageService))
	}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1.Use(middleware.AuthMiddleware())
	v1.GET("/websocket", handlers.WebSocketHandler(connectionHub, repos.ChatRepository))
}
