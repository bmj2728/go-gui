package testutil

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"
)

// CreateTestImage generates a test image with specified dimensions and format
func CreateTestImage(width, height int, format string) (image.Image, error) {
	img := CreateColorImage(width, height, 128, 128, 128)

	// Encode and decode to ensure proper format
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return nil, err
		}
		return jpeg.Decode(&buf)
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
		return png.Decode(&buf)
	case "gif":
		if err := gif.Encode(&buf, img, nil); err != nil {
			return nil, err
		}
		return gif.Decode(&buf)
	default:
		return img, nil
	}
}

// CreateColorImage generates a solid color image
func CreateColorImage(width, height int, r, g, b uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := color.RGBA{R: r, G: g, B: b, A: 255}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	return img
}

// CreateGradientImage creates an image with a horizontal gradient
func CreateGradientImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create gradient from black to white
			intensity := uint8(float64(x) / float64(width) * 255)
			c := color.RGBA{R: intensity, G: intensity, B: intensity, A: 255}
			img.Set(x, y, c)
		}
	}

	return img
}

// ImagesEqual compares two images for equality
func ImagesEqual(img1, img2 image.Image) bool {
	if img1 == nil || img2 == nil {
		return img1 == img2
	}

	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1 != bounds2 {
		return false
	}

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return false
			}
		}
	}

	return true
}

// GetAspectRatio calculates the aspect ratio of an image
func GetAspectRatio(img image.Image) float64 {
	if img == nil {
		return 0.0
	}

	bounds := img.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	if height == 0 {
		return 0.0
	}

	return width / height
}

// AssertAspectRatio verifies the aspect ratio is within tolerance
func AssertAspectRatio(t *testing.T, img image.Image, expected float64) {
	t.Helper()

	actual := GetAspectRatio(img)
	tolerance := 0.01 // 1% tolerance for floating point comparison

	if abs(actual-expected) > tolerance {
		t.Errorf("Aspect ratio mismatch: expected %.4f, got %.4f", expected, actual)
	}
}

// AssertImageDimensions verifies image has expected dimensions
func AssertImageDimensions(t *testing.T, img image.Image, expectedWidth, expectedHeight int) {
	t.Helper()

	if img == nil {
		t.Fatal("Image is nil")
	}

	bounds := img.Bounds()
	actualWidth := bounds.Dx()
	actualHeight := bounds.Dy()

	if actualWidth != expectedWidth || actualHeight != expectedHeight {
		t.Errorf("Image dimensions mismatch: expected %dx%d, got %dx%d",
			expectedWidth, expectedHeight, actualWidth, actualHeight)
	}
}

// GetImageDimensions returns the width and height of an image
func GetImageDimensions(img image.Image) (width, height int) {
	if img == nil {
		return 0, 0
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// EncodeImage encodes an image to bytes in the specified format
func EncodeImage(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
		return buf.Bytes(), err
	case "png":
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	case "gif":
		err := gif.Encode(&buf, img, nil)
		return buf.Bytes(), err
	default:
		return nil, nil
	}
}

// CreateTestImageBytes creates image bytes for testing
func CreateTestImageBytes(width, height int, format string) ([]byte, error) {
	img := CreateColorImage(width, height, 100, 150, 200)
	return EncodeImage(img, format)
}
