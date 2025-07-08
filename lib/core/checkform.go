package core

import (
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type Checkform struct {
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

	Properties   map[string]Feature `json:"-" bexpr:"properties"`
	FeatureOrder []string           `json:"feature_order"`

	Widget struct {
		FormlyConfig map[string]any `json:"formlyConfig"`
	} `json:"widget"`
}

func (feature *Checkform) Validate(tag string, i any) (changed bool, err error) {
	var result any
	var b bool

	if feature.Tag != "" && feature.Tag == tag {
		return
	}

	for _, f := range feature.GetChildren() {
		subchanged, _ := f.Validate(tag, i)
		changed = changed || subchanged
	}
	/*
		if feature.Tag != "" && feature.valueEvaluator != nil {
			result, err := feature.valueEvaluator.Evaluate(i)
			if err != nil {
				return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, err)
			}
			b, _ := bexpr.CoerceBool(result)
			//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
			changed = changed || feature.Value != b
			feature.Value = b
		}
	*/
	if len(feature.DisabledOn) > 0 {
		result, err = feature.disabledEvaluator.Evaluate(i)
		if err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.DisabledOn, err)
		}
		if b, err = bexpr.CoerceBool(result); err != nil {
			//fmt.Printf("%s %v %T\n", feature.DisabledOn, err, result)
			b = false
		}
		//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
		changed = changed || feature.Disabled != b
		feature.Disabled = b
	}
	return
}

func (feature *Checkform) ApplyDefaults() error {
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

func (feature *Checkform) Set(tag string, value any) error {
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

func (feature *Checkform) ApplicableFeatures() []Feature {
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

func (feature *Checkform) GetChildren() []Feature {
	/*
		featsLabels := make([]string, 0)
		for k := range feature.Properties {
			featsLabels = append(featsLabels, k)
		}
		order := func(a, b string) int {
			posA := IndexOf(feature.FeatureOrder, a)
			posB := IndexOf(feature.FeatureOrder, b)
			if posA != -1 && posB != -1 {
				if IndexOf(feature.FeatureOrder, a) > IndexOf(feature.FeatureOrder, b) {
					return 1
				}
				return -1
			} else if posA != -1 {
				return -1
			} else if posB != -1 {
				return 1
			}
			return 0
		}
		slices.SortFunc(featsLabels, order)

		for i, l := range featsLabels {
			feats[i] = feature.Properties[l]
		}
	*/
	feats := make([]Feature, 0)
	if len(feature.FeatureOrder) > 0 {
		for _, k := range feature.FeatureOrder {
			if v, ok := feature.Properties[k]; ok {
				feats = append(feats, v)
			}
		}
	} else {
		for _, v := range feature.Properties {
			feats = append(feats, v)
		}
	}
	return feats
}

func (feature *Checkform) GetValue() any {
	values := make(map[string]any)
	//if feature.Disabled {
	//	return values
	//}
	for k, control := range feature.Properties {
		values[k] = control.GetValue()
	}
	return values
}

func (feature *Checkform) SetValue(value any) error {
	switch t := value.(type) {
	case map[string]any:
		for k, v := range t {
			if _, ok := feature.Properties[k]; ok {
				//fmt.Printf("setting prop %s = %v\n", k, v)
				if err := feature.Properties[k].SetValue(v); err != nil {
					return fmt.Errorf("failed to set value for %q: %v", k, err)
				}
			}
		}
	default:
		return fmt.Errorf("unknown type %T for %v", t, t)
	}
	return nil
}

func (feature *Checkform) GetTag() string {
	return feature.Tag
}
func (feature *Checkform) SetTag(tag string) {
	feature.Tag = tag
}

func (feature *Checkform) IsDisabled() bool {
	return feature.Disabled
}

func (feature *Checkform) GetType() string {
	return feature.Type
}
func (feature *Checkform) GetTitle() string {
	return feature.Title
}

func (feature *Checkform) UnmarshalJSON(bytes []byte) (err error) {
	marshaller := &CheckformFormly{feature}
	return marshaller.UnmarshalJSON(bytes)
}

func (p *Checkform) MarshalJSON() ([]byte, error) {
	marshaller := &CheckformFormly{p}
	return marshaller.MarshalJSON()
}
