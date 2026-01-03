package ui

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"time"

	"go-gui/pkg/shared/api"
)

func HandleButtonClick() (image.Image, *api.CatMetadata, error) {
	img, metadata, err := api.RequestRandomCat(30 * time.Second)
	if err != nil {
		log.Printf("Error fetching image: %v", err)
		return nil, nil, err
	}

	return img, metadata, nil
}
