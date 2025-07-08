package widgets

import (
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LabellerFunc func(any) string
type OptionSelectedFunc func(any) bool

type SelectEntryWithData struct {
	widget.BaseWidget

	obj              *widget.Select
	labeller         LabellerFunc
	selectionChecker OptionSelectedFunc

	items          []any
	selectedOption string
	selectedItem   any

	OnChanged func(item any)
}

func NewSelectWithData(
	items []any,
	labeller LabellerFunc,
	selectionChecker OptionSelectedFunc) *SelectEntryWithData {
	c := &SelectEntryWithData{
		labeller:         labeller,
		selectionChecker: selectionChecker,
		items:            items,
		selectedOption:   "",
		selectedItem:     nil,
	}

	c.ExtendBaseWidget(c)

	options := make([]string, 0)
	for _, item := range items {
		option := labeller(item)
		options = append(options, option)
		if selectionChecker(item) {
			c.selectedItem = item
			c.selectedOption = labeller(item)
		}
	}
	c.obj = widget.NewSelect(options, c.onChanged)
	c.obj.SetSelected(c.selectedOption)

	return c
}

func (c *SelectEntryWithData) SetOnChanged(callback func(item any)) {
	c.OnChanged = callback
}

func (c *SelectEntryWithData) SetSelected(s string) {
	c.obj.SetSelected(s)
}

func (c *SelectEntryWithData) onChanged(s string) {
	for _, item := range c.items {
		option := c.labeller(item)
		if s == option {
			c.selectedItem = item
			break
		}
	}
	if c.selectedItem != nil {
		c.selectedOption = s
	} else {
		c.selectedOption = ""
	}
	if c.OnChanged != nil {
		c.OnChanged(c.selectedItem)
	}
}

func (c *SelectEntryWithData) SelectedItem() any {
	return c.selectedItem
}

// Move implements fyne.Widget
func (c *SelectEntryWithData) Move(p fyne.Position) {
	c.obj.Move(p)
}

// Refresh implements fyne.Widget
func (c *SelectEntryWithData) Refresh() {
	c.obj.Refresh()
}

// CreateRenderer implements fyne.Widget
func (c *SelectEntryWithData) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.obj)
}
