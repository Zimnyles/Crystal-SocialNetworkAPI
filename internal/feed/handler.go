package feed

import (
	"net/http"
	"strconv"
	"time"
	"zimniyles/fibergo/pkg/generator"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"
	"zimniyles/fibergo/views/components"

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
	h.router.Get("/createpost", h.postCreate)
	h.router.Post("api/createpost", h.apiPostCreate)
}

func (h *FeedHandler) feed(c *fiber.Ctx) error {
	component := views.FeedPage()
	return tadapter.Render(c, component, 200)
}

func (h *FeedHandler) postCreate(c *fiber.Ctx) error {
	_, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	component := views.CreatePostPage()
	return tadapter.Render(c, component, 200)
}

func (h *FeedHandler) apiPostCreate(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login")
	if login == nil {
		c.Response().Header.Add("Hx-Redirect", "/login")
		return c.Redirect("/login", http.StatusOK)
	}

	authedLogin := sess.Get("login").(string)

	content := c.FormValue("content")
	h.customLogger.Info().Msg(content)

	image, err := c.FormFile("image")

	if err != nil || image == nil {
		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)
	}

	uniqueFilenameCode := generator.GenerateFilename()
	uniqueFilename := "postimage_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".jpg"
	filepath := "static/images/postimages/" + uniqueFilename

	if err := c.SaveFile(image, filepath); err != nil {
		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)
	}

	err = h.repository.AddNewFeedPost(authedLogin, content, filepath)
		if err != nil {
			h.customLogger.Error().Msg(err.Error())
			component := components.Notification("Произошла ошибка на сервере, попробуйте повторить попытку позже", components.NotificationFail)
			return tadapter.Render(c, component, http.StatusBadRequest)
		}
	
	msg := "Всё получилось! Пост можно увидеть на странице Новости или в своем профиле."
	component := components.Notification(msg, components.NotificationSuccess)
	return tadapter.Render(c, component, http.StatusOK)
}
