package people

import (
	"fmt"
	"math"
	"net/http"
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"
	"zimniyles/fibergo/views/widgets"

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
	h.router.Post("api/findpeople", middleware.AuthRequired(h.store), h.apiFindPeople)
	h.router.Get("api/findpeople", middleware.AuthRequired(h.store), h.apiFindPeople)
	profileGroup.Get("/", middleware.AuthRequired(h.store), h.people)
}

func (h *PeopleHandler) apiFindPeople(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	fmt.Println(login)

	content := c.FormValue("content")
	h.customLogger.Info().Msg(content)

	PAGE_ITEMS := 100
	page := c.QueryInt("page", 1)
	count := h.repository.CountAll()

	users, err := h.repository.GetAll(PAGE_ITEMS, (page-1)*PAGE_ITEMS, content)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	component := widgets.PeopleList(users, int(math.Ceil(float64(count/PAGE_ITEMS))), page, "/api/findpeople?page=%d", login)
	return tadapter.Render(c, component, http.StatusOK)
}

func (h *PeopleHandler) people(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	fmt.Println(login)

	PAGE_ITEMS := 10
	page := c.QueryInt("page", 1)
	count := h.repository.CountAll()

	users, err := h.repository.GetAll(PAGE_ITEMS, (page-1)*PAGE_ITEMS, "")
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	component := views.PeoplePage(users, int(math.Ceil(float64(count/PAGE_ITEMS))), page, "/people?page=%d", login)
	return tadapter.Render(c, component, http.StatusOK)

}
