package gui

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

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
	"terra9.it/vadovia/core/features"
	"terra9.it/vadovia/gui/widgets"
)

type application struct {
	app fyne.App
}

type mainWindow struct {
	window  fyne.Window
	content *fyne.Container

	response               *canvas.Text
	suitableFeatures       *widgets.MultiSelectList
	suitableFeaturesGroup  *container.TabItem
	DM2010Features         *widgets.MultiSelectList
	DM2010FeaturesGroup    *container.TabItem
	screeningFeatures      *widgets.MultiSelectList
	screeningFeaturesGroup *container.TabItem

	tabContainer *container.AppTabs

	app fyne.App
}

func newMainWindow(app fyne.App) *mainWindow {
	return &mainWindow{
		app:    app,
		window: app.NewWindow("vadoVIA"),
	}
}

func (w *mainWindow) Update(project *core.Project) {
	if project.CapacityMW < (core.ScreeningThreshold/2) || project.Art22bisArea() {
		w.HideSuitableFeatures()
		w.HideDM2010Features()
	} else {
		w.ShowSuitableFeatures()
		if project.SuitableArea() {
			w.HideDM2010Features()
		} else {
			w.ShowDM2010Features()
		}
	}
	EIAProcedure := project.EnvironmentalProcedures[0]
	if !EIAProcedure.IsApplicable() && project.CapacityMW >= core.ScreeningThreshold*0.5 {
		w.ShowScreeningCriteria()
	} else {
		w.HideScreeningCriteria()
	}

	w.tabContainer.Refresh()

	w.response.Text = project.Evaluate()
	w.response.Refresh()

}

func (w *mainWindow) ShowSuitableFeatures() {
	w.suitableFeatures.Enable()
	w.tabContainer.EnableItem(w.suitableFeaturesGroup)
}

func (w *mainWindow) ShowDM2010Features() {
	w.DM2010Features.Enable()
	w.tabContainer.EnableItem(w.DM2010FeaturesGroup)
}

func (w *mainWindow) ShowScreeningCriteria() {
	w.screeningFeatures.Enable()
	w.tabContainer.EnableItem(w.screeningFeaturesGroup)
}

func (w *mainWindow) HideSuitableFeatures() {
	w.suitableFeatures.Disable()
	w.tabContainer.DisableItem(w.suitableFeaturesGroup)
}

func (w *mainWindow) HideDM2010Features() {
	w.DM2010Features.Disable()
	w.tabContainer.DisableItem(w.DM2010FeaturesGroup)
}

func (w *mainWindow) HideScreeningCriteria() {
	w.screeningFeatures.Disable()
	w.tabContainer.DisableItem(w.screeningFeaturesGroup)
}

func (w *mainWindow) Create(project *core.Project) {

	projectCapacityMW := widget.NewEntry()
	projectCapacityMW.SetText("1")
	projectCapacityMW.OnChanged = func(s string) {
		power, _ := strconv.ParseFloat(s, 32)
		project.SetPower(math.Round(power*100) / 100)
		w.Update(project)
	}

	localizationFeatures := make([]interface{}, 0)
	for i, c := range project.Features[features.LOCATION].(features.LocationFeatures).BaseFeatures {
		if i == 0 {
			c.Checked = true
		}
		localizationFeatures = append(localizationFeatures, c)
	}

	localizationArea := widgets.NewSelectWithData(localizationFeatures,
		func(i interface{}) string {
			c := i.(*features.Feature)
			if len(c.Reference) > 0 {
				return fmt.Sprintf("%v (%v)", c.Title, c.Reference)
			}
			return c.Title
		},
		func(i interface{}) bool {
			c := i.(*features.Feature)
			return c != nil && c.Checked
		},
	)

	localizationArea.OnChanged = func(item interface{}) {
		for _, c := range localizationFeatures {
			c.(*features.Feature).Checked = false
		}
		condition := item.(*features.Feature)
		condition.Checked = true
		w.Update(project)
	}

	options := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range project.Features[features.SUITABILITY].(features.SuitabilityFeatures).BaseFeatures {
		options = append(options, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c,
			Checked:  c.Checked,
			Disabled: c.Disabled,
		})
	}
	w.suitableFeatures = widgets.NewMultiSelectList(options,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*features.Feature).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*features.Feature)
			condition.Checked = item.Checked
			w.Update(project)
		},
	)
	w.suitableFeaturesGroup = container.NewTabItem("Aree idonee d.lgs. 199/2021", w.suitableFeatures)

	unsuitableOptions := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range project.Features[features.NATIONAL_GUIDELINES].(features.NationalGuidelinesFeatures).BaseFeatures {
		unsuitableOptions = append(unsuitableOptions, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c,
			Checked:  c.Checked,
			Disabled: c.Disabled,
		})
	}
	w.DM2010Features = widgets.NewMultiSelectList(unsuitableOptions,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*features.Feature).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*features.Feature)
			condition.Checked = item.Checked
			switch condition.Tags {
			case features.EUAP, features.Ramsar, features.ReteNatura2000, features.Art142Dlgs2004, features.Art142Dlgs2004Coste, features.Art142Dlgs2004MontagneBoschi:
				for i, c := range project.Features[features.SCREENING_GUIDELINES].(features.ScreeningGuidelinesFeatures).BaseFeatures {
					if c.Tags == condition.Tags {
						w.screeningFeatures.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
						continue
					}
					if c.Tags == features.Art142Dlgs2004 && (condition.Tags == features.Art142Dlgs2004Coste || condition.Tags == features.Art142Dlgs2004MontagneBoschi) {
						w.screeningFeatures.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
					}
				}
			}
			w.Update(project)
		},
	)
	w.DM2010FeaturesGroup = container.NewTabItem("Aree non idonee Linee Guida Nazionali", w.DM2010Features)

	screeningCriteria := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range project.Features[features.SCREENING_GUIDELINES].(features.ScreeningGuidelinesFeatures).BaseFeatures {
		screeningCriteria = append(screeningCriteria, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c,
			Checked:  c.Checked,
			Disabled: c.Disabled,
		})
	}
	w.screeningFeatures = widgets.NewMultiSelectList(screeningCriteria,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*features.Feature).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*features.Feature)
			condition.Checked = item.Checked
			switch condition.Tags {
			case features.EUAP, features.Ramsar, features.ReteNatura2000, features.Art142Dlgs2004, features.Art142Dlgs2004Coste, features.Art142Dlgs2004MontagneBoschi:
				for i, c := range project.Features[features.NATIONAL_GUIDELINES].(features.NationalGuidelinesFeatures).BaseFeatures {
					if c.Tags == condition.Tags {
						w.DM2010Features.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
						continue
					}
					if c.Tags == features.Art142Dlgs2004 && (condition.Tags == features.Art142Dlgs2004Coste || condition.Tags == features.Art142Dlgs2004MontagneBoschi) {
						w.DM2010Features.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
					}
				}
			}
			w.Update(project)
		},
	)
	w.screeningFeaturesGroup = container.NewTabItem("Linee Guida Screening DM 30/03/2015", w.screeningFeatures)

	powerLabel := widget.NewLabel("Potenza Impianto FV (MW):")
	powerLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	powerLabel.Alignment = fyne.TextAlignTrailing

	locLabel := widget.NewLabel("Localizzazione:")
	locLabel.Alignment = fyne.TextAlignTrailing
	locLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	form := container.NewGridWithColumns(2,
		powerLabel, projectCapacityMW,
		locLabel, localizationArea,
	)
	w.tabContainer = container.NewAppTabs(w.suitableFeaturesGroup, w.DM2010FeaturesGroup, w.screeningFeaturesGroup)

	w.response = canvas.NewText("", color.Black)
	w.response.TextSize = 16

	report := widget.NewButtonWithIcon("Dettagli", w.app.Settings().Theme().Icon(theme.IconNameDocument), func() {

		details := widget.NewRichTextFromMarkdown(project.Describe())
		details.Wrapping = fyne.TextWrapWord
		d := dialog.NewCustom("Dettagli", "Chiudi", details, w.window)
		//d.Resize(w.window.Canvas().Size())
		d.Resize(fyne.NewSize(600, 400))
		d.Show()
	})

	about := widget.NewButtonWithIcon("Informazioni", w.app.Settings().Theme().Icon(theme.IconNameInfo), func() {

		d := dialog.NewCustom("Informazioni", "Chiudi", aboutView(), w.window)
		d.Resize(fyne.NewSize(600, 400))
		d.Show()
	})

	image := canvas.NewImageFromResource(assets.ResourceLogovadoviacompactPng)
	image.FillMode = canvas.ImageFillOriginal

	topBar := container.NewBorder(nil, nil, container.NewCenter(image), nil, form)
	bottomBar := container.NewBorder(nil, nil, nil, container.NewHBox(about, report), w.response)

	w.content = container.NewBorder(topBar, bottomBar, nil, nil,
		w.tabContainer,
	)
}

func (w *mainWindow) CreateAndShow(project *core.Project) {
	w.Create(project)
	w.Update(project)
	w.window.SetContent(container.New(layout.NewMaxLayout(), w.content))
	w.window.Resize(fyne.NewSize(1200, 800))
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
