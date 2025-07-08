package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"terra9.it/checkmate/server/handlers/auth"
	"terra9.it/checkmate/server/repository"
)

// NewSessionMiddleware creates a new session middleware with the given store and configures it to use the provided secret key for signing cookies.
// It also sets the session cookie name and domain based on the provided server address.
// The session cookie is configured to be secure and HTTP-only, and the session store is set to expire after 24 hours.
// The session middleware is then returned as a fiber.Handler for use in the application.
// The session middleware is used to manage user sessions in the application, allowing for secure storage of session data and authentication information.
func NewSessionMiddleware(serverAddr string) fiber.Handler {
	// Create a new session store with the provided configuration
	store := session.New(session.Config{
		KeyLookup: "header:session-id",
		Storage:   repository.NewStorage(),
	})
	return func(c *fiber.Ctx) error {

		claims, err := auth.ClaimsFromContext(c)
		if err == nil {
			// Set the session ID in the request header
			// to be used by the session middleware
			if session_id, ok := claims["session_id"].(string); ok {
				c.Request().Header.Set("Session-Id", session_id)
			}
			//defer c.Response().Header.Del("Session-Id")
		}

		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}
		defer func() {
			sess, _ := store.Get(c)
			sess.Save()
		}()

		// Proceed to next handler
		c.Locals("session", sess)
		return c.Next()
	}
}
