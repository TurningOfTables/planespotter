package main

import (
	"fmt"
	"planespotter/helpers/types"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// InitUi takes the url, savePath and Save Data to form the main sections of the Fyne UI.
// Returns the Fyne App and Fyne Window
func InitUi(url, savePath string, saveData types.SaveData) (fyne.App, fyne.Window) {
	icon, _ := fyne.LoadResourceFromPath("assets/plane.png")
	app := app.NewWithID("Planespotter")
	app.SetIcon(icon)
	window := WindowSetup(app, icon)

	title := widget.NewLabelWithStyle("Configuration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	settingsForm := FormSetup(savePath, saveData)

	startButton := widget.NewButton("Start", func() {
		if !started {
			startUpdateLoop(url, saveData)
		}
	})
	stopButton := widget.NewButton("Stop", func() {
		if started {
			stopUpdateLoop()
		}
	})

	if started {
		status.Set(StartedText)
	} else {
		status.Set(StoppedText)
	}
	statusLabel := widget.NewLabelWithData(status)

	window.SetContent(container.NewVBox(title, settingsForm, startButton, stopButton, statusLabel))
	return app, window
}

// WindowSetup takes the Fyne App and icon (Fyne Resource) and creates the window with correct size and title
// Returns the Fyne Window
func WindowSetup(app fyne.App, icon fyne.Resource) fyne.Window {
	window := app.NewWindow("Planespotter")
	window.Resize(fyne.NewSize(500, 200))
	window.Hide()

	return window
}

// FormSetup takes the savePath and a SaveData, and creates the settings form for the user
// to update their config with save button. Returns a Fyne Container for inclusion in a Fyne Window.
func FormSetup(savePath string, saveData types.SaveData) *fyne.Container {
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
			{Text: "Latitude", HintText: "Decimal degrees", Widget: uiLatitude},
			{Text: "Longitude", HintText: "Decimal degrees", Widget: uiLongitude},
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
			SaveConfig(savePath, newConfig)
			if started {
				stopUpdateLoop()
			}
			url, saveData := InitSaveData(savePath)
			startUpdateLoop(url, saveData)

		},
	}

	c := container.NewVBox(settingsForm)
	return c
}
