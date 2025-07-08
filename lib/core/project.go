package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"terra9.it/checkmate/core/converter"
	"terra9.it/checkmate/loader"
)

const (
	CACHE_FEATURES = true
)

type ResourceLoader interface {
	Name() string
	Get(filename string) ([]byte, bool)
	Set(filename string, data []byte) error
	SaveAs(filename string) error
}

type FeatureDef struct {
	Lang      string   `json:"lang"`
	Filenames []string `json:"filenames"`
}

type TemplateDef struct {
	Name         string             `json:"name"`
	Filenames    []string           `json:"filenames"`
	Format       string             `json:"format"`
	ReferenceDoc string             `json:"reference_doc"`
	Template     *template.Template `json:"-"`
}

type Project struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	License string `json:"license"`

	Features     []Feature      `json:"-" bexpr:"features"`
	Tags         map[string]any `json:"tags" bexpr:"tags"`
	TemplateDefs []*TemplateDef `json:"templates"`
	Loader       ResourceLoader `json:"-"`

	ProjectFile string `json:"-"`
	isDirty     bool   `json:"-"`
}

type ProjectExport struct {
	Tags   []string       `json:"tags"`
	Values map[string]any `json:"values"`
}

func NewProject(loader ResourceLoader) *Project {
	p := &Project{
		Tags:         make(map[string]any),
		TemplateDefs: make([]*TemplateDef, 0),
		Loader:       loader,
		ProjectFile:  "",
	}

	if err := p.LoadFeatures(); err != nil {
		panic(err)
	}

	for _, f := range p.Features {
		f.ApplyDefaults()
	}
	if err := p.LoadProjectDataFromFile("data.json"); err != nil {
		if err.Error() != "file data.json not found" {
			panic(err)
		}
	}

	return p
}

func LoadProject(filename string) (*Project, error) {
	var pkgLoader *loader.ResourceLoader
	var err error
	if pkg, found := strings.CutSuffix(filename, loader.CHECKLIST_EXT); found {
		pkgLoader, err = loader.NewLoader(pkg)
	}
	if err != nil {
		return nil, err
	}
	project := NewProject(pkgLoader)
	project.ProjectFile = pkgLoader.Name()
	return project, nil
}

type _Project Project

func (p *Project) UnmarshalJSON(bytes []byte) (err error) {
	var isCached, ok bool
	var valueBytes []byte

	foo := _Project{}

	type _ProjectJSON struct {
		Name         string           `json:"name"`
		Author       string           `json:"author"`
		License      string           `json:"license"`
		TemplateDefs []*TemplateDef   `json:"templates"`
		Features     []map[string]any `json:"features"`
	}
	pfi := _ProjectJSON{}
	if CACHE_FEATURES {
		if cached, ok := p.Loader.Get("cachedFeatures.json"); ok {
			err = json.Unmarshal(cached, &pfi)
			if err != nil {
				return err
			}
			//foo.Name = pfi.Name
			//foo.Author = pfi.Author
			//foo.License = pfi.License
			//foo.TemplateDefs = pfi.TemplateDefs
			isCached = true
		}
	}
	if !isCached {
		if err := json.Unmarshal(bytes, &pfi); err != nil {
			return err
		}
	}

	foo.Name = pfi.Name
	foo.Author = pfi.Author
	foo.License = pfi.License
	foo.TemplateDefs = pfi.TemplateDefs
	foo.Tags = make(map[string]any)
	foo.ProjectFile = p.ProjectFile
	foo.Loader = p.Loader

	foo.Features = make([]Feature, len(pfi.Features))
	jsonFeatures := make([]map[string]any, len(pfi.Features))

	for i, feature := range pfi.Features {

		value := reflect.New(knownTypes[feature["type"].(string)]).Interface().(Feature)

		if feature["$ref"] == nil || feature["$ref"] == "" {
			valueBytes, err = json.Marshal(feature)
			if err != nil {
				return err
			}

			if err = json.Unmarshal(valueBytes, &value); err != nil {
				return err
			}
			if CACHE_FEATURES && !isCached {
				jsonFeatures[i] = feature
			}
		} else {
			valueBytes, ok = p.Loader.Get(feature["$ref"].(string))
			if !ok {
				return fmt.Errorf("ref file %s not found", feature["$ref"].(string))
			}
			if err = json.Unmarshal(valueBytes, &value); err != nil {
				return err
			}
			value.SetTag(feature["tag"].(string))
			if CACHE_FEATURES && !isCached {
				jsonFeature := make(map[string]any)
				err = json.Unmarshal(valueBytes, &jsonFeature)
				if err != nil {
					return err
				}
				jsonFeature["tag"] = value.GetTag()
				jsonFeatures[i] = jsonFeature
			}
		}

		foo.Features[i] = value
		//fmt.Println(key, value)
	}

	if CACHE_FEATURES && !isCached {
		onlyFeatures := _ProjectJSON{
			Name:         foo.Name,
			Author:       foo.Author,
			License:      foo.License,
			TemplateDefs: foo.TemplateDefs,
			Features:     jsonFeatures,
		}
		bYaml, _ := json.Marshal(onlyFeatures)
		p.Loader.Set("cachedFeatures.json", bYaml)
	}

	*p = Project(foo)

	return nil

}

func (p *Project) MarshalJSON() ([]byte, error) {
	type jsonProject struct {
		Type     string             `json:"type"`
		Name     string             `json:"name"`
		Author   string             `json:"author"`
		License  string             `json:"license"`
		Features map[string]Feature `json:"properties"`
	}
	foo := jsonProject{
		Type:    "object",
		Name:    p.Name,
		Author:  p.Author,
		License: p.License,
	}

	foo.Features = make(map[string]Feature)

	for _, feature := range p.Features {
		key := feature.GetTag()
		if strings.HasPrefix(key, "_") {
			continue
		}
		foo.Features[key] = feature
		//foo.Features[key].SetTag(key)
	}

	valueBytes, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}

	return valueBytes, nil
}

func (p *Project) LoadFeatures() error {
	if content, ok := p.Loader.Get("config.json"); ok {
		if err := json.Unmarshal(content, p); err != nil {
			return err
		}
	}
	/*
		if p.FeatureDefs != nil {
			for _, f := range p.FeatureDefs {
				for _, tmplFile := range f.Filenames {
					var i18nFeatures _Project
					if content, ok := p.Loader.Get(tmplFile); ok {
						if err := json.Unmarshal(content, &i18nFeatures); err != nil {
							return err
						}
						for k, v := range i18nFeatures.Features {
							p.Features[k+"."+f.Lang] = v
						}
					} else {
						return fmt.Errorf("file %s not found", tmplFile)
					}
				}
			}
		}
		if p.StatusDefs != nil {
			tmpl := template.New("status")
			for _, tmplFile := range p.StatusDefs.Filenames {
				if content, ok := p.Loader.Get(tmplFile); ok {
					tmpl = template.Must(tmpl.Parse(string(content)))
				} else {
					return fmt.Errorf("file %s not found", tmplFile)
				}
			}
			p.StatusDefs.Template = tmpl
		}
	*/

	for _, t := range p.TemplateDefs {
		t.Template = template.New("template")
		for _, tmplFile := range t.Filenames {
			if content, ok := p.Loader.Get(tmplFile); ok {
				t.Template = template.Must(t.Template.Parse(string(content)))
			} else {
				return fmt.Errorf("file %s not found", tmplFile)
			}
		}
	}
	return nil
}

func (p *Project) UpdateTags() {
	for k := range p.Tags {
		delete(p.Tags, k)
	}

	for _, f := range p.Features {
		for _, t := range f.ApplicableFeatures() {
			p.Tags[t.GetTag()] = t.GetValue()
		}
	}
}

func (p *Project) SetFeature(tag string, value any) error {
	for _, f := range p.Features {
		if err := f.Set(tag, value); err != nil {
			return err
		}
	}
	if _, err := p.Validate(tag); err != nil {
		return err
	}
	p.SetDirty(true)
	return nil
}

func (p *Project) Validate(tag string) (changed bool, err error) {
	var count int
	for count = 0; count < 100; count++ {
		changed = false
		p.UpdateTags()
		for _, f := range p.Features {
			var fchanged bool
			fchanged, err = f.Validate(tag, p)
			if err != nil {
				return
			}
			//if fchanged {
			//	fmt.Println("tag:", tag, "against", f.GetTag(), "changed")
			//}

			changed = changed || fchanged
		}
		if err != nil {
			return
		}
		if !changed {
			p.UpdateTags()
			for _, f := range p.Features {
				var fchanged bool
				fchanged, err = f.Validate("", p)
				if err != nil {
					return
				}
				//if fchanged {
				//	fmt.Println("tag:", tag, "against", f.GetTag(), "changed")
				//}
				changed = changed || fchanged
			}
			if !changed {
				break
			}
		}
		count++
	}
	if count == 100 {
		err = fmt.Errorf("too many iterations in validation")
	}
	return
}

func (p *Project) Render(t *TemplateDef) (output string, err error) {

	var buf bytes.Buffer

	if t == nil || t.Template == nil {
		return "", fmt.Errorf("template not found")
	}

	err = t.Template.Execute(&buf, p)
	if err != nil {
		return
	}
	output = strings.Trim(regexp.MustCompile("\r\n[\r\n]+").ReplaceAllString(buf.String(), "\r\n\r\n"), "\r\n")

	if t.Format == "docx" {
		f := converter.File{}
		tmpFile, _ := f.TempFile("output.docx")
		var fname, fref string
		if fname, err = f.WriteToTempFile("output.md", []byte(output)); err != nil {
			return
		}

		args := make([]string, 0)
		args = append(args, "-f markdown+pipe_tables", "--from=markdown", "--to=docx",
			"--columns=43", "--wrap=preserve", "--output="+tmpFile)

		if t.ReferenceDoc != "" {
			var refBytes []byte
			var ok bool
			refBytes, ok = p.Loader.Get(t.ReferenceDoc)
			if ok {
				if fref, err = f.WriteToTempFile("refdoc.docx", refBytes); err != nil {
					return
				}
				args = append(args, "--reference-doc="+fref)
			}
		}
		args = append(args, fname)
		if _, err = converter.ExecCommand(time.Second*120, "pandoc", args...); err != nil {
			return
		}
		if err = converter.FixDocxTableStyle(tmpFile, "output.docx"); err != nil {
			return
		}
		return "output.docx", nil

	}

	return
}

func (p *Project) Evaluate() string {
	// Files are provided as a slice of strings.
	if len(p.TemplateDefs) == 0 {
		return ""
	}

	output, _ := p.Render(p.TemplateDefs[0])
	return output
}

func (p *Project) SaveProject() error {
	return p.SaveProjectAs(p.ProjectFile)
}

func (p *Project) SaveProjectAs(filename string) error {
	export := p.ExportData()

	data, err := json.Marshal(&export)
	if err != nil {
		return err
	}

	if err := p.Loader.Set("data.json", data); err != nil {
		//return saveToFile(Settings.StoragePath(), data)
		return err
	}
	p.ProjectFile = filename
	p.SetDirty(false)
	return p.Loader.SaveAs(filename)
}

func (p *Project) ExportData() ProjectExport {
	export := ProjectExport{Tags: make([]string, 0), Values: make(map[string]any)}

	p.Validate("")
	for t, v := range p.Tags {
		switch value := v.(type) {
		case bool:
			export.Tags = append(export.Tags, t)
		default:
			export.Values[t] = value
		}
	}
	return export
}

func (p *Project) LoadProjectData(export ProjectExport) error {

	p.ResetFeatures()

	for k, v := range export.Values {
		if err := p.SetFeature(k, v); err != nil {
			return err
		}
	}

	for _, v := range export.Tags {
		if err := p.SetFeature(v, true); err != nil {
			return err
		}
	}
	if _, err := p.Validate(""); err != nil {
		return err
	}
	p.SetDirty(false)
	return nil
}

func (p *Project) LoadProjectDataFromFile(filename string) error {
	export := ProjectExport{Tags: make([]string, 0), Values: make(map[string]any)}

	data, ok := p.Loader.Get(filename)
	if !ok {
		return fmt.Errorf("file %s not found", filename)
	}

	err := json.Unmarshal(data, &export)
	if err != nil {
		return err
	}

	return p.LoadProjectData(export)
}

func (p *Project) SetDirty(b bool) {
	p.isDirty = b
}

func (p *Project) Dirty() bool {
	return p.isDirty
}

func (p *Project) ResetFeatures() {
	for k := range p.Tags {
		delete(p.Tags, k)
	}

	for _, f := range p.Features {
		for _, t := range f.ApplicableFeatures() {
			switch t.GetValue().(type) {
			case bool:
				t.SetValue(false)
			case int64:
				t.SetValue(0)
			case string:
				t.SetValue("")
			default:
				t.SetValue(nil)
			}
		}
	}
	p.ProjectFile = ""
	p.SetDirty(false)
}

func (p *Project) GetValue() any {
	values := make(map[string]any)
	for _, feature := range p.Features {
		key := feature.GetTag()
		if key == "_" {
			continue
		}
		values[key] = feature.GetValue()
		//if value := feature.GetValue(); value != nil {
		//	values[k] = value
		//}
	}
	return values
}

func (p *Project) SetValue(value any) error {
	switch t := value.(type) {
	case map[string]any:
		for k, v := range t {
			for _, f := range p.Features {
				if f.GetTag() == k {
					f.SetValue(v)
					break
				}
			}
		}
	default:
		return fmt.Errorf("unknown type %T for %v", t, t)
	}
	return nil
}
