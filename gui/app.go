package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"terra9.it/vadovia/assets"
	"terra9.it/vadovia/core"
	"terra9.it/vadovia/gui/wizard"
)

type application struct {
	app fyne.App
}

type mainWindow struct {
	window  fyne.Window
	content *fyne.Container

	response *canvas.Text

	wizard *wizard.Wizard

	app fyne.App
}

func newMainWindow(app fyne.App) *mainWindow {
	return &mainWindow{
		app:    app,
		window: app.NewWindow("vadoVIA"),
	}
}

func (w *mainWindow) Update(project *core.Project) {
	w.response.Text = project.Evaluate()
	w.response.Refresh()
	w.wizard.Refresh()
	//w.wizard.VBox.Refresh()
}

func (w *mainWindow) Create(project *core.Project) {

	wc := &wizardConfig{project: project}

	steps := make([]wizard.WizardStep, len(project.FeatureOrder))
	//steps[0] = NewProjectStep(project, w)
	for i, tag := range project.FeatureOrder {
		steps[i] = NewMultiselectStep(project, w, project.Features[tag])
	}
	w.wizard = wizard.NewWizard(wc, steps)

	//w.tabContainer = container.NewAppTabs(w.screeningFeaturesGroup)

	w.response = canvas.NewText("", color.Black)
	w.response.TextSize = 16

	report := widget.NewButtonWithIcon("Dettagli", w.app.Settings().Theme().Icon(theme.IconNameDocument), func() {

		details := widget.NewRichTextFromMarkdown(project.Describe())
		details.Wrapping = fyne.TextWrapWord
		d := dialog.NewCustom("Dettagli", "Chiudi", details, w.window)
		//d.Resize(w.window.Canvas().Size())
		d.Resize(w.window.Canvas().Size().Subtract(fyne.NewDelta(50, 50)))
		d.Show()
	})

	about := widget.NewButtonWithIcon("Informazioni", w.app.Settings().Theme().Icon(theme.IconNameInfo), func() {

		d := dialog.NewCustom("Informazioni", "Chiudi", aboutView(), w.window)
		d.Resize(fyne.NewSize(600, 400))
		d.Show()
	})

	image := canvas.NewImageFromResource(assets.ResourceLogovadoviacompactPng)
	image.FillMode = canvas.ImageFillOriginal

	topBar := container.NewBorder(nil, nil, container.NewCenter(image), nil, nil)
	bottomBar := container.NewBorder(nil, nil, nil, container.NewHBox(about, report), w.response)

	w.content = container.NewBorder(topBar, bottomBar, nil, nil,
		//w.tabContainer,
		w.wizard.GetContainer(),
	)
}

func (w *mainWindow) CreateAndShow(project *core.Project) {
	w.Create(project)
	project.Validate("")
	w.Update(project)
	w.window.SetContent(container.New(layout.NewMaxLayout(), w.content))
	w.window.Resize(fyne.NewSize(900, 700))
	w.window.SetIcon(assets.ResourceIconPng)
	w.window.CenterOnScreen()
	w.window.Show()
}

func NewApp() *application {

	a := &application{
		app: app.NewWithID("it.terra9.vadovia"),
	}

	newMainWindow(a.app).CreateAndShow(core.NewProject())
	return a
}

func (a *application) Run() {
	a.app.Run()
}

// -- START of Wizard configuration struct --
type wizardConfig struct {
	wizard.BaseWizardConfig
	project *core.Project
}

// Callback called on 3 cases:
// 1. Cancel Button clicked
// 2. Close button clicked after Finished
// 3. Window instance closed
func (c *wizardConfig) OnClose() {
	//createInstance.wizard.Close()
}

// Stores new sites created and returns finish view
func (c *wizardConfig) OnFinish() *fyne.Container {
	//_, err := createInstance.site.StoreSite()
	//if err != nil {
	//	return view.NewErrorSiteCreationLayout(err)
	//}
	//return view.NewSuccessSiteCreationLayout()
	details := widget.NewRichTextFromMarkdown(c.project.Describe())
	details.Wrapping = fyne.TextWrapWord

	return container.NewPadded(details)
}

func (c *wizardConfig) CanClose() bool {
	return false
}

func (c *wizardConfig) CanFinish() bool {
	return false
}
