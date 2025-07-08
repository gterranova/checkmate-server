package wizard

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type navLabelRenderer struct {
	icon  *canvas.Image
	label *canvas.Text

	objects []fyne.CanvasObject
	button  *navLabelButton
}

const navLabelSize = 10
const iconSize = 16

// MinSize calculates the minimum size of a navLabel button. A fixed amount.
func (b *navLabelRenderer) MinSize() fyne.Size {
	textMin := fyne.MeasureText(b.label.Text, navLabelSize, fyne.TextStyle{Bold: true})
	if b.button.icon == nil {
		return textMin
	}
	return textMin.AddWidthHeight(iconSize, 0)
}

// Layout the components of the widget
func (b *navLabelRenderer) Layout(size fyne.Size) {
	//inner := size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	//b.icon.Resize(inner)
	b.icon.Move(fyne.NewPos(0, (size.Height-iconSize)/2))

	//textSize := size.Height * .67
	//textMin := fyne.MeasureText(b.label.Text, textSize, fyne.TextStyle{Bold: true})
	textMin := b.MinSize()

	//b.label.TextSize = 10 //textSize
	if b.button.icon == nil {
		b.label.Resize(fyne.NewSize(size.Width, textMin.Height))
		b.label.Move(fyne.NewPos(0, (size.Height-textMin.Height)/2))
	} else {
		b.label.Resize(fyne.NewSize(size.Width-iconSize-theme.Padding(), textMin.Height))
		b.label.Move(fyne.NewPos(iconSize+theme.Padding(), (size.Height-textMin.Height)/2))
	}
}

// ApplyTheme is called when the navLabelButton may need to update it's look
func (b *navLabelRenderer) ApplyTheme() {
	b.label.Color = theme.ForegroundColor()
	b.Refresh()
}

func (b *navLabelRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

func (b *navLabelRenderer) Refresh() {
	b.label.Text = b.button.text

	b.icon.Hidden = b.button.icon == nil
	if b.button.icon != nil {
		b.icon.Resource = b.button.icon
	}

	b.Layout(b.button.Size())
	canvas.Refresh(b.button)
}

func (b *navLabelRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *navLabelRenderer) Destroy() {
}

// navLabelButton widget is a scalable button that has a text label and icon and triggers an event func when clicked
type navLabelButton struct {
	widget.BaseWidget
	text  string
	icon  fyne.Resource
	color color.Color

	tap func(bool)
}

// Tapped is called when a regular tap is reported
func (b *navLabelButton) Tapped(ev *fyne.PointEvent) {
	b.tap(true)
}

// TappedSecondary is called when an alternative tap is reported
func (b *navLabelButton) TappedSecondary(ev *fyne.PointEvent) {
	b.tap(false)
}

func (b *navLabelButton) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(b.text, theme.ForegroundColor())
	text.TextSize = navLabelSize
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle.Bold = true
	text.Color = b.color

	icon := canvas.NewImageFromResource(b.icon)
	icon.Resize(fyne.NewSize(iconSize, iconSize))
	icon.FillMode = canvas.ImageFillStretch

	objects := []fyne.CanvasObject{
		text,
		icon,
	}

	return &navLabelRenderer{icon, text, objects, b}
	//return widget.NewSimpleRenderer(text)
}

// SetText allows the button label to be changed
func (b *navLabelButton) SetText(text string) {
	b.text = text

	b.Refresh()
}

// SetIcon updates the icon on a label - pass nil to hide an icon
func (b *navLabelButton) SetIcon(icon fyne.Resource) {
	b.icon = icon

	b.Refresh()
}

// newNavigationLabel creates a new button widget with the specified label, themed icon and tap handler
func newNavigationLabel(label string, color color.Color, icon fyne.Resource, tap func(bool)) *navLabelButton {
	button := &navLabelButton{text: label, icon: icon, color: color, tap: tap}
	button.ExtendBaseWidget(button)
	return button
}
