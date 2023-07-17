package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"terra9.it/vadovia/assets"
	"terra9.it/vadovia/core/features"
)

type Project struct {
	Features     map[string]*features.Feature `json:"features" bexpr:"features"`
	Tags         map[string]any               `json:"tags" bexpr:"tags"`
	FeatureOrder []string                     `json:"feature_order"`
	template     *template.Template           `json:"-"`
}

func NewProject() *Project {
	paths := []string{
		"output.tmpl",
	}

	p := &Project{
		Features: make(map[string]*features.Feature),
		Tags:     make(map[string]any),
		template: template.Must(template.New("output.tmpl").ParseFiles(paths...)),
	}

	if err := json.Unmarshal(assets.Features, p); err != nil {
		panic(err)
	}

	for _, f := range p.Features {
		f.ApplyDefaults()
	}

	return p
}

func (p *Project) UpdateTags() {
	for k := range p.Tags {
		delete(p.Tags, k)
	}

	for _, f := range p.Features {
		for _, t := range f.ApplicableFeatures() {
			p.Tags[t.Tag] = t.Value
		}
	}
}

func (p *Project) SetFeature(tag string, value interface{}) {
	for _, f := range p.Features {
		f.Set(tag, value)
	}
	p.Validate(tag)
}

func (p *Project) Validate(tag string) {
	for {
		p.UpdateTags()
		changed := false
		for _, f := range p.Features {
			changed = changed || f.Validate(tag, p)
		}
		if !changed {
			p.UpdateTags()
			for _, f := range p.Features {
				changed = changed || f.Validate("", p)
			}
			if !changed {
				break
			}
		}
	}
	/*
		for t := range p.Tags {
			fmt.Printf(`%v=%v `, t, p.Tags[t])
		}
		fmt.Println("")
	*/
}

func (p *Project) Describe() string {

	var buf bytes.Buffer

	err := p.template.Execute(&buf, p)
	if err != nil {
		fmt.Println(err)
	}
	output := strings.Trim(regexp.MustCompile("\r\n[\r\n]+").ReplaceAllString(buf.String(), "\r\n\r\n"), "\r\n")

	return output
}

func (p *Project) Evaluate() string {
	// Files are provided as a slice of strings.
	var buf bytes.Buffer

	err := p.template.Execute(&buf, p)
	if err != nil {
		fmt.Println(err)
	}
	output := strings.Trim(regexp.MustCompile("\r\n[\r\n]+").ReplaceAllString(buf.String(), "\r\n\r\n"), "\r\n")
	//fmt.Println(output)

	return strings.Split(output, "\n")[0]
}
