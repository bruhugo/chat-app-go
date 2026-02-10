package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/grongoglongo/chatter-go/internal/exceptions"
	"github.com/grongoglongo/chatter-go/internal/messenger"
	"github.com/grongoglongo/chatter-go/internal/repositories"
	"github.com/grongoglongo/chatter-go/internal/utils"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func WebSocketHandler(h *messenger.ConnectionHub, repo repositories.ChatRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Client %v failed to upgrade to websocket: %s", c.Value("userId"), err.Error())
			c.Error(exceptions.NewHttpError("Failed to upgrade to websocket.", http.StatusBadRequest))
			return
		}

		userId, ok := utils.ConvertAnyToInt64(c.Value("userId"))
		if !ok {
			c.Error(exceptions.BadRequestError)
			return
		}

		log.Printf("Websocket connection opened for user %d", userId)

		chats, err := repo.FindByUser(userId)
		if err != nil {
			log.Print("Failed to convert user id to int64.")
			c.Error(exceptions.BadRequestError)
			return
		}

		chatIds := []int64{}
		for _, chat := range chats {
			chatIds = append(chatIds, chat.ID)
		}

		//	write
		go func() {
			cha := make(chan messenger.Event)
			connectionId := h.Subscribe(cha, userId, chatIds)
			for {
				select {
				case event := <-cha:
					err = conn.WriteJSON(event)
					if err != nil {
						conn.Close()
						h.Unsubscribe(connectionId)
						log.Printf("Failed to sent json to user %d via websocket. Closing connection", userId)
						return
					}
				case <-c.Done():
					conn.Close()
					h.Unsubscribe(connectionId)
					log.Printf("Websocket connection close for user %d", userId)
					return
				}
			}
		}()
	}
}
