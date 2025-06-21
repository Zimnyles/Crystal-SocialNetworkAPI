package messenger

import (
	"net/http"
	"zimniyles/fibergo/internal/models"
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"
	"zimniyles/fibergo/views/components"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type MessengerHandler struct {
	router           fiber.Router
	customLogger     *zerolog.Logger
	repository       MessagesRepo
	globalRepository models.GlobalRepo
	store            *session.Store
}

type MessagesRepo interface {
	GetUserChats(userID int) ([]models.ChatPreview, error)
}

// var (
// 	chat    = models.NewChat()
// 	mutex   = sync.Mutex{}
// 	clients = make(map[*websocket.Conn]bool)
// )

func NewMessengerHandler(router fiber.Router, customLogger *zerolog.Logger, messengerRepository MessagesRepo, globalRepository models.GlobalRepo, store *session.Store) {
	h := &MessengerHandler{
		router:           router,
		customLogger:     customLogger,
		repository:       messengerRepository,
		globalRepository: globalRepository,
		store:            store,
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
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	userLogin := sess.Get("login").(string)
	userID, err := h.globalRepository.GetIDfromLogin(userLogin)
	if err != nil {
		h.customLogger.Error().Err(err).Msg("cannot get userID from login(messagesHandler)")
	}

	userChats, err := h.repository.GetUserChats(userID)
	if err != nil {
		h.customLogger.Error().Err(err).Msg("cannot get userChats (messagesHandler)")
	}

	component := views.MessagesPage(userChats)
	return tadapter.Render(c, component, 200)

}

func (h *MessengerHandler) dialog(c *fiber.Ctx) error {

	component := components.ErrorComponent(200, "ok")
	return tadapter.Render(c, component, http.StatusOK)

}
