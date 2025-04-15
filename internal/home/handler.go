package home

import (
	"net/http"
	"strings"
	"zimniyles/fibergo/config"
	"zimniyles/fibergo/internal/post"
	"zimniyles/fibergo/pkg/jwt"
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
	"golang.org/x/crypto/bcrypt"
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

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *post.PostRepository, store *session.Store, repositoryH *UsersRepository, authConnfig *config.AuthConfig) {
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
	h.router.Get("api/logout", h.apiLogout)
	//POST
	h.router.Post("api/login", h.apiLogin)
	h.router.Post("api/registration", h.apiRegistration)
}

func (h *HomeHandler) apiLogout(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Delete("login")

	if err := sess.Save(); err != nil {
		panic(err)
	}

	c.Response().Header.Add("Hx-Redirect", "/")
	return c.Redirect("/", http.StatusOK)

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

	msg, err := h.repositoryH.addUser(form, h.customLogger)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		component = components.Notification(msg, components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	sess.Set("login", strings.ToLower(form.Login))
	if err := sess.Save(); err != nil {
		panic(err)
	}

	c.Response().Header.Add("Hx-Redirect", "/")
	return c.Redirect("/", http.StatusOK)

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

	emailIsExists, _ := h.repositoryH.IsEmailExistsForLogin(form, h.customLogger)

	if !emailIsExists {
		component := components.Notification("Пользователся с такой почтой не существует", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	UserCredentials, _ := h.repositoryH.GetPasswordByEmail(form, h.customLogger)
	if UserCredentials == nil {
		h.customLogger.Info().Msg("ошбика сервера 1")
		component := components.Notification("Ошибка сервера, попробуйте позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	err := bcrypt.CompareHashAndPassword([]byte(UserCredentials.PasswordHash), []byte(form.Password))
	if err != nil {
		h.customLogger.Info().Msg("ошбика сервера 2")
		component := components.Notification("Неверный пароль", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	jwt := jwt.NewJWT(config.NewAuthConfig().Secret)
	jwtToken, err := jwt.Create(form.Email)
	if err != nil {
		component := components.Notification("jwt токен не был сгенерирован, ошибка на сервере", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	c.SendString(jwtToken)

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Set("login", strings.ToLower(form.Login))
	if err := sess.Save(); err != nil {
		panic(err)
	}

	c.Response().Header.Add("Hx-Redirect", "/")
	return c.Redirect("/", http.StatusOK)

}

func (h *HomeHandler) home(c *fiber.Ctx) error {

	return c.Redirect("/feed", 302)

}

func (h *HomeHandler) login(c *fiber.Ctx) error {
	component := views.Login()

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
