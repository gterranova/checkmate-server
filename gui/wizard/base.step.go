package wizard

import (
	"fyne.io/fyne/v2"
)

// Base widget is a basic, unvalidated implementation of WizardStep
// Can to be composed in the final Steps
type BaseWizardStep struct {
	Caption string
	Title   string
	Content *fyne.Container
}

func (s *BaseWizardStep) GetTitle() string {
	return s.Title
}

func (s *BaseWizardStep) GetCaption() string {
	return s.Caption
}

func (s *BaseWizardStep) OnLoad() *fyne.Container {
	return s.Content
}

func (s *BaseWizardStep) OnValidate() bool {
	return true
}

func (s *BaseWizardStep) OnLeave() {
}
