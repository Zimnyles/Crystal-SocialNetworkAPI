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
	"zimniyles/fibergo/views/widgets"

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
	h.router.Get("api/module-content", middleware.AuthRequired(h.store), h.apiGetModule)
	//Post
	h.router.Post("api/upload-avatar", middleware.AuthRequired(h.store), h.apiUploadAvatar)
}

func (h *ProfileHandler) apiGetModule(c *fiber.Ctx) error {
	chosenModule := c.Query("module")
	//?module=friend ?module=photos ?module=groups
	login := c.Query("login")
	userID, err := h.repository.GetIDfromLogin(login)
	if err != nil {
		panic(err)
	}

	switch chosenModule {
    case "photo":
        photos, err := h.repository.GetUserPhotos(userID, h.customLogger)
		if err != nil {
			h.customLogger.Error().Err(err).Msg("cannot get userID from login(profileHandler/getModule)")
		}
		component := widgets.UserPhotosList(photos, login)
		return tadapter.Render(c, component, http.StatusOK)
		
    case "friends":
        
    
    case "groups":
        
    
    default:
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid module parameter",
            "valid_modules": []string{"photos", "friends", "groups"},
        })
    }

	return nil
	
}

func (h *ProfileHandler) apiUploadAvatar(c *fiber.Ctx) error {
	
	file, err := c.FormFile("avatar")

	if err != nil || file == nil {
		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)
	}

	uniqueFilenameCode := generator.GenerateFilename()
	uniqueFilename := "avatar_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".jpg"
	uniqueFilenameToPhotoBin := "photo_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".jpg"
	filepath := "static/images/useravatars/" + uniqueFilename
	filepathToPhotoBin := "static/images/userphotos/" + uniqueFilenameToPhotoBin
	newAvatarPath := "static/images/useravatars/" + uniqueFilename

	if err := c.SaveFile(file, filepath); err != nil {
		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)
	}
	if err := c.SaveFile(file, filepathToPhotoBin); err != nil {
		component := views.ErrorPage(500, "в разработке1")
		return tadapter.Render(c, component, 500)
	}


	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login").(string)
	userID, err := h.repository.GetIDfromLogin(login)
	if err != nil {
		h.customLogger.Error().Err(err).Msg("Failed to get user id from login (profilehandler)")
	}

	err = h.repository.UpdateUserAvatar(userID, newAvatarPath)
	if err != nil {
		h.customLogger.Error().Err(err).Msg("Failed to update user pfp (profilehandler)")
	}
	err = h.repository.AddUserAvatarToPhotos(userID, filepathToPhotoBin, true)
	if err != nil {
		h.customLogger.Error().Err(err).Msg("Failed to add user pfp to global photos bin (profilehandler)")
	}

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

	//ProfileData------------------------------------------------------------------------------//
	UserData, err := h.repository.GetUserDataFromLogin(username, h.customLogger)
	if err != nil {
		component := views.ErrorPage(http.StatusInternalServerError, "Ошибка сервера, попробуйте позже")
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}
	userPosts, err := h.repository.GetAllUserPosts(username, PAGE_ITEMS, (page-1)*PAGE_ITEMS)
	if err != nil {
		component := views.ErrorPage(http.StatusInternalServerError, "Ошибка сервера, попробуйте позже")
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}
	//----------------------------------------------------------------------------------------//

	totalPages := int(math.Ceil(float64(count) / float64(PAGE_ITEMS)))

	if page > totalPages && totalPages > 0 {
		return c.Redirect(fmt.Sprintf("/profile/"+username+"?page=%d", totalPages), fiber.StatusSeeOther)
	}

	component := views.ProfilePage(*UserData, userPosts, totalPages, page, "/profile/"+username+"?page=%d")

	return tadapter.Render(c, component, http.StatusOK)

}
