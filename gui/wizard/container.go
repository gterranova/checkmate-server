package wizard

import (
	"fyne.io/fyne/v2"
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
	//titleContainer        *fyne.Container
	stepsContainer        *fyne.Container
	taskContainer         *fyne.Container
	stepsButtonsContainer *fyne.Container
}

// Instantiates a new container group
func (w *Wizard) buildContainers() *WizardContainers {
	return &WizardContainers{
		//titleContainer:        w.buildTitleContainer(),
		stepsContainer:        w.buildStepsContainer(),
		taskContainer:         w.buildTaskContainer(),
		stepsButtonsContainer: w.buildStepButtonsContainer(),
	}

}

// Builds the header container, with title and caption
//func (w *Wizard) buildTitleContainer() *fyne.Container {
//	title := canvas.NewText(w.config.GetTitle(), theme.ForegroundColor())
//	title.TextSize = 20
//	title.Alignment = fyne.TextAlignCenter

//	caption := canvas.NewText(w.Steps[w.currentStep].GetCaption(), theme.ForegroundColor())
//	caption.TextSize = 10
//	caption.Alignment = fyne.TextAlignCenter

//	return container.NewVBox(
//		title,
//		caption,
//		widget.NewSeparator(),
//	)
//}

// Rebuilds the header container
//func (w *Wizard) rebuildTitleContainer() {
//	tc := w.buildTitleContainer()
//	w.VBox.Objects[1] = tc
//}

// Builds steps lists for showing current task
func (w *Wizard) buildStepsContainer() *fyne.Container {
	vb := container.NewVBox()
	for i, step := range w.Steps {
		if step.Disabled() {
			continue
		}
		color := theme.DisabledColor()
		if i == w.currentStep {
			color = theme.ForegroundColor()
		}
		current := i
		t := newNavigationLabel(step.GetTitle(), color, nil, func(b bool) {
			//fmt.Println(current)
			w.SelectStep(current)
		})

		vb.Add(t)
	}
	vb.Resize(fyne.NewSize(200, 200))
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
	w.HBox.Objects[1] = w.containers.taskContainer
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

	container := container.NewHBox(
		layout.NewSpacer(),
		w.BackButton,
		w.NextButton,
	)
	if w.config.CanFinish() {
		container.Add(w.FinishButton)
	}
	if w.config.CanClose() {
		container.Add(w.CancelButton)
	}
	return container
}

// Generates and assign to Wizard all the buttons of the bottom lane
func (w *Wizard) buildStepButtons() {
	w.BackButton = w.buildBackButton()
	w.NextButton = w.buildNextButton()
	if w.config.CanFinish() {
		w.FinishButton = w.buildFinishButton()
	}
	if w.config.CanClose() {
		w.CancelButton = w.buildCancelButton()
	}
}

// Builds the back button, this set the wizard to the previous step
func (w *Wizard) buildBackButton() *widget.Button {
	backButton := widget.NewButtonWithIcon(
		BACK, theme.NavigateBackIcon(), func() {
			w.Back()
		})
	backButton.IconPlacement = widget.ButtonIconLeadingText
	return backButton
}

// Builds the next button, this set the wizard to the forward step
func (w *Wizard) buildNextButton() *widget.Button {
	nextButton := widget.NewButtonWithIcon(
		NEXT, theme.NavigateNextIcon(), func() {
			w.Next()
		})
	nextButton.IconPlacement = widget.ButtonIconTrailingText
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

// Recreates steps container and replaces in the container
func (w *Wizard) RebuildStepsContainer() {
	w.rebuildStepsContainer()
}
