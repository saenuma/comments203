package main

import (
	"log"

	"github.com/sqweek/dialog"
)

func PickImageFile() string {
	filename, err := dialog.File().Filter("PNG Image", "png").Filter("JPEG Image", "jpg").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}
