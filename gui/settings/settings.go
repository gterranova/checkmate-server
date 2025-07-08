package settings

import (
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const (
	SystemThemeName = "system default"
)

type ResourceLoader interface {
	Get(filename string) ([]byte, bool)
	Set(filename string, data []byte) error
	Save() error
}

// Initilize this variable to access the env values
var Settings *settings
var Loader ResourceLoader

type ThemeChangeListener func()

var listeners []ThemeChangeListener

func RemoveThemeChangeListeners() {
	listeners = listeners[:0]
}

func AddThemeChangeListener(fn ThemeChangeListener) {
	listeners = append(listeners, fn)
}

// We will call this in main.go to load the env variables
func InitSettings(loader ResourceLoader) {
	if Settings == nil {
		Settings = &settings{}
		Loader = loader
		load()
	}
}

// Settings gives access to user interfaces to control Fyne settings
type settings struct {
	fyne.Settings `json:"-"`

	ThemeName string `json:"theme"`
	Language  string `json:"lang"`

	InstalledProjects []string `json:"installed_projects"`

	userTheme fyne.Theme
	logo      *canvas.Image
}

//func (sc *settings) StoragePath() string {
//	return "settings.json"
//}

func load() {
	if settings, ok := Loader.Get("settings.json"); ok {
		//err := loadFromFile(Settings.StoragePath())
		//if err != nil {
		//	fyne.LogError("Settings load error:", err)
		//}
		//} else {
		err := json.Unmarshal(settings, &Settings)
		if err != nil {
			fyne.LogError("Settings load error:", err)
		}
	} 
	//else {
	//
	//	}
	setupTheme()
}

func setupTheme() {

	switch Settings.ThemeName {
	case "light":
		Settings.userTheme = &myThemeLight{}
		//Settings.logo = canvas.NewImageFromResource(Resource(ResourceNameLogoLight))

	default:
		Settings.userTheme = &myThemeDark{}
		//Settings.logo = canvas.NewImageFromResource(Resource(ResourceNameLogoDark))
	}
	//Settings.logo.FillMode = canvas.ImageFillOriginal
}

/*
	func loadFromFile(path string) error {
		file, err := os.Open(path) // #nosec
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(filepath.Dir(path), 0700)
				if err != nil {
					return err
				}
				return nil
			}
			return err
		}
		decode := json.NewDecoder(file)

		return decode.Decode(&Settings.SettingsSchema)
	}
*/
func save() error {
	data, err := json.Marshal(&Settings)
	if err != nil {
		return err
	}

	if err := Loader.Set("settings.json", data); err != nil {
		//return saveToFile(Settings.StoragePath(), data)
		return err
	}
	return Loader.Save()
}

/*
func saveToFile(path string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	return os.WriteFile(path, data, 0644)
}
*/

func ThemeName() string {
	return Settings.ThemeName
}

func Language() string {
	return Settings.Language
}

func InstalledProjects() []string {
	return Settings.InstalledProjects
}

func ThemeVariant() fyne.ThemeVariant {
	if Settings.ThemeName == "dark" {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func SetTheme(name string) {
	if name == SystemThemeName {
		name = ""
	}
	Settings.ThemeName = name
	setupTheme()

	err := save()
	if err != nil {
		fyne.LogError("Failed on saving", err)
	}
}

func SetLanguage(name string) {
	Settings.Language = name

	err := save()
	if err != nil {
		fyne.LogError("Failed on saving", err)
	}
}

func ApplyTheme() {
	app := fyne.CurrentApp()
	if app != nil {
		settingsChanged := make(chan fyne.Settings)

		app.Settings().SetTheme(Settings.userTheme)
		app.Settings().AddChangeListener(settingsChanged)
		go func() {
			<-settingsChanged
			setupTheme()
			for _, l := range listeners {
				l()
			}
		}()

	}
}

func Theme() fyne.Theme {
	return Settings.userTheme
}

func LogoImage() *canvas.Image {
	return Settings.logo
}

func SetRememberCredentials(v bool) {
	app := fyne.CurrentApp()
	if app != nil {
		fyne.CurrentApp().Preferences().SetBool("RememberMe", v)
	}
}

func RememberCredentials() bool {
	app := fyne.CurrentApp()
	if app != nil {
		return fyne.CurrentApp().Preferences().BoolWithFallback("RememberMe", false)
	}
	return false
}

func SetCredentials(user string, password string) {
	app := fyne.CurrentApp()
	if app != nil {
		fyne.CurrentApp().Preferences().SetString("User", user)
		fyne.CurrentApp().Preferences().SetString("Password", password)
	}
}

func Credentials() (user string, password string) {
	app := fyne.CurrentApp()
	if app != nil {
		return fyne.CurrentApp().Preferences().StringWithFallback("User", ""),
			fyne.CurrentApp().Preferences().StringWithFallback("Password", "")
	}
	return "", ""
}
