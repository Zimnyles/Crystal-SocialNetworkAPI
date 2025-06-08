package friends

import (
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type FriendsHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *FriendsRepository
	store        *session.Store
}

func NewFriendsHandler(router fiber.Router, customLogger *zerolog.Logger, feedRepository *FriendsRepository, store *session.Store) {
	h := &FriendsHandler{
		router:       router,
		customLogger: customLogger,
		repository:   feedRepository,
		store:        store,
	}
	profileGroup := h.router.Group("/friends")
	profileGroup.Get("/", middleware.AuthRequired(h.store), h.friends)
}

func (h *FriendsHandler) friends(c *fiber.Ctx) error {
	

	component := views.FriendsPage()
	return tadapter.Render(c, component, 200)
}
