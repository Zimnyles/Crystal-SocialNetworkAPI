package profile

import (
	"net/http"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type ProfileHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *ProfileRepository
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *ProfileRepository) {
	h := &ProfileHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
	}
	profileGroup := h.router.Group("/profile")
	profileGroup.Get("/:username", h.profile)
}

func (h *ProfileHandler) profile(c *fiber.Ctx) error {

	username := c.Params("username")
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

	component := views.ProfilePage(*UserData)

	return tadapter.Render(c, component, http.StatusOK)

}
