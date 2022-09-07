package pdftpl

import (
	"fmt"
	"strings"

	"github.com/signintech/gopdf"
)

type tag struct {
	x, y, w, h float64
	fontSize   int
	align      int
}

func (t tag) move(x float64, y float64) tag {
	t.x += x
	t.y += y
	return t
}

func parseTag(tagStr string) (*tag, error) {
	t := tag{}
	var align string

	if _, err := fmt.Sscanf(tagStr, "%f,%f,%f,%f,%d,%s", &t.x, &t.y, &t.w, &t.h, &t.fontSize, &align); err != nil {
		return nil, fmt.Errorf("pdftplタグが不正です: %v", tagStr)
	}

	for _, m := range strings.Split(align, ",") {
		switch m {
		case "top":
			t.align |= gopdf.Top
		case "middle":
			t.align |= gopdf.Middle
		case "bottom":
			t.align |= gopdf.Bottom
		case "left":
			t.align |= gopdf.Left
		case "center":
			t.align |= gopdf.Center
		case "right":
			t.align |= gopdf.Right
		}
	}

	return &t, nil
}

func parseRelTag(tagStr string) (float64, float64, error) {
	var x, y float64
	if _, err := fmt.Sscanf(tagStr, "%f,%f", &x, &y); err != nil {
		return 0, 0, fmt.Errorf("pdftpl相対タグが不正です: %v", tagStr)
	}

	return x, y, nil
}
