package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (v *viewProvider) GlossaryPage() *fyne.Container {
	gtext := `
# GGAPCMON 

A monitor for UPS's under the management of the APCUPSD applicaiton. 

Application Under Construction

Send comments to skoona at gmail dot com

`

	rtext := widget.NewRichTextFromMarkdown("")
	rtext.ParseMarkdown(gtext)
	rtext.Wrapping = fyne.TextWrapWord
	return container.NewMax(
		container.NewVScroll(
			rtext,
		),
	)
}
