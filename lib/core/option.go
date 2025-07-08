package core

import (
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type Option struct {
	Type              string           `json:"type"`
	Title             string           `json:"title"`
	Value             bool             `json:"value" bexpr:"value"`
	Disabled          bool             `json:"disabled,omitempty"`
	Default           bool             `json:"default,omitempty"`
	Tag               string           `json:"tag,omitempty"`
	Condition         string           `json:"condition,omitempty"`
	DisabledOn        string           `json:"disabled_on,omitempty"`
	valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`
}

func (feature *Option) Validate(tag string, i any) (changed bool, err error) {
	if feature.Tag != "" && feature.valueEvaluator != nil {
		result, eval_err := feature.valueEvaluator.Evaluate(i)
		if eval_err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, eval_err)
		}
		b, _ := bexpr.CoerceBool(result)
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

func (feature *Option) ApplyDefaults() error {
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

func (feature *Option) Set(tag string, value any) error {
	if feature.Tag == tag {
		return feature.SetValue(value)
	}
	return nil
}

func (feature *Option) ApplicableFeatures() []Feature {
	f := make([]Feature, 0)
	if !feature.Disabled {
		if t := feature.Value; t {
			f = append(f, feature)
		}
	}
	return f
}

func (feature *Option) GetChildren() []Feature {
	return []Feature{}
}

func (feature *Option) GetValue() any {
	return feature.Value
}
func (feature *Option) SetValue(value any) error {
	feature.Value = value.(bool)
	return nil
}
func (feature *Option) GetTag() string {
	return feature.Tag
}
func (feature *Option) SetTag(tag string) {
	feature.Tag = tag
}
func (feature *Option) IsDisabled() bool {
	return feature.Disabled
}

func (feature *Option) GetType() string {
	return feature.Type
}
func (feature *Option) GetTitle() string {
	return feature.Title
}
