package pdftpl

import (
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

// Template はGeneratorに登録されたテンプレートです
type Template struct {
	pdf *gopdf.GoPdf
	id  int
}

// AddPage はページを追加します
func (tpl *Template) AddPage(vars interface{}) error {
	return tpl.addPage(vars, false)
}

// AddPageDebug はデバッグ情報付きでページを追加します
func (tpl *Template) AddPageDebug(vars interface{}) error {
	return tpl.addPage(vars, true)
}

func (tpl *Template) addPage(vars interface{}, debug bool) error {
	// ページ追加
	tpl.pdf.AddPage()
	tpl.pdf.UseImportedTemplate(tpl.id, 0, 0, gopdf.PageSizeA4.W, gopdf.PageSizeA4.H)

	cb := func(text string, t tag) {
		if debug {
			tpl.pdf.SetStrokeColor(255, 128, 128)
			tpl.pdf.SetLineWidth(2)
			tpl.pdf.Rectangle(t.x, t.y, t.x+t.w, t.y+10, "D", 3, 10)
		}

		tpl.pdf.SetFont(t.fontFace, "", t.fontSize)
		tpl.pdf.SetX(t.x)
		tpl.pdf.SetY(t.y)
		tpl.pdf.MultiCellWithOption(&gopdf.Rect{W: t.w, H: 1000}, text, gopdf.CellOption{
			Align: t.align,
			// Align: gopdf.Right | gopdf.Center,
		})
	}

	if err := parseVars(vars, cb); err != nil {
		return errors.Wrap(err, "parseVars")
	}
	return nil
}
