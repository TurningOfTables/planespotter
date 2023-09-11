package ui

import (
	"fmt"
	"planespotter/helpers/config"
	"planespotter/helpers/types"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func InitUi(savePath string, c types.Config) fyne.App {
	icon, _ := fyne.LoadResourceFromPath("assets/plane.png")
	ui := app.New()
	ui.SetIcon(icon)
	uiWindow := ui.NewWindow("Planespotter")
	uiWindow.SetCloseIntercept(func() {
		uiWindow.Hide()
	})

	if desk, ok := ui.(desktop.App); ok {
		m := fyne.NewMenu("Planespotter", fyne.NewMenuItem("Settings", func() {
			uiWindow.Show()
		}))
		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(m)
	}
	uiWindow.Resize(fyne.NewSize(500, 200))
	uiWindow.Hide()

	title := widget.NewLabelWithStyle("Configuration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	uiLatitude := widget.NewEntry()
	uiLatitude.SetText(fmt.Sprintf("%v", c.Position.Latitude))

	uiLongitude := widget.NewEntry()
	uiLongitude.SetText(fmt.Sprintf("%v", c.Position.Longitude))

	uiSpotDistance := widget.NewEntry()
	uiSpotDistance.SetText(strconv.Itoa(c.SpotDistanceKm))

	uiCheckFreq := widget.NewEntry()
	uiCheckFreq.SetText(strconv.Itoa(c.CheckFreqSeconds))

	uiUsername := widget.NewEntry()
	uiUsername.SetText(c.ApiAuth.Username)

	uiPassword := widget.NewPasswordEntry()
	uiPassword.SetText(c.ApiAuth.Password)

	settingsForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Latitude", Widget: uiLatitude},
			{Text: "Longitude", Widget: uiLongitude},
			{Text: "OpenSky username", Widget: uiUsername},
			{Text: "OpenSky password", Widget: uiPassword},
			{Text: "Spot distance (km)", Widget: uiSpotDistance},
			{Text: "Check frequency (seconds)", Widget: uiCheckFreq},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			var newConfig types.Config
			newConfig.Position.Latitude, _ = strconv.ParseFloat(uiLatitude.Text, 64)
			newConfig.Position.Longitude, _ = strconv.ParseFloat(uiLongitude.Text, 64)
			newConfig.ApiAuth.Username = uiUsername.Text
			newConfig.ApiAuth.Password = uiPassword.Text
			newConfig.SpotDistanceKm, _ = strconv.Atoi(uiSpotDistance.Text)
			newConfig.CheckFreqSeconds, _ = strconv.Atoi(uiCheckFreq.Text)
			config.SaveConfig(savePath, newConfig)
		},
	}

	uiWindow.SetContent(container.NewVBox(title, settingsForm))
	return ui
}
