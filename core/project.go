package core

import (
	"fmt"
	"strings"

	"terra9.it/vadovia/core/features"
)

type Project struct {
	CapacityMW              float64                        `json:"capacityMW"`
	AuthorizationProcedures []AuthorizationProcedure       `json:"-"`
	EnvironmentalProcedures []EnvironmentalProcedure       `json:"-"`
	Features                map[string]features.FeatureSet `json:"features"`
}

func NewProject() *Project {
	p := &Project{
		CapacityMW:              1,
		AuthorizationProcedures: make([]AuthorizationProcedure, 0),
		EnvironmentalProcedures: make([]EnvironmentalProcedure, 0),
		Features:                make(map[string]features.FeatureSet),
	}
	p.Features[features.SOURCE] = features.NewSourceFeatures()
	p.Features[features.LOCATION] = features.NewLocationFeatures()
	p.Features[features.SUITABILITY] = features.NewSuitabilityFeatures()
	p.Features[features.NATIONAL_GUIDELINES] = features.NewNationalGuidelinesFeatures()
	p.Features[features.SCREENING_GUIDELINES] = features.NewScreeningGuidelinesFeatures()

	p.AddAuthorizationProcedure(
		NewFreeBuildingRequirements(p),
		NewDILARequirements(p),
		NewPASRequirements(p),
		NewAURequirements(p),
	)
	p.AddEnvironmentalProcedure(NewEIA(p), NewScreening(p))
	return p
}

func (p *Project) SetPower(capacityMW float64) {
	p.CapacityMW = capacityMW
}

func (p *Project) AddAuthorizationProcedure(procedure ...AuthorizationProcedure) {
	p.AuthorizationProcedures = append(p.AuthorizationProcedures, procedure...)
}

func (p *Project) AddEnvironmentalProcedure(procedure ...EnvironmentalProcedure) {
	p.EnvironmentalProcedures = append(p.EnvironmentalProcedures, procedure...)
}

func (p *Project) Art22bisArea() bool {
	return p.Features[features.LOCATION].Has(features.Art22bis)
}

func (p *Project) SuitableArea() bool {
	return p.Features[features.SUITABILITY].Any()
}

func (p *Project) UnsuitableDM2010Area() bool {
	return p.Features[features.NATIONAL_GUIDELINES].Any()
}

func (p *Project) CanElevateThresholds() bool {
	return p.Art22bisArea() || p.SuitableArea() || !p.UnsuitableDM2010Area()
}

func (p *Project) ApplicableAuthorizationProcedure() AuthorizationProcedure {
	for _, c := range p.AuthorizationProcedures {
		if c.IsApplicable() {
			return c
		}
	}
	return nil
}

func (p *Project) ApplicableEnvironmentalProcedure() EnvironmentalProcedure {
	for _, c := range p.EnvironmentalProcedures {
		if c.IsApplicable() {
			return c
		}
	}
	return nil
}

func (p *Project) Describe() string {

	ret := make([]string, 0)
	if loc := p.Features[features.LOCATION].ApplicableFeatures(); len(loc) > 0 {
		ret = append(ret, fmt.Sprintf("Progetto della potenza nominale di **%v MW**, situato in **%v**.  ", p.CapacityMW, loc[0].Title))
	} else {
		ret = append(ret, fmt.Sprintf("Progetto della potenza nominale di **%v MW**.  ", p.CapacityMW))
	}

	if p.CanElevateThresholds() {

		switch {
		case p.Art22bisArea():
			ret = append(ret, p.Features[features.LOCATION].Describe()...)

			ret = append(ret, `Il Progetto **può beneficiare dell'elevazione delle soglie** per lo Screening e la VIA ai sensi dell'art. 47, comma 11-bis del decreto 13/2023 secondo cui le soglie per la VIA statale e verifica di assoggettabilità/screening " *sono rispettivamente fissati a 20 MW e 10 MW, purché:*`)
			ret = append(ret, `[...] *b)  l'impianto si trovi nelle aree di cui all'articolo 22-bis del decreto legislativo 8 novembre 2021, n. 199;* **[Nota: aree industriali, discariche e cave]** [...]"`)

		case p.SuitableArea():
			ret = append(ret, p.Features[features.SUITABILITY].Describe()...)

			ret = append(ret, `Il Progetto **può beneficiare dell'elevazione delle soglie** per lo Screening e la VIA ai sensi dell'art. 47, comma 11-bis del decreto 13/2023 secondo cui le soglie per la VIA statale e verifica di assoggettabilità/screening " *sono rispettivamente fissati a 20 MW e 10 MW, purché:*`)
			ret = append(ret, `*a) l'impianto si trovi nelle aree classificate idonee ai sensi dell'articolo 20 del decreto legislativo 8 novembre 2021, n. 199, ivi comprese le aree di cui al comma 8 del medesimo articolo 20;* [...]"`)

		case !p.UnsuitableDM2010Area():
			ret = append(ret, p.Features[features.NATIONAL_GUIDELINES].Describe()...)

			ret = append(ret, `Il Progetto **può beneficiare dell'elevazione delle soglie** per lo Screening e la VIA ai sensi dell'art. 47, comma 11-bis del decreto 13/2023 secondo cui le soglie per la VIA statale e verifica di assoggettabilità/screening " *sono rispettivamente fissati a 20 MW e 10 MW, purché:*`)

			ret = append(ret, `[...] *c)  fuori dei casi di cui alle lettere a) e b), l'impianto non sia situato all'interno di aree comprese tra quelle specificamente elencate e individuate ai sensi della lettera f) dell'allegato 3 annesso al decreto del Ministro dello sviluppo economico 10 settembre 2010, pubblicato nella Gazzetta Ufficiale n. 219 del 18 settembre 2010*".`)
		}
	} else {
		ret = append(ret, p.Features[features.SUITABILITY].Describe()...)
		if p.UnsuitableDM2010Area() {
			ret = append(ret, p.Features[features.NATIONAL_GUIDELINES].Describe()...)
		}

		ret = append(ret, `Il Progetto non soddisfa alcuno dei casi previsti dall'art. 47, comma 11-bis del decreto 13/2023 per i quali è prevista l'elevazione delle soglie per lo Screening e per la VIA.`)
	}

	procedure := p.ApplicableEnvironmentalProcedure()
	screeningThreshold := p.EnvironmentalProcedures[1].ApplicableThreshold()
	screeningFeatures := p.Features[features.SCREENING_GUIDELINES].(features.ScreeningGuidelinesFeatures)
	if procedure != nil && procedure.Name() != EIAName {
		ret = append(ret, p.Features[features.SCREENING_GUIDELINES].Describe()...)
		ret = append(ret, fmt.Sprintf("**%v**", p.Evaluate()))

	} else if procedure == nil && !screeningFeatures.Any() && p.CapacityMW >= screeningThreshold*0.5 {
		ret = append(ret, fmt.Sprintf("**%v**", p.Evaluate()))

		ret = append(ret, fmt.Sprintf(`**Si raccomanda di verificare che non ricorra alcuno dei criteri previsti dalle Linee Guida per la verifica di assoggettabilità a VIA, Allegate al DM 30/03/2015, in applicazione dei quali la soglia per lo screening sarebbe soggetta a dimezzamento. In tal caso, infatti, per effetto del dimezzamento, la soglia per lo screening diventerebbe %v MW, anziché %v MW, ed il Progetto andrebbe sottoposto a verifica di assoggettabilità a VIA.**`, screeningThreshold*0.5, screeningThreshold))
	} else {
		ret = append(ret, fmt.Sprintf("**%v**", p.Evaluate()))
	}

	return strings.Join(ret, "\n\n")
}

func (p *Project) Evaluate() string {
	if procedure := p.ApplicableEnvironmentalProcedure(); procedure != nil {
		return fmt.Sprintf("Richiede %v, soglia applicabile %v MW.", procedure.Name(), procedure.ApplicableThreshold())
	}
	return "Non richiede procedure ambientali."
}
