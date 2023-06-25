package features

type ConditionValue uint

const (
	Photovoltaic  ConditionValue = iota // impianto solare fotovoltaico
	SolarThermal                        // impianto solare termico
	Wind                                // impianto eolico
	Hydroelectric                       // impianto idroelettrico

	NotApplicable
	Art22bis
	Art20
	DM10sett2010

	EUAP
	Ramsar
	ReteNatura2000
	Art142Dlgs2004
	Art142Dlgs2004Coste
	Art142Dlgs2004MontagneBoschi
)

const (
	SOURCE               = "source"
	LOCATION             = "location"
	SUITABILITY          = "suitability"
	NATIONAL_GUIDELINES  = "national_guidelines"
	SCREENING_GUIDELINES = "screening_guidelines"
)

type Feature struct {
	Title       string         `json:"title"`
	Description string         `json:"-"`
	Reference   string         `json:"-"`
	Checked     bool           `json:"checked"`
	Disabled    bool           `json:"disabled"`
	Tags        ConditionValue `json:"value"`
}

type FeatureSet interface {
	Any() bool
	Has(f ConditionValue) bool
	ApplicableFeatures() []*Feature
	Describe() []string
}

type BaseFeatures []*Feature

func (features BaseFeatures) ApplicableFeatures() []*Feature {
	f := make([]*Feature, 0)
	for _, feature := range features {
		if !feature.Disabled && feature.Checked {
			f = append(f, feature)
		}
	}
	return f
}

func (features BaseFeatures) Has(f ConditionValue) bool {
	for _, feature := range features {
		if !feature.Disabled && feature.Checked && feature.Tags == f {
			return true
		}
	}
	return false
}

func (features BaseFeatures) Any() bool {
	for _, feature := range features {
		if !feature.Disabled && feature.Checked {
			return true
		}
	}
	return false
}
