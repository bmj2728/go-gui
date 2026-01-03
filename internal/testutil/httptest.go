package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// ServerConfig configures the behavior of the mock HTTP server
type ServerConfig struct {
	MetadataResponse   string        // JSON to return for metadata
	ImageData          []byte        // Image bytes to return
	MetadataDelay      time.Duration // Artificial delay for metadata
	ImageDelay         time.Duration // Artificial delay for image
	MetadataStatus     int           // HTTP status for metadata
	ImageStatus        int           // HTTP status for image
	ShouldFailMetadata bool          // Simulate metadata fetch failure
	ShouldFailImage    bool          // Simulate image fetch failure
}

// NewMockCatAPIServer creates a test HTTP server that simulates the cat API
// Returns metadata server URL and image server URL
func NewMockCatAPIServer(config ServerConfig) (*httptest.Server, *httptest.Server) {
	// Create image server first
	imageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.ImageDelay > 0 {
			time.Sleep(config.ImageDelay)
		}

		if config.ShouldFailImage {
			http.Error(w, "Image server error", http.StatusInternalServerError)
			return
		}

		status := config.ImageStatus
		if status == 0 {
			status = http.StatusOK
		}

		w.WriteHeader(status)
		if status == http.StatusOK && config.ImageData != nil {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(config.ImageData)
		}
	}))

	// Create metadata server that references image server
	metadataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.MetadataDelay > 0 {
			time.Sleep(config.MetadataDelay)
		}

		if config.ShouldFailMetadata {
			http.Error(w, "Metadata server error", http.StatusInternalServerError)
			return
		}

		status := config.MetadataStatus
		if status == 0 {
			status = http.StatusOK
		}

		w.WriteHeader(status)
		if status == http.StatusOK {
			// Inject image server URL into metadata
			metadata := config.MetadataResponse
			if metadata == "" {
				metadata = ValidMetadataJSON()
			}
			// Replace placeholder with actual image server URL
			metadata = fmt.Sprintf(metadata, imageServer.URL)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(metadata))
		}
	}))

	return metadataServer, imageServer
}

// ValidMetadataJSON returns a valid CatMetadata JSON string
// Use %s placeholder for URL that will be replaced with image server URL
func ValidMetadataJSON() string {
	return `{
		"id": "test_cat_001",
		"tags": ["cute", "fluffy"],
		"created_at": "2025-01-01T12:00:00Z",
		"url": "%s",
		"mimetype": "image/jpeg"
	}`
}

// ValidMetadataJSONWithURL returns valid metadata with a specific URL
func ValidMetadataJSONWithURL(url string) string {
	metadata := map[string]interface{}{
		"id":         "test_cat_001",
		"tags":       []string{"cute", "fluffy"},
		"created_at": "2025-01-01T12:00:00Z",
		"url":        url,
		"mimetype":   "image/jpeg",
	}
	bytes, _ := json.Marshal(metadata)
	return string(bytes)
}

// MalformedMetadataJSON returns JSON with syntax errors
func MalformedMetadataJSON() string {
	return `{
		"id": "test_cat_002",
		"tags": ["cute", "fluffy"
		"created_at": "2025-01-01T12:00:00Z"
	`
}

// MissingFieldsMetadataJSON returns JSON with missing required fields
func MissingFieldsMetadataJSON() string {
	return `{
		"id": "test_cat_003"
	}`
}

// WrongTypesMetadataJSON returns JSON with incorrect field types
func WrongTypesMetadataJSON() string {
	return `{
		"id": 12345,
		"tags": "not_an_array",
		"created_at": "invalid_date",
		"url": 123,
		"mimetype": true
	}`
}

// EmptyMetadataJSON returns empty JSON object
func EmptyMetadataJSON() string {
	return `{}`
}

// ValidJPEGBytes returns minimal valid JPEG image bytes
func ValidJPEGBytes() []byte {
	// Minimal 1x1 pixel JPEG image
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x01, 0x00, 0x48, 0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
		0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x00, 0x0B,
		0x08, 0x00, 0x01, 0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xFF, 0xC4, 0x00,
		0x14, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xDA, 0x00, 0x08, 0x01,
		0x01, 0x00, 0x00, 0x3F, 0x00, 0x7F, 0xFF, 0xD9,
	}
}

// ValidPNGBytes returns minimal valid PNG image bytes
func ValidPNGBytes() []byte {
	// Minimal 1x1 pixel PNG image
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}

// ValidGIFBytes returns minimal valid GIF image bytes
func ValidGIFBytes() []byte {
	// Minimal 1x1 pixel GIF image
	return []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00,
		0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x21, 0xF9, 0x04, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
		0x00, 0x02, 0x02, 0x44, 0x01, 0x00, 0x3B,
	}
}

// CorruptedImageBytes returns invalid image data
func CorruptedImageBytes() []byte {
	// Starts like JPEG but is corrupted
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0xDE, 0xAD, 0xBE, 0xEF, 0xCA, 0xFE, 0xBA, 0xBE,
	}
}

// PartialImageBytes returns truncated JPEG data
func PartialImageBytes() []byte {
	return []byte{0xFF, 0xD8, 0xFF, 0xE0}
}
