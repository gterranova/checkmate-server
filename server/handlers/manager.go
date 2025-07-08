package handlers

import "github.com/gofiber/fiber/v2"

type HandlerFactoryFunc func(config *HandlerConfig) Handler

type HandlersManager struct {
	handlers map[string]HandlerFactoryFunc
}

var Manager = &HandlersManager{}

func (ds *HandlersManager) AddHandler(name string, handlerFactory HandlerFactoryFunc) {
	if ds.handlers == nil {
		ds.handlers = make(map[string]HandlerFactoryFunc)
	}
	if _, exists := ds.handlers[name]; exists {
		return
	}
	ds.handlers[name] = handlerFactory
}

func (ds *HandlersManager) GetHandler(config *HandlerConfig) Handler {
	if _, exists := ds.handlers[config.HandlerName]; exists {
		return ds.handlers[config.HandlerName](config)
	}
	return nil
}

func (ds *HandlersManager) Call(
	ctx *fiber.Ctx,
	name, arguments string,
) error {
	return nil
}
