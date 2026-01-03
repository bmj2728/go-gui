package api

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	caasBaseURL        = "https://cataas.com/"
	caasCatEndpoint    = "cat"
	caasCatGIFEndpoint = "cat/gif"
)

type CatMetadata struct {
	ID        string    `json:"id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
	MIMEType  string    `json:"mimetype"`
}

func (cm *CatMetadata) GetID() string {
	return cm.ID
}

func (cm *CatMetadata) GetTags() []string {
	return cm.Tags
}

func (cm *CatMetadata) GetCreatedAt() time.Time {
	return cm.CreatedAt
}

func (cm *CatMetadata) GetURL() string {
	return cm.URL
}

func (cm *CatMetadata) GetMIMEType() string {
	return cm.MIMEType
}

func RequestRandomCat(timeout time.Duration) (image.Image, *CatMetadata, error) {
	// make some stuff
	bodyReader := bytes.NewReader(make([]byte, 0))
	reqURL := caasBaseURL + caasCatEndpoint + "?json=true" //first get the metadata in json format
	client := &http.Client{Timeout: timeout}
	var meta CatMetadata

	req, err := http.NewRequest(http.MethodGet, reqURL, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	// make the req
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	// clean up when done
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {

		}
	}(resp.Body)

	//unmarshall into a metadata struct
	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Fetching image: %v", meta)

	// now get the actual image
	imgResp, err := http.Get(meta.URL)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error fetching image: %v", err)
		}
	}(imgResp.Body)

	// Read in the data
	respBody, err := io.ReadAll(imgResp.Body)
	if err != nil {
		return nil, nil, err
	}

	// decode the image
	img, format, err := image.Decode(bytes.NewReader(respBody))
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return nil, nil, err
	}

	mFormat := "image/" + format

	if mFormat == meta.MIMEType {
		log.Printf("Expected format registered - %s:%s", mFormat, meta.MIMEType)
	} else {
		log.Printf("Unexpected format registered: %s:%s", mFormat, meta.MIMEType)
	}

	return img, &meta, nil
}
