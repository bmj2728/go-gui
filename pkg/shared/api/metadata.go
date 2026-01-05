package api

import "time"

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
