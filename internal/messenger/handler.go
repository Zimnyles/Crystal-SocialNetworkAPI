package messenger

import (
	"net/http"
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"
	"zimniyles/fibergo/views/components"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type MessengerHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *MessengerRepository
	store        *session.Store
}

// var (
// 	chat    = models.NewChat()
// 	mutex   = sync.Mutex{}
// 	clients = make(map[*websocket.Conn]bool)
// )

func NewMessengerHandler(router fiber.Router, customLogger *zerolog.Logger, messengerRepository *MessengerRepository, store *session.Store) {
	h := &MessengerHandler{
		router:       router,
		customLogger: customLogger,
		repository:   messengerRepository,
		store:        store,
	}
	messagesGroup := h.router.Group("/messages")
	messagesGroup.Get("/", h.messages)
	h.router.Get("/messages/:username", middleware.AuthRequired(h.store), h.dialog)

	// messagesGroup.Get("/ws", func(c *fiber.Ctx) error {
	// 	if websocket.IsWebSocketUpgrade(c) {
	// 		return c.Next()
	// 	}
	// 	return c.SendStatus(fiber.StatusUpgradeRequired)
	// }, websocket.New(h.WebSocketHandler))

}

func (h *MessengerHandler) messages(c *fiber.Ctx) error {

	component := views.MessagesPage()
	return tadapter.Render(c, component, 200)

}

func (h *MessengerHandler) dialog(c *fiber.Ctx) error {

	component := components.ErrorComponent(200, "ok")
	return tadapter.Render(c, component, http.StatusOK)

}
