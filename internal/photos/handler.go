package photos

import (
	"zimniyles/fibergo/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type PhotosHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *PhotosRepository
	store        *session.Store
}

func NewPhotosHandler(router fiber.Router, customLogger *zerolog.Logger, repository *PhotosRepository, store *session.Store) {
	h := &PhotosHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}
	//Post
	h.router.Post("api/upload-photo", middleware.AuthRequired(h.store), h.apiUploadPhoto)
}

func (h *PhotosHandler) apiUploadPhoto(c*fiber.Ctx) error {
	return nil
}
