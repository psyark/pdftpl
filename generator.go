package pdftpl

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

// NewGenerator は新しいGeneratorを返します
func NewGenerator() *Generator {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	return &Generator{pdf: pdf}
}

type Generator struct {
	pdf *gopdf.GoPdf
}

func (gen *Generator) RegisterFont(name string, ttfData []byte) error {
	return gen.pdf.AddTTFFontData(name, ttfData)
}

// RegisterTemplate はテンプレートとなるPDFをページ単位で登録します
func (gen *Generator) RegisterTemplate(src io.ReadSeeker, pageNumber int) (tmpl *Template, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("panic: %v", err)
		}
	}()
	id := gen.pdf.ImportPageStream(&src, pageNumber, "/MediaBox")
	tmpl = &Template{id}
	return
}

// AddPage はページを追加します
func (gen *Generator) AddPage(vars interface{}, tpl *Template) error {
	return gen.addPage(vars, tpl, false)
}

// AddPageDebug はデバッグ情報付きでページを追加します
func (gen *Generator) AddPageDebug(vars interface{}, tpl *Template) error {
	return gen.addPage(vars, tpl, true)
}

func (gen *Generator) addPage(vars interface{}, tpl *Template, debug bool) error {
	// ページ追加
	gen.pdf.AddPage()
	if tpl != nil {
		gen.pdf.UseImportedTemplate(tpl.id, 0, 0, gopdf.PageSizeA4.W, gopdf.PageSizeA4.H)
	}

	texts, err := parseVars(vars)
	if err != nil {
		return errors.Wrap(err, "parseVars")
	}

	for _, t := range texts {
		if debug {
			gen.pdf.SetStrokeColor(255, 128, 128)
			gen.pdf.SetLineWidth(2)
			gen.pdf.Rectangle(t.X, t.Y, t.X+t.W, t.Y+10, "D", 3, 10)
		}

		gen.pdf.SetFont(t.FontFace, "", t.FontSize)
		gen.pdf.SetX(t.X)
		gen.pdf.SetY(t.Y)

		if t.Text != "" {
			texts, err := gen.pdf.SplitTextWithWordWrap(t.Text, t.W)
			if err != nil {
				return err
			}
			for _, text := range texts {
				gen.pdf.MultiCellWithOption(&gopdf.Rect{W: t.W, H: gopdf.PageSizeA4.H}, text, gopdf.CellOption{Align: t.Align})
			}
		}
	}

	return nil
}

// Generate はPDFのバイト列を返却します
func (gen *Generator) Generate() ([]byte, error) {
	return gen.pdf.GetBytesPdfReturnErr()
}

// Template はGeneratorに登録されたテンプレートです
type Template struct {
	id int
}
