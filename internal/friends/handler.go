package friends

import (
	"fmt"
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

type FriendsHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   FriendsRepo
	store        *session.Store
}

type FriendsRepo interface {
	GetAllFriendRequests(userID int) ([]models.FriendRequestList, error)
	GetAcceptedFriends(userID int) ([]models.FriendList, error)
	GetIDfromLogin(login string) (int, error)
	AcceptFriendship(userId int, friendId int) bool
	DeclineFriendship(userId int, friendId int) bool
}

func NewFriendsHandler(router fiber.Router, customLogger *zerolog.Logger, friendsRepository FriendsRepo, store *session.Store) {
	h := &FriendsHandler{
		router:       router,
		customLogger: customLogger,
		repository:   friendsRepository,
		store:        store,
	}

	friendsGroup := h.router.Group("/friends")
	friendsGroup.Get("/", middleware.AuthRequired(h.store), h.friends)

	h.router.Post("api/acceptfriendship/:username", middleware.AuthRequired(h.store), h.apiAcceptFriendship)
	h.router.Post("api/declinefriendship/:username", middleware.AuthRequired(h.store), h.apiDeclineFriendship)
}

func (h *FriendsHandler) apiAcceptFriendship(c *fiber.Ctx) error {
	friendLogin := c.Params("username")

	userID := getUserID(h.store, c, h.repository)
	friendID, err := h.repository.GetIDfromLogin(friendLogin)
	if err != nil {
		panic(err)
	}

	isAccepted := h.repository.AcceptFriendship(userID, friendID)
	if !isAccepted {
		component := components.Notification("Произошла ошибка на сервере, попробуйте повторить попытку позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}


		friends, err := h.repository.GetAcceptedFriends(userID)
	if err != nil {
		fmt.Println(err)
	}

	component := widgets.FriendList(friends)
	return tadapter.Render(c, component, 200)
}

func (h *FriendsHandler) apiDeclineFriendship(c *fiber.Ctx) error {

	friendLogin := c.Params("username")

	userID := getUserID(h.store, c, h.repository)
	friendID, err := h.repository.GetIDfromLogin(friendLogin)
	if err != nil {
		panic(err)
	}

	isAccepted := h.repository.DeclineFriendship(userID, friendID)
	if !isAccepted {
		component := components.Notification("Произошла ошибка на сервере, попробуйте повторить попытку позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	return nil

}

func (h *FriendsHandler) friends(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	userID, _ := h.repository.GetIDfromLogin(login)

	friends, err := h.repository.GetAcceptedFriends(userID)
	if err != nil {
		fmt.Println(err)
	}
	requests, err := h.repository.GetAllFriendRequests(userID)
	if err != nil {
		fmt.Println(err)
	}

	friendsData := models.FriendPageCredentials{
		Friends:        friends,
		FriendRequests: requests,
	}

	component := views.FriendsPage(friendsData)
	return tadapter.Render(c, component, 200)
}

func getUserID(store *session.Store, c *fiber.Ctx, friendsRepository FriendsRepo) int {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	login := sess.Get("login").(string)
	userID, _ := friendsRepository.GetIDfromLogin(login)

	return userID
}
