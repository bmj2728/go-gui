package metadata

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

func (cm *CatMetadata) ToCatDBMetadata() *CatDBMetadata {
	return &CatDBMetadata{
		ID:        cm.ID,
		Tags:      cm.Tags,
		CreatedAt: cm.CreatedAt,
		MIMEType:  cm.MIMEType,
	}
}

type CatDBMetadata struct {
	ID        string    `json:"id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	MIMEType  string    `json:"mimetype"`
}

func (cdm *CatDBMetadata) GetID() string {
	return cdm.ID
}

func (cdm *CatDBMetadata) GetTags() []string {
	return cdm.Tags
}

func (cdm *CatDBMetadata) GetCreatedAt() time.Time {
	return cdm.CreatedAt
}

func (cdm *CatDBMetadata) GetMIMEType() string {
	return cdm.MIMEType
}

func (cdm *CatDBMetadata) ToCatMetadata(url string) *CatMetadata {
	return &CatMetadata{
		ID:        cdm.ID,
		Tags:      cdm.Tags,
		CreatedAt: cdm.CreatedAt,
		URL:       url,
		MIMEType:  cdm.MIMEType,
	}
}
