// Package pdftpl はテンプレートのPDFから値の埋め込まれたPDFを出力するためのパッケージです
package pdftpl

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

const fontIPAexGothic = "IPAexGothic"

// Generator はPDFを出力するためのインターフェースです
// 名前を変えるかもしれません
type Generator interface {
	RegisterTemplate(io.ReadSeeker, int) (*Template, error)
	Generate(io.Writer) error
}

// NewGenerator は新しいGeneratorを返します
func NewGenerator() (Generator, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := pdf.AddTTFFont(fontIPAexGothic, "./fonts/ipaexg.ttf"); err != nil {
		return nil, errors.Wrap(err, "AddTTFFont")
	}
	return &generatorImpl{pdf: pdf}, nil
}

type generatorImpl struct {
	pdf *gopdf.GoPdf
}

// RegisterTemplate はテンプレートとなるPDFをページ単位で登録します
func (gen *generatorImpl) RegisterTemplate(reader io.ReadSeeker, page int) (tmpl *Template, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("panic: %v", err)
		}
	}()
	id := gen.pdf.ImportPageStream(&reader, page, "/MediaBox")
	tmpl = &Template{pdf: gen.pdf, id: id}
	return
}

// Generate はPDFを出力します
func (gen *generatorImpl) Generate(writer io.Writer) error {
	if err := gen.pdf.Write(writer); err != nil {
		gen.pdf.Close()
		return err
	}
	return gen.pdf.Close()
}

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

		tpl.pdf.SetFont(fontIPAexGothic, "", t.fontSize)
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
