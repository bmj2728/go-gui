package utils

import (
	"image"
	"sync"

	"gioui.org/layout"
	"gioui.org/op/paint"
)

type CurrentImage struct {
	img       image.Image
	mu        sync.Mutex
	isLoading bool
}

func DrawImage(gtx layout.Context, img image.Image) layout.Dimensions {
	// Convert to paint.ImageOp
	imgOp := paint.NewImageOp(img)
	imgOp.Filter = paint.FilterLinear

	// Scale to fit available space while maintaining aspect ratio
	bounds := img.Bounds()
	imgW, imgH := float32(bounds.Dx()), float32(bounds.Dy())
	maxW, maxH := float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y)

	scale := min(maxW/imgW, maxH/imgH)
	finalW, finalH := int(imgW*scale), int(imgH*scale)

	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Pt(finalW, finalH)}
}
