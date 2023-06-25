package wizard

import "fyne.io/fyne/v2"

// Wizard configuration parameters
const (
	ContentTitleTextSize = 14
	ContentTextSize      = 11
)

type WizardConfig interface {
	GetTitle() string
	OnClose()
	OnFinish() *fyne.Container
}

type BaseWizardConfig struct {
	Title string
}

func (wc *BaseWizardConfig) GetTitle() string {
	return wc.Title
}
