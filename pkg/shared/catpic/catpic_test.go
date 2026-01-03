package catpic

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/bmj2728/catfetch/internal/testutil"
)

// TestNewCatImage tests the constructor
func TestNewCatImage(t *testing.T) {
	t.Run("with_valid_image", func(t *testing.T) {
		img := testutil.CreateColorImage(100, 100, 255, 0, 0)
		catPic := NewCatImage(img)

		testutil.AssertNotNil(t, catPic, "CatPic should not be nil")
		testutil.AssertFalse(t, catPic.IsLoading(), "should not be loading initially")

		retrievedImg := catPic.GetImage()
		testutil.AssertNotNil(t, retrievedImg, "retrieved image should not be nil")
	})

	t.Run("with_nil_image", func(t *testing.T) {
		catPic := NewCatImage(nil)

		testutil.AssertNotNil(t, catPic, "CatPic should not be nil")
		testutil.AssertFalse(t, catPic.IsLoading(), "should not be loading initially")

		retrievedImg := catPic.GetImage()
		testutil.AssertNil(t, retrievedImg, "retrieved image should be nil")
	})
}

// TestCatPic_GetImage tests the GetImage method
func TestCatPic_GetImage(t *testing.T) {
	t.Run("returns_correct_image", func(t *testing.T) {
		img := testutil.CreateColorImage(50, 50, 100, 150, 200)
		catPic := NewCatImage(img)

		retrieved := catPic.GetImage()
		testutil.AssertNotNil(t, retrieved, "retrieved image should not be nil")

		// Verify dimensions
		bounds := retrieved.Bounds()
		testutil.AssertEqual(t, 50, bounds.Dx(), "width")
		testutil.AssertEqual(t, 50, bounds.Dy(), "height")
	})

	t.Run("returns_nil_for_nil_image", func(t *testing.T) {
		catPic := NewCatImage(nil)
		retrieved := catPic.GetImage()
		testutil.AssertNil(t, retrieved, "should return nil")
	})
}

// TestCatPic_SetImage tests the SetImage method
func TestCatPic_SetImage(t *testing.T) {
	t.Run("sets_image_correctly", func(t *testing.T) {
		catPic := NewCatImage(nil)

		img := testutil.CreateColorImage(100, 100, 255, 128, 0)
		catPic.SetImage(img)

		retrieved := catPic.GetImage()
		testutil.AssertNotNil(t, retrieved, "image should be set")
	})

	t.Run("replaces_existing_image", func(t *testing.T) {
		img1 := testutil.CreateColorImage(50, 50, 255, 0, 0)
		catPic := NewCatImage(img1)

		img2 := testutil.CreateColorImage(100, 100, 0, 255, 0)
		catPic.SetImage(img2)

		retrieved := catPic.GetImage()
		testutil.AssertNotNil(t, retrieved, "image should be set")

		bounds := retrieved.Bounds()
		testutil.AssertEqual(t, 100, bounds.Dx(), "should have new image width")
		testutil.AssertEqual(t, 100, bounds.Dy(), "should have new image height")
	})

	t.Run("can_set_to_nil", func(t *testing.T) {
		img := testutil.CreateColorImage(50, 50, 255, 0, 0)
		catPic := NewCatImage(img)

		catPic.SetImage(nil)

		retrieved := catPic.GetImage()
		testutil.AssertNil(t, retrieved, "image should be nil")
	})
}

// TestCatPic_IsLoading tests the IsLoading method
func TestCatPic_IsLoading(t *testing.T) {
	t.Run("initially_false", func(t *testing.T) {
		catPic := NewCatImage(nil)
		testutil.AssertFalse(t, catPic.IsLoading(), "should not be loading initially")
	})

	t.Run("true_after_set_loading", func(t *testing.T) {
		catPic := NewCatImage(nil)
		catPic.SetLoading()
		testutil.AssertTrue(t, catPic.IsLoading(), "should be loading after SetLoading")
	})

	t.Run("false_after_clear_loading", func(t *testing.T) {
		catPic := NewCatImage(nil)
		catPic.SetLoading()
		catPic.ClearLoading()
		testutil.AssertFalse(t, catPic.IsLoading(), "should not be loading after ClearLoading")
	})
}

// TestCatPic_SetLoading tests the SetLoading method
func TestCatPic_SetLoading(t *testing.T) {
	catPic := NewCatImage(nil)

	testutil.AssertFalse(t, catPic.IsLoading(), "initially false")

	catPic.SetLoading()
	testutil.AssertTrue(t, catPic.IsLoading(), "should be true after SetLoading")

	// Multiple calls should work
	catPic.SetLoading()
	testutil.AssertTrue(t, catPic.IsLoading(), "should still be true")
}

// TestCatPic_ClearLoading tests the ClearLoading method
func TestCatPic_ClearLoading(t *testing.T) {
	catPic := NewCatImage(nil)

	catPic.SetLoading()
	testutil.AssertTrue(t, catPic.IsLoading(), "should be loading")

	catPic.ClearLoading()
	testutil.AssertFalse(t, catPic.IsLoading(), "should not be loading after clear")

	// Multiple calls should work
	catPic.ClearLoading()
	testutil.AssertFalse(t, catPic.IsLoading(), "should still be false")
}

// TestCatPic_LoadingStateTransitions tests loading state transitions
func TestCatPic_LoadingStateTransitions(t *testing.T) {
	catPic := NewCatImage(nil)

	// Test multiple transitions
	testutil.AssertFalse(t, catPic.IsLoading(), "initial state")

	catPic.SetLoading()
	testutil.AssertTrue(t, catPic.IsLoading(), "after first set")

	catPic.ClearLoading()
	testutil.AssertFalse(t, catPic.IsLoading(), "after first clear")

	catPic.SetLoading()
	testutil.AssertTrue(t, catPic.IsLoading(), "after second set")

	catPic.ClearLoading()
	testutil.AssertFalse(t, catPic.IsLoading(), "after second clear")
}

// TestCatPic_Draw_NilImage tests Draw with nil image
func TestCatPic_Draw_NilImage(t *testing.T) {
	catPic := NewCatImage(nil)

	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Constraints: layout.Constraints{
			Min: image.Pt(0, 0),
			Max: image.Pt(400, 500),
		},
	}

	dims := catPic.Draw(gtx)

	// Should return minimum dimensions
	testutil.AssertEqual(t, 0, dims.Size.X, "nil image width")
	testutil.AssertEqual(t, 0, dims.Size.Y, "nil image height")
}

// TestCatPic_Draw_ImageFitsConstraints tests Draw when image fits within constraints
func TestCatPic_Draw_ImageFitsConstraints(t *testing.T) {
	// Image smaller than constraints
	img := testutil.CreateColorImage(100, 100, 255, 0, 0)
	catPic := NewCatImage(img)

	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Constraints: layout.Constraints{
			Min: image.Pt(0, 0),
			Max: image.Pt(400, 400),
		},
	}

	dims := catPic.Draw(gtx)

	// Image should be scaled up to fit, or stay same size
	// Based on the code, it scales to fit the constraints
	testutil.AssertTrue(t, dims.Size.X > 0, "width should be positive")
	testutil.AssertTrue(t, dims.Size.Y > 0, "height should be positive")
	testutil.AssertTrue(t, dims.Size.X <= 400, "width should not exceed max")
	testutil.AssertTrue(t, dims.Size.Y <= 400, "height should not exceed max")
}

// TestCatPic_Draw_ImageScaling tests image scaling for various scenarios
func TestCatPic_Draw_ImageScaling(t *testing.T) {
	tests := []struct {
		name        string
		imageWidth  int
		imageHeight int
		maxWidth    int
		maxHeight   int
		checkAspect bool
	}{
		{
			name:        "landscape_image_wider_than_constraints",
			imageWidth:  1920,
			imageHeight: 1080,
			maxWidth:    400,
			maxHeight:   400,
			checkAspect: true,
		},
		{
			name:        "portrait_image_taller_than_constraints",
			imageWidth:  800,
			imageHeight: 1200,
			maxWidth:    400,
			maxHeight:   400,
			checkAspect: true,
		},
		{
			name:        "square_image_with_square_constraints",
			imageWidth:  500,
			imageHeight: 500,
			maxWidth:    400,
			maxHeight:   400,
			checkAspect: true,
		},
		{
			name:        "very_wide_image",
			imageWidth:  1000,
			imageHeight: 100,
			maxWidth:    400,
			maxHeight:   400,
			checkAspect: true,
		},
		{
			name:        "very_tall_image",
			imageWidth:  100,
			imageHeight: 1000,
			maxWidth:    400,
			maxHeight:   400,
			checkAspect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := testutil.CreateColorImage(tt.imageWidth, tt.imageHeight, 128, 128, 128)
			catPic := NewCatImage(img)

			var ops op.Ops
			gtx := layout.Context{
				Ops: &ops,
				Constraints: layout.Constraints{
					Min: image.Pt(0, 0),
					Max: image.Pt(tt.maxWidth, tt.maxHeight),
				},
			}

			dims := catPic.Draw(gtx)

			// Verify dimensions are within constraints
			testutil.AssertTrue(t, dims.Size.X <= tt.maxWidth, "width within max")
			testutil.AssertTrue(t, dims.Size.Y <= tt.maxHeight, "height within max")
			testutil.AssertTrue(t, dims.Size.X > 0, "width positive")
			testutil.AssertTrue(t, dims.Size.Y > 0, "height positive")

			// Verify aspect ratio is preserved (within tolerance)
			if tt.checkAspect {
				originalAspect := float64(tt.imageWidth) / float64(tt.imageHeight)
				scaledAspect := float64(dims.Size.X) / float64(dims.Size.Y)
				tolerance := 0.05 // 5% tolerance for integer rounding

				diff := originalAspect - scaledAspect
				if diff < 0 {
					diff = -diff
				}
				if diff > tolerance {
					t.Errorf("Aspect ratio not preserved: original=%.4f, scaled=%.4f", originalAspect, scaledAspect)
				}
			}
		})
	}
}

// TestCatPic_Draw_AspectRatioPreservation tests aspect ratio preservation
func TestCatPic_Draw_AspectRatioPreservation(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		height      int
		aspectRatio float64
	}{
		{"16:9", 1920, 1080, 16.0 / 9.0},
		{"4:3", 800, 600, 4.0 / 3.0},
		{"1:1", 500, 500, 1.0},
		{"21:9", 2560, 1080, 21.0 / 9.0},
		{"9:16", 1080, 1920, 9.0 / 16.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := testutil.CreateColorImage(tt.width, tt.height, 200, 100, 50)
			catPic := NewCatImage(img)

			var ops op.Ops
			gtx := layout.Context{
				Ops: &ops,
				Constraints: layout.Constraints{
					Min: image.Pt(0, 0),
					Max: image.Pt(400, 500),
				},
			}

			dims := catPic.Draw(gtx)

			// Calculate scaled aspect ratio
			scaledAspect := float64(dims.Size.X) / float64(dims.Size.Y)
			tolerance := 0.05 // 5% tolerance for integer rounding

			diff := tt.aspectRatio - scaledAspect
			if diff < 0 {
				diff = -diff
			}

			if diff > tolerance {
				t.Errorf("Aspect ratio not preserved: expected %.4f, got %.4f", tt.aspectRatio, scaledAspect)
			}
		})
	}
}

// TestCatPic_Draw_ExtremeAspectRatios tests very wide and very tall images
func TestCatPic_Draw_ExtremeAspectRatios(t *testing.T) {
	t.Run("very_wide_10_to_1", func(t *testing.T) {
		img := testutil.CreateColorImage(1000, 100, 255, 0, 0)
		catPic := NewCatImage(img)

		var ops op.Ops
		gtx := layout.Context{
			Ops: &ops,
			Constraints: layout.Constraints{
				Min: image.Pt(0, 0),
				Max: image.Pt(400, 400),
			},
		}

		dims := catPic.Draw(gtx)

		testutil.AssertTrue(t, dims.Size.X <= 400, "width within constraints")
		testutil.AssertTrue(t, dims.Size.Y <= 400, "height within constraints")

		// Should be constrained by width
		testutil.AssertEqual(t, 400, dims.Size.X, "should use full width")
	})

	t.Run("very_tall_1_to_10", func(t *testing.T) {
		img := testutil.CreateColorImage(100, 1000, 0, 255, 0)
		catPic := NewCatImage(img)

		var ops op.Ops
		gtx := layout.Context{
			Ops: &ops,
			Constraints: layout.Constraints{
				Min: image.Pt(0, 0),
				Max: image.Pt(400, 400),
			},
		}

		dims := catPic.Draw(gtx)

		testutil.AssertTrue(t, dims.Size.X <= 400, "width within constraints")
		testutil.AssertTrue(t, dims.Size.Y <= 400, "height within constraints")

		// Should be constrained by height
		testutil.AssertEqual(t, 400, dims.Size.Y, "should use full height")
	})
}

// TestCatPic_Draw_ZeroDimensions tests edge cases with zero dimensions
func TestCatPic_Draw_ZeroDimensions(t *testing.T) {
	t.Run("zero_width_constraint", func(t *testing.T) {
		img := testutil.CreateColorImage(100, 100, 255, 255, 255)
		catPic := NewCatImage(img)

		var ops op.Ops
		gtx := layout.Context{
			Ops: &ops,
			Constraints: layout.Constraints{
				Min: image.Pt(0, 0),
				Max: image.Pt(0, 400),
			},
		}

		// This should not panic
		testutil.AssertNoPanic(t, func() {
			dims := catPic.Draw(gtx)
			_ = dims
		}, "should not panic with zero width constraint")
	})

	t.Run("zero_height_constraint", func(t *testing.T) {
		img := testutil.CreateColorImage(100, 100, 255, 255, 255)
		catPic := NewCatImage(img)

		var ops op.Ops
		gtx := layout.Context{
			Ops: &ops,
			Constraints: layout.Constraints{
				Min: image.Pt(0, 0),
				Max: image.Pt(400, 0),
			},
		}

		// This should not panic
		testutil.AssertNoPanic(t, func() {
			dims := catPic.Draw(gtx)
			_ = dims
		}, "should not panic with zero height constraint")
	})
}

// TestCatPic_Draw_LargeConstraints tests with very large constraints
func TestCatPic_Draw_LargeConstraints(t *testing.T) {
	img := testutil.CreateColorImage(100, 100, 128, 128, 128)
	catPic := NewCatImage(img)

	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Constraints: layout.Constraints{
			Min: image.Pt(0, 0),
			Max: image.Pt(10000, 10000),
		},
	}

	dims := catPic.Draw(gtx)

	// Image should be scaled up significantly
	testutil.AssertTrue(t, dims.Size.X > 100, "width should be scaled up")
	testutil.AssertTrue(t, dims.Size.Y > 100, "height should be scaled up")
	testutil.AssertTrue(t, dims.Size.X <= 10000, "width within constraints")
	testutil.AssertTrue(t, dims.Size.Y <= 10000, "height within constraints")

	// Aspect ratio should still be preserved
	scaledAspect := float64(dims.Size.X) / float64(dims.Size.Y)
	expectedAspect := 1.0 // 100x100 is square
	tolerance := 0.05     // 5% tolerance for integer rounding

	diff := expectedAspect - scaledAspect
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("Aspect ratio not preserved with large constraints: expected %.4f, got %.4f", expectedAspect, scaledAspect)
	}
}
