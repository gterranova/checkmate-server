package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	xwidget "fyne.io/x/fyne/widget"

	"terra9.it/checkmate/core"
	"terra9.it/checkmate/gui/widgets"
	"terra9.it/checkmate/gui/wizard"
)

// MultiselectStep implement a wizard step based on default BasewizardStep
type MultiselectStep struct {
	wizard.BaseWizardStep
	project   *core.Project
	w         *mainWindow
	feature   core.Feature
	canvasObj fyne.CanvasObject
}

// Creates and returns the introduction step pane
func NewMultiselectStep(project *core.Project, w *mainWindow, feature core.Feature) *MultiselectStep {
	step := &MultiselectStep{project: project, w: w, feature: feature}
	step.Content = step.CreateContent()
	step.Caption = "" //feature.GetCaption()
	step.Title = feature.GetTitle()
	return step
}

func (step *MultiselectStep) CreateContent() *fyne.Container {
	label := step.feature.GetTitle()
	titleLabel := canvas.NewText(label, theme.ForegroundColor())
	titleLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	titleLabel.TextSize = wizard.ContentTitleTextSize

	captionLabel := canvas.NewText(step.feature.GetTitle(), theme.DisabledColor())
	captionLabel.TextStyle = fyne.TextStyle{
		Bold: false,
	}
	captionLabel.TextSize = wizard.ContentCaptionTextSize
	if len(step.feature.GetTitle()) == 0 {
		captionLabel.Hide()
	}

	switch step.feature.(type) {
	case *core.Checklist:
		step.canvasObj = step.CreateChecklistContainer()
	case *core.Checkform:
		step.canvasObj = step.CreateFormContainer()
	default:
		step.canvasObj = nil
	}
	vb := container.NewBorder(container.NewPadded(container.NewVBox(titleLabel, captionLabel, widget.NewSeparator())),
		nil, nil, nil, container.NewPadded(step.canvasObj))
	return vb
}

func (step *MultiselectStep) CreateChecklistContainer() *widgets.MultiSelectList {
	options := make([]*widgets.MultiSelectListItem, 0)
	for i, c := range step.feature.GetChildren() {
		options = append(options, &widgets.MultiSelectListItem{
			Index:    i,
			Value:    c.(*core.Checkbox),
			Checked:  c.GetValue().(bool),
			Disabled: c.IsDisabled(),
		})
	}
	return widgets.NewMultiSelectList(options,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(core.Feature).GetTitle())
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(core.Feature)
			step.project.SetFeature(condition.GetTag(), item.Checked)
			//condition.Checked = item.Checked
			step.w.Update(step.project)
		},
	)
}

func (step *MultiselectStep) CreateFormContainer() fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, child := range step.feature.GetChildren() {
		switch t := child.(type) {
		case *core.String:

			inputFeature := child
			control := widget.NewEntry()
			if inputFeature.GetValue() != nil {
				control.SetText(inputFeature.GetValue().(string))
			}
			control.OnChanged = func(s string) {
				step.project.SetFeature(inputFeature.GetTag(), s)
				step.w.Update(step.project)
			}
			controlLabel := widget.NewLabel(inputFeature.GetTitle())
			controlLabel.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			controlLabel.Alignment = fyne.TextAlignLeading

			objects = append(objects, controlLabel, control)
		case *core.Number:
			inputFeature := t
			control := xwidget.NewNumericalEntry()
			control.SetText(fmt.Sprintf("%v", inputFeature.Value))
			control.OnChanged = func(s string) {
				value, _ := strconv.ParseInt(s, 10, 64)
				step.project.SetFeature(inputFeature.Tag, value)
				step.w.Update(step.project)
			}
			controlLabel := widget.NewLabel(inputFeature.Title)
			controlLabel.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			controlLabel.Alignment = fyne.TextAlignLeading
			objects = append(objects, controlLabel, control)
		case *core.Select:
			selectFeature := t
			options := make([]any, 0)
			for _, option := range selectFeature.GetChildren() {
				if !option.IsDisabled() {
					options = append(options, option)
				}
			}

			control := widgets.NewSelectWithData(options,
				func(i any) string {
					c := i.(core.Feature)
					//if len(c.Reference) > 0 {
					//	return fmt.Sprintf("%v (%v)", c.Title, c.Reference)
					//}
					return c.GetTitle()
				},
				func(i any) bool {
					c := i.(core.Feature)
					return c != nil && c.GetValue() != nil && c.GetValue().(bool)
				},
			)

			control.OnChanged = func(item any) {
				for _, c := range options {
					c.(core.Feature).SetValue(false)
					step.project.SetFeature(c.(core.Feature).GetTag(), false)
				}
				feature := item.(core.Feature)
				step.project.SetFeature(feature.GetTag(), true)
				//fmt.Println(step.project.Tags, control.SelectedItem())
				step.w.Update(step.project)
			}
			controlLabel := widget.NewLabel(selectFeature.Title)
			controlLabel.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			controlLabel.Alignment = fyne.TextAlignLeading

			objects = append(objects, controlLabel, control)
		default:
			panic("Unknown type")
		}
	}
	return container.New(layout.NewFormLayout(), objects...)
}

func (step *MultiselectStep) Disabled() bool {
	return step.feature.IsDisabled()
}

func (step *MultiselectStep) OnLeave() {
	step.w.wizard.RebuildStepsContainer()
}

func (step *MultiselectStep) Refresh() {
	switch step.feature.(type) {
	case *core.Checklist:
		obj := step.canvasObj.(*widgets.MultiSelectList)
		options := make([]*widgets.MultiSelectListItem, 0)
		for i, c := range step.feature.GetChildren() {
			options = append(options, &widgets.MultiSelectListItem{
				Index:    i,
				Value:    c,
				Checked:  c.GetValue().(bool),
				Disabled: c.IsDisabled(),
			})
		}
		obj.SetItems(options)
	case *core.Checkform:
		objects := step.canvasObj.(*fyne.Container).Objects
		for i, child := range step.feature.GetChildren() {
			if child.IsDisabled() {
				objects[i*2].Hide()
				objects[i*2+1].Hide()
			} else {
				objects[i*2].Show()
				objects[i*2+1].Show()
			}
		}
		//default:
		//	step.Content.Refresh()
	}
}
