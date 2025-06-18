package profile

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"zimniyles/fibergo/pkg/generator"
	"zimniyles/fibergo/pkg/middleware"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type ProfileHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *ProfileRepository
	store        *session.Store
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *ProfileRepository, store *session.Store) {
	h := &ProfileHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}

	profileGroup := h.router.Group("/profile")
	//Get
	profileGroup.Get("/:username", h.profile)
	//Post
	h.router.Post("api/upload-avatar", middleware.AuthRequired(h.store), h.apiUploadAvatar)
}

func (h *ProfileHandler) apiUploadAvatar(c *fiber.Ctx) error {

	file, err := c.FormFile("avatar")

	if err != nil || file == nil {
		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)
	}

	uniqueFilenameCode := generator.GenerateFilename()
	uniqueFilename := "avatar_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".jpg"
	filepath := "static/images/useravatars/" + uniqueFilename
	newAvatarPath := "static/images/useravatars/" + uniqueFilename

	if err := c.SaveFile(file, filepath); err != nil {

		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)

	}

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	if login == "" {

	}

	err = h.repository.UpdateUserAvatar(login, newAvatarPath)

	if err != nil {
		redirectLink := "/profile/" + login
		return c.Redirect(redirectLink, 302)
	}

	redirectLink := "/profile/" + login
	return c.Redirect(redirectLink, 302)

}

func (h *ProfileHandler) profile(c *fiber.Ctx) error {
	username := c.Params("username")

	PAGE_ITEMS := 10
	page := c.QueryInt("page", 1)

	count, err := h.repository.CountUserPosts(username)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	isLoginExists, _ := h.repository.IsLoginExistsForString(username, h.customLogger)
	if !isLoginExists {
		component := views.ErrorPage(http.StatusNotFound, "Страница не найдена")
		return tadapter.Render(c, component, http.StatusNotFound)
	}

	UserData, err := h.repository.GetUserDataFromLogin(username, h.customLogger)
	if err != nil {
		component := views.ErrorPage(http.StatusInternalServerError, "Ошибка сервера, попробуйте позже")
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}

	userPosts, err := h.repository.GetAllUserPosts(username, PAGE_ITEMS, (page-1)*PAGE_ITEMS)

	totalPages := int(math.Ceil(float64(count) / float64(PAGE_ITEMS)))

	if page > totalPages && totalPages > 0 {
		return c.Redirect(fmt.Sprintf("/profile/" + username + "?page=%d", totalPages), fiber.StatusSeeOther)
	}

	component := views.ProfilePage(*UserData, userPosts, totalPages, page, "/profile/" + username + "?page=%d")

	return tadapter.Render(c, component, http.StatusOK)

}
