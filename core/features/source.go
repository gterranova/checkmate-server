package features

import "encoding/json"

type SourceFeatures struct {
	BaseFeatures
}

func NewSourceFeatures() SourceFeatures {
	return SourceFeatures{
		BaseFeatures{
			&Feature{
				Title: "Impianto Fotovoltaico",
				Tags:  Photovoltaic,
			},
			&Feature{
				Title: "Impianto Solare Termico",
				Tags:  SolarThermal,
			},
			&Feature{
				Title: "Impianto Eolico",
				Tags:  Wind,
			},
			&Feature{
				Title: "Impianto Idroelettrico",
				Tags:  Hydroelectric,
			},
		},
	}
}

func (features SourceFeatures) Describe() []string {
	ret := make([]string, 0)
	for _, cond := range features.ApplicableFeatures() {
		ret = append(ret, cond.Title)
	}
	return ret
}

func (features SourceFeatures) MarshalJSON() ([]byte, error) {
	return json.Marshal(features.BaseFeatures)
}
