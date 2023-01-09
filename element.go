package pdftpl

import (
	"image"

	"github.com/signintech/gopdf"
)

var (
	_ element = &textElement{}
	_ element = &imageElement{}
)

type element interface {
	draw(pdf *gopdf.GoPdf, debug bool) error
}

type textElement struct {
	textTag
	text string
}

func (e *textElement) draw(pdf *gopdf.GoPdf, debug bool) error {
	if debug {
		pdf.SetStrokeColor(255, 128, 128)
		pdf.SetLineWidth(2)
		pdf.Rectangle(e.X, e.Y, e.X+e.W, e.Y+10, "D", 3, 10)
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

func (e *imageElement) draw(pdf *gopdf.GoPdf, debug bool) error {
	if e.image != nil {
		pdf.ImageFrom(e.image, e.X, e.Y, &gopdf.Rect{W: e.W, H: e.H})
	}
	return nil
}
