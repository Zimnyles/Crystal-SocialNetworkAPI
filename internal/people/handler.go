package people

import (
	"fmt"
	"math"
	"net/http"
	"zimniyles/fibergo/internal/models"
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"
	"zimniyles/fibergo/views/components"
	"zimniyles/fibergo/views/widgets"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type PeopleHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   PeopleRepo
	store        *session.Store
}

type PeopleRepo interface {
	GetIDfromLogin(login string) (int, error)
	IsLoginExists(login string, logger *zerolog.Logger) (bool, error)
	AddFriend(userID int, friendID int) (*friendshipStatus, error)
	GetAll(userID, limit, offset int, search string) ([]models.PeopleProfileCredentials, error)
	CountAll() int
	CountNonFriends(userID int) (int, error)
}

func NewPeopleHandler(router fiber.Router, customLogger *zerolog.Logger, feedRepository PeopleRepo, store *session.Store) {
	h := &PeopleHandler{
		router:       router,
		customLogger: customLogger,
		repository:   feedRepository,
		store:        store,
	}
	peopleGroup := h.router.Group("/people")
	h.router.Post("api/findpeople", middleware.AuthRequired(h.store), h.apiFindPeople)
	h.router.Post("api/addfriend/:username", middleware.AuthRequired(h.store), h.apiAddFriend)
	h.router.Get("api/findpeople", middleware.AuthRequired(h.store), h.apiFindPeople)
	peopleGroup.Get("/", middleware.AuthRequired(h.store), h.people)
}

func (h *PeopleHandler) apiAddFriend(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	userLogin := sess.Get("login").(string)
	friendLogin := c.Params("username")

	isFriendExists, err := h.repository.IsLoginExists(friendLogin, h.customLogger)
	if err != nil {
		panic(err)
	}
	if !isFriendExists {
		c.Response().Header.Add("Hx-Redirect", "/404")
		return c.Redirect("/404", http.StatusOK)
	}

	userID, _ := h.repository.GetIDfromLogin(userLogin)
	friendID, _ := h.repository.GetIDfromLogin(friendLogin)
	if userID == 0 || friendID == 0 {
		c.Response().Header.Add("Hx-Redirect", "/404")
		return c.Redirect("/404", http.StatusOK)
	}

	friendshipStatus, err := h.repository.AddFriend(userID, friendID)
	h.customLogger.Info().Int("id", friendshipStatus.Id).Str("status", friendshipStatus.Status).Time("created_at", friendshipStatus.CreatedAt).Int("originator_id", friendshipStatus.OriginatorID).Int("recipient_id", friendshipStatus.RecipientID).Msg("New friend request")

	component := components.AddFriendResponse()
	return tadapter.Render(c, component, http.StatusOK)

}

func (h *PeopleHandler) apiFindPeople(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	userID, _ := h.repository.GetIDfromLogin(login)

	content := c.FormValue("content")
	h.customLogger.Info().Msg(content)

	PAGE_ITEMS := 100
	page := c.QueryInt("page", 1)
	count := h.repository.CountAll()

	users, err := h.repository.GetAll(userID, PAGE_ITEMS, (page-1)*PAGE_ITEMS, content)
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
	userID, _ := h.repository.GetIDfromLogin(login)

	PAGE_ITEMS := 10
	page := c.QueryInt("page", 1)

	count, err := h.repository.CountNonFriends(userID)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	users, err := h.repository.GetAll(userID, PAGE_ITEMS, (page-1)*PAGE_ITEMS, "")
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	totalPages := int(math.Ceil(float64(count) / float64(PAGE_ITEMS)))

	if page > totalPages && totalPages > 0 {
		return c.Redirect(fmt.Sprintf("/people?page=%d", totalPages), fiber.StatusSeeOther)
	}

	component := views.PeoplePage(users, totalPages, page, "/people?page=%d", login)
	return tadapter.Render(c, component, http.StatusOK)
}
