package markdown

import (
	"bytes"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
	"terra9.it/checkmate/server/handlers"
)

var YAML_HEADER_SEPARATOR = []byte{'-', '-', '-'}

type MDHandler struct {
	config *handlers.HandlerConfig
}

var _ handlers.Handler = &MDHandler{}

func (t *MDHandler) GetInfo() error {
	if len(t.config.Data) == 0 {
		t.config.Params = fiber.Map{
			"type": "empty",
		}
		return nil
	}
	t.config.Params["type"] = "markdown"

	if bytes.HasPrefix(t.config.Data, YAML_HEADER_SEPARATOR) {
		parts := bytes.Split(t.config.Data, YAML_HEADER_SEPARATOR)
		if len(parts) == 3 {
			t.config.Data = parts[2]
			if err := yaml.Unmarshal(parts[1], &t.config.Params); err != nil {
				return err
			}
		}
	}
	if _, ok := t.config.Params["title"]; !ok {
		t.config.Params["title"] = t.extractTitleFromBytes(t.config.Data)
	}
	t.config.Params["content"] = string(t.config.Data)
	return nil
}

// Paginate implements handlers.Handler.
func (t *MDHandler) Paginate() (*handlers.Pagination, error) {
	return handlers.DefaultPagination(t.config)
}

func (t *MDHandler) Call(ctx *fiber.Ctx) error {
	return handlers.DefaultCallHandler(ctx, t, t.config)
}

// ExtractTitleFromBytes extracts the title from Markdown content provided as []byte.
func (t *MDHandler) extractTitleFromBytes(content []byte) string {
	lines := regexp.MustCompile(`\r\n|\n|\r`).Split(string(content), -1)

	for _, line := range lines {
		// Check for level 1 heading
		match := regexp.MustCompile(`^#\s+(.+)`).FindStringSubmatch(line)
		if len(match) > 1 {
			return match[1]
		}
	}

	// Return an empty string if no title is found
	return ""
}

func NewMDHandler(config *handlers.HandlerConfig) *MDHandler {
	t := MDHandler{
		config: config,
	}
	if err := t.GetInfo(); err != nil {
		return nil
	}
	return &t
}

func init() {
	handlers.Manager.AddHandler(".md", func(config *handlers.HandlerConfig) handlers.Handler {
		return NewMDHandler(config)
	})
}
