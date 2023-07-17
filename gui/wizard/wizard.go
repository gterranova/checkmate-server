package wizard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Represents a wizard pane content with logic to navigate between steps
// with validation and completition logic
type Wizard struct {
	// Configuration
	config WizardConfig

	// Container properties
	VBox       *fyne.Container
	HBox       *fyne.Container
	containers *WizardContainers

	// Step buttons
	BackButton   *widget.Button
	NextButton   *widget.Button
	FinishButton *widget.Button
	CancelButton *widget.Button

	// Phases properties
	currentStep int
	Steps       []WizardStep
}

// Instantiates a new wizard
func NewWizard(config WizardConfig, steps []WizardStep) *Wizard {
	w := &Wizard{
		config:      config,
		Steps:       steps,
		currentStep: 0,
	}
	w.containers = w.buildContainers()
	w.buildBorderLayout()
	w.rebuildStepsContainer()
	return w
}

// Build the main layout where all parts of wizard resides
func (w *Wizard) buildBorderLayout() {
	w.HBox = container.New(layout.NewFormLayout(),
		container.NewPadded(w.containers.stepsContainer),
		w.containers.taskContainer,
	)
	w.VBox = container.NewBorder(
		nil, //w.containers.titleContainer,
		container.NewCenter(w.containers.stepsButtonsContainer),
		nil, nil,
		w.HBox,
	)
}

// Implements gui.Layout, this makes possible to set the content on window
func (w *Wizard) GetContainer() *fyne.Container {
	return w.VBox
}

func (w *Wizard) Refresh() {
	for _, s := range w.Steps {
		s.Refresh()
	}
	w.refreshContainers()
}
