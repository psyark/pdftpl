package pdftpl

import "image/color"

type addPageOptions struct {
	DebugBorderColor color.Color
}
type AddPageOption func(o *addPageOptions)

func Debug(borderColor color.Color) AddPageOption {
	return func(o *addPageOptions) {
		o.DebugBorderColor = borderColor
	}
}
