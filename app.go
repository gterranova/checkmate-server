package main

import (
	"fmt"
	"image/color"
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
	"terra9.it/vadovia/widgets"
)

type application struct {
	app fyne.App
}

type mainWindow struct {
	window  fyne.Window
	content *fyne.Container

	response                *canvas.Text
	suitableConditions      *widgets.MultiSelectList
	suitableConditionsGroup *container.TabItem
	DM2010Conditions        *widgets.MultiSelectList
	DM2010ConditionsGroup   *container.TabItem
	screeningCriteria       *widgets.MultiSelectList
	screeningCriteriaGroup  *container.TabItem

	tabContainer *container.AppTabs

	app fyne.App
}

func newMainWindow(app fyne.App) *mainWindow {
	return &mainWindow{
		app:    app,
		window: app.NewWindow("vadoVIA"),
	}
}

func (w *mainWindow) Update(project *Project) {
	if project.capacityMW < ScreeningThreshold || project.locationConditions.Art22bisArea() == Art22bis {
		w.HideSuitableConditions()
		w.HideDM2010Conditions()
	} else {
		w.ShowSuitableConditions()
		if project.SuitableArea() != NotApplicable {
			w.HideDM2010Conditions()
		} else {
			w.ShowDM2010Conditions()
		}
	}
	EIAProcedure := project.environmentalProcedures[0]
	if !EIAProcedure.IsApplicable() && project.capacityMW >= ScreeningThreshold*0.5 {
		w.ShowScreeningCriteria()
	} else {
		w.HideScreeningCriteria()
	}

	w.tabContainer.Refresh()

	w.response.Text = project.Evaluate()
	w.response.Refresh()

}

func (w *mainWindow) ShowSuitableConditions() {
	w.suitableConditions.Enable()
	w.tabContainer.EnableItem(w.suitableConditionsGroup)
}

func (w *mainWindow) ShowDM2010Conditions() {
	w.DM2010Conditions.Enable()
	w.tabContainer.EnableItem(w.DM2010ConditionsGroup)
}

func (w *mainWindow) ShowScreeningCriteria() {
	w.screeningCriteria.Enable()
	w.tabContainer.EnableItem(w.screeningCriteriaGroup)
}

func (w *mainWindow) HideSuitableConditions() {
	w.suitableConditions.Disable()
	w.tabContainer.DisableItem(w.suitableConditionsGroup)
}

func (w *mainWindow) HideDM2010Conditions() {
	w.DM2010Conditions.Disable()
	w.tabContainer.DisableItem(w.DM2010ConditionsGroup)
}

func (w *mainWindow) HideScreeningCriteria() {
	w.screeningCriteria.Disable()
	w.tabContainer.DisableItem(w.screeningCriteriaGroup)
}

func (w *mainWindow) Create(project *Project) {

	projectCapacityMW := widget.NewEntry()
	projectCapacityMW.SetText("1")
	projectCapacityMW.OnChanged = func(s string) {
		power, _ := strconv.ParseFloat(s, 32)
		project.SetPower(power)
		w.Update(project)
	}

	localizationConditions := make([]interface{}, 0)
	for i, c := range project.locationConditions {
		if i == 0 {
			c.Checked = true
		}
		localizationConditions = append(localizationConditions, c)
	}

	localizationArea := widgets.NewSelectWithData(localizationConditions,
		func(i interface{}) string {
			c := i.(*Condition)
			if len(c.Reference) > 0 {
				return fmt.Sprintf("%v (%v)", c.Title, c.Reference)
			}
			return c.Title
		},
		func(i interface{}) bool {
			c := i.(*Condition)
			return c != nil && c.Checked
		},
	)

	localizationArea.OnChanged = func(item interface{}) {
		for _, c := range localizationConditions {
			c.(*Condition).Checked = false
		}
		condition := item.(*Condition)
		condition.Checked = true
		w.Update(project)
	}

	options := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range project.suitabilityConditions {
		options = append(options, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c,
			Checked:  c.Checked,
			Disabled: c.Disabled,
		})
	}
	w.suitableConditions = widgets.NewMultiSelectList(options,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*Condition).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*Condition)
			condition.Checked = item.Checked
			w.Update(project)
		},
	)
	w.suitableConditionsGroup = container.NewTabItem("Aree idonee d.lgs. 199/2021", w.suitableConditions)
	//w.suitableConditionsGroup = widgets.NewGroup("Aree idonee ex art. 20, d.lgs. 199/2021", w.suitableConditions)

	unsuitableOptions := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range project.nationalGuidelinesConditions {
		unsuitableOptions = append(unsuitableOptions, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c,
			Checked:  c.Checked,
			Disabled: c.Disabled,
		})
	}
	w.DM2010Conditions = widgets.NewMultiSelectList(unsuitableOptions,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*Condition).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*Condition)
			condition.Checked = item.Checked
			switch condition.Value {
			case EUAP, Ramsar, ReteNatura2000, Art142Dlgs2004, Art142Dlgs2004Coste, Art142Dlgs2004MontagneBoschi:
				for i, c := range project.screeningGuidelinesConditions {
					if c.Value == condition.Value {
						w.screeningCriteria.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
						continue
					}
					if c.Value == Art142Dlgs2004 && (condition.Value == Art142Dlgs2004Coste || condition.Value == Art142Dlgs2004MontagneBoschi) {
						w.screeningCriteria.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
					}
				}
			}
			w.Update(project)
		},
	)
	w.DM2010ConditionsGroup = container.NewTabItem("Aree non idonee Linee Guida Nazionali", w.DM2010Conditions)

	screeningCriteria := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range project.screeningGuidelinesConditions {
		screeningCriteria = append(screeningCriteria, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c,
			Checked:  c.Checked,
			Disabled: c.Disabled,
		})
	}
	w.screeningCriteria = widgets.NewMultiSelectList(screeningCriteria,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*Condition).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*Condition)
			condition.Checked = item.Checked
			switch condition.Value {
			case EUAP, Ramsar, ReteNatura2000, Art142Dlgs2004, Art142Dlgs2004Coste, Art142Dlgs2004MontagneBoschi:
				for i, c := range project.nationalGuidelinesConditions {
					if c.Value == condition.Value {
						w.DM2010Conditions.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
						continue
					}
					if c.Value == Art142Dlgs2004 && (condition.Value == Art142Dlgs2004Coste || condition.Value == Art142Dlgs2004MontagneBoschi) {
						w.DM2010Conditions.GetItems()[i].Checked = condition.Checked
						c.Checked = condition.Checked
					}
				}
			}
			w.Update(project)
		},
	)
	w.screeningCriteriaGroup = container.NewTabItem("Linee Guida Screening DM 30/03/2015", w.screeningCriteria)

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
	w.tabContainer = container.NewAppTabs(w.suitableConditionsGroup, w.DM2010ConditionsGroup, w.screeningCriteriaGroup)

	w.response = canvas.NewText("", color.Black)
	w.response.TextSize = 16

	report := widget.NewButtonWithIcon("Dettagli", w.app.Settings().Theme().Icon(theme.IconNameDocument), func() {

		details := widget.NewRichTextFromMarkdown(project.Describe())
		details.Wrapping = fyne.TextWrapWord
		d := dialog.NewCustom("Dettagli", "Chiudi", details, w.window)
		d.Resize(w.window.Canvas().Size())
		d.Show()
	})

	image := canvas.NewImageFromResource(assets.ResourceLogovadoviacompactPng)
	image.FillMode = canvas.ImageFillOriginal

	topBar := container.NewBorder(nil, nil, container.NewCenter(image), nil, form)
	bottomBar := container.NewBorder(nil, nil, nil, report, w.response)

	w.content = container.NewBorder(topBar, bottomBar, nil, nil,
		w.tabContainer,
	)
}

func (w *mainWindow) CreateAndShow(project *Project) {
	w.Create(project)
	w.Update(project)
	w.window.SetContent(container.New(layout.NewMaxLayout(), w.content))
	w.window.Resize(fyne.NewSize(1200, 600))
	w.window.SetIcon(assets.ResourceIconPng)
	w.window.CenterOnScreen()
	w.window.Show()
}

func NewApp() *application {

	a := &application{
		app: app.NewWithID("it.terra9.vadovia"),
	}

	newMainWindow(a.app).CreateAndShow(NewProject())
	return a
}

func (a *application) Run() {
	a.app.Run()
}
