package feed

import (
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type FeedHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *FeedRepository
}

func NewFeedHandler(router fiber.Router, customLogger *zerolog.Logger, feedRepository *FeedRepository) {
	h := &FeedHandler{
		router:       router,
		customLogger: customLogger,
		repository:   feedRepository,
	}
	profileGroup := h.router.Group("/feed")
	profileGroup.Get("/", h.feed)
}

func (h *FeedHandler) feed(c *fiber.Ctx) error {
	component := views.FeedPage()
	return tadapter.Render(c, component, 200)
}
