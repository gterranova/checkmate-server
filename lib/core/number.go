package core

import (
	"fmt"
	"strconv"

	"github.com/gterranova/go-bexpr"
)

type Number struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Value      int64  `json:"value" bexpr:"value"`
	Disabled   bool   `json:"disabled,omitempty"`
	Default    int64  `json:"default,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Condition  string `json:"condition,omitempty"`
	DisabledOn string `json:"disabled_on,omitempty"`
	InfoUrl    string `json:"info_url,omitempty"`

	valueEvaluator    *bexpr.Evaluator `json:"-" bexpr:"-"`
	disabledEvaluator *bexpr.Evaluator `json:"-" bexpr:"-"`
}

func (feature *Number) Validate(tag string, i any) (changed bool, err error) {
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
	if changed {
		eval, err := bexpr.CreateEvaluator(feature.Condition)
		if err != nil {
			return changed, fmt.Errorf("failed to create evaluator for expression %q: %v", feature.Condition, err)
		} else {
			feature.valueEvaluator = eval
		}
	}
	if !feature.Disabled && feature.Tag != "" && feature.valueEvaluator != nil {
		result, eval_err := feature.valueEvaluator.Evaluate(i)
		if eval_err != nil {
			return changed, fmt.Errorf("failed to run evaluation of expression %q: %v", feature.Condition, eval_err)
		}
		b, _ := bexpr.CoerceInt64(result)
		//fmt.Printf("Result of expression %q evaluation: %t\n", feature.Condition, result)
		changed = changed || feature.Value != b
		feature.Value = b
	}
	return
}

func (feature *Number) ApplyDefaults() error {
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

func (feature *Number) Set(tag string, value any) error {
	if feature.Tag == tag {
		return feature.SetValue(value)
	}
	return nil
}

func (feature *Number) ApplicableFeatures() []Feature {
	f := make([]Feature, 0)
	if !feature.Disabled {
		f = append(f, feature)
	}
	return f
}

func (feature *Number) GetChildren() []Feature {
	return []Feature{}
}

func (feature *Number) GetValue() any {
	return feature.Value
}
func (feature *Number) SetValue(value any) error {
	switch v := value.(type) {
	case int64:
		feature.Value = v
	case int:
		if i64, err := strconv.ParseInt(strconv.Itoa(v), 10, 64); err == nil {
			feature.Value = i64
		} else {
			fmt.Println(i64, "is not an integer.")
		}
	case float64:
		if i64, err := strconv.ParseInt(fmt.Sprintf("%.0f", v), 10, 64); err == nil {
			feature.Value = i64
		}
	default:
		fmt.Printf("cannot convert number to int64 %v %T", value, value)
	}
	return nil
}
func (feature *Number) GetTag() string {
	return feature.Tag
}
func (feature *Number) SetTag(tag string) {
	feature.Tag = tag
}
func (feature *Number) IsDisabled() bool {
	return feature.Disabled
}

func (feature *Number) GetType() string {
	return feature.Type
}
func (feature *Number) GetTitle() string {
	return feature.Title
}

// InfoUrl implements FeatureWithInfoUrl.
func (feature *Number) GetInfoUrl() string {
	return feature.InfoUrl
}

func (p *Number) MarshalJSON() ([]byte, error) {
	marshaller := &NumberFormly{p}
	return marshaller.MarshalJSON()
}
