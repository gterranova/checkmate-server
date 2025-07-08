package json

import (
	"encoding/json"

	"terra9.it/checkmate/server/handlers"

	"github.com/gofiber/fiber/v2"
)

type JSONHandler struct {
	config *handlers.HandlerConfig
}

var _ handlers.Handler = &JSONHandler{}

func (t *JSONHandler) GetInfo() error {
	if len(t.config.Data) == 0 {
		t.config.Params = fiber.Map{
			"type": "empty",
		}
		return nil
	}
	if err := json.Unmarshal(t.config.Data, &t.config.Params); err != nil {
		return err
	}
	return nil
}

// Paginate implements handlers.Handler.
func (t *JSONHandler) Paginate() (*handlers.Pagination, error) {
	return handlers.DefaultPagination(t.config)
}

func (t *JSONHandler) Call(ctx *fiber.Ctx) error {
	return handlers.DefaultCallHandler(ctx, t, t.config)
}

func NewJSONHandler(config *handlers.HandlerConfig) *JSONHandler {
	t := JSONHandler{
		config: config,
	}
	if err := t.GetInfo(); err != nil {
		return nil
	}
	return &t
}

func init() {
	handlers.Manager.AddHandler(".json", func(config *handlers.HandlerConfig) handlers.Handler {
		return NewJSONHandler(config)
	})
}
