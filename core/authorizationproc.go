package core

const (
	FreeBuildingName = "Comunicazione relativa alle attività in edilizia libera (art. 6, comma 11 D.lgs.28/2011)"
	DILAName         = "Dichiarazione di Inizio Lavori Asseverata (art. 6-bis D.lgs. 28/2011)"
	PASName          = "Procedura Abilitativa Semplificata (art. 6 D.lgs. 28/2011)"
	AUName           = "Autorizzazione Unica (art. 5 D.lgs. 28/2011)"
)

type AuthorizationProcedure interface {
	Name() string
	IsApplicable() bool
}

type FreeBuildingRequirements struct{ project *Project }

type DILARequirements struct{ project *Project }

type PASRequirements struct{ project *Project }

type AURequirements struct{ project *Project }

func NewFreeBuildingRequirements(p *Project) *FreeBuildingRequirements {
	return &FreeBuildingRequirements{project: p}
}

func NewDILARequirements(p *Project) *DILARequirements {
	return &DILARequirements{project: p}
}

func NewPASRequirements(p *Project) *PASRequirements {
	return &PASRequirements{project: p}
}

func NewAURequirements(p *Project) *AURequirements {
	return &AURequirements{project: p}
}

func (p *FreeBuildingRequirements) Name() string {
	return FreeBuildingName
}

func (p *DILARequirements) Name() string {
	return DILAName
}

func (p *PASRequirements) Name() string {
	return PASName
}

func (p *AURequirements) Name() string {
	return AUName
}

func (p *FreeBuildingRequirements) IsApplicable() bool {
	panic("Not implemented")
}

func (p *DILARequirements) IsApplicable() bool {
	panic("Not implemented")
}

func (p *PASRequirements) IsApplicable() bool {
	panic("Not implemented")
}

func (p *AURequirements) IsApplicable() bool {
	panic("Not implemented")
}
