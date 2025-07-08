package html

import (
	"bytes"

	"terra9.it/checkmate/server/handlers"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

var HTML_HEADER_SEPARATOR = []byte{'-', '-', '-'}

type HTMLHandler struct {
	config *handlers.HandlerConfig
}

var _ handlers.Handler = &HTMLHandler{}

func (t *HTMLHandler) GetInfo() error {
	if len(t.config.Data) == 0 {
		t.config.Params = fiber.Map{
			"type": "empty",
		}
		return nil
	}
	t.config.Params["type"] = "html"

	if bytes.HasPrefix(t.config.Data, HTML_HEADER_SEPARATOR) {
		parts := bytes.Split(t.config.Data, HTML_HEADER_SEPARATOR)
		if len(parts) == 3 {
			t.config.Data = parts[2]
			if err := yaml.Unmarshal(parts[1], &t.config.Params); err != nil {
				return err
			}
		}
	}
	//if _, ok := t.config.Params["title"]; !ok {
	//	t.config.Params["title"] = ExtractTitleFromBytes(t.config.Data)
	//}
	t.config.Params["content"] = string(t.config.Data)
	return nil
}

// Paginate implements handlers.Handler.
func (t *HTMLHandler) Paginate() (*handlers.Pagination, error) {
	return handlers.DefaultPagination(t.config)
}

func (t *HTMLHandler) Call(ctx *fiber.Ctx) error {
	return handlers.DefaultCallHandler(ctx, t, t.config)
}

func NewHTMLHandler(config *handlers.HandlerConfig) *HTMLHandler {
	t := HTMLHandler{
		config: config,
	}
	if err := t.GetInfo(); err != nil {
		return nil
	}
	return &t
}

func init() {
	handlers.Manager.AddHandler(".html", func(config *handlers.HandlerConfig) handlers.Handler {
		return NewHTMLHandler(config)
	})
}
