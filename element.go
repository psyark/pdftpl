package pdftpl

import (
	"image"
	"image/color"

	"github.com/signintech/gopdf"
)

var (
	_ element = &textElement{}
	_ element = &imageElement{}
)

type element interface {
	draw(pdf *gopdf.GoPdf, options *addPageOptions) error
}

type textElement struct {
	textTag
	text string
}

func (e *textElement) draw(pdf *gopdf.GoPdf, options *addPageOptions) error {
	if options.DebugBorderColor != color.Transparent {
		r, g, b, _ := options.DebugBorderColor.RGBA()
		pdf.SetStrokeColor(uint8(r>>8), uint8(g>>8), uint8(b>>8))
		pdf.SetLineWidth(2)
		if err := pdf.Rectangle(e.X, e.Y, e.X+e.W, e.Y+10, "D", 3, 10); err != nil {
			return err
		}
	}

	pdf.SetFont(e.FontFace, "", e.FontSize)
	pdf.SetX(e.X)
	pdf.SetY(e.Y)

	if e.text != "" {
		texts, err := pdf.SplitTextWithWordWrap(e.text, e.W)
		if err != nil {
			return err
		}
		for _, text := range texts {
			pdf.MultiCellWithOption(&gopdf.Rect{W: e.W, H: gopdf.PageSizeA4.H}, text, gopdf.CellOption{Align: e.Align})

			if e.LineHeight != 0 {
				pdf.SetY(pdf.GetY() + e.FontSize*(e.LineHeight-1))
			}
		}
	}

	return nil
}

type imageElement struct {
	imageTag
	image image.Image
}

func (e *imageElement) draw(pdf *gopdf.GoPdf, options *addPageOptions) error {
	if e.image != nil {
		pdf.ImageFrom(e.image, e.X, e.Y, &gopdf.Rect{W: e.W, H: e.H})
	}
	return nil
}
