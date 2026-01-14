package ui

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/bmj2728/catfetch/pkg/shared/api"
	"github.com/bmj2728/catfetch/pkg/shared/metadata"
)

func HandleButtonClick(client *api.CatClient) (image.Image, *metadata.CatMetadata, error) {
	img, md, err := client.RequestRandomCat()
	if err != nil {
		log.Printf("Error fetching image: %v", err)
		return nil, nil, err
	}

	return img, md, nil
}
