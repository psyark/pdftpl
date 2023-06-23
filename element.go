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
	pdf.SetFont(e.FontFace, "", e.FontSize)
	pdf.SetX(e.X)
	pdf.SetY(e.Y)

	if options.DebugBorderColor != color.Transparent {
		lines := 1
		if e.text != "" {
			texts, err := pdf.SplitTextWithWordWrap(e.text, e.W)
			if err != nil {
				return err
			}
			lines = len(texts)
		}

		box := gopdf.Box{
			Left:   e.X,
			Top:    e.Y,
			Right:  e.X + e.W,
			Bottom: e.Y + float64(lines)*e.FontSize,
		}
		if e.LineHeight != 0 {
			box.Bottom = e.Y + (float64(lines-1)*e.LineHeight+1)*e.FontSize
		}

		if err := drawDebugBorder(pdf, box, options.DebugBorderColor); err != nil {
			return err
		}
	}

	if e.text != "" {
		texts, err := pdf.SplitTextWithWordWrap(e.text, e.W)
		if err != nil {
			return err
		}

		for _, text := range texts {
			err := pdf.MultiCellWithOption(&gopdf.Rect{W: e.W, H: gopdf.PageSizeA4.H}, text, gopdf.CellOption{Align: e.Align})
			if err != nil {
				return err
			}

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
	if options.DebugBorderColor != color.Transparent {
		box := gopdf.Box{Left: e.X - 1, Top: e.Y - 1, Right: e.X + e.W + 2, Bottom: e.Y + e.H + 2}
		if err := drawDebugBorder(pdf, box, options.DebugBorderColor); err != nil {
			return err
		}
	}

	if e.image != nil {
		return pdf.ImageFrom(e.image, e.X, e.Y, &gopdf.Rect{W: e.W, H: e.H})
	}
	return nil
}

func drawDebugBorder(pdf *gopdf.GoPdf, box gopdf.Box, color color.Color) error {
	r, g, b, _ := color.RGBA()
	pdf.SetStrokeColor(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	pdf.SetLineWidth(2)
	return pdf.Rectangle(box.Left, box.Top, box.Right, box.Bottom, "D", 3, 10)
}
