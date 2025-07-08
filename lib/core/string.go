package core

import (
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type String struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Value      string `json:"value" bexpr:"value"`
	Disabled   bool   `json:"disabled,omitempty"`
	Default    string `json:"default,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Condition  string `json:"condition,omitempty"`
	DisabledOn string `json:"disabled_on,omitempty"`
	InfoUrl    string `json:"info_url,omitempty"`

	valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`
}

func (feature *String) Validate(tag string, i any) (changed bool, err error) {
	if feature.Tag != "" && feature.valueEvaluator != nil {
		result, eval_err := feature.valueEvaluator.Evaluate(i)
		if eval_err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, eval_err)
		}
		b, _ := result.(string)
		//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
		changed = changed || feature.Value != b
		feature.Value = b
	}
	if len(feature.DisabledOn) > 0 {
		result, eval_err := feature.disabledEvaluator.Evaluate(i)
		if eval_err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.DisabledOn, eval_err)
		}
		b, _ := bexpr.CoerceBool(result)
		//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
		changed = changed || feature.Disabled != b
		feature.Disabled = b
	}
	return
}

func (feature *String) ApplyDefaults() error {
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

func (feature *String) Set(tag string, value any) error {
	if feature.Tag == tag {
		return feature.SetValue(value)
	}
	return nil
}

func (feature *String) ApplicableFeatures() []Feature {
	f := make([]Feature, 0)
	if !feature.Disabled {
		f = append(f, feature)
	}
	return f
}

func (feature *String) GetChildren() []Feature {
	return []Feature{}
}

func (feature *String) GetValue() any {
	return feature.Value
}
func (feature *String) SetValue(value any) error {
	feature.Value = value.(string)
	return nil
}
func (feature *String) GetTag() string {
	return feature.Tag
}
func (feature *String) SetTag(tag string) {
	feature.Tag = tag
}
func (feature *String) IsDisabled() bool {
	return feature.Disabled
}

func (feature *String) GetType() string {
	return feature.Type
}
func (feature *String) GetTitle() string {
	return feature.Title
}

// InfoUrl implements FeatureWithInfoUrl.
func (feature *String) GetInfoUrl() string {
	return feature.InfoUrl
}

func (p *String) MarshalJSON() ([]byte, error) {
	marshaller := &StringFormly{p}
	return marshaller.MarshalJSON()
}
