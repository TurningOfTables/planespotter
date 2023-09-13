package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitUi(t *testing.T) {
	url := "https://www.example.com"
	err := CreateSaveIfNotExists(testSavePath)
	if err != nil {
		t.Error(err)
	}

	saveData, err := GetSave(testSavePath)
	if err != nil {
		t.Error(err)
	}
	_, window := InitUi(url, testSavePath, saveData)

	assert.NotNil(t, window.Content()) // weak assertion

	err = os.Remove(testSavePath)
	if err != nil {
		t.Error(err)
	}
}
