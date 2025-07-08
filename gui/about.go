package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"terra9.it/checkmate/internal"
)

// AboutView displays the logo and a version link for application information.
func aboutView(app *application) fyne.CanvasObject {
	imageRes := &fyne.StaticResource{
		StaticContent: app.loader.MustGet("icon.png"),
	}

	logo := canvas.NewImageFromResource(imageRes)
	logo.FillMode = canvas.ImageFillOriginal
	content := widget.NewRichTextFromMarkdown(internal.Version.Info())
	content.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		container.NewCenter(logo),
		nil, nil, nil,
		content,
	)
}
