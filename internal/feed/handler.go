package feed

import (
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type FeedHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *FeedRepository
	store        *session.Store
}

func NewFeedHandler(router fiber.Router, customLogger *zerolog.Logger, feedRepository *FeedRepository, store *session.Store) {
	h := &FeedHandler{
		router:       router,
		customLogger: customLogger,
		repository:   feedRepository,
		store:        store,
	}
	profileGroup := h.router.Group("/feed")
	profileGroup.Get("/", h.feed)
	h.router.Get("/createpost", h.postcreate)
}

func (h *FeedHandler) feed(c *fiber.Ctx) error {
	component := views.FeedPage()
	return tadapter.Render(c, component, 200)
}

func (h *FeedHandler) postcreate(c *fiber.Ctx) error {
	_, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	component := views.FeedPage()
	return tadapter.Render(c, component, 200)
}
