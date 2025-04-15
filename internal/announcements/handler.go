package announcements

import (
	"math"
	"net/http"
	"zimniyles/fibergo/internal/post"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/views"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AnnouncementsHandler struct {
	router         fiber.Router
	customLogger   *zerolog.Logger
	repository     *AnnouncementsRepository
	postRepository *post.PostRepository
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *AnnouncementsRepository, postRepository *post.PostRepository) {
	h := &AnnouncementsHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		postRepository: postRepository,
	}
	profileGroup := h.router.Group("/announcements")
	profileGroup.Get("/", h.announcements)
}

func (h *AnnouncementsHandler) announcements(c *fiber.Ctx) error {
	PAGE_ITEMS := 2
	page := c.QueryInt("page", 1)

	count := h.postRepository.CountAll()
	posts, err := h.postRepository.GetAll(PAGE_ITEMS, (page-1)*PAGE_ITEMS)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	component := views.AnnoncementPage(posts, int(math.Ceil(float64(count/PAGE_ITEMS))), page)
	return tadapter.Render(c, component, http.StatusOK)
}
