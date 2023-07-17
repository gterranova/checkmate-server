package features

import (
	"fmt"

	"github.com/gterranova/go-bexpr"
)

type Feature struct {
	Type       string     `json:"type"`
	Title      string     `json:"title"`
	Caption    string     `json:"caption"`
	Reference  string     `json:"reference"`
	Value      any        `json:"value" bexpr:"value"`
	Disabled   bool       `json:"disabled"`
	Default    any        `json:"default"`
	Tag        string     `json:"tag"`
	Condition  string     `json:"condition"`
	DisabledOn string     `json:"disabled_on"`
	Features   []*Feature `json:"features" bexpr:"features"`

	valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`
}

func (feature *Feature) Validate(tag string, i interface{}) (changed bool) {
	if feature.Tag != "" && feature.Tag == tag {
		return
	}

	for _, f := range feature.Features {
		changed = changed || f.Validate(tag, i)
	}
	if feature.Tag != "" && feature.valueEvaluator != nil {
		result, err := feature.valueEvaluator.Evaluate(i)
		if err != nil {
			fmt.Printf("failed to run evaluation of expression %q: %v\n", feature.Condition, err)
			return
		}
		b, _ := bexpr.CoerceBool(result)
		//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
		changed = changed || feature.Value != b
		feature.Value = b
	}
	if len(feature.DisabledOn) > 0 {
		result, err := feature.disabledEvaluator.Evaluate(i)
		if err != nil {
			fmt.Printf("failed to run evaluation of expression %q: %v\n", feature.DisabledOn, err)
			return
		}
		b, _ := bexpr.CoerceBool(result)
		//fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
		changed = changed || feature.Disabled != b
		feature.Disabled = b
	}
	return
}

func (feature *Feature) ApplyDefaults() {
	// apply defaults
	if feature.Condition != "" {
		eval, err := bexpr.CreateEvaluator(feature.Condition)
		if err != nil {
			fmt.Printf("failed to create evaluator for expression %q: %v", feature.Condition, err)
		} else {
			feature.valueEvaluator = eval
		}
	}

	if feature.DisabledOn != "" {
		eval, err := bexpr.CreateEvaluator(feature.DisabledOn)
		if err != nil {
			fmt.Printf("failed to create evaluator for expression %q: %v", feature.DisabledOn, err)
		} else {
			feature.disabledEvaluator = eval
		}
	}

	if feature.Type == "select" && feature.Default != nil {
		for _, f := range feature.Features {
			if f.Tag == feature.Default.(string) {
				f.Value = true
			} else {
				f.Value = false
			}
		}
	} else if feature.Type == "input" && feature.Default != nil {
		feature.Value = feature.Default
	} else {
		for _, f := range feature.Features {
			f.ApplyDefaults()
		}
	}
}

func (feature *Feature) Set(tag string, value interface{}) {
	for _, f := range feature.Features {
		f.Set(tag, value)
	}

	if feature.Tag == tag {
		//fmt.Println(feature.Tag, "found")
		feature.Value = value
	}
}

func (feature *Feature) ApplicableFeatures() []*Feature {
	f := make([]*Feature, 0)
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
	for _, child := range feature.Features {
		if !child.Disabled {
			f = append(f, child.ApplicableFeatures()...)
		}
	}
	return f
}

func (feature *Feature) GetChildren() []*Feature {
	return feature.Features
}

func (feature *Feature) Any() bool {
	if !feature.Disabled && feature.Value.(bool) {
		return true
	}
	for _, child := range feature.Features {
		if !child.Disabled && child.Value != nil && child.Value.(bool) {
			return true
		}
	}
	return false
}
