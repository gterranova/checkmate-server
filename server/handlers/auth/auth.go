package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/golang-jwt/jwt/v5"

	"github.com/spf13/viper"
	"terra9.it/checkmate/server/handlers"
	"terra9.it/checkmate/server/models"
	"terra9.it/checkmate/server/repository"
)

type AuthHandler struct {
	config *handlers.HandlerConfig
}

var _ handlers.Handler = &AuthHandler{}

func (a *AuthHandler) GetInfo() error {
	return nil
}

// Paginate implements handlers.Handler.
func (t *AuthHandler) Paginate() (*handlers.Pagination, error) {
	return handlers.DefaultPagination(t.config)
	/*
		pagination, err := handlers.DefaultPagination(t.config)
		if err != nil {
			return nil, err
		}
		pagination.PathParts = append(pagination.PathParts, &handlers.PageItem{
			Title: "Auth",
			Href:  "",
		})
		return pagination, nil
	*/
}

// Call implements handlers.Handler.
func (a *AuthHandler) Call(ctx *fiber.Ctx) error {
	claims, err := ClaimsFromContext(ctx)
	if err != nil {
		// not logged in
		if strings.HasSuffix(a.config.Path, "/login") {
			return a.Login(ctx)
		}

		if strings.HasSuffix(a.config.Path, "/settings") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	if strings.HasSuffix(a.config.Path, "/logout") {
		return a.Logout(ctx)
	}

	if strings.HasSuffix(a.config.Path, "/token") {
		return a.Token(ctx)
	}

	if strings.HasSuffix(a.config.Path, "/user") {
		return ctx.JSON(fiber.Map{
			"ID":    claims["ID"],
			"email": claims["email"],
		})
	}

	if strings.HasSuffix(a.config.Path, "/settings") {
		return a.SettingsPage(ctx)
	}

	return handlers.CallFileHandler(a.config, ctx)
}

func (a *AuthHandler) LoginPage(ctx *fiber.Ctx) error {
	return handlers.CallFileHandler(&handlers.HandlerConfig{
		Host:   a.config.Host,
		Path:   a.config.Host.DocumentFolder + "/auth/login",
		Params: a.config.Params,
	}, ctx)
}

func (a *AuthHandler) SettingsPage(ctx *fiber.Ctx) error {
	return handlers.CallFileHandler(&handlers.HandlerConfig{
		Host:   a.config.Host,
		Path:   a.config.Host.DocumentFolder + "/auth/setting",
		Params: a.config.Params,
	}, ctx)
}

// Login route
func (a *AuthHandler) Login(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return a.LoginPage(c)
	}

	// Extract the credentials from the request body
	loginRequest := new(models.FormResponse[models.LoginRequest])
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Find the user by credentials
	user, err := repository.FindByCredentials(loginRequest.Changes.Email, loginRequest.Changes.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	day := time.Hour * 24
	// Create the JWT claims, which includes the user ID and expiry time
	claims := jwt.MapClaims{
		"ID":         user.ID,
		"email":      user.Email,
		"session_id": user.SessionID,
		"exp":        time.Now().Add(day * 1).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	key := viper.GetString("AUTH_SECRET")

	t, err := token.SignedString([]byte(key))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Response().Header.Set("Authorization", t)

	a.config.Params["user"] = user
	// Return the token
	return c.JSON(fiber.Map{
		"user": user,
	})
	//return c.JSON(fiber.Map{
	//	"type": "redirect",
	//	"url":  "auth/settings",
	//	"user": user,
	//})
	//return a.SettingsPage(c)
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	s := c.Locals("session")
	if s != nil {
		if sess := s.(*session.Session); sess != nil {
			sess.Destroy()
		}
	}
	c.Locals("user", nil)
	return c.JSON(fiber.Map{
		"user": nil,
	})
	//return c.Next()
}

// Token route
func (a *AuthHandler) Token(c *fiber.Ctx) error {

	claims, err := ClaimsFromContext(c)
	if err != nil {
		return err
	}
	//day := time.Hour * 24
	// Create the JWT claims, which includes the user ID and expiry time
	newclaims := jwt.MapClaims{
		"ID":         claims["ID"],
		"email":      claims["email"],
		"session_id": claims["session_id"],
		//"exp":        time.Now().Add(day * 1).Unix(),
		"exp": time.Now().Add(time.Second * 60).Unix(),
	}
	// Create token
	newtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, newclaims)
	// Generate encoded token and send it as response.
	key := viper.GetString("AUTH_SECRET")
	t, err := newtoken.SignedString([]byte(key))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Return the token
	c.Response().Header.Set("Authorization", t)
	return c.JSON(models.LoginResponse{
		Token: t,
	})
}

func ClaimsFromContext(c *fiber.Ctx) (jwt.MapClaims, error) {
	u := c.Locals("user")
	if u == nil {
		//fmt.Println("No user found in context")
		return nil, fmt.Errorf("unauthorized")
	}
	// Check if the token is valid
	token := u.(*jwt.Token)
	if token == nil || !token.Valid {
		//fmt.Println("No valid token found in context")
		return nil, fmt.Errorf("invalid token")
	}
	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		//fmt.Println("Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func init() {
	handlers.Manager.AddHandler("auth", func(config *handlers.HandlerConfig) handlers.Handler {
		return &AuthHandler{
			config: config,
		}
	})
}
