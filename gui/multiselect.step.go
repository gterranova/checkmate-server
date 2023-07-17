package gui

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"terra9.it/vadovia/core"
	"terra9.it/vadovia/core/features"
	"terra9.it/vadovia/gui/widgets"
	"terra9.it/vadovia/gui/wizard"
)

// MultiselectStep implement a wizard step based on default BasewizardStep
type MultiselectStep struct {
	wizard.BaseWizardStep
	project   *core.Project
	w         *mainWindow
	feature   *features.Feature
	canvasObj fyne.CanvasObject
}

// Creates and returns the introduction step pane
func NewMultiselectStep(project *core.Project, w *mainWindow, feature *features.Feature) *MultiselectStep {
	step := &MultiselectStep{project: project, w: w, feature: feature}
	step.Content = step.CreateContent()
	step.Caption = feature.Caption
	step.Title = feature.Title
	return step
}

func (step *MultiselectStep) CreateContent() *fyne.Container {
	titleLabel := canvas.NewText(step.feature.Title, color.Black)
	titleLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	titleLabel.TextSize = wizard.ContentTitleTextSize

	captionLabel := canvas.NewText(step.feature.Caption, color.Gray16{0x8888})
	captionLabel.TextStyle = fyne.TextStyle{
		Bold: false,
	}
	captionLabel.TextSize = wizard.ContentCaptionTextSize
	if len(step.feature.Caption) == 0 {
		captionLabel.Hide()
	}

	switch step.feature.Type {
	case "checklist":
		step.canvasObj = step.CreateChecklistContainer()
	case "form":
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
			Value:    c,
			Checked:  c.Value != nil && c.Value.(bool),
			Disabled: c.Disabled,
		})
	}
	return widgets.NewMultiSelectList(options,
		func(item *widgets.MultiSelectListItem) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item *widgets.MultiSelectListItem, cnvObj fyne.CanvasObject) {
			label := cnvObj.(*widget.Label)
			label.SetText(item.Value.(*features.Feature).Title)
		},
		func(item *widgets.MultiSelectListItem) {
			condition := item.Value.(*features.Feature)
			step.project.SetFeature(condition.Tag, item.Checked)
			//condition.Checked = item.Checked
			step.w.Update(step.project)
		},
	)
}

func (step *MultiselectStep) CreateFormContainer() fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, child := range step.feature.GetChildren() {
		switch child.Type {
		case "input":
			inputFeature := child
			control := widget.NewEntry()
			if inputFeature.Value != nil {
				control.SetText(fmt.Sprintf("%v", inputFeature.Value))
			}
			control.OnChanged = func(s string) {
				power, _ := strconv.ParseFloat(s, 32)
				step.project.SetFeature(inputFeature.Tag, power)
				step.w.Update(step.project)
			}
			controlLabel := widget.NewLabel(inputFeature.Title)
			controlLabel.TextStyle = fyne.TextStyle{
				Bold: true,
			}
			controlLabel.Alignment = fyne.TextAlignLeading

			objects = append(objects, controlLabel, control)
		case "select":
			selectFeature := child
			options := make([]interface{}, 0)
			for _, option := range selectFeature.GetChildren() {
				if !option.Disabled {
					options = append(options, option)
				}
			}

			control := widgets.NewSelectWithData(options,
				func(i interface{}) string {
					c := i.(*features.Feature)
					if len(c.Reference) > 0 {
						return fmt.Sprintf("%v (%v)", c.Title, c.Reference)
					}
					return c.Title
				},
				func(i interface{}) bool {
					c := i.(*features.Feature)
					return c != nil && c.Value != nil && c.Value.(bool)
				},
			)

			control.OnChanged = func(item interface{}) {
				for _, c := range options {
					c.(*features.Feature).Value = false
					//step.project.SetFeature(c.(*features.Feature).Tag, false)
				}
				feature := item.(*features.Feature)
				step.project.SetFeature(feature.Tag, true)
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
	return step.feature.Disabled
}

func (step *MultiselectStep) OnLeave() {
	step.w.wizard.RebuildStepsContainer()
}

func (step *MultiselectStep) Refresh() {
	switch step.feature.Type {
	case "checklist":
		obj := step.canvasObj.(*widgets.MultiSelectList)
		options := make([]*widgets.MultiSelectListItem, 0)
		for i, c := range step.feature.GetChildren() {
			options = append(options, &widgets.MultiSelectListItem{
				Index:    i,
				Value:    c,
				Checked:  c.Value != nil && c.Value.(bool),
				Disabled: c.Disabled,
			})
		}
		obj.SetItems(options)
	default:
		step.Content.Refresh()
	}
}
