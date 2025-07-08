package core

import (
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type Select struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Value      any    `json:"value" bexpr:"value"`
	Disabled   bool   `json:"disabled,omitempty"`
	Default    any    `json:"default,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Condition  string `json:"condition,omitempty"`
	DisabledOn string `json:"disabled_on,omitempty"`
	InfoUrl    string `json:"info_url,omitempty"`

	valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`

	//Properties []any `json:"properties" bexpr:"properties"`
	Enum []*Option `json:"enum" bexpr:"enum"`
}

func (feature *Select) Validate(tag string, i any) (changed bool, err error) {
	if feature.Tag != "" && feature.Tag == tag {
		return
	}

	for _, f := range feature.GetChildren() {
		subchanged, _ := f.Validate(tag, i)
		changed = changed || subchanged
	}
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

func (feature *Select) ApplyDefaults() error {
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

	for _, f := range feature.GetChildren() {
		if f.GetTag() == feature.Default.(string) {
			if err := f.SetValue(true); err != nil {
				return err
			}
		} else {
			if err := f.SetValue(false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (feature *Select) Set(tag string, value any) error {
	for _, f := range feature.GetChildren() {
		if err := f.Set(tag, value); err != nil {
			return err
		}
	}

	if feature.Tag == tag {
		//fmt.Println(feature.Tag, "found")
		return feature.SetValue(value)
	}
	return nil
}

func (feature *Select) ApplicableFeatures() []Feature {
	f := make([]Feature, 0)
	/*
		if !feature.Disabled {
			switch t := feature.Value.(type) {
			case nil:
			case bool:
				if t {
					f = append(f, feature)
				}
			default:
				f = append(f, feature)
			}
		}
	*/
	for _, child := range feature.GetChildren() {
		if !child.IsDisabled() {
			f = append(f, child.ApplicableFeatures()...)
		}
	}
	return f
}

func (feature *Select) GetChildren() []Feature {
	feats := make([]Feature, 0)
	for _, v := range feature.Enum {
		feats = append(feats, v)
	}
	return feats
}

func (feature *Select) GetValue() any {
	for _, option := range feature.Enum {
		if option.GetValue().(bool) {
			//fmt.Printf("getting select %s = %v (%T)\n", option.GetTag(), option.GetValue(), option.GetValue())
			return option.GetTag()
		}
	}
	return nil
}

func (feature *Select) SetValue(value any) error {
	switch t := value.(type) {
	case string:
		for _, option := range feature.Enum {
			if err := option.SetValue(option.GetTag() == t); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown type %T for %v", t, t)
	}
	return nil
}
func (feature *Select) GetTag() string {
	return feature.Tag
}
func (feature *Select) SetTag(tag string) {
	feature.Tag = tag
}
func (feature *Select) IsDisabled() bool {
	return feature.Disabled
}

func (feature *Select) GetType() string {
	return feature.Type
}
func (feature *Select) GetTitle() string {
	return feature.Title
}

// InfoUrl implements FeatureWithInfoUrl.
func (feature *Select) GetInfoUrl() string {
	return feature.InfoUrl
}

func (feature *Select) UnmarshalJSON(bytes []byte) (err error) {
	marshaller := &SelectFormly{feature}
	return marshaller.UnmarshalJSON(bytes)
}

func (p *Select) MarshalJSON() ([]byte, error) {
	marshaller := &SelectFormly{p}
	return marshaller.MarshalJSON()
}
