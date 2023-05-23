package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"terra9.it/vadovia/assets"
	"terra9.it/vadovia/internal"
)

// AboutView displays the logo and a version link for application information.
func aboutView() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(assets.ResourceLogovadoviacompactPng)
	logo.FillMode = canvas.ImageFillOriginal
	content := widget.NewRichTextFromMarkdown(internal.Version.Info())
	content.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		container.NewCenter(logo),
		nil, nil, nil,
		content,
	)
}
