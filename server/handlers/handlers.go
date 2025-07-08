package handlers

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gopkg.in/yaml.v3"
	"terra9.it/checkmate/loader"
)

type Host struct {
	DocumentFolder string
}

type HandlerConfig struct {
	Host        *Host
	HandlerName string
	Path        string
	Data        []byte
	Params      map[string]any
}

type Handler interface {
	GetInfo() error
	Paginate() (*Pagination, error)
	Call(ctx *fiber.Ctx) error
}

func MetaForPath(config *HandlerConfig, requestedPath string) error {
	var info fs.FileInfo
	var metaFile string
	var data []byte
	var err error

	if requestedPath == "" || requestedPath == "." {
		return nil
	}

	parent := filepath.Dir(requestedPath)
	//fmt.Println(requestedPath, "not valid, trying parent", parent)
	MetaForPath(config, parent)

	if info, err = os.Stat(requestedPath); err != nil {
		return nil
	}

	if !info.IsDir() {
		return nil
	}

	metaFile = path.Join(requestedPath, "_meta.yaml")
	if _, err = os.Stat(metaFile); err != nil {
		return nil
	}

	if data, err = os.ReadFile(metaFile); err == nil {
		err = yaml.Unmarshal(data, config.Params)
		//fmt.Println(requestedPath, "found", metaFile, "with data", config.Params, err)
		return err
	}

	return nil

}

func FileForPath(requestedPath string) (filename string, err error) {
	var info fs.FileInfo

	extensions := []string{loader.CHECKLIST_EXT, loader.JSON_EXT, loader.YAML_EXT, loader.MD_EXT, loader.HTML_EXT}

	if info, err = os.Stat(requestedPath); err == nil && info.IsDir() {
		requestedPath = path.Join(requestedPath, "index")
	}
	for _, ext := range extensions {
		if _, err = os.Stat(requestedPath + ext); err == nil {
			return requestedPath + ext, nil
		}
	}
	return "", err
}

func FileHandler(config *HandlerConfig) (Handler, error) {
	var filename string
	var err error

	if filename, err = FileForPath(config.Path); err != nil {
		return nil, err
	}
	if config.Data, err = os.ReadFile(filename); err != nil {
		return nil, err
	}

	config.HandlerName = path.Ext(filename)
	if handler := Manager.GetHandler(config); handler != nil {
		return handler, nil
	}
	return nil, nil
}

func HandlerForPath(host *Host, pathUrl string) (Handler, error) {
	var mode string
	var err error

	if host == nil {
		return nil, fmt.Errorf("host is nil")
	}

	config := &HandlerConfig{
		Host:   host,
		Path:   strings.Replace(pathUrl, "/api/v1", host.DocumentFolder, 1),
		Data:   nil,
		Params: make(map[string]any),
	}

	if parts := strings.SplitN(config.Path, ":", 2); len(parts) > 1 {
		config.Path, mode = parts[0], parts[1]
		config.Params["mode"] = mode
	}

	//fmt.Println("REQ PATH:", config.Path)
	if err = MetaForPath(config, config.Path); err != nil {
		return nil, err
	}

	if handlerName, ok := config.Params["handler"]; ok {
		config.HandlerName = handlerName.(string)
		delete(config.Params, "handler")
		if handler := Manager.GetHandler(config); handler != nil {
			return handler, nil
		}
	}

	return FileHandler(config)
}

func Page(hosts map[string]*Host) fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		var host *Host
		var ok bool

		if host, ok = hosts[ctx.Hostname()]; !ok {
			return ctx.Next()
		}
		log.Println(ctx.Hostname(), ctx.Path(), host)

		handler, err := HandlerForPath(host, ctx.Path())
		if err != nil {
			return PageNotFound(hosts)(ctx)
		}

		if handler != nil {
			return handler.Call(ctx)
		}
		return PageNotFound(hosts)(ctx)
	}
}

func DefaultPagination(config *HandlerConfig) (*Pagination, error) {
	parts := make([]*PageItem, 0)

	parent := filepath.Dir(config.Path)
	parentHandler, err := HandlerForPath(config.Host, parent)
	if err == nil && parentHandler != nil {
		// If the parent handler has pagination, get the path parts
		pagination, err := parentHandler.Paginate()
		if err == nil && len(pagination.PathParts) > 0 {
			// If the pagination is nil or has no path parts, return an empty slice
			parts = append(parts, pagination.PathParts...)
		}
	}
	urlPath := strings.TrimPrefix(config.Path, path.Dir(config.Host.DocumentFolder))
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimPrefix(urlPath, "\\")

	if title, ok := config.Params["title"]; ok {
		// If the title is already set in params, use it
		parts = append(parts, &PageItem{
			Href:  urlPath,
			Title: title.(string),
		})
	}
	return &Pagination{
		PathParts: parts,
	}, nil
}

func DefaultCallHandler(ctx *fiber.Ctx, handler Handler, config *HandlerConfig) error {
	pagination, err := handler.Paginate()
	if err != nil {
		return err
	}
	config.Params["pagination"] = pagination
	return ctx.JSON(config.Params)
}

func CallFileHandler(config *HandlerConfig, ctx *fiber.Ctx) error {
	if handler, err := FileHandler(config); err == nil {
		return DefaultCallHandler(ctx, handler, config)
	}
	return ctx.Next()
}

func PageNotFound(hosts map[string]*Host) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var host *Host
		var ok bool

		//log.Println(ctx.Hostname(), ctx.Path())
		if host, ok = hosts[ctx.Hostname()]; !ok {
			return ctx.Next()
		}
		return CallFileHandler(&HandlerConfig{
			Host:   host,
			Path:   host.DocumentFolder + "/_404",
			Params: make(map[string]any),
		}, ctx)
	}
}

func SessionFromContext(c *fiber.Ctx) (*session.Session, error) {
	sess := c.Locals("session")
	if sess == nil {
		return nil, fiber.ErrUnauthorized
	}
	session, ok := sess.(*session.Session)
	//fmt.Println("Session keys", session.Keys(), ok)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}
	return session, nil
}
