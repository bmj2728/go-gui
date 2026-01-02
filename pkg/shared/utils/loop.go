package utils

import (
	//"image"
	"log"

	"go-gui/pkg/shared/handlers"

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
	var currentImage CurrentImage
	// Theme for material widgets
	th := material.NewTheme()

	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle button click
			if fetchButton.Clicked(gtx) && !currentImage.isLoading {
				go func(wind *app.Window) {
					img, err := handlers.HandleButtonClick()
					currentImage.mu.Lock()
					currentImage.isLoading = true
					if err != nil {
						log.Printf("Error handling button click: %v", err)
					} else {
						currentImage.img = img
					}
					currentImage.isLoading = false
					currentImage.mu.Unlock()
					wind.Invalidate()
				}(w)

			}

			// Layout
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				// Button at top
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &fetchButton, "Fetch Image")
						return btn.Layout(gtx)
					})
				}),

				// Image display area
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if currentImage.img == nil {
							// Show placeholder
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}
						return DrawImage(gtx, currentImage.img)
					})
				}),
			)

			e.Frame(gtx.Ops)

		}
	}
}
