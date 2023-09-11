package ui

import (
	"fmt"
	"planespotter/helpers/save"
	"planespotter/helpers/types"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func InitUi(savePath string, saveData types.SaveData) fyne.App {
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
			saveData.Position.Latitude, _ = strconv.ParseFloat(uiLatitude.Text, 64)
			saveData.Position.Longitude, _ = strconv.ParseFloat(uiLongitude.Text, 64)
			saveData.ApiAuth.Username = uiUsername.Text
			saveData.ApiAuth.Password = uiPassword.Text
			saveData.SpotDistanceKm, _ = strconv.Atoi(uiSpotDistance.Text)
			saveData.CheckFreqSeconds, _ = strconv.Atoi(uiCheckFreq.Text)
			save.SaveConfig(savePath, newConfig)
		},
	}

	uiWindow.SetContent(container.NewVBox(title, settingsForm))
	return ui
}
