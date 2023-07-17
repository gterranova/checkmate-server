package wizard

import (
	"fyne.io/fyne/v2"
)

// Defines the minimum behavior to a step
type WizardStep interface {
	GetTitle() string
	GetCaption() string
	OnLoad() *fyne.Container
	OnValidate() bool
	OnLeave()
	Disabled() bool
	Refresh()
}

// This method controls the back button behaviour
func (w *Wizard) Back() {
	if w.currentStep <= 0 {
		return
	}
	w.Steps[w.currentStep].OnLeave()
	w.currentStep--
	if w.Steps[w.currentStep].Disabled() {
		w.Back()
		return
	}
	w.refreshContainers()
}

// This method controls the next button behaviour
func (w *Wizard) Next() {
	totalSteps := len(w.Steps)
	if w.currentStep >= totalSteps-1 {
		return
	}
	if !w.Steps[w.currentStep].OnValidate() {
		return
	}
	w.Steps[w.currentStep].OnLeave()
	w.currentStep++
	if w.Steps[w.currentStep].Disabled() {
		w.Next()
		return
	}
	w.refreshContainers()
}

// This method controls the finish button behaviour
func (w *Wizard) Finish() {
	if w.currentStep != len(w.Steps)-1 {
		return
	}
	if !w.Steps[w.currentStep].OnValidate() {
		return
	}
	w.Steps[w.currentStep].OnLeave()
	w.currentStep++
	w.refreshButtonStatus()
	w.buildFinishContainer(w.config.OnFinish())
	w.CancelButton.SetText(CLOSE)
}

// Refresh all the containers
func (w *Wizard) refreshContainers() {
	w.refreshButtonStatus()
	w.rebuildTaskContainer()
	w.rebuildStepsContainer()
	//w.rebuildTitleContainer()
}

// Refresh all steps buttons based on current step
func (w *Wizard) refreshButtonStatus() {
	w.refreshBackButton()
	w.refreshNextButton()
	if w.config.CanFinish() {
		w.refreshFinishButton()
	}
}

// Enable or Disable Back button based on current step
func (w *Wizard) refreshBackButton() {
	if w.currentStep <= 0 || w.currentStep > len(w.Steps)-1 {
		w.BackButton.Disable()
	} else {
		w.BackButton.Enable()
	}
}

// Enable or Disable Next button based on current step
func (w *Wizard) refreshNextButton() {
	if w.currentStep >= len(w.Steps)-1 {
		w.NextButton.Disable()
	} else {
		w.NextButton.Enable()
	}
}

// Enable or Disable Finish button based on current step
func (w *Wizard) refreshFinishButton() {
	if w.currentStep == len(w.Steps)-1 {
		w.FinishButton.Enable()
	} else {
		w.FinishButton.Disable()
	}
}
