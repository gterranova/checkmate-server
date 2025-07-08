package main

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"terra9.it/checkmate/gui/settings"
)

var langToISO = map[string]string{"Italiano": "it", "English": "en"}
var ISOToLang = map[string]string{"it": "Italiano", "en": "English"}

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	IconName     fyne.ThemeIconName
	View         fyne.CanvasObject
}

// SettingsView contains the gui information for the settings screen.
func settingsView(aw *mainWindow) fyne.CanvasObject {

	// TODO: Add setting for changing language.

	settingsContentView := container.New(layout.NewVBoxLayout())

	// Make it possible for the user to switch themes.
	def := settings.ThemeName()
	themeNames := []string{"dark", "light"}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		themeNames = append(themeNames, settings.SystemThemeName)
		if settings.ThemeName() == "" {
			def = "light"
		}
	}
	themes := widget.NewSelect(themeNames, func(theme string) {
		settings.SetTheme(theme)
		settings.ApplyTheme()
		aw.window.Canvas().Refresh(settingsContentView)
	})
	themes.SetSelected(def)
	// Add the theme switcher next to a label.

	themeChanger := container.New(layout.NewGridLayout(2), widget.NewLabel("Tema"), themes)

	// userInterfaceSettings is a group holding widgets related to user interface settings such as theme.
	userInterfaceSettings := container.NewVBox(themeChanger)

	// languages
	defLang := ISOToLang[settings.Language()]
	languages := widget.NewSelect([]string{"Italiano", "English"}, func(lang string) {
		settings.SetLanguage(langToISO[lang])
		settings.ApplyTheme()
		aw.window.Canvas().Refresh(settingsContentView)
	})
	languages.SetSelected(defLang)

	langChanger := container.New(layout.NewGridLayout(2), widget.NewLabel("Lingue"), languages)

	// userInterfaceSettings is a group holding widgets related to user interface settings such as theme.
	langSettings := container.NewVBox(langChanger)

	content := container.NewMax()

	var (
		// panels defines the metadata for each tutorial
		panels = map[string]Tutorial{
			"appearance": {"Aspetto", "", "colorPalette", userInterfaceSettings},
			"language":   {"Lingua", "", "document", langSettings},
		}

		// panelIndex  defines how our panels should be laid out in the index tree
		panelIndex = map[string][]string{
			"": {"appearance", "language"},
		}
	)

	setContent := func(t Tutorial) {
		if fyne.CurrentDevice().IsMobile() {
			child := aw.app.app.NewWindow(t.Title)
			child.SetContent(t.View)
			child.Show()
			child.SetOnClosed(func() {
			})
			return
		}

		content.Objects = []fyne.CanvasObject{t.View}
		content.Refresh()
	}

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return panelIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := panelIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			//return widget.NewButtonWithIcon("Settings panel node", &fyne.StaticResource{}, func() {})
			return container.NewHBox(widget.NewIcon(settings.Theme().Icon(theme.IconNameCheckButton)), widget.NewLabel("Settings panel node"))
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := panels[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			b := obj.(*fyne.Container)
			b.Objects[0].(*widget.Icon).SetResource(settings.Theme().Icon(t.IconName))
			b.Objects[1].(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := panels[uid]; ok {
				setContent(t)
			}
		},
	}
	tree.Select("appearance")
	return container.NewBorder(nil, nil, tree, nil, content)
}
