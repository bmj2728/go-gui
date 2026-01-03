package ui

import (
	//"image"
	"log"

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
	var currentImage catpic.CatPic
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
			if fetchButton.Clicked(gtx) && !currentImage.IsLoading() {
				currentImage.SetLoading()
				go func(wind *app.Window) {
					img, _, err := HandleButtonClick()
					if err != nil {
						log.Printf("Error handling button click: %v", err)
					} else {
						currentImage.SetImage(img)
						//fmt.Println(meta.Tags)
					}
					currentImage.ClearLoading()
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
						if currentImage.GetImage() == nil {
							// Show placeholder
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}
						return currentImage.Draw(gtx)
					})
				}),
			)

			e.Frame(gtx.Ops)

		}
	}
}
