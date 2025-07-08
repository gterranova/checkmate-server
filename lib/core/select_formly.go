package core

import (
	"encoding/json"
)

type SelectFormly struct {
	*Select
}

type SelectProps struct {
	Disabled     bool               `json:"disabled,omitempty"`
	Multiple     bool               `json:"multiple,omitempty"`
	DefaultValue string             `json:"defaultValue,omitempty"`
	HideDisabled bool               `json:"hide_disabled,omitempty"`
	InfoUrl      string             `json:"info_url,omitempty"`
	Options      []OptionLabelValue `json:"options,omitempty"`
}

func (feature *SelectFormly) UnmarshalJSON(bytes []byte) (err error) {
	type _Select Select
	foo := _Select{}
	if err := json.Unmarshal(bytes, &foo); err != nil {
		return err
	}

	type _FooChildren struct {
		Enum []*Option `json:"enum" bexpr:"enum"`
	}

	fooitems := _FooChildren{}
	if err := json.Unmarshal(bytes, &fooitems); err != nil {
		return err
	}
	for i := range foo.Enum {
		foo.Enum[i] = fooitems.Enum[i]
	}

	*feature.Select = Select(foo)
	return nil

}

func (feature *SelectFormly) MarshalJSON() ([]byte, error) {

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
		Type:     feature.Type,
		Title:    feature.Title,
		Disabled: feature.Disabled,
		Enum:     make([]string, 0),
	}

	foo.Widget.FormlyConfig = make(map[string]any)

	props := SelectProps{Options: make([]OptionLabelValue, 0)}
	for _, opt := range feature.Enum {
		foo.Enum = append(foo.Enum, opt.Tag)
		props.Options = append(props.Options, OptionLabelValue{Label: opt.Title, Value: opt.Tag, Disabled: opt.Disabled})
	}
	if feature.Default != nil {
		props.DefaultValue = feature.Default.(string)
	}

	if feature.InfoUrl != "" {
		props.InfoUrl = feature.InfoUrl
	}

	foo.Widget.FormlyConfig["type"] = "select"
	foo.Widget.FormlyConfig["multiple"] = "false"
	foo.Widget.FormlyConfig["props"] = props

	valueBytes, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}

	return valueBytes, nil
}
