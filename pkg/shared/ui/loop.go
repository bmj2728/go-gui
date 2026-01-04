package ui

import (
	"image"
	"image/color"
	//"image"
	"log"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/bmj2728/catfetch/pkg/shared/catpic"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Run(w *app.Window) error {
	// button
	var fetchButton widget.Clickable
	// thread-safe image wrapper
	var currentImage catpic.CatPic //threadsafe wrapper for image.Image
	// Ops list
	var ops op.Ops

	newBg := color.NRGBA{R: 40, G: 42, B: 54, A: 255}

	// Theme for material widgets
	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Draw background
			winRect := clip.Rect{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
			}
			paint.FillShape(&ops, newBg, winRect.Op())

			// Handle button click
			if fetchButton.Clicked(gtx) && !currentImage.IsLoading() {
				currentImage.SetLoading()
				go func(wind *app.Window) {
					img, _, err := HandleButtonClick()
					if err != nil {
						log.Printf("Error handling button click: %v", err)
					} else {
						currentImage.SetImage(img)
					}
					currentImage.ClearLoading()
					wind.Invalidate()
				}(w)
			}

			// Layout UI components
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layoutButton(gtx, th, &fetchButton, 12)
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layoutImageDisplay(gtx, &currentImage, 24)
				}),
			)

			e.Frame(gtx.Ops)

		}
	}
}

// layoutButton renders the fetch button with padding and styling
func layoutButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, insetPixels unit.Dp) layout.Dimensions {
	inset := layout.UniformInset(insetPixels)

	dims := layoutButtonDims(gtx, inset, th, btn)

	return dims

}

func layoutButtonDims(gtx layout.Context, inset layout.Inset, th *material.Theme, btn *widget.Clickable) layout.Dimensions {
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Create button with styling
		button := material.Button(th, btn, "Fetch a Cat")
		button.CornerRadius = unit.Dp(16)
		button.Background = color.NRGBA{R: 189, G: 147, B: 249, A: 255}
		button.Color = color.NRGBA{R: 248, G: 248, B: 242, A: 255}

		// Set fixed button size
		gtx.Constraints.Min.X = gtx.Dp(120)
		gtx.Constraints.Max.X = gtx.Dp(120)
		gtx.Constraints.Min.Y = gtx.Dp(40)
		gtx.Constraints.Max.Y = gtx.Dp(40)

		return button.Layout(gtx)
	})
}

// layoutImageDisplay renders the image display area with padding
func layoutImageDisplay(gtx layout.Context, img *catpic.CatPic, insetPixels unit.Dp) layout.Dimensions {
	// Create the inset
	inset := layout.UniformInset(insetPixels)

	dims := layoutImageDisplayDims(gtx, img, inset)

	return dims

}

func layoutImageDisplayDims(gtx layout.Context, img *catpic.CatPic, inset layout.Inset) layout.Dimensions {
	return inset.Layout(gtx, img.Draw)
}
