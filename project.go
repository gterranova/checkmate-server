package main

import (
	"fmt"
	"strings"
)

type Project struct {
	capacityMW                    float64
	environmentalProcedures       []EnvironmentalProcedure
	locationConditions            LocationConditions
	suitabilityConditions         SuitabilityConditions
	nationalGuidelinesConditions  NationalGuidelinesConditions
	screeningGuidelinesConditions ScreeningGuidelinesConditions
}

func NewProject() *Project {
	p := &Project{
		capacityMW:                    1,
		environmentalProcedures:       make([]EnvironmentalProcedure, 0),
		locationConditions:            NewLocationConditions(),
		suitabilityConditions:         NewSuitabilityConditions(),
		nationalGuidelinesConditions:  NewNationalGuidelinesConditions(),
		screeningGuidelinesConditions: NewScreeningGuidelinesConditions(),
	}
	p.AddEnvironmentalProcedure(NewEIA(p), NewScreening(p))
	return p
}

func (p *Project) SetPower(capacityMW float64) {
	p.capacityMW = capacityMW
}

func (p *Project) AddEnvironmentalProcedure(procedure ...EnvironmentalProcedure) {
	p.environmentalProcedures = append(p.environmentalProcedures, procedure...)
}

func (p *Project) Art22bisArea() ConditionValue {
	return p.locationConditions.Art22bisArea()
}

func (p *Project) SuitableArea() ConditionValue {
	return p.suitabilityConditions.SuitableArea()
}

func (p *Project) NonUnsuitableDM2010Area() ConditionValue {
	return p.nationalGuidelinesConditions.NonUnsuitableDM2010Area()
}

func (p *Project) CanElevateThresholds() bool {
	return p.Art22bisArea() != NotApplicable || p.SuitableArea() != NotApplicable || p.NonUnsuitableDM2010Area() == NotApplicable
}

func (p *Project) ApplicableEnvironmentalProcedure() EnvironmentalProcedure {
	for _, c := range p.environmentalProcedures {
		if c.IsApplicable() {
			return c
		}
	}
	return nil
}

func (p *Project) Describe() string {

	ret := make([]string, 0)
	if loc := p.locationConditions.ApplicableCondition(); loc != nil {
		ret = append(ret, fmt.Sprintf("Progetto della potenza nominale di **%v MW**, situato in **%v**.  ", p.capacityMW, loc.Title))
	} else {
		ret = append(ret, fmt.Sprintf("Progetto della potenza nominale di **%v MW**.  ", p.capacityMW))
	}

	if p.CanElevateThresholds() {

		switch {
		case p.Art22bisArea() == Art22bis:
			ret = append(ret, p.locationConditions.Describe()...)

			ret = append(ret, `Il Progetto **può beneficiare dell'elevazione delle soglie** per lo Screening e la VIA ai sensi dell'art. 47, comma 11-bis del decreto 13/2023 secondo cui le soglie per la VIA statale e verifica di assoggettabilità/screening " *sono rispettivamente fissati a 20 MW e 10 MW, purché:*`)
			ret = append(ret, `[...] *b)  l'impianto si trovi nelle aree di cui all'articolo 22-bis del decreto legislativo 8 novembre 2021, n. 199;* **[Nota: aree industriali, discariche e cave]** [...]"`)

		case p.SuitableArea() != NotApplicable:
			ret = append(ret, p.suitabilityConditions.Describe()...)

			ret = append(ret, `Il Progetto **può beneficiare dell'elevazione delle soglie** per lo Screening e la VIA ai sensi dell'art. 47, comma 11-bis del decreto 13/2023 secondo cui le soglie per la VIA statale e verifica di assoggettabilità/screening " *sono rispettivamente fissati a 20 MW e 10 MW, purché:*`)
			ret = append(ret, `*a) l'impianto si trovi nelle aree classificate idonee ai sensi dell'articolo 20 del decreto legislativo 8 novembre 2021, n. 199, ivi comprese le aree di cui al comma 8 del medesimo articolo 20;* [...]"`)

		case p.NonUnsuitableDM2010Area() == NotApplicable:
			ret = append(ret, p.nationalGuidelinesConditions.Describe()...)

			ret = append(ret, `Il Progetto **può beneficiare dell'elevazione delle soglie** per lo Screening e la VIA ai sensi dell'art. 47, comma 11-bis del decreto 13/2023 secondo cui le soglie per la VIA statale e verifica di assoggettabilità/screening " *sono rispettivamente fissati a 20 MW e 10 MW, purché:*`)

			ret = append(ret, `[...] *c)  fuori dei casi di cui alle lettere a) e b), l'impianto non sia situato all'interno di aree comprese tra quelle specificamente elencate e individuate ai sensi della lettera f) dell'allegato 3 annesso al decreto del Ministro dello sviluppo economico 10 settembre 2010, pubblicato nella Gazzetta Ufficiale n. 219 del 18 settembre 2010*".`)
		}
	} else {
		ret = append(ret, p.suitabilityConditions.Describe()...)
		if p.NonUnsuitableDM2010Area() != NotApplicable {
			ret = append(ret, p.nationalGuidelinesConditions.Describe()...)
		}

		ret = append(ret, `Il Progetto non soddisfa alcuno dei casi previsti dall'art. 47, comma 11-bis del decreto 13/2023 per i quali è prevista l'elevazione delle soglie per lo Screening e per la VIA.`)
	}

	procedure := p.ApplicableEnvironmentalProcedure()
	screeningThreshold := p.environmentalProcedures[1].ApplicableThreshold()
	if procedure != nil && procedure.Name() != EIAName {
		ret = append(ret, p.screeningGuidelinesConditions.Describe()...)
		ret = append(ret, fmt.Sprintf("**%v**", p.Evaluate()))

	} else if procedure == nil && p.capacityMW >= screeningThreshold*0.5 {
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
