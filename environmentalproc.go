package main

type Threshold float64

const (
	EIAName       = "Valutazione di Impatto Ambientale (VIA) statale"
	ScreeningName = "Verifica di assoggettabilità a VIA (screening)"

	EIAThreshold       = 10.00
	ScreeningThreshold = 1.00

	ElevatedEIAThreshold       = 20.00
	ElevatedScreeningThreshold = 10.00
)

type EnvironmentalProcedure interface {
	Name() string
	ApplicableThreshold() float64
	IsApplicable() bool
}

type EIARequirements struct {
	project *Project
}

func NewEIA(p *Project) *EIARequirements {
	return &EIARequirements{project: p}
}

func (e *EIARequirements) Name() string {
	return EIAName
}

func (e *EIARequirements) ApplicableThreshold() float64 {
	if e.project.CanElevateThresholds() {
		return ElevatedEIAThreshold
	}
	return EIAThreshold
}

func (e *EIARequirements) IsApplicable() bool {
	return e.project.capacityMW >= e.ApplicableThreshold()
}

type ScreeningRequirements struct {
	project *Project
}

func NewScreening(p *Project) *ScreeningRequirements {
	return &ScreeningRequirements{project: p}
}

func (e *ScreeningRequirements) Name() string {
	return ScreeningName
}

func (e *ScreeningRequirements) ApplicableThreshold() float64 {
	refThreshold := ScreeningThreshold
	if e.project.CanElevateThresholds() {
		refThreshold = ElevatedScreeningThreshold
	}
	if e.project.screeningGuidelinesConditions.ShouldHalveThreshold() != NotApplicable {
		refThreshold *= 0.5
	}
	return refThreshold
}

func (e *ScreeningRequirements) IsApplicable() bool {
	return e.project.capacityMW >= e.ApplicableThreshold()
}
