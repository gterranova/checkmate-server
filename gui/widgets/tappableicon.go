package widgets

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type BehaviourType int

const (
	Tappable BehaviourType = iota
	Toggleable
)

type TappableIcon struct {
	widget.BaseWidget
	Resource   fyne.Resource
	Importance widget.ButtonImportance
	OnTapped   func()

	disabled, focused, hovered, toggled, inverted bool

	size         float32
	beaviourType BehaviourType
	propertyLock sync.RWMutex
	cachedRes    fyne.Resource
	renderer     fyne.WidgetRenderer
}

func NewTappableIcon(res fyne.Resource, size float32, beaviourType BehaviourType, onTapped func(), inverted bool) *TappableIcon {
	icon := &TappableIcon{Resource: res, size: size, OnTapped: onTapped, beaviourType: beaviourType, inverted: inverted}
	icon.ExtendBaseWidget(icon)
	return icon
}

func NewTappableIconWithMenu(res fyne.Resource, size float32, menu *fyne.Menu, inverted bool) *TappableIcon {
	icon := &TappableIcon{Resource: res, size: size, beaviourType: Toggleable, inverted: inverted}
	icon.OnTapped = func() {
		h := icon.MinSize().Height
		//holder := aw.app.Driver().CanvasForObject(icon)
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(icon).AddXY(-h, h+theme.InnerPadding()*2)
		//holder.Focus(icon)
		popup := widget.NewPopUpMenu(menu, fyne.CurrentApp().Driver().CanvasForObject(icon))
		popup.ShowAtPosition(pos)
		popup.OnDismiss = func() {
			icon.UnToggle()
			popup.Hide()
		}
	}
	icon.ExtendBaseWidget(icon)
	return icon
}

func (o *TappableIcon) SetResource(res fyne.Resource) {
	o.Resource = res
	o.renderer.(*iconRenderer).SetResource(res)
	o.Refresh()
}

// Tapped implements fyne.Tappable
func (o *TappableIcon) Tapped(*fyne.PointEvent) {
	if o.Disabled() {
		return
	}
	if o.beaviourType == Toggleable {
		o.toggled = !o.toggled
		if !o.toggled {
			return
		}
	}

	if o.OnTapped != nil {
		o.OnTapped()
	}
	o.Refresh()
}

func (o *TappableIcon) Toggle() {
	if o.beaviourType == Toggleable {
		o.toggled = true
		o.Refresh()
	}
}

func (o *TappableIcon) UnToggle() {
	if o.beaviourType == Toggleable {
		o.toggled = false
		o.Refresh()
	}
}

func (o *TappableIcon) Toggled() bool {
	return o.toggled
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (i *TappableIcon) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	i.propertyLock.RLock()
	defer i.propertyLock.RUnlock()

	i.renderer = NewTappableIconRenderer(i)
	return i.renderer
}

func (l *TappableIcon) Enable() {
	l.disabled = false
	l.Refresh()
}

func (l *TappableIcon) Disable() {
	l.disabled = true
	l.Refresh()
}
func (l *TappableIcon) Disabled() bool {
	return l.disabled
}

// Cursor returns the cursor type of this widget
func (b *TappableIcon) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (b *TappableIcon) FocusGained() {
	b.focused = true
	b.Refresh()
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (b *TappableIcon) FocusLost() {
	b.focused = false
	b.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (b *TappableIcon) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (b *TappableIcon) MouseIn(*desktop.MouseEvent) {
	b.hovered = true

	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *TappableIcon) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (b *TappableIcon) MouseOut() {
	b.hovered = false

	b.Refresh()
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (b *TappableIcon) TypedRune(rune) {
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (b *TappableIcon) TypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeySpace {
		b.Tapped(nil)
	}
}

func (b *TappableIcon) Refresh() {
	if b.renderer != nil {
		b.renderer.Refresh()
	}
}

func (b *TappableIcon) backgroundColor() color.Color {
	switch {
	case b.Disabled():
		if b.Importance == widget.LowImportance {
			return color.Transparent
		}
		return theme.DisabledButtonColor()
	case b.focused:
		return blendColor(theme.ButtonColor(), theme.FocusColor())
	case b.hovered, b.toggled:
		bg := theme.ButtonColor()
		switch b.Importance {
		case widget.HighImportance:
			bg = theme.PrimaryColor()
		case widget.DangerImportance:
			bg = theme.ErrorColor()
		case widget.WarningImportance:
			bg = theme.WarningColor()
		}

		return blendColor(bg, theme.HoverColor())
	case b.Importance == widget.HighImportance:
		return theme.PrimaryColor()
	case b.Importance == widget.LowImportance:
		return color.Transparent
	case b.Importance == widget.DangerImportance:
		return theme.ErrorColor()
	case b.Importance == widget.WarningImportance:
		return theme.WarningColor()
	default:
		return theme.ButtonColor()
	}
}

type iconRenderer struct {
	raster  *canvas.Image
	objects []fyne.CanvasObject

	tappableIcon *TappableIcon
	background   *canvas.Rectangle
}

func NewTappableIconRenderer(i *TappableIcon) *iconRenderer {
	img := canvas.NewImageFromResource(i.Resource)
	img.FillMode = canvas.ImageFillContain
	r := &iconRenderer{
		tappableIcon: i,
		raster:       img,
	}
	if i.inverted {
		r.SetObjects([]fyne.CanvasObject{r.raster})
	} else {
		r.background = canvas.NewRectangle(theme.ButtonColor())
		r.SetObjects([]fyne.CanvasObject{r.background, r.raster})
	}
	r.applyTheme()
	return r
}

func (r *iconRenderer) SetResource(res fyne.Resource) {
	r.tappableIcon.propertyLock.RLock()
	r.raster = canvas.NewImageFromResource(r.tappableIcon.Resource)
	r.raster.FillMode = canvas.ImageFillContain
	if r.tappableIcon.inverted {
		r.SetObjects([]fyne.CanvasObject{r.raster})
	} else {
		r.SetObjects([]fyne.CanvasObject{r.background, r.raster})
	}
	r.tappableIcon.propertyLock.RUnlock()
}

// Destroy does nothing in the base implementation.
//
// Implements: fyne.WidgetRenderer
func (r *iconRenderer) Destroy() {
}

// Objects returns the objects that should be rendered.
//
// Implements: fyne.WidgetRenderer
func (r *iconRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// SetObjects updates the objects of the renderer.
func (r *iconRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

func (r *iconRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.InnerPadding()*2, theme.InnerPadding()*2)
}

func (i *iconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(i.tappableIcon.size, i.tappableIcon.size).Add(i.padding())
}

func (i *iconRenderer) Layout(size fyne.Size) {
	if len(i.Objects()) == 0 {
		return
	}

	iconSize := fyne.NewSize(i.tappableIcon.size, i.tappableIcon.size)
	padding := i.padding()
	pos := fyne.Position{
		X: padding.Width / 2,
		Y: (size.Height - iconSize.Height) / 2,
	}

	if i.background == nil {
		i.objects[0].Move(pos)
		i.objects[0].Resize(iconSize)
	} else {
		i.objects[0].Move(pos)
		i.objects[0].Resize(iconSize)
		i.objects[1].Move(pos)
		i.objects[1].Resize(iconSize)
	}
	//for _, obj := range i.objects {
	//	obj.Resize(size)
	//}
}

func (i *iconRenderer) Refresh() {
	i.tappableIcon.propertyLock.RLock()
	i.applyTheme()
	i.tappableIcon.propertyLock.RUnlock()

	for _, obj := range i.objects {
		obj.Refresh()
	}

}

// applyTheme updates this button to match the current theme
func (r *iconRenderer) applyTheme() {
	if !r.tappableIcon.Visible() {
		r.tappableIcon.Hide()
		return
	}
	if !r.tappableIcon.inverted && r.background != nil {
		r.background.FillColor = r.tappableIcon.backgroundColor()
		r.background.Refresh()
	}
	if r.tappableIcon.Resource != nil {
		switch res := r.tappableIcon.Resource.(type) {
		case *theme.ThemedResource:
			if r.tappableIcon.inverted {
				var cln fyne.ThemeColorName
				switch {
				case r.tappableIcon.disabled:
					cln = theme.ColorNameDisabled
					r.tappableIcon.cachedRes = theme.NewDisabledResource(r.tappableIcon.Resource)

				case r.tappableIcon.toggled:
					cln = theme.ColorNamePrimary
				case r.tappableIcon.focused:
					cln = theme.ColorNameFocus
				case r.tappableIcon.hovered:
					cln = theme.ColorNamePlaceHolder

				}
				//if r.tappableIcon.Importance == widget.HighImportance || r.tappableIcon.Importance == widget.DangerImportance || r.tappableIcon.Importance == widget.WarningImportance {
				res.ColorName = cln
				r.tappableIcon.cachedRes = theme.NewInvertedThemedResource(res)
			} else if r.tappableIcon.Importance != widget.HighImportance && r.tappableIcon.Importance != widget.DangerImportance && r.tappableIcon.Importance != widget.WarningImportance {
				r.tappableIcon.cachedRes = res
			}
			//}
		case *theme.InvertedThemedResource:

			if r.tappableIcon.inverted && r.tappableIcon.Importance != widget.HighImportance && r.tappableIcon.Importance != widget.DangerImportance && r.tappableIcon.Importance != widget.WarningImportance {
				r.tappableIcon.cachedRes = theme.NewInvertedThemedResource(res)
			} else {
				r.tappableIcon.cachedRes = res.Original()
			}
		}

		if r.tappableIcon.Resource == r.tappableIcon.cachedRes {
			return
		}
		r.raster = canvas.NewImageFromResource(r.tappableIcon.cachedRes)
		r.tappableIcon.Show()
	}
}

func blendColor(under, over color.Color) color.Color {
	// This alpha blends with the over operator, and accounts for RGBA() returning alpha-premultiplied values
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()

	srcAlpha := float32(srcA) / 0xFFFF
	dstAlpha := float32(dstA) / 0xFFFF

	outAlpha := srcAlpha + dstAlpha*(1-srcAlpha)
	outR := srcR + uint32(float32(dstR)*(1-srcAlpha))
	outG := srcG + uint32(float32(dstG)*(1-srcAlpha))
	outB := srcB + uint32(float32(dstB)*(1-srcAlpha))
	// We create an RGBA64 here because the color components are already alpha-premultiplied 16-bit values (they're just stored in uint32s).
	return color.RGBA64{R: uint16(outR), G: uint16(outG), B: uint16(outB), A: uint16(outAlpha * 0xFFFF)}

}
