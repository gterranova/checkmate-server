package checklist

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"terra9.it/checkmate/core"
	"terra9.it/checkmate/server/handlers"
	"terra9.it/checkmate/server/models"
)

func doPaginate(config *handlers.HandlerConfig, project *core.Project) (*handlers.Pagination, error) {
	var featureName string
	var prevFeature core.Feature
	var basePath string
	var err error

	params := config.Params

	basePath, err = ProjectBasePath(config.Path)
	if err != nil {
		return nil, err
	}

	urlPath := strings.TrimPrefix(basePath, config.Host.DocumentFolder)
	urlPath = strings.TrimPrefix(urlPath, "/")

	paginate := handlers.Pagination{
		Total:     0,
		Count:     0,
		PathParts: make([]*handlers.PageItem, 0),
	}

	// path paths for breadcrumbs
	paginate.PathParts = append(paginate.PathParts, doBreadcrumbs(config, basePath, project, params)...)

	if len(project.Features) == 0 {
		// If there are no features or the feature is not specified, return
		return &paginate, nil
	}

	if f, ok := params["feature"]; ok {
		featureName = f.(string)
	}

	if len(featureName) == 0 {
		for i, feat := range project.Features {
			if !feat.IsDisabled() && !strings.HasPrefix(feat.GetTag(), "_") {
				featureName = project.Features[i].GetTag()
				break
			}
		}
	}

	//fmt.Println(config.Path, feature)
	featureIsDisabled := func(f core.Feature) bool {
		if f.IsDisabled() {
			return true
		}
		allChildrenDisabled := true
		for _, child := range f.GetChildren() {
			allChildrenDisabled = allChildrenDisabled && child.IsDisabled()
			if !allChildrenDisabled {
				return false
			}
		}
		return allChildrenDisabled
	}

	// Get the current page from the query parameters
	for i, curFeature := range project.Features {
		if curFeature == nil || featureIsDisabled(curFeature) || strings.HasPrefix(curFeature.GetTag(), "_") {
			curFeature = nil
			continue
		}
		paginate.Total++
		if paginate.Total == 1 {
			paginate.First = &handlers.PageItem{
				Href:  fmt.Sprintf("%s/%s", urlPath, project.Features[i].GetTag()),
				Title: curFeature.GetTitle(),
			}
		}
		if curFeature.GetTag() == featureName {
			paginate.Count = paginate.Total
			if paginate.Count > 1 && prevFeature != nil {
				paginate.Prev = &handlers.PageItem{
					Href:  fmt.Sprintf("%s/%s", urlPath, prevFeature.GetTag()),
					Title: prevFeature.GetTitle(),
				}
			}
		}
		if paginate.Count > 0 && paginate.Count+1 == paginate.Total {
			paginate.Next = &handlers.PageItem{
				Href:  fmt.Sprintf("%s/%s", urlPath, project.Features[i].GetTag()),
				Title: curFeature.GetTitle(),
			}
		}
		paginate.Last = &handlers.PageItem{
			Href:  fmt.Sprintf("%s/%s", urlPath, project.Features[i].GetTag()),
			Title: curFeature.GetTitle(),
		}
		prevFeature = curFeature
	}
	switch paginate.Count {

	case 0:
		// If the feature is not found, return nil to indicate no pagination
		return &paginate, nil
	case 1:
		// If the prev page is the same as the first page, remove it
		paginate.First = nil
		paginate.Prev = nil
	case paginate.Total:
		// If the next page is the same as the last page, remove it
		paginate.Next = nil
		paginate.Last = nil
	}

	return &paginate, nil
}

func doBreadcrumbs(config *handlers.HandlerConfig, basePath string, project *core.Project, params map[string]any) []*handlers.PageItem {
	parts := make([]*handlers.PageItem, 0)

	parent := filepath.Dir(basePath)
	parentHandler, err := handlers.HandlerForPath(config.Host, parent)
	if err == nil && parentHandler != nil {
		// If the parent handler has pagination, get the path parts
		pagination, err := parentHandler.Paginate()
		if err == nil && len(pagination.PathParts) > 0 {
			// If the pagination is nil or has no path parts, return an empty slice
			parts = append(parts, pagination.PathParts...)
		}
	}
	urlPath := strings.TrimPrefix(basePath, config.Host.DocumentFolder)
	urlPath = strings.TrimPrefix(urlPath, "/")

	parts = append(parts, &handlers.PageItem{
		Href:  urlPath,
		Title: project.Name,
	})
	var featureName string
	if f, ok := params["feature"]; ok {
		featureName = f.(string)
	}

	if len(featureName) > 0 {
		for _, feature := range project.Features {
			if feature.GetTag() == featureName {
				parts = append(parts, &handlers.PageItem{
					Href:  fmt.Sprintf("%s/%s", urlPath, featureName),
					Title: feature.GetTitle(),
				})
			}
		}
	}
	return parts
}

func doSendForm(ctx *fiber.Ctx, config *handlers.HandlerConfig, project *core.Project) error {
	var singlePage bool

	params := config.Params

	if params == nil {
		params = make(map[string]any)
	}
	params["type"] = "form"
	params["title"] = project.Name

	if c, ok := config.Params["class"]; ok {
		params["class"] = c
	}

	var featureName string
	if f, ok := params["single_page"]; ok {
		singlePage = f.(bool)
	}

	if !singlePage {

		if f, ok := params["feature"]; ok {
			featureName = f.(string)
		}

		if len(featureName) == 0 {
			for i, feat := range project.Features {
				if !feat.IsDisabled() && !strings.HasPrefix(feat.GetTag(), "_") {
					featureName = project.Features[i].GetTag()
					break
				}
			}
		}
	}

	if !singlePage && len(featureName) > 0 {
		var feature core.Feature
		var found bool

		for _, feature = range project.Features {
			if feature.GetTag() == featureName {
				found = true
				break
			}
		}
		if !found {
			return ctx.Next()
		}

		//fmt.Println("Project feature", feature.GetTag(), feature.GetValue())

		wrapper := map[string]any{
			"type": "object",
			"properties": map[string]any{
				"feature": feature,
			},
		}
		value_wrapper := map[string]any{
			"feature": feature.GetValue(),
		}
		params["schema"] = wrapper
		params["model"] = value_wrapper
	} else {
		params["schema"] = project
		params["model"] = project.GetValue()
	}

	if !singlePage {
		if pagination, err := doPaginate(config, project); err != nil {
			return err
		} else {
			params["pagination"] = pagination
		}
	}
	//fmt.Println("Project features")
	//for k, v := range project.Tags {
	//	fmt.Println(k, v)
	//}
	return ctx.JSON(params)

}

func doChecklist(ctx *fiber.Ctx, config *handlers.HandlerConfig, project *core.Project) error {
	var singlePage bool

	sess, _ := handlers.SessionFromContext(ctx)
	params := config.Params

	if f, ok := params["single_page"]; ok {
		singlePage = f.(bool)
	}

	if sess != nil {
		data := sess.Get("project")
		if data != nil {
			//var export core.ProjectExport
			//err := json.Unmarshal(data.([]byte), &export)
			values := make(map[string]any)
			err := json.Unmarshal(data.([]byte), &values)
			if err != nil {
				return err
			}
			//if err := project.LoadProjectData(export); err != nil {
			if err := project.SetValue(values); err != nil {
				//return ctx.Next()
				return err
			}
		}
		defer func() {
			//d := project.ExportData()
			buf, _ := json.Marshal(project.GetValue())
			//fmt.Println("Saving project data", sess.ID())
			sess.Set("project", buf)
			sess.Save()
		}()
	}

	if ctx.Method() == "PUT" {
		var export models.FormResponse[map[string]any]
		export.Changes = make(map[string]any)
		export.Model = make(map[string]any)

		if err := ctx.BodyParser(&export); err != nil {
			return fiber.NewError(500, err.Error())
		}
		var featureName string
		var found bool

		if f, ok := params["feature"]; ok {
			featureName = f.(string)
		}

		if !singlePage && featureName != "" {
			//fmt.Println("Updating feature", feature, export)
			for _, feature := range project.Features {
				if feature.GetTag() == featureName {
					found = true
					break
				}
			}
			if !found {
				return ctx.RedirectToRoute("/", nil)
			}
			//value := export.Model["feature"]
			//if err := schema.SetValue(value); err != nil {
			//	return err
			//}
			changed := export.Changes["feature"]
			if err := project.SetFeature(featureName, changed); err != nil {
				return err
			}
		} else {
			//if err := project.SetValue(export.Model); err != nil {
			//	return err
			//}
			for k, v := range export.Changes {
				if err := project.SetFeature(k, v); err != nil {
					return err
				}
			}
		}
	}
	//project.UpdateTags()
	_, err := project.Validate("")
	if err != nil {
		return err
	}
	if feedback, ok := params["feedback"]; ok {
		output := project.Evaluate()
		config.Params[feedback.(string)] = strings.TrimSpace(output)
	}
	/*
		if config.Params["mode"] != "render" {
			err = project.StatusDefs.Template.Execute(&buf, project)
			if err == nil {
				//fmt.Println(buf.String())
			} else {
				fmt.Println(err)
			}
			fmt.Println("Project features:", project.Tags)
			return doSendForm(ctx, config, project)
		}

		var buf bytes.Buffer

		err = project.StatusDefs.Template.Execute(&buf, project)
		if err == nil {
			//fmt.Println(buf.String())
		} else {
			fmt.Println(err)
		}

		config.Data = buf.Bytes()
		return handlers.Manager.GetHandler(&handlers.HandlerConfig{
			HandlerName: ".md",
			Data:        config.Data,
			Params:      config.Params,
		}).Call(ctx)
	*/
	return doSendForm(ctx, config, project)
}

func ProjectBasePath(requestedPath string) (basePath string, err error) {
	var info fs.FileInfo

	if requestedPath == "" || requestedPath == "." || requestedPath == ".." {
		return
	}

	if info, err = os.Stat(requestedPath); err != nil {
		goto next
	}

	if !info.IsDir() {
		goto next
	}

	return strings.ReplaceAll(requestedPath, "\\", "/"), nil

next:
	parent := filepath.Dir(requestedPath)
	//fmt.Println(requestedPath, "not valid, trying parent", parent)
	return ProjectBasePath(parent)
}

func init() {
	handlers.Manager.AddHandler(".chlx", func(config *handlers.HandlerConfig) handlers.Handler {
		return NewCHLXHandler(config)
	})
	handlers.Manager.AddHandler("checklist", func(config *handlers.HandlerConfig) handlers.Handler {
		return NewChecklistHandler(config)
	})
}
