package checklist

import (
	"github.com/gofiber/fiber/v2"
	"terra9.it/checkmate/core"
	"terra9.it/checkmate/loader"
	"terra9.it/checkmate/server/handlers"
)

type CHLXHandler struct {
	config  *handlers.HandlerConfig
	project *core.Project
}

var _ handlers.Handler = &CHLXHandler{}

func (t *CHLXHandler) GetInfo() error {
	return nil
}

// Paginate implements handlers.Handler.
func (t *CHLXHandler) Paginate() (*handlers.Pagination, error) {
	return doPaginate(t.config, t.project)
}

func (t *CHLXHandler) Call(ctx *fiber.Ctx) error {
	var pkgLoader *loader.ResourceLoader
	var err error

	pkgLoader, err = loader.NewEmptyLoader("checklist").FromBuffer(t.config.Data)
	if err != nil {
		return err
	}

	project := core.NewProject(pkgLoader)
	t.project = project

	return doChecklist(ctx, t.config, t.project)
}

func NewCHLXHandler(config *handlers.HandlerConfig) *CHLXHandler {
	t := CHLXHandler{
		config: config,
	}
	if err := t.GetInfo(); err != nil {
		return nil
	}
	return &t
}
