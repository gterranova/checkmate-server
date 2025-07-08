package checklist

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"terra9.it/checkmate/core"
	"terra9.it/checkmate/loader"
	"terra9.it/checkmate/server/handlers"
)

type ChecklistHandler struct {
	basePath string
	config   *handlers.HandlerConfig
	project  *core.Project
}

var _ handlers.Handler = &ChecklistHandler{}

func (t *ChecklistHandler) GetInfo() error {
	basePath, err := ProjectBasePath(t.config.Path)
	if err != nil {
		return err
	}

	feature := strings.TrimPrefix(t.config.Path, basePath)
	feature = strings.TrimPrefix(feature, "/")
	t.basePath = basePath
	t.config.Params["feature"] = feature

	return nil
}

// Paginate implements handlers.Handler.
func (t *ChecklistHandler) Paginate() (*handlers.Pagination, error) {
	return doPaginate(t.config, t.project)
}

func (t *ChecklistHandler) Call(ctx *fiber.Ctx) error {
	project := core.NewProject(loader.NewEmptyLoader(t.basePath))
	t.project = project

	return doChecklist(ctx, t.config, t.project)
}

func NewChecklistHandler(config *handlers.HandlerConfig) *ChecklistHandler {
	t := ChecklistHandler{
		config: config,
	}
	if err := t.GetInfo(); err != nil {
		return nil
	}
	return &t
}
