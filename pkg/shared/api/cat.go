package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bmj2728/catfetch/pkg/shared/catdb"
	"github.com/bmj2728/catfetch/pkg/shared/metadata"
)

type CatClient struct {
	db             *catdb.CatDB
	defaultTimeout time.Duration
}

func NewCatClient(db *catdb.CatDB, defaultTimeout time.Duration) *CatClient {
	return &CatClient{
		db:             db,
		defaultTimeout: defaultTimeout,
	}
}

// TODO: Modularize RequestRandomCat, helper funcs retrieveMetadata
// TODO: RequestCatByTag(tag string), RequestCatByID(id string), RequestCatByMetadata(md *CatMetadata), RequestCatByURL(url string)

func (cc *CatClient) RequestRandomCat() (image.Image, *metadata.CatMetadata, error) {
	// make some stuff
	bodyReader := bytes.NewReader(make([]byte, 0))
	// first get the metadata in JSON format
	// the NewCatURL provides a CatURL struct using the caas base - https://cataas.com/cat
	// AsJSON adds the json=true param to the CatURL's param slice
	// Generate validates and constructs the URL, returning an error if not valid
	reqURL, err := NewCatURL().AsJSON().Generate()
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(reqURL)
	client := &http.Client{Timeout: cc.defaultTimeout}
	var meta metadata.CatMetadata

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

	err = cc.db.AddCatVersion(&meta, respBody)
	if err != nil {
		log.Printf("Error adding cat version: %v", err)
	}

	mFormat := "image/" + format

	if mFormat == meta.MIMEType {
		log.Printf("Expected format registered - %s:%s", mFormat, meta.MIMEType)
	} else {
		log.Printf("Unexpected format registered: %s:%s", mFormat, meta.MIMEType)
	}

	return img, &meta, nil
}
