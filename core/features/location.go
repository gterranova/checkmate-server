package features

import (
	"encoding/json"
	"fmt"
)

type LocationFeatures struct {
	BaseFeatures
}

func NewLocationFeatures() LocationFeatures {
	return LocationFeatures{
		BaseFeatures{
			&Feature{
				Reference:   "",
				Title:       "Aree agricole",
				Description: ``,
				Tags:        NotApplicable,
			},
			&Feature{
				Reference:   "Art. 22-bis, comma 1, D.lgs. 199/2021",
				Title:       "Aree a destinazione industriale, artigianale e commerciale",
				Description: `*L'installazione, con qualunque modalità, di impianti fotovoltaici su terra e delle relative opere connesse e infrastrutture necessarie, ubicati nelle zone e nelle aree a destinazione industriale, artigianale e commerciale* [...] *è considerata attività di manutenzione ordinaria e non è subordinata acquisizione di permessi, autorizzazioni o atti di assenso comunque denominati, fatte salve le valutazioni ambientali di cui al titolo III della parte seconda del decreto legislativo 3 aprile 2006, n. 152, ove previste.*`,
				Tags:        Art22bis,
			},
			&Feature{
				Reference:   "Art. 22-bis, comma 1, D.lgs. 199/2021",
				Title:       "Discariche o lotti di discarica chiusi e ripristinati",
				Description: `*L'installazione, con qualunque modalità, di impianti fotovoltaici su terra e delle relative opere connesse e infrastrutture necessarie, ubicati* [...] *in discariche o lotti di discarica chiusi e ripristinati* [...] *è considerata attività di manutenzione ordinaria e non è subordinata acquisizione di permessi, autorizzazioni o atti di assenso comunque denominati, fatte salve le valutazioni ambientali di cui al titolo III della parte seconda del decreto legislativo 3 aprile 2006, n. 152, ove previste.*`,
				Tags:        Art22bis,
			},
			&Feature{
				Reference:   "Art. 22-bis, comma 1, D.lgs. 199/2021",
				Title:       "Cave o lotti o porzioni di cave non suscettibili di ulteriore sfruttamento",
				Description: `*L'installazione, con qualunque modalità, di impianti fotovoltaici su terra e delle relative opere connesse e infrastrutture necessarie, ubicati* [...] *in cave o lotti o porzioni di cave non suscettibili di ulteriore sfruttamento, è considerata attività di manutenzione ordinaria e non è subordinata acquisizione di permessi, autorizzazioni o atti di assenso comunque denominati, fatte salve le valutazioni ambientali di cui al titolo III della parte seconda del decreto legislativo 3 aprile 2006, n. 152, ove previste.*`,
				Tags:        Art22bis,
			},
			&Feature{
				Reference:   "",
				Title:       "Altro",
				Description: ``,
				Tags:        NotApplicable,
			},
		},
	}
}

func (features LocationFeatures) Describe() []string {
	ret := make([]string, 0)
	for _, feature := range features.ApplicableFeatures() {
		if len(feature.Reference) > 0 {
			ret = append(ret, fmt.Sprintf("Ai sensi dell'%v, \"%v\".", feature.Reference, feature.Description))
		}
	}
	return ret
}

func (features LocationFeatures) MarshalJSON() ([]byte, error) {
	return json.Marshal(features.BaseFeatures)
}
