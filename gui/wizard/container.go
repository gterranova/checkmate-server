package wizard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	BACK   = "Indietro"
	NEXT   = "Avanti"
	FINISH = "Fine"
	CANCEL = "Annulla"
	CLOSE  = "Chiudi"
)

// Represents all de parts for the border layout of the wizard
type WizardContainers struct {
	titleContainer        *fyne.Container
	stepsContainer        *fyne.Container
	taskContainer         *fyne.Container
	stepsButtonsContainer *fyne.Container
}

// Instantiates a new container group
func (w *Wizard) buildContainers() *WizardContainers {
	return &WizardContainers{
		titleContainer:        w.buildTitleContainer(),
		stepsContainer:        w.buildStepsContainer(),
		taskContainer:         w.buildTaskContainer(),
		stepsButtonsContainer: w.buildStepButtonsContainer(),
	}

}

// Builds the header container, with title and caption
func (w *Wizard) buildTitleContainer() *fyne.Container {
	title := canvas.NewText(w.config.GetTitle(), theme.ForegroundColor())
	title.TextSize = 20
	title.Alignment = fyne.TextAlignCenter

	caption := canvas.NewText(w.Steps[w.currentStep].GetCaption(), theme.ForegroundColor())
	caption.TextSize = 10
	caption.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		title,
		caption,
		widget.NewSeparator(),
	)
}

// Rebuilds the header container
func (w *Wizard) rebuildTitleContainer() {
	tc := w.buildTitleContainer()
	w.VBox.Objects[0] = tc
}

// Builds steps lists for showing current task
func (w *Wizard) buildStepsContainer() *fyne.Container {
	vb := container.NewVBox()
	for i, step := range w.Steps {
		t := canvas.NewText(step.GetTitle(), theme.DisabledColor())
		t.TextSize = 10
		if i == w.currentStep {
			t.Alignment = fyne.TextAlignTrailing
			t.Color = theme.ForegroundColor()
			t.TextStyle = fyne.TextStyle{
				Bold: true,
			}
		}
		vb.Add(t)
	}
	hb := container.NewHBox(vb, widget.NewSeparator())
	return hb
}

// Recreates steps container and replaces in the container
func (w *Wizard) rebuildStepsContainer() {
	sc := w.buildStepsContainer()
	w.HBox.Objects[0] = sc
}

// Defaults the task container with first task element
func (w *Wizard) buildTaskContainer() *fyne.Container {
	return w.Steps[w.currentStep].OnLoad()
}

// Rebuilds current task container, tipically when step changes
func (w *Wizard) rebuildTaskContainer() {
	w.containers.taskContainer = w.buildTaskContainer()
	w.HBox.Objects[1] = container.NewPadded(
		container.NewPadded(w.containers.taskContainer),
	)
}

// Builds on finish wizard steps
func (w *Wizard) buildFinishContainer(content *fyne.Container) {
	w.HBox.Objects[1] = content
	w.HBox.Objects[0] = container.NewVBox()
}

// Builds all the button to navigate through the wizard
func (w *Wizard) buildStepButtonsContainer() *fyne.Container {
	w.buildStepButtons()
	w.refreshButtonStatus()

	return container.NewHBox(
		layout.NewSpacer(),
		w.BackButton,
		w.NextButton,
		w.FinishButton,
		w.CancelButton,
	)
}

// Generates and assign to Wizard all the buttons of the bottom lane
func (w *Wizard) buildStepButtons() {
	w.BackButton = w.buildBackButton()
	w.NextButton = w.buildNextButton()
	w.FinishButton = w.buildFinishButton()
	w.CancelButton = w.buildCancelButton()
}

// Builds the back button, this set the wizard to the previous step
func (w *Wizard) buildBackButton() *widget.Button {
	backButton := widget.NewButton(
		BACK, func() {
			w.Back()
		})
	return backButton
}

// Builds the next button, this set the wizard to the forward step
func (w *Wizard) buildNextButton() *widget.Button {
	nextButton := widget.NewButton(
		NEXT, func() {
			w.Next()
		})
	return nextButton
}

// Builds finish button, this commits the wizard steps
func (w *Wizard) buildFinishButton() *widget.Button {
	finishButton := widget.NewButton(
		FINISH, func() {
			w.Finish()
		})
	return finishButton
}

// Builds the cancel button, this must close and free all wizards resources
func (w *Wizard) buildCancelButton() *widget.Button {
	cancelButton := widget.NewButton(
		CANCEL, func() {
			w.config.OnClose()
		})
	return cancelButton
}
