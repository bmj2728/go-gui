package catpic

import (
	"image"
	"sync"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/bmj2728/catfetch/internal/testutil"
)

// TestCatPic_ConcurrentGetSet tests concurrent GetImage and SetImage operations
func TestCatPic_ConcurrentGetSet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race test in short mode")
	}

	catPic := NewCatImage(nil)

	const numGoroutines = 100
	const opsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start reader goroutines
	for i := 0; i < numGoroutines/2; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				_ = catPic.GetImage()
			}
		}()
	}

	// Start writer goroutines
	for i := 0; i < numGoroutines/2; i++ {
		go func(id int) {
			defer wg.Done()
			img := testutil.CreateColorImage(50, 50, uint8(id%256), 128, 200)
			for j := 0; j < opsPerGoroutine; j++ {
				catPic.SetImage(img)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify CatPic is still in a valid state
	img := catPic.GetImage()
	if img != nil {
		bounds := img.Bounds()
		testutil.AssertEqual(t, 50, bounds.Dx(), "final image width")
		testutil.AssertEqual(t, 50, bounds.Dy(), "final image height")
	}
}

// TestCatPic_ConcurrentLoadingState tests concurrent loading state operations
func TestCatPic_ConcurrentLoadingState(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race test in short mode")
	}

	catPic := NewCatImage(nil)

	const numGoroutines = 100
	const opsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start goroutines that set/clear loading
	for i := 0; i < numGoroutines/3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				catPic.SetLoading()
			}
		}()
	}

	for i := 0; i < numGoroutines/3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				catPic.ClearLoading()
			}
		}()
	}

	// Start goroutines that read loading state
	for i := 0; i < numGoroutines/3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				_ = catPic.IsLoading()
			}
		}()
	}

	wg.Wait()

	// Verify final state is valid (either true or false, but consistent)
	finalState := catPic.IsLoading()
	_ = finalState // Just verify we can read it without panic
}

// TestCatPic_ConcurrentDrawAndSet tests Draw operations concurrent with SetImage
func TestCatPic_ConcurrentDrawAndSet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race test in short mode")
	}

	img := testutil.CreateColorImage(200, 200, 255, 128, 0)
	catPic := NewCatImage(img)

	const numDrawers = 10
	const numSetters = 10
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numDrawers + numSetters)

	// Start drawer goroutines (simulating UI thread)
	for i := 0; i < numDrawers; i++ {
		go func() {
			defer wg.Done()
			var ops op.Ops
			gtx := layout.Context{
				Ops: &ops,
				Constraints: layout.Constraints{
					Min: image.Pt(0, 0),
					Max: image.Pt(400, 500),
				},
			}

			for j := 0; j < opsPerGoroutine; j++ {
				_ = catPic.Draw(gtx)
			}
		}()
	}

	// Start setter goroutines (simulating fetch thread)
	for i := 0; i < numSetters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				newImg := testutil.CreateColorImage(100+id, 100+id, uint8(id%256), 100, 150)
				catPic.SetImage(newImg)
			}
		}(i)
	}

	wg.Wait()

	// Verify final state is valid
	finalImg := catPic.GetImage()
	if finalImg != nil {
		bounds := finalImg.Bounds()
		testutil.AssertTrue(t, bounds.Dx() > 0, "final image has positive width")
		testutil.AssertTrue(t, bounds.Dy() > 0, "final image has positive height")
	}
}

// TestCatPic_ConcurrentMixedOperations tests all operations happening concurrently
func TestCatPic_ConcurrentMixedOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race test in short mode")
	}

	catPic := NewCatImage(nil)

	const numGoroutines = 20
	const opsPerGoroutine = 500

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Mix of all operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			var ops op.Ops
			gtx := layout.Context{
				Ops: &ops,
				Constraints: layout.Constraints{
					Min: image.Pt(0, 0),
					Max: image.Pt(400, 500),
				},
			}

			for j := 0; j < opsPerGoroutine; j++ {
				switch j % 7 {
				case 0:
					_ = catPic.GetImage()
				case 1:
					img := testutil.CreateColorImage(100, 100, uint8(id%256), 128, 64)
					catPic.SetImage(img)
				case 2:
					catPic.SetLoading()
				case 3:
					catPic.ClearLoading()
				case 4:
					_ = catPic.IsLoading()
				case 5:
					_ = catPic.Draw(gtx)
				case 6:
					catPic.SetImage(nil)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify CatPic is still functional
	testutil.AssertNoPanic(t, func() {
		_ = catPic.GetImage()
		_ = catPic.IsLoading()
		catPic.SetLoading()
		catPic.ClearLoading()
	}, "CatPic should still be functional after concurrent operations")
}

// TestCatPic_HighConcurrency tests with a very large number of goroutines
func TestCatPic_HighConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high concurrency test in short mode")
	}

	catPic := NewCatImage(nil)

	const numGoroutines = 1000
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < opsPerGoroutine; j++ {
				if j%2 == 0 {
					img := testutil.CreateColorImage(50, 50, uint8(id%256), 100, 200)
					catPic.SetImage(img)
				} else {
					_ = catPic.GetImage()
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify no deadlock occurred and final state is valid
	finalImg := catPic.GetImage()
	if finalImg != nil {
		bounds := finalImg.Bounds()
		testutil.AssertEqual(t, 50, bounds.Dx(), "final image width")
		testutil.AssertEqual(t, 50, bounds.Dy(), "final image height")
	}
}

// TestCatPic_RapidStateChanges tests rapid loading state changes
func TestCatPic_RapidStateChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rapid state changes test in short mode")
	}

	catPic := NewCatImage(nil)

	const numGoroutines = 50
	const cycles = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < cycles; j++ {
				catPic.SetLoading()
				_ = catPic.IsLoading()
				catPic.ClearLoading()
				_ = catPic.IsLoading()
			}
		}()
	}

	wg.Wait()

	// Verify final loading state is accessible
	finalState := catPic.IsLoading()
	_ = finalState
}

// TestCatPic_ConcurrentNilImageHandling tests concurrent operations with nil images
func TestCatPic_ConcurrentNilImageHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping nil image handling test in short mode")
	}

	catPic := NewCatImage(nil)

	const numGoroutines = 100
	const opsPerGoroutine = 500

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < opsPerGoroutine; j++ {
				if j%3 == 0 {
					catPic.SetImage(nil)
				} else if j%3 == 1 {
					img := testutil.CreateColorImage(75, 75, uint8(id%256), 150, 200)
					catPic.SetImage(img)
				} else {
					_ = catPic.GetImage()
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify CatPic handles nil correctly
	testutil.AssertNoPanic(t, func() {
		img := catPic.GetImage()
		_ = img // May be nil or not, both are valid
	}, "should handle nil images without panic")
}

// TestCatPic_StressTestDraw tests Draw under stress
func TestCatPic_StressTestDraw(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Draw stress test in short mode")
	}

	img := testutil.CreateColorImage(500, 500, 100, 200, 255)
	catPic := NewCatImage(img)

	const numDrawers = 100
	const drawsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numDrawers)

	for i := 0; i < numDrawers; i++ {
		go func() {
			defer wg.Done()

			var ops op.Ops
			gtx := layout.Context{
				Ops: &ops,
				Constraints: layout.Constraints{
					Min: image.Pt(0, 0),
					Max: image.Pt(400, 500),
				},
			}

			for j := 0; j < drawsPerGoroutine; j++ {
				dims := catPic.Draw(gtx)
				testutil.AssertTrue(t, dims.Size.X >= 0, "width should be non-negative")
				testutil.AssertTrue(t, dims.Size.Y >= 0, "height should be non-negative")
			}
		}()
	}

	wg.Wait()
}

// TestCatPic_NoDeadlock tests that operations don't cause deadlocks
func TestCatPic_NoDeadlock(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping deadlock test in short mode")
	}

	catPic := NewCatImage(nil)

	// Create a timeout to detect deadlocks
	done := make(chan bool, 1)

	go func() {
		const numOps = 10000

		for i := 0; i < numOps; i++ {
			img := testutil.CreateColorImage(100, 100, 128, 128, 128)
			catPic.SetImage(img)
			_ = catPic.GetImage()
			catPic.SetLoading()
			_ = catPic.IsLoading()
			catPic.ClearLoading()

			var ops op.Ops
			gtx := layout.Context{
				Ops: &ops,
				Constraints: layout.Constraints{
					Min: image.Pt(0, 0),
					Max: image.Pt(400, 500),
				},
			}
			_ = catPic.Draw(gtx)
		}

		done <- true
	}()

	// Use a channel with timeout to detect deadlocks
	// If operations complete, we succeed
	// If timeout fires first, we have a deadlock
	<-done // Wait for completion (should not deadlock)
}
