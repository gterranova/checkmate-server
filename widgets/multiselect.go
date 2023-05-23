package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type itemObjCreateFunc func(item *MultiSelectListItem) fyne.CanvasObject
type itemObjUpdateFunc func(item *MultiSelectListItem, cnvObj fyne.CanvasObject)

type MultiSelectList struct {
	widget.BaseWidget

	container *fyne.Container
	items     []*MultiSelectListItem
	layout    *MultiSelectListLayout

	callback func(*MultiSelectListItem)
	disabled bool

	objCreateFunc itemObjCreateFunc
	objUpdateFunc itemObjUpdateFunc
}

type MultiSelectListItem struct {
	Index    int
	Checked  bool
	Disabled bool
	Value    interface{}
}

type MultiSelectListObject struct {
	widget.BaseWidget

	obj   fyne.CanvasObject
	image *canvas.Image

	item *MultiSelectListItem
	list *MultiSelectList
}

// Tapped implements fyne.Tappable
func (o *MultiSelectListObject) Tapped(*fyne.PointEvent) {
	if !o.list.Disabled() && !o.item.Disabled {
		o.item.Checked = !o.item.Checked
		if o.item.Checked {
			o.image.Resource = theme.CheckButtonCheckedIcon()
		} else {
			o.image.Resource = theme.CheckButtonIcon()
		}
		o.list.OnChange(o.item)
		o.Refresh()
	}
}

// Tapped implements fyne.Tappable
func (l *MultiSelectList) OnChange(item *MultiSelectListItem) {
	if l.callback != nil {
		l.callback(item)
	}
	//l.Refresh()
	l.SetItems(l.GetItems())
}

// CreateRenderer implements fyne.Widget
func (o *MultiSelectListObject) CreateRenderer() fyne.WidgetRenderer {
	// this encapsulate the image and the button
	box := container.New(layout.NewHBoxLayout(), o.image, o.obj)
	return widget.NewSimpleRenderer(box)
}

func (o *MultiSelectListObject) SetItem(item *MultiSelectListItem) {
	o.item = item
	if o.item.Checked {
		o.image.Resource = theme.CheckButtonCheckedIcon()
	} else {
		o.image.Resource = theme.CheckButtonIcon()
	}
	if o.list.Disabled() || o.item.Disabled {
		o.image.Translucency = 0.5
	}
	o.list.objUpdateFunc(item, o.obj)
	o.Refresh()
}

func NewMultiSelectListObject(list *MultiSelectList, item *MultiSelectListItem) *MultiSelectListObject {
	obj := &MultiSelectListObject{item: item, list: list}
	obj.image = canvas.NewImageFromResource(theme.CheckButtonIcon())
	obj.image.SetMinSize(fyne.NewSize(20, 20))
	obj.image.FillMode = canvas.ImageFillContain

	obj.obj = list.objCreateFunc(item)

	obj.ExtendBaseWidget(obj)
	return obj
}

func NewMultiSelectList(items []*MultiSelectListItem, itemCreateFunc itemObjCreateFunc, itemUpdateFunc itemObjUpdateFunc, OnChange func(*MultiSelectListItem)) *MultiSelectList {
	list := &MultiSelectList{callback: OnChange, objCreateFunc: itemCreateFunc, objUpdateFunc: itemUpdateFunc}
	list.layout = &MultiSelectListLayout{d: list}
	list.items = items
	list.container = container.New(list.layout, list.makeList())
	list.ExtendBaseWidget(list)

	return list
}

func (l *MultiSelectList) GetValue() []*MultiSelectListItem {
	ret := make([]*MultiSelectListItem, 0)
	for _, item := range l.items {
		if item.Checked {
			ret = append(ret, item)
		}
	}
	return ret
}

func (l *MultiSelectList) GetItems() []*MultiSelectListItem {
	return l.items
}

func (l *MultiSelectList) SetItems(items []*MultiSelectListItem) {
	l.items = items
	l.container.RemoveAll()
	l.container.Add(l.makeList())
	l.container.Refresh()
	l.Refresh()
}

func (l *MultiSelectList) makeList() *widget.List {
	return widget.NewList(
		func() int {
			return len(l.items)
		},
		func() fyne.CanvasObject {
			return NewMultiSelectListObject(l, &MultiSelectListItem{})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*MultiSelectListObject).SetItem(l.items[i])
		})
}

func (l *MultiSelectList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.container)
}

func (l *MultiSelectList) Enable() {
	l.disabled = false
	l.SetItems(l.GetItems())
}

func (l *MultiSelectList) Disable() {
	l.disabled = true
	l.SetItems(l.GetItems())
}
func (l *MultiSelectList) Disabled() bool {
	return l.disabled
}

type MultiSelectListLayout struct {
	d *MultiSelectList
}

func (l *MultiSelectListLayout) Layout(obj []fyne.CanvasObject, size fyne.Size) {

	// icon
	if len(obj) > 0 {
		//height := fyne.Min(float32(len(l.d.items)*int(size.Height))+theme.Padding()*2, 34*6)
		obj[0].Resize(fyne.NewSize(size.Width, size.Height))
	}
}

func (l *MultiSelectListLayout) MinSize(obj []fyne.CanvasObject) fyne.Size {
	if obj == nil {
		return fyne.NewSize(0, 0)
	}
	contentMin := obj[0].MinSize()

	width := contentMin.Width
	height := fyne.Min(float32(len(l.d.items)*int(contentMin.Height))+theme.Padding()*2, 34*6)

	return fyne.NewSize(width, height)
}
