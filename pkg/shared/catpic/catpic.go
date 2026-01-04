package catpic

import (
	"image"
	"sync"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

type CatPic struct {
	img       image.Image
	mu        sync.Mutex
	isLoading bool
}

func NewCatImage(img image.Image) *CatPic {
	return &CatPic{
		img: img,
	}
}

func (p *CatPic) IsLoading() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isLoading
}

func (p *CatPic) GetImage() image.Image {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.img
}

func (p *CatPic) SetImage(img image.Image) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.img = img
}

func (p *CatPic) SetLoading() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isLoading = true
}

func (p *CatPic) ClearLoading() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isLoading = false
}

func (p *CatPic) Draw(gtx layout.Context) layout.Dimensions {
	img := p.GetImage()
	if img == nil {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}

	return widget.Image{
		Src:      paint.NewImageOp(img),
		Fit:      widget.Contain,
		Position: layout.Center,
	}.Layout(gtx)
}
