package core

import (
	"encoding/json"
)

type StringFormly struct {
	*String
}

func (feature *StringFormly) MarshalJSON() ([]byte, error) {
	type jsonFeature struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Value    string `json:"value,omitempty"`
		Disabled bool   `json:"disabled,omitempty"`
		Default  string `json:"default,omitempty"`
		Widget   struct {
			FormlyConfig map[string]any `json:"formlyConfig,omitempty"`
		} `json:"widget,omitempty"`
	}
	foo := jsonFeature{
		Type:     feature.Type,
		Title:    feature.Title,
		Value:    feature.Value,
		Disabled: feature.Disabled,
		Default:  feature.Default,
	}

	if feature.InfoUrl != "" {
		foo.Widget.FormlyConfig = make(map[string]any)
		foo.Widget.FormlyConfig["props"] = map[string]any{
			"info_url": feature.InfoUrl,
		}
	}

	valueBytes, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}

	return valueBytes, nil
}
