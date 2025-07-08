package core

import (
	"encoding/json"
	"reflect"
)

type CheckformFormly struct {
	*Checkform
}

func (p *CheckformFormly) MarshalJSON() ([]byte, error) {

	type jsonCheckform struct {
		Type     string `json:"type"`
		Title    string `json:"title,omitempty"`
		Disabled bool   `json:"disabled,omitempty"`
		Default  any    `json:"default,omitempty"`

		Properties   map[string]Feature `json:"properties" bexpr:"properties"`
		FeatureOrder []string           `json:"feature_order"`
		Widget       struct {
			FormlyConfig map[string]any `json:"formlyConfig,omitempty"`
		} `json:"widget,omitempty"`
	}
	foo := jsonCheckform{
		Type:         "object",
		Title:        p.Title,
		Disabled:     p.Disabled,
		Properties:   p.Properties,
		FeatureOrder: p.FeatureOrder,
	}
	foo.Widget.FormlyConfig = make(map[string]any)
	foo.Widget.FormlyConfig["type"] = p.Type
	foo.Widget.FormlyConfig["props"] = map[string]any{
		"hide_disabled": p.HideDisabled,
	}

	valueBytes, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}

	return valueBytes, nil
}

func (feature *CheckformFormly) UnmarshalJSON(bytes []byte) (err error) {
	type _Checkform Checkform
	foo := _Checkform{}
	if err := json.Unmarshal(bytes, &foo); err != nil {
		return err
	}

	type _ObjectProperties struct {
		Properties map[string]struct {
			Type string `json:"type"`
		} `json:"properties"`
	}
	pf := _ObjectProperties{}
	if err := json.Unmarshal(bytes, &pf); err != nil {
		return err
	}

	type _ObjectPropertiesInterface struct {
		Properties map[string]any `json:"properties"`
	}
	pfi := _ObjectPropertiesInterface{}
	if err := json.Unmarshal(bytes, &pfi); err != nil {
		return err
	}

	foo.Properties = make(map[string]Feature)

	for key, feature := range pf.Properties {
		value := reflect.New(knownTypes[feature.Type]).Interface().(Feature)

		valueBytes, err := json.Marshal(pfi.Properties[key])
		if err != nil {
			return err
		}

		if err = json.Unmarshal(valueBytes, &value); err != nil {
			return err
		}
		foo.Properties[key] = value
		// foo.Properties[key].SetTag(key)
	}
	foo.Widget.FormlyConfig = make(map[string]any)

	props := ChecklistProps{
		Disabled:     feature.Disabled,
		HideDisabled: true, //feature.HideDisabled,
	}
	foo.Widget.FormlyConfig["type"] = "object"
	foo.Widget.FormlyConfig["props"] = props

	if foo.FeatureOrder != nil {
		// Integrate FeatureOrder
		for k := range foo.Properties {
			found := false
			for _, l := range foo.FeatureOrder {
				if l == k {
					found = true
					break
				}
			}
			if !found {
				foo.FeatureOrder = append(foo.FeatureOrder, k)
			}
		}
		//fmt.Println(foo.FeatureOrder)
	}
	*feature.Checkform = Checkform(foo)

	return nil

}
