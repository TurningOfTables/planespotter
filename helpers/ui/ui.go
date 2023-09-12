package ui

import (
	"fmt"
	"planespotter/helpers/save"
	"planespotter/helpers/types"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func WindowSetup(app fyne.App, icon fyne.Resource) fyne.Window {
	window := app.NewWindow("Planespotter")
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	if desk, ok := app.(desktop.App); ok {
		m := fyne.NewMenu("Planespotter", fyne.NewMenuItem("Settings", func() {
			window.Show()
		}))
		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(m)
	}
	window.Resize(fyne.NewSize(500, 200))
	window.Hide()

	return window
}

func GenerateSettingsForm(savePath string, saveData types.SaveData) *fyne.Container {
	uiLatitude := widget.NewEntry()
	uiLatitude.SetText(fmt.Sprintf("%v", saveData.Position.Latitude))

	uiLongitude := widget.NewEntry()
	uiLongitude.SetText(fmt.Sprintf("%v", saveData.Position.Longitude))

	uiSpotDistance := widget.NewEntry()
	uiSpotDistance.SetText(strconv.Itoa(saveData.SpotDistanceKm))

	uiCheckFreq := widget.NewEntry()
	uiCheckFreq.SetText(strconv.Itoa(saveData.CheckFreqSeconds))

	uiUsername := widget.NewEntry()
	uiUsername.SetText(saveData.ApiAuth.Username)

	uiPassword := widget.NewPasswordEntry()
	uiPassword.SetText(saveData.ApiAuth.Password)

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
			save.SaveConfig(savePath, newConfig)
		},
	}

	c := container.NewVBox(settingsForm)
	return c
}
