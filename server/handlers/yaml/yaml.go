package yaml

import (
	"gopkg.in/yaml.v3"
	"terra9.it/checkmate/server/handlers"

	"github.com/gofiber/fiber/v2"
)

type YAMLHandler struct {
	config *handlers.HandlerConfig
}

var _ handlers.Handler = &YAMLHandler{}

func (t *YAMLHandler) GetInfo() error {
	if len(t.config.Data) == 0 {
		t.config.Params = fiber.Map{
			"type": "empty",
		}
		return nil
	}
	if err := yaml.Unmarshal(t.config.Data, &t.config.Params); err != nil {
		return err
	}
	return nil
}

// Paginate implements handlers.Handler.
func (t *YAMLHandler) Paginate() (*handlers.Pagination, error) {
	return handlers.DefaultPagination(t.config)
}

func (t *YAMLHandler) Call(ctx *fiber.Ctx) error {
	return handlers.DefaultCallHandler(ctx, t, t.config)
}

func NewYAMLHandler(config *handlers.HandlerConfig) *YAMLHandler {
	t := YAMLHandler{
		config: config,
	}
	if err := t.GetInfo(); err != nil {
		return nil
	}
	return &t
}

func init() {
	handlers.Manager.AddHandler(".yaml", func(config *handlers.HandlerConfig) handlers.Handler {
		return NewYAMLHandler(config)
	})
}
