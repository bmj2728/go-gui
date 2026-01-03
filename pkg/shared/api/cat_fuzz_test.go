package api

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/bmj2728/catfetch/internal/testutil"
)

// FuzzMetadataJSON fuzzes JSON unmarshaling for CatMetadata
func FuzzMetadataJSON(f *testing.F) {
	// Add seed corpus with various JSON inputs
	f.Add(`{"id":"test","tags":[],"created_at":"2025-01-01T00:00:00Z","url":"http://example.com","mimetype":"image/jpeg"}`)
	f.Add(`{}`)
	f.Add(`{"id":""}`)
	f.Add(`{"tags":"not_an_array"}`)
	f.Add(`{"id":123}`)
	f.Add(`{"created_at":"invalid"}`)
	f.Add(testutil.ValidMetadataJSON())
	f.Add(testutil.MalformedMetadataJSON())
	f.Add(testutil.EmptyMetadataJSON())
	f.Add(testutil.MissingFieldsMetadataJSON())
	f.Add(`{"id":"test","tags":["a","b","c","d","e","f","g","h","i","j"]}`)
	f.Add(`{"url":"https://very-long-url.example.com/with/many/path/segments/and/query/params?key1=value1&key2=value2"}`)

	f.Fuzz(func(t *testing.T, jsonData string) {
		// Fuzzing should never panic
		var meta CatMetadata

		// Try to unmarshal - this should handle all malformed input gracefully
		err := json.Unmarshal([]byte(jsonData), &meta)

		if err != nil {
			// Malformed JSON is expected and acceptable
			return
		}

		// If unmarshaling succeeds, all getters should work without panic
		_ = meta.GetID()
		_ = meta.GetTags()
		_ = meta.GetCreatedAt()
		_ = meta.GetURL()
		_ = meta.GetMIMEType()

		// Test that we can marshal it back
		_, marshalErr := json.Marshal(meta)
		if marshalErr != nil {
			// Some edge cases might not marshal back, that's okay
			return
		}
	})
}

// FuzzImageData fuzzes image decoding
func FuzzImageData(f *testing.F) {
	// Add seed corpus with various image formats
	f.Add(testutil.ValidJPEGBytes())
	f.Add(testutil.ValidPNGBytes())
	f.Add(testutil.ValidGIFBytes())
	f.Add(testutil.CorruptedImageBytes())
	f.Add(testutil.PartialImageBytes())
	f.Add([]byte{}) // Empty data
	f.Add([]byte{0x00, 0x00, 0x00, 0x00})
	f.Add([]byte{0xFF, 0xFF, 0xFF, 0xFF})

	// Add some known image format headers
	f.Add([]byte{0xFF, 0xD8, 0xFF})                               // JPEG header start
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) // PNG header
	f.Add([]byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61})             // GIF89a header

	f.Fuzz(func(t *testing.T, imageData []byte) {
		// Image decoding should never panic, even with malformed data
		_, _, err := image.Decode(bytes.NewReader(imageData))

		if err != nil {
			// Errors are expected for invalid image data
			return
		}

		// If decoding succeeded, that's fine too
		// The image package should handle all edge cases
	})
}

// FuzzImageDataWithFormat fuzzes image decoding and format detection
func FuzzImageDataWithFormat(f *testing.F) {
	// Seed with valid images
	f.Add(testutil.ValidJPEGBytes(), "jpeg")
	f.Add(testutil.ValidPNGBytes(), "png")
	f.Add(testutil.ValidGIFBytes(), "gif")
	f.Add(testutil.CorruptedImageBytes(), "unknown")

	f.Fuzz(func(t *testing.T, imageData []byte, expectedFormat string) {
		// Should never panic
		img, format, err := image.Decode(bytes.NewReader(imageData))

		if err != nil {
			// Decoding errors are acceptable
			return
		}

		// If decoding succeeded, verify format is reasonable
		if format != "jpeg" && format != "png" && format != "gif" {
			// Unknown format detected - this is possible
		}

		// If we got an image, verify it has reasonable bounds
		if img != nil {
			bounds := img.Bounds()
			width := bounds.Dx()
			height := bounds.Dy()

			// Reasonable image dimensions check
			if width < 0 || height < 0 {
				t.Errorf("Invalid image dimensions: %dx%d", width, height)
			}

			// Very large dimensions might indicate an issue
			if width > 100000 || height > 100000 {
				// This is suspicious but not necessarily wrong
				// Just make sure it doesn't panic
				_ = width
				_ = height
			}
		}
	})
}

// FuzzCatMetadataRoundTrip fuzzes JSON marshal/unmarshal round trip
func FuzzCatMetadataRoundTrip(f *testing.F) {
	// Seed with structured data
	f.Add("test_id", `["tag1","tag2"]`, "2025-01-01T00:00:00Z", "http://example.com", "image/jpeg")
	f.Add("", `[]`, "", "", "")
	f.Add("very_long_id_with_many_characters_that_might_cause_issues", `["a","b","c"]`, "2025-12-31T23:59:59Z", "https://example.com", "image/png")

	f.Fuzz(func(t *testing.T, id, tagsJSON, createdAt, url, mimetype string) {
		// Create metadata struct
		meta := CatMetadata{
			ID:       id,
			URL:      url,
			MIMEType: mimetype,
		}

		// Try to parse tags
		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err == nil {
			meta.Tags = tags
		}

		// Try to parse time
		// We'll skip time parsing in fuzz test as it's complex
		// and we're mainly testing JSON handling

		// Marshal the metadata
		data, err := json.Marshal(meta)
		if err != nil {
			// Marshal should rarely fail for our struct
			return
		}

		// Unmarshal it back
		var meta2 CatMetadata
		err = json.Unmarshal(data, &meta2)
		if err != nil {
			t.Errorf("Failed to unmarshal after marshal: %v", err)
			return
		}

		// Basic consistency checks
		if meta.ID != meta2.ID {
			t.Errorf("ID mismatch after round trip: %s != %s", meta.ID, meta2.ID)
		}

		if meta.URL != meta2.URL {
			t.Errorf("URL mismatch after round trip: %s != %s", meta.URL, meta2.URL)
		}

		if meta.MIMEType != meta2.MIMEType {
			t.Errorf("MIMEType mismatch after round trip: %s != %s", meta.MIMEType, meta2.MIMEType)
		}
	})
}

// FuzzImageBoundsCalculation fuzzes image bounds calculations
func FuzzImageBoundsCalculation(f *testing.F) {
	// Seed with various image dimensions
	f.Add(testutil.ValidJPEGBytes())
	f.Add(testutil.ValidPNGBytes())
	f.Add(testutil.ValidGIFBytes())

	f.Fuzz(func(t *testing.T, imageData []byte) {
		img, _, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			return
		}

		if img == nil {
			return
		}

		// Perform bounds calculations that might be done in real code
		bounds := img.Bounds()

		// These operations should never panic
		_ = bounds.Dx()
		_ = bounds.Dy()
		_ = bounds.Min.X
		_ = bounds.Min.Y
		_ = bounds.Max.X
		_ = bounds.Max.Y

		// Calculate aspect ratio (division by zero protection)
		width := float64(bounds.Dx())
		height := float64(bounds.Dy())

		if height != 0 {
			aspectRatio := width / height
			_ = aspectRatio
		}

		if width != 0 {
			inverseRatio := height / width
			_ = inverseRatio
		}
	})
}
