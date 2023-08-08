package commons

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func SknSelectResource(alias string) fyne.Resource {
	return sknImageByName(alias, false, false).Resource
}
func SknSelectThemedResource(alias string) fyne.Resource {
	return sknImageByName(alias, true, false).Resource
}

func SknSelectImage(alias string) *canvas.Image {
	return sknImageByName(alias, false, false)
}
func SknSelectThemedImage(alias string) *canvas.Image {
	return sknImageByName(alias, true, false)
}
func SknSelectThemedInvertedImage(alias string) *canvas.Image {
	return sknImageByName(alias, true, true)
}

func sknImageByName(alias string, themed bool, inverted bool) *canvas.Image {
	var selected fyne.Resource

	switch alias {
	case "apcupsd":
		selected = resourceApcupsdPng
	case "charging":
		selected = resourceChargingPng
	case "preferences":
		selected = resourceGapcprefsPng
	case "onbattery":
		selected = resourceOnbattPng
	case "online":
		selected = resourceOnlinePng
	case "unknown":
		selected = resourceUnpluggedPng
	default:
		selected = resourceApcupsdPng
	}

	image := canvas.NewImageFromResource(selected)
	if themed {
		image = canvas.NewImageFromResource(theme.NewThemedResource(selected))
	}
	if inverted {
		image = canvas.NewImageFromResource(theme.NewInvertedThemedResource(selected))
	}

	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScaleSmooth
	return image
}
