package core

import (
	"encoding/json"
)

type ChecklistFormly struct {
	*Checklist
}

func (feature *ChecklistFormly) MarshalJSON() ([]byte, error) {

	type jsonChecklist struct {
		Type     string `json:"type"`
		Title    string `json:"title,omitempty"`
		Disabled bool   `json:"disabled,omitempty"`

		Enum   []string `json:"enum" bexpr:"enum"`
		Widget struct {
			FormlyConfig map[string]any `json:"formlyConfig"`
		} `json:"widget"`
	}
	foo := jsonChecklist{
		Type:     "array",
		Title:    feature.Title,
		Disabled: feature.Disabled,
		Enum:     make([]string, 0),
	}

	foo.Widget.FormlyConfig = make(map[string]any)

	props := ChecklistProps{Options: make([]OptionLabelValue, 0), InfoUrls: make([]string, 0)}
	for _, opt := range feature.Enum {
		foo.Enum = append(foo.Enum, opt.Tag)
		props.Options = append(props.Options, OptionLabelValue{
			Label:    opt.Title,
			Value:    opt.Tag,
			Disabled: opt.Disabled,
		})
		props.InfoUrls = append(props.InfoUrls, opt.InfoUrl)
	}
	props.Multiple = true
	props.Disabled = feature.Disabled
	props.HideDisabled = feature.HideDisabled
	if feature.Default != nil {
		props.DefaultValue = feature.Default.(string)
	}

	foo.Widget.FormlyConfig["type"] = "checklist"
	foo.Widget.FormlyConfig["props"] = props

	valueBytes, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}

	return valueBytes, nil
}
