package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*myThemeLight)(nil)
var _ fyne.Theme = (*myThemeDark)(nil)

type myThemeDark struct{}
type myThemeLight struct{ myThemeDark }

type ThemeResourceName string

func (m myThemeDark) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return color.NRGBA{0x18, 0x0C, 0x27, 0xFF}
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNameInputBackground:
		return color.NRGBA{0x00, 0x00, 0x00, 0x00}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{0xFF, 0xFF, 0xFF, 0x40}
	case theme.ColorNameFocus:
		return theme.HoverColor()
	case theme.ColorNameError, theme.ColorNameWarning:
		return color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x39, G: 0x39, B: 0x3a, A: 0xff}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0f}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 0x39, G: 0x39, B: 0x3a, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x66}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x99}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{A: 0x66}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m myThemeDark) Font(style fyne.TextStyle) fyne.Resource {

	if style.Bold && style.Italic {
		//Heavy with italics
		return theme.DefaultTextBoldItalicFont()
	} else if style.Bold {
		//heavy
		return theme.DefaultTextBoldFont()
	} else if style.Monospace {
		//Spaced out smaller font
		return theme.DefaultTextMonospaceFont()
	}
	//standard bold
	return theme.DefaultTextFont()

}

func (m myThemeDark) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myThemeDark) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m myThemeLight) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	//case theme.ColorNamePrimary:
	//	return color.NRGBA{0xF5, 0x24, 0x7D, 0xFF}
	case theme.ColorNameError, theme.ColorNameWarning:
		return color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	}

	return theme.DefaultTheme().Color(name, theme.VariantLight)
}
