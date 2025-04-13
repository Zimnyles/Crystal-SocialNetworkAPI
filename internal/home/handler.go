package home

import (
	"math"
	"net/http"
	"zimniyles/fibergo/internal/post"
	"zimniyles/fibergo/pkg/tadapter"
	"zimniyles/fibergo/pkg/validator"
	"zimniyles/fibergo/views"
	"zimniyles/fibergo/views/components"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type HomeHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *post.PostRepository
	repositoryH  *UsersRepository
	store        *session.Store
}

type User struct {
	Id   int
	Name string
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *post.PostRepository, store *session.Store, repositoryH *UsersRepository) {
	h := &HomeHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		repositoryH:  repositoryH,
		store:        store,
	}
	//GET
	h.router.Get("/", h.home)
	h.router.Get("login", h.login)
	h.router.Get("register", h.register)
	h.router.Get("/404", h.error)
	//POST
	h.router.Post("api/login", h.apiLogin)
	h.router.Post("api/registration", h.apiRegistration)
}

func (h *HomeHandler) apiRegistration(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	form := UserCreateForm{
		Login:    c.FormValue("login"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	errors := validate.Validate(
		&validators.EmailIsPresent{Name: "Email", Field: form.Email, Message: "Email не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Password", Field: form.Password, Message: "Пароль не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Login", Field: form.Login, Message: "Логин не задан или задан неверно"},
	)

	var component templ.Component
	if len(errors.Error()) > 0 {
		component = components.Notification(validator.FormatErrors(errors), components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	err = h.repositoryH.addUser(form, h.customLogger)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		component = components.Notification("Произошла ошибка на сервере, попробуйте повторить попытку позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	msg := "Аккаунт зарегестрирован"
	
	sess.Set("login", form.Login)
	if err := sess.Save(); err != nil {
		panic(err)
	}

	component = components.Notification(msg, components.NotificationSuccess)
	return tadapter.Render(c, component, http.StatusOK)

}

func (h *HomeHandler) register(c *fiber.Ctx) error {
	component := views.Register()

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	userLogin := ""
	if login, ok := sess.Get("login").(string); ok {
		userLogin = login
		h.customLogger.Info().Msg("1")
	}

	c.Locals("login", userLogin)

	return tadapter.Render(c, component, http.StatusOK)

}

func (h *HomeHandler) apiLogin(c *fiber.Ctx) error {
	form := LoginForm{
		Login:    c.FormValue("login"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	if form.Login == "1" && form.Email == "2" && form.Password == "3" {
		sess, err := h.store.Get(c)
		if err != nil {
			panic(err)
		}
		sess.Set("login", form.Login)
		if err := sess.Save(); err != nil {
			panic(err)
		}

		c.Response().Header.Add("Hx-Redirect", "/")
		return c.Redirect("/", http.StatusOK)
	}

	component := components.Notification("Неверный логин или пароль", components.NotificationFail)
	return tadapter.Render(c, component, http.StatusBadRequest)

}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	PAGE_ITEMS := 2
	page := c.QueryInt("page", 1)
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	userLogin := ""
	if login, ok := sess.Get("login").(string); ok {
		h.customLogger.Info().Msg("1")
		userLogin = login
	}

	c.Locals("login", userLogin)

	count := h.repository.CountAll()
	posts, err := h.repository.GetAll(PAGE_ITEMS, (page-1)*PAGE_ITEMS)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	component := views.Main(posts, int(math.Ceil(float64(count/PAGE_ITEMS))), page)
	return tadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) login(c *fiber.Ctx) error {
	component := views.Login()

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	userLogin := ""
	if login, ok := sess.Get("login").(string); ok {
		userLogin = login
		h.customLogger.Info().Msg("1")
	}

	c.Locals("login", userLogin)

	return tadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) error(c *fiber.Ctx) error {

	h.customLogger.Info().
		Bool("isAdmin", true).
		Str("email", "a@a.ru").
		Int("id", 10).
		Msg("инфо")
	return c.SendString("error")
}
