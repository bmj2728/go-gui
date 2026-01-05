package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	caasTags = "https://cataas.com/api/tags?json=true" //will return valid tags
)

var AvailableTags = CAASTags{}

type CAASTags []string

func FetchCAASTags(timeout time.Duration) {
	bodyReader := bytes.NewReader(make([]byte, 0))
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, caasTags, bodyReader)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	// clean up when done
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&AvailableTags)
	if err != nil {
		fmt.Println(err)
	}

}
