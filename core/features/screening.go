package features

import (
	"encoding/json"
	"fmt"
)

type ScreeningGuidelinesFeatures struct {
	BaseFeatures
}

func NewScreeningGuidelinesFeatures() ScreeningGuidelinesFeatures {
	return ScreeningGuidelinesFeatures{
		BaseFeatures{
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.1. Cumulo con altri progetti",
				Description: ``,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.2. Rischio di incidenti, per quanto riguarda, in particolare, le sostanze o le tecnologie utilizzate",
				Description: ``,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.1. Zone umide",
				Description: `zone umide di importanza internazionale (Ramsar)`,
				Tags:        Ramsar,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.2. Zone costiere",
				Description: `vincoli di cui al Codice dei beni culturali e del paesaggio (art. 142) - Aree di rispetto coste e corpi idrici`,
				Tags:        Art142Dlgs2004Coste,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.3. Zone montuose e forestali",
				Description: `vincoli di cui al Codice dei beni culturali e del paesaggio (art. 142) - Montagne oltre 1600 o 1200 metri; piano forestale regionale/provinciale; in assenza di piano forestale vedi vincoli di cui al Codice dei beni culturali e del paesaggio (art. 142) - Boschi`,
				Tags:        Art142Dlgs2004MontagneBoschi,
			},

			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.4. Riserve e parchi naturali, zone classificate o protette ai sensi della normativa nazionale",
				Description: `Elenco ufficiale aree naturali protette (EUAP)`,
				Tags:        EUAP,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.5. Zone protette speciali designate ai sensi delle direttive 2009/147/CE e 92/43/CEE",
				Description: `Siti di importanza comunitaria (SIC), Zone di protezione speciale (ZPS)`,
				Tags:        ReteNatura2000,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.6. Zone nelle quali gli standard di qualità ambientale fissati dalla normativa dell'Unione europea sono già stati superati",
				Description: `dati di qualità delle acque superficiali e sotterranee`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.7. Zone a forte densità demografica",
				Description: `zone a forte densità demografica si intendono i centri abitati, così come delimitati dagli strumenti urbanistici comunali, posti all'interno dei territori comunali con densità superiore a 500 abitanti per km2 e popolazione di almeno 50.000 abitanti (EUROSTAT)`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Linee Guida Screening DM 30/03/2015",
				Title:       "4.3.8. Zone di importanza storica, culturale o archeologica",
				Description: `immobili e le aree di cui all'art. 136 del Codice dei beni culturali e del paesaggio di cui al decreto legislativo n. 42/2004 dichiarati di notevole interesse pubblico ai sensi dell'art. 140 del medesimo decreto e gli immobili e le aree di interesse artistico, storico, archeologico o etnoantropologico di cui all'art. 10, comma 3, lettera a), del medesimo decreto`,
				Tags:        DM10sett2010,
			},
		},
	}
}

func (features ScreeningGuidelinesFeatures) Describe() []string {
	ret := make([]string, 0)
	applicableFeatures := features.ApplicableFeatures()
	if len(applicableFeatures) == 0 {
		return ret
	}
	ret = append(ret, "Ai sensi delle Linee Guida per la verifica di assoggettabilità a VIA, Allegate al DM 30/03/2015, la **soglia per lo screening è soggetta a dimezzamento** in applicazione dei seguenti criteri:")

	for _, cond := range applicableFeatures {
		if len(cond.Description) > 0 {
			ret = append(ret, fmt.Sprintf("- %v (%v)", cond.Title, cond.Description))
		} else {
			ret = append(ret, fmt.Sprintf("- %v", cond.Title))
		}
	}
	return ret
}

func (features ScreeningGuidelinesFeatures) MarshalJSON() ([]byte, error) {
	return json.Marshal(features.BaseFeatures)
}
