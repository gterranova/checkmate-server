package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/skratchdot/open-golang/open"

	"terra9.it/checkmate/core"
	"terra9.it/checkmate/gui/assets"
	"terra9.it/checkmate/gui/settings"
	"terra9.it/checkmate/gui/wizard"
	"terra9.it/checkmate/loader"
)

type application struct {
	app    fyne.App
	loader *loader.ResourceLoader
}

type mainWindow struct {
	window   fyne.Window
	response *canvas.Text
	wizard   *wizard.Wizard
	app      *application
}

func (app *application) newMainWindow() *mainWindow {
	w := &mainWindow{
		app:    app,
		window: app.app.NewWindow("checkmate"),
	}

	var resLoader *loader.ResourceLoader
	var err error

	if _, ok := w.app.app.(desktop.App); ok {
		fmt.Println("Start desktop app")
		if resLoader, err = loader.NewLoader("settings"); err != nil {
			resLoader, err = loader.NewEmptyLoader("settings").FromBuffer(assets.ResourceSettingsChlx.StaticContent)
		}
	} else {
		resLoader, err = loader.NewEmptyLoader("settings").FromBuffer(assets.ResourceSettingsChlx.StaticContent)
	}
	if err != nil {
		dialog.ShowError(err, w.window)
	}
	app.loader = resLoader

	settings.InitSettings(app.loader)
	settings.ApplyTheme()

	w.window.Resize(fyne.NewSize(900, 700))
	iconRes := &fyne.StaticResource{
		StaticContent: w.app.loader.MustGet("icon.png"),
	}
	w.window.SetIcon(iconRes)
	w.window.CenterOnScreen()
	return w
}

func (w *mainWindow) Update(project *core.Project) {
	title, _ := strings.CutSuffix(w.window.Title(), "*")
	if project.Dirty() {
		w.window.SetTitle(title + "*")
	} else {
		w.window.SetTitle(title)
	}
	w.response.Text = project.Evaluate()
	w.response.Refresh()
	w.wizard.Refresh()
	//w.wizard.VBox.Refresh()
}

// ToolbarLabel is a blank, stretchable space for a toolbar.
// This is typically used to assist layout if you wish some left and some right aligned items.
// Space will be split evebly amongst all the spacers on a toolbar.
type ToolbarLabel struct {
	label string `json:"-"`
}

// ToolbarObject gets the actual spacer object for this ToolbarLabel
func (t *ToolbarLabel) ToolbarObject() fyne.CanvasObject {
	return widget.NewLabelWithStyle(t.label, fyne.TextAlignLeading, fyne.TextStyle{
		Bold: true,
	})
}

// NewToolbarLabel returns a new spacer item for a Toolbar to assist with ToolbarItem alignment
func NewToolbarLabel(label string) *ToolbarLabel {
	return &ToolbarLabel{label}
}

func (w *mainWindow) topBar(project *core.Project) fyne.CanvasObject {
	//imageRes := &fyne.StaticResource{
	//	StaticContent: w.app.loader.MustGet("logo.png"),
	//}

	//image := canvas.NewImageFromResource(imageRes)
	//image.FillMode = canvas.ImageFillOriginal

	menuItem1 := fyne.NewMenuItem("Configurazione", func() {
		settingsDialog := dialog.NewCustom("Settings", "CLOSE", settingsView(w), w.window)
		settingsDialog.Resize(fyne.NewSize(750, 500))
		settingsDialog.Show()
	})
	menuItem3 := fyne.NewMenuItem("Informazioni", func() {
		d := dialog.NewCustom("Informazioni", "Chiudi", aboutView(w.app), w.window)
		d.Resize(fyne.NewSize(600, 400))
		d.Show()
	})
	menuItem2 := fyne.NewMenuItem("Esci", func() {
		w.window.Close()
	})
	menu := fyne.NewMenu("", menuItem1, menuItem3, fyne.NewMenuItemSeparator(), menuItem2)

	//settingsButton := widgets.NewTappableIconWithMenu(theme.SettingsIcon(), 40, menu, true)
	//settingsButton.Importance = widget.LowImportance

	var toolbar fyne.CanvasObject
	var settingAction *widget.ToolbarAction
	settingAction = widget.NewToolbarAction(theme.SettingsIcon(), func() {
		//h := settingAction.ToolbarObject().MinSize().Height
		//holder := aw.app.Driver().CanvasForObject(icon)
		popup := widget.NewPopUpMenu(menu, w.window.Canvas())
		pos := w.app.app.Driver().AbsolutePositionForObject(settingAction.ToolbarObject()).
			AddXY(toolbar.Size().Width-popup.MinSize().Width+theme.Padding(), toolbar.Size().Height+theme.Padding())
		//holder.Focus(icon)
		popup.ShowAtPosition(pos)
		popup.OnDismiss = func() {
			popup.Hide()
		}
	})
	if project != nil {
		var resetAction, openAction, saveAction, saveAsAction *widget.ToolbarAction
		settingAction = widget.NewToolbarAction(theme.SettingsIcon(), func() {
			//h := settingAction.ToolbarObject().MinSize().Height
			//holder := aw.app.Driver().CanvasForObject(icon)
			popup := widget.NewPopUpMenu(menu, w.window.Canvas())
			pos := w.app.app.Driver().AbsolutePositionForObject(settingAction.ToolbarObject()).
				AddXY(toolbar.Size().Width-popup.MinSize().Width+theme.Padding(), toolbar.Size().Height+theme.Padding())
			//holder.Focus(icon)
			popup.ShowAtPosition(pos)
			popup.OnDismiss = func() {
				popup.Hide()
			}
		})

		resetAction = widget.NewToolbarAction(theme.FolderNewIcon(), func() {
			project.ResetFeatures()
			for _, f := range project.Features {
				f.ApplyDefaults()
			}
			steps := make([]wizard.WizardStep, len(project.Features))
			//steps[0] = NewProjectStep(project, w)
			for i, feat := range project.Features {
				steps[i] = NewMultiselectStep(project, w, feat)
			}
			w.wizard.Steps = steps
			w.Update(project)
		})

		openAction = widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
				if err == nil && uc != nil {
					uc.Close()
					project, err = core.LoadProject(uc.URI().Path())
					if err != nil {
						dialog.ShowError(err, w.window)
						return
					}
					w.ChecklistPage(project)
				}
			}, w.window)
			d.Show()
		})

		saveAction = widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			if project.ProjectFile != "" {
				if err := project.SaveProject(); err != nil {
					panic(err)
				}
			} else {
				d := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
					if err == nil && uc != nil {
						uc.Close()
						os.Remove(uc.URI().Path())
						if err = project.SaveProjectAs(uc.URI().Path()); err != nil {
							panic(err)
						}
						w.window.SetTitle(project.ProjectFile)
					}
				}, w.window)
				d.Show()
			}
		})

		saveAsAction = widget.NewToolbarAction(theme.UploadIcon(), func() {
			menuItem1 := fyne.NewMenuItem("Salva con nome...", func() {
				d := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
					if err == nil && uc != nil {
						uc.Close()
						os.Remove(uc.URI().Path())
						if err = project.SaveProjectAs(uc.URI().Path()); err != nil {
							panic(err)
						}
						w.window.SetTitle(project.ProjectFile)
					}
				}, w.window)
				d.Show()
			})
			outputOptions := make([]*fyne.MenuItem, 0)
			outputOptions = append(outputOptions, menuItem1)
			outputOptions = append(outputOptions, fyne.NewMenuItemSeparator())
			for _, t := range project.TemplateDefs {
				template_def := t
				m := fyne.NewMenuItem(t.Name, func() {
					output, err := project.Render(template_def)
					if err != nil {
						dialog.ShowError(err, w.window)
						return
					}
					if template_def.Format != "docx" {
						details := widget.NewRichTextFromMarkdown(output)
						details.Wrapping = fyne.TextWrapWord
						d := dialog.NewCustom(t.Name, "Chiudi", details, w.window)
						//d.Resize(w.window.Canvas().Size())
						d.Resize(w.window.Canvas().Size().Subtract(fyne.NewDelta(50, 50)))
						d.Show()
					} else {
						dialog.ShowConfirm("Report", "Report generato come "+output+". Vuoi aprirlo?", func(b bool) {
							if b {
								open.Start(output)
							}
						}, w.window)
					}
				})
				outputOptions = append(outputOptions, m)
			}
			menu := fyne.NewMenu("", outputOptions...)

			popup := widget.NewPopUpMenu(menu, w.window.Canvas())
			pos := w.app.app.Driver().AbsolutePositionForObject(saveAsAction.ToolbarObject()).
				AddXY(toolbar.Size().Width-popup.MinSize().Width+theme.Padding()-32-16, toolbar.Size().Height+theme.Padding())
			//holder.Focus(icon)
			popup.ShowAtPosition(pos)
			popup.OnDismiss = func() {
				popup.Hide()
			}
		})

		toolbar = widget.NewToolbar(
			widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
				w.LandingPage()
			}),
			NewToolbarLabel(project.Name),
			widget.NewToolbarSpacer(),
			resetAction,
			openAction,
			saveAction,
			saveAsAction,
			widget.NewToolbarSeparator(),
			settingAction,
		)
	} else {
		iconRes := &fyne.StaticResource{
			StaticContent: w.app.loader.MustGet("icon.png"),
		}

		toolbar = widget.NewToolbar(
			widget.NewToolbarAction(iconRes, func() {}),
			NewToolbarLabel("SmartCheck"),
			widget.NewToolbarSpacer(),
			settingAction,
		)
	}

	return container.NewBorder(toolbar, nil, nil, nil, nil)
}

func (w *mainWindow) ChecklistContent(project *core.Project) fyne.CanvasObject {

	wc := &wizardConfig{project: project}

	steps := make([]wizard.WizardStep, len(project.Features))
	//steps[0] = NewProjectStep(project, w)
	for i, feat := range project.Features {
		steps[i] = NewMultiselectStep(project, w, feat)
	}
	w.wizard = wizard.NewWizard(wc, steps)

	settings.RemoveThemeChangeListeners()
	settings.AddThemeChangeListener(func() {
		if w.wizard != nil {
			steps := make([]wizard.WizardStep, len(project.Features))
			//steps[0] = NewProjectStep(project, w)
			for i, feat := range project.Features {
				steps[i] = NewMultiselectStep(project, w, feat)
			}
			w.wizard.Steps = steps
			w.wizard.Refresh()
		}
		if w.response != nil {
			w.response.Color = theme.ForegroundColor()
		}
	})

	//w.tabContainer = container.NewAppTabs(w.screeningFeaturesGroup)

	w.response = canvas.NewText("", theme.ForegroundColor())
	w.response.TextSize = 16

	bottomBar := w.response

	project.Validate("")
	w.Update(project)

	return container.NewBorder(w.topBar(project), bottomBar, nil, nil,
		//w.tabContainer,
		w.wizard.GetContainer(),
	)
}

func (w *mainWindow) ChecklistPage(project *core.Project) {
	title := project.ProjectFile
	if title == "" {
		title = project.Name
	}
	w.window.SetTitle(title)
	w.window.SetContent(container.New(layout.NewMaxLayout(), w.ChecklistContent(project)))
	w.window.SetCloseIntercept(func() {
		w.LandingPage()
	})
	w.window.Show()
}

func (w *mainWindow) LandingContent() fyne.CanvasObject {

	//bottomBar := container.NewBorder(nil, nil, nil, container.NewHBox(about, report), w.response)

	cards := make([]fyne.CanvasObject, 0)
	for _, pkg := range settings.InstalledProjects() {
		var pkgLoader *loader.ResourceLoader
		var err error

		if embeddedZip, ok := w.app.loader.Get(pkg + loader.CHECKLIST_EXT); ok {
			pkgLoader, err = loader.NewEmptyLoader(pkg).FromBuffer(embeddedZip)
		} else {
			pkgLoader, err = loader.NewLoader(pkg)
		}
		if err != nil {
			dialog.ShowError(err, w.window)
			return nil
		}

		project := core.NewProject(pkgLoader)
		card := widget.NewCard(project.Name, "Subtitle", widget.NewButton("Run", func() {
			if embeddedZip, ok := w.app.loader.Get(pkg + loader.CHECKLIST_EXT); ok {
				pkgLoader, err = loader.NewEmptyLoader(pkg).FromBuffer(embeddedZip)
			} else {
				pkgLoader, err = loader.NewLoader(pkg)
			}
			project := core.NewProject(pkgLoader)
			w.ChecklistPage(project)
		}))
		logo := &fyne.StaticResource{
			StaticName:    "logo.png",
			StaticContent: pkgLoader.MustGet("logo.png"),
		}

		image2 := canvas.NewImageFromResource(logo)
		image2.FillMode = canvas.ImageFillContain
		card.SetImage(image2)

		cards = append(cards, card)
	}

	return container.NewBorder(w.topBar(nil), nil /*bottomBar*/, nil, nil,
		//w.tabContainer,
		container.NewVBox(cards...),
	)
}

func (w *mainWindow) LandingPage() {

	w.wizard = nil
	w.window.SetTitle("SmartCheck")

	settings.RemoveThemeChangeListeners()

	w.window.SetCloseIntercept(func() {
		w.app.app.Quit()
	})

	w.window.SetContent(w.LandingContent())
	w.window.Show()
}

func NewApp() *application {

	a := &application{
		app: app.NewWithID("it.terra9.checkmate"),
	}
	mainWin := a.newMainWindow()

	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		project, err := core.LoadProject(args[0])

		if err != nil {
			return nil
		}
		mainWin.ChecklistPage(project)
	} else {
		mainWin.LandingPage()
	}
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
	return nil
}

func (c *wizardConfig) CanClose() bool {
	return false
}

func (c *wizardConfig) CanFinish() bool {
	return false
}
