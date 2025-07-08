package core

import (
	"reflect"
)

type Feature interface {
	GetChildren() []Feature
	GetValue() any
	SetValue(value any) error
	GetTag() string
	SetTag(tag string)
	GetType() string
	GetTitle() string
	IsDisabled() bool
	Validate(tag string, i any) (changed bool, err error)
	ApplyDefaults() error
	Set(tag string, value any) error
	ApplicableFeatures() []Feature
}

type FeatureWithInfoUrl interface {
	GetInfoUrl() string
}

var knownTypes = map[string]reflect.Type{
	"checkform": reflect.TypeOf(Checkform{}),
	"checklist": reflect.TypeOf(Checklist{}),
	"checkbox":  reflect.TypeOf(Checkbox{}),
	"select":    reflect.TypeOf(Select{}),
	"option":    reflect.TypeOf(Option{}),
	"string":    reflect.TypeOf(String{}),
	"number":    reflect.TypeOf(Number{}),
}

var _ Feature = (*Checklist)(nil)
var _ Feature = (*Checkform)(nil)
var _ Feature = (*Checkbox)(nil)
var _ FeatureWithInfoUrl = (*Checkbox)(nil)

var _ Feature = (*Select)(nil)
var _ FeatureWithInfoUrl = (*Select)(nil)

var _ Feature = (*Option)(nil)

var _ Feature = (*String)(nil)
var _ FeatureWithInfoUrl = (*String)(nil)

var _ Feature = (*Number)(nil)
var _ FeatureWithInfoUrl = (*Number)(nil)

type OptionLabelValue struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled,omitempty"`
	InfoUrl  string `json:"info_url,omitempty"`
}

type ChecklistProps struct {
	Disabled     bool               `json:"disabled,omitempty"`
	Multiple     bool               `json:"multiple,omitempty"`
	DefaultValue string             `json:"defaultValue,omitempty"`
	HideDisabled bool               `json:"hide_disabled,omitempty"`
	InfoUrls     []string           `json:"info_urls,omitempty"`
	Options      []OptionLabelValue `json:"options,omitempty"`
}

/*
type _Checklist Checklist

	func (feature *Checklist) UnmarshalJSON(bytes []byte) (err error) {
		foo := _Checklist{}
		if err := json.Unmarshal(bytes, &foo); err != nil {
			return err
		}

		type _FooChildren struct {
			Enum []*Checkbox `json:"enum" bexpr:"enum"`
		}

		fooitems := _FooChildren{}
		if err := json.Unmarshal(bytes, &fooitems); err != nil {
			return err
		}
		for i := range foo.Enum {
			foo.Enum[i] = fooitems.Enum[i]
		}

		*feature = Checklist(foo)

		return nil

}
*/
