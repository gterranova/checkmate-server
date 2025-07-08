package core

import (
	"encoding/json"
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type Checkbox struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Value      bool   `json:"value" bexpr:"value"`
	Disabled   bool   `json:"disabled,omitempty"`
	Default    bool   `json:"default,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Condition  string `json:"condition,omitempty"`
	DisabledOn string `json:"disabled_on,omitempty"`
	InfoUrl    string `json:"info_url,omitempty"`

	valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`
}

func (feature *Checkbox) Validate(tag string, i any) (changed bool, err error) {
	var result any
	var b bool

	if feature.Tag != "" && feature.valueEvaluator != nil {
		result, err = feature.valueEvaluator.Evaluate(i)
		if err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, err)
		}
		b, _ = bexpr.CoerceBool(result)
		changed = changed || feature.Value != b
		feature.Value = b
		//if changed {
		//	fmt.Printf("%s: Result of expression %q evaluation: %t\n", feature.Tag, feature.Condition, b)
		//}
	}
	if len(feature.DisabledOn) > 0 {
		if feature.disabledEvaluator == nil {
			fmt.Printf("%s: failed to create evaluator for expression %q\n", feature.Tag, feature.DisabledOn)
			return changed, fmt.Errorf("%s: failed to create evaluator for expression %q", feature.Tag, feature.DisabledOn)
		}
		result, err = feature.disabledEvaluator.Evaluate(i)
		if err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, err)
		}
		if b, err = bexpr.CoerceBool(result); err != nil {
			//fmt.Printf("%s %v %T\n", feature.DisabledOn, err, result)
			//b = false
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, err)
		}
		//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
		changed = changed || feature.Disabled != b
		feature.Disabled = b
	}
	return changed, nil
}

func (feature *Checkbox) ApplyDefaults() error {
	// apply defaults
	if feature.Condition != "" {
		eval, err := bexpr.CreateEvaluator(feature.Condition)
		if err != nil {
			return fmt.Errorf("failed to create evaluator for expression %q: %v", feature.Condition, err)
		} else {
			feature.valueEvaluator = eval
		}
	}

	if feature.DisabledOn != "" {
		eval, err := bexpr.CreateEvaluator(feature.DisabledOn)
		if err != nil {
			return fmt.Errorf("failed to create evaluator for expression %q: %v", feature.DisabledOn, err)
		} else {
			feature.disabledEvaluator = eval
		}
	}

	return feature.SetValue(feature.Default)
}

func (feature *Checkbox) Set(tag string, value any) error {
	if feature.Tag == tag {
		return feature.SetValue(value)
	}
	return nil
}

func (feature *Checkbox) ApplicableFeatures() []Feature {
	f := make([]Feature, 0)
	if !feature.Disabled {
		if t := feature.Value; t {
			f = append(f, feature)
		}
	}
	return f
}

func (feature *Checkbox) GetChildren() []Feature {
	return []Feature{}
}

func (feature *Checkbox) GetValue() any {
	//if feature.Disabled {
	//	return false
	//}
	return feature.Value
}
func (feature *Checkbox) SetValue(value any) error {
	feature.Value = value.(bool)
	return nil
}
func (feature *Checkbox) GetTag() string {
	return feature.Tag
}
func (feature *Checkbox) SetTag(tag string) {
	feature.Tag = tag
}
func (feature *Checkbox) IsDisabled() bool {
	return feature.Disabled
}

func (feature *Checkbox) GetType() string {
	return feature.Type
}
func (feature *Checkbox) GetTitle() string {
	return feature.Title
}

func (feature *Checkbox) GetInfoUrl() string {
	return feature.InfoUrl
}

func (feature *Checkbox) MarshalJSON() ([]byte, error) {
	type jsonFeature struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Value    bool   `json:"value,omitempty"`
		Disabled bool   `json:"disabled,omitempty"`
		Default  bool   `json:"default,omitempty"`
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
