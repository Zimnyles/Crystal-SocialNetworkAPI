package messenger

import (
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

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

func NewMessengerHandler(router fiber.Router, customLogger *zerolog.Logger, messengerRepository *MessengerRepository, store *session.Store) {
	h := &MessengerHandler{
		router:       router,
		customLogger: customLogger,
		repository:   messengerRepository,
		store:        store,
	}
	profileGroup := h.router.Group("/messages")
	profileGroup.Get("/", h.messages)
}

func (h *MessengerHandler) messages(c *fiber.Ctx) error {
	component := views.MessagesPage()
	return tadapter.Render(c, component, 200)

}
