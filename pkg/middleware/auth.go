package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		userLogin := ""
		if login, ok := sess.Get("login").(string); ok {
			userLogin = login
		}

		c.Locals("login", userLogin)
		return c.Next()
	}
}
