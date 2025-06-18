package post

import (
	"net/http"
	"strconv"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/pkg/validator"
	"zimniyles/fibergo/views/components"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

//NOT USED 

type PostHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *PostRepository
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *PostRepository) {
	h := &PostHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
	}
	postGroup := h.router.Group("/post")
	postGroup.Post("/", h.createPost)
}


func (h *PostHandler) createPost(c *fiber.Ctx) error {
	form := PostCreateForm{
		Name:        c.FormValue("name"),
		Breed:       c.FormValue("breed"),
		Price:       c.FormValue("price"),
		Location:    c.FormValue("location"),
		Description: c.FormValue("description"),
		Email:       c.FormValue("email"),
	}
	errors := validate.Validate(
		&validators.StringIsPresent{Name: "Name", Field: form.Name, Message: "Кличка не указана или указана неверно"},
		&validators.EmailIsPresent{Name: "Email", Field: form.Email, Message: "Email не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Location", Field: form.Location, Message: "Город не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Breed", Field: form.Breed, Message: "Порода не задана или задана неверно"},
		&validators.StringIsPresent{Name: "Price", Field: form.Price, Message: "Цена не задана или задана неверно"},
		&validators.StringIsPresent{Name: "Description", Field: form.Description, Message: "Описание не задано или задано неверно"},
	)
	var component templ.Component
	if len(errors.Error()) > 0 {
		component = components.Notification(validator.FormatErrors(errors), components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}
	id, err := h.repository.addPost(form, h.customLogger)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		component = components.Notification("Произошла ошибка на сервере, попробуйте повторить попытку позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}
	msg := "Всё получилось! Номер вашего объявления: " + strconv.Itoa(id)
	component = components.Notification(msg, components.NotificationSuccess)
	return tadapter.Render(c, component, http.StatusOK)
}
