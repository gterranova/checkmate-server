package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Group widget is list of widgets that contains a visual border around the list and a group title at the top.
type Group struct {
	widget.BaseWidget

	Text    string
	box     *fyne.Container
	content fyne.CanvasObject
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (g *Group) Resize(size fyne.Size) {
	g.BaseWidget.Resize(size)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (g *Group) Move(pos fyne.Position) {
	g.BaseWidget.Move(pos)
}

// MinSize returns the smallest size this widget can shrink to
func (g *Group) MinSize() fyne.Size {
	return g.BaseWidget.MinSize()
}

// Show this widget, if it was previously hidden
func (g *Group) Show() {
	g.BaseWidget.Show()
}

// Hide this widget, if it was previously visible
func (g *Group) Hide() {
	g.BaseWidget.Hide()
}

// Append adds a new CanvasObject to the end of the group
func (g *Group) Append(object fyne.CanvasObject) {
	g.box.Add(object)

	g.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (g *Group) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(g.Text)
	label.TextStyle = fyne.TextStyle{Bold: true}

	labelBg := canvas.NewRectangle(theme.BackgroundColor())
	line := canvas.NewRectangle(theme.ButtonColor())
	objects := []fyne.CanvasObject{line, labelBg, label, g.content}
	return &groupRenderer{label: label, line: line, labelBg: labelBg,
		objects: objects, group: g}
}

// NewGroup creates a new grouped list widget with a title and the specified list of child objects.
func NewGroup(title string, children ...fyne.CanvasObject) *Group {
	box := container.NewVBox(children...)
	group := &Group{widget.BaseWidget{}, title, box, box}
	group.ExtendBaseWidget(group)
	//Renderer(group).Layout(group.MinSize())
	return group
}

// NewGroupWithScroller creates a new grouped list widget with a title and the specified list of child objects.
// This group will scroll when the available space is less than needed to display the items it contains.
func NewGroupWithScroller(title string, children ...fyne.CanvasObject) *Group {
	box := container.NewVBox(children...)
	group := &Group{widget.BaseWidget{}, title, box, container.NewScroll(box)}
	group.ExtendBaseWidget(group)
	//Renderer(group).Layout(group.MinSize())
	return group
}

type groupRenderer struct {
	label         *widget.Label
	line, labelBg *canvas.Rectangle

	objects []fyne.CanvasObject
	group   *Group
}

func (g *groupRenderer) MinSize() fyne.Size {
	labelMin := g.label.MinSize()
	groupMin := g.group.content.MinSize()

	return fyne.NewSize(fyne.Max(labelMin.Width, groupMin.Width),
		labelMin.Height+groupMin.Height+theme.Padding())
}

func (g *groupRenderer) Layout(size fyne.Size) {
	labelWidth := g.label.MinSize().Width
	labelHeight := g.label.MinSize().Height

	g.line.Move(fyne.NewPos(0, (labelHeight-theme.Padding())/2))
	g.line.Resize(fyne.NewSize(size.Width, theme.Padding()))

	g.labelBg.Move(fyne.NewPos(size.Width/2-labelWidth/2, 0))
	g.labelBg.Resize(g.label.MinSize())
	g.label.Move(fyne.NewPos(size.Width/2-labelWidth/2, 0))
	g.label.Resize(g.label.MinSize())

	g.group.content.Move(fyne.NewPos(0, labelHeight+theme.Padding()))
	g.group.content.Resize(fyne.NewSize(size.Width, size.Height-labelHeight-theme.Padding()))
}

func (g *groupRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (g *groupRenderer) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *groupRenderer) Refresh() {
	g.line.FillColor = theme.ButtonColor()
	g.labelBg.FillColor = theme.BackgroundColor()

	g.label.TextStyle = fyne.TextStyle{Bold: true}
	g.label.SetText(g.group.Text)

	g.Layout(g.group.Size())

	canvas.Refresh(g.group)
}

func (g *groupRenderer) Destroy() {
}
