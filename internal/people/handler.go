package people

import (
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type PeopleHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *PeopleRepository
	store        *session.Store
}

func NewPeopleHandler(router fiber.Router, customLogger *zerolog.Logger, feedRepository *PeopleRepository, store *session.Store) {
	h := &PeopleHandler{
		router:       router,
		customLogger: customLogger,
		repository:   feedRepository,
		store:        store,
	}
	profileGroup := h.router.Group("/people")
	profileGroup.Get("/", middleware.AuthRequired(h.store), h.people)
}

func (h *PeopleHandler) people(c *fiber.Ctx) error {
	component := views.PeoplePage()
	return tadapter.Render(c, component, 200)
}
