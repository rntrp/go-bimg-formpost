package img

type Resize int

const (
	Fit Resize = iota
	FitBlack
	FitWhite
	FitUpscale
	FitUpscaleBlack
	FitUpscaleWhite
	Fill
	FillTopLeft
	FillTop
	FillTopRight
	FillLeft
	FillRight
	FillBottomLeft
	FillBottom
	FillBottomRight
	Stretch
)
