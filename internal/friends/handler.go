package friends

import (
	"zimniyles/fibergo/internal/models"
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

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	userID, _ := h.repository.GetIDfromLogin(login)

	friends, _ := h.repository.GetAcceptedFriends(userID)
	requests, _ := h.repository.GetAllFriendRequests(userID)

	friendsData := models.FriendPageCredentials{
		Friends:        friends,
		FriendRequests: requests,
	}

	component := views.FriendsPage(friendsData)
	return tadapter.Render(c, component, 200)
}
