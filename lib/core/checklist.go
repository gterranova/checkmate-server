package core

import (
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type Checklist struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Value    any    `json:"value" bexpr:"value"`
	Disabled bool   `json:"disabled,omitempty"`
	Default  any    `json:"default,omitempty"`
	Tag      string `json:"tag,omitempty"`
	//Condition         string           `json:"condition,omitempty"`
	DisabledOn   string `json:"disabled_on,omitempty"`
	HideDisabled bool   `json:"hide_disabled,omitempty"`
	//valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`

	//Properties []any `json:"properties" bexpr:"properties"`
	Enum []*Checkbox `json:"enum" bexpr:"enum"`
}

func (feature *Checklist) Validate(tag string, i any) (changed bool, err error) {
	if feature.Tag != "" && feature.Tag == tag {
		return
	}

	for _, f := range feature.GetChildren() {
		subchanged, _ := f.Validate(tag, i)
		changed = changed || subchanged
	}
	/*
		if feature.Tag != "" && feature.valueEvaluator != nil {
			result, eval_err := feature.valueEvaluator.Evaluate(i)
			if eval_err != nil {
				return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, eval_err)
			}
			b, _ := bexpr.CoerceBool(result)
			fmt.Printf("Result of expression %q evaluation: %s=%t\n", feature.Condition, feature.Tag, result)
			changed = changed || feature.Value != b
			feature.Value = b
		}
	*/
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
	//fmt.Printf("Checklist %s: changed=%t, disabled=%t, childrenDisabled=%t\n", feature.Tag, changed, feature.Disabled, allChildrenDisabled)
	return
}

func (feature *Checklist) ApplyDefaults() error {
	// apply defaults
	/*
		if feature.Condition != "" {
			eval, err := bexpr.CreateEvaluator(feature.Condition)
			if err != nil {
				return fmt.Errorf("failed to create evaluator for expression %q: %v", feature.Condition, err)
			} else {
				feature.valueEvaluator = eval
			}
		}
	*/
	if feature.DisabledOn != "" {
		eval, err := bexpr.CreateEvaluator(feature.DisabledOn)
		if err != nil {
			return fmt.Errorf("failed to create evaluator for expression %q: %v", feature.DisabledOn, err)
		} else {
			feature.disabledEvaluator = eval
		}
	}

	for _, f := range feature.GetChildren() {
		if err := f.ApplyDefaults(); err != nil {
			return fmt.Errorf("failed to apply defaults for %q: %v", f.GetTag(), err)
		}
	}
	return nil
}

func (feature *Checklist) Set(tag string, value any) error {
	for _, f := range feature.GetChildren() {
		if err := f.Set(tag, value); err != nil {
			return err
		}
	}

	if feature.Tag == tag {
		return feature.SetValue(value)
	}
	return nil
}

func (feature *Checklist) ApplicableFeatures() []Feature {
	f := make([]Feature, 0)
	if feature.IsDisabled() {
		return f
	}
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

func (feature *Checklist) GetChildren() []Feature {
	feats := make([]Feature, 0)
	for _, v := range feature.Enum {
		feats = append(feats, v)
	}
	return feats
}

func (feature *Checklist) GetValue() any {
	values := make([]any, 0)
	//if feature.Disabled {
	//	return values
	//}
	for _, checkbox := range feature.GetChildren() {
		if checkbox.GetValue().(bool) {
			values = append(values, checkbox.GetTag())
		}
	}
	return values
}

func (feature *Checklist) SetValue(value any) error {
	switch t := value.(type) {
	case []any:
		for _, checkbox := range feature.Enum {
			if err := checkbox.SetValue(false); err != nil {
				return err
			}
		}
		for _, v := range t {
			tag := v.(string)
			for _, checkbox := range feature.Enum {
				if checkbox.GetTag() == tag {
					//fmt.Printf("setting check %s = %v\n", checkbox.GetTag(), true)
					if err := checkbox.SetValue(true); err != nil {
						return err
					}
				}
			}
		}
	default:
		return fmt.Errorf("unknown type %T for %v", t, t)
	}
	return nil
}

func (feature *Checklist) GetTag() string {
	return feature.Tag
}
func (feature *Checklist) SetTag(tag string) {
	feature.Tag = tag
}

func (feature *Checklist) IsDisabled() bool {
	return feature.Disabled
}

func (feature *Checklist) GetType() string {
	return feature.Type
}
func (feature *Checklist) GetTitle() string {
	return feature.Title
}

//func (feature *Checklist) UnmarshalJSON(bytes []byte) (err error) {
//	marshaller := &ChecklistFormly{feature}
//	return marshaller.UnmarshalJSON(bytes)
//}

func (p *Checklist) MarshalJSON() ([]byte, error) {
	marshaller := &ChecklistFormly{p}
	return marshaller.MarshalJSON()
}
