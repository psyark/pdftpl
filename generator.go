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
func (gen *Generator) AddPage(tpl *Template, vars interface{}) error {
	return gen.addPage(tpl, vars, false)
}

// AddPageDebug はデバッグ情報付きでページを追加します
func (gen *Generator) AddPageDebug(tpl *Template, vars interface{}) error {
	return gen.addPage(tpl, vars, true)
}

func (gen *Generator) addPage(tpl *Template, vars interface{}, debug bool) error {
	// ページ追加
	gen.pdf.AddPage()
	gen.pdf.UseImportedTemplate(tpl.id, 0, 0, gopdf.PageSizeA4.W, gopdf.PageSizeA4.H)

	cb := func(text string, t tag) {
		if debug {
			gen.pdf.SetStrokeColor(255, 128, 128)
			gen.pdf.SetLineWidth(2)
			gen.pdf.Rectangle(t.x, t.y, t.x+t.w, t.y+10, "D", 3, 10)
		}

		gen.pdf.SetFont(t.fontFace, "", t.fontSize)
		gen.pdf.SetX(t.x)
		gen.pdf.SetY(t.y)
		gen.pdf.MultiCellWithOption(&gopdf.Rect{W: t.w, H: gopdf.PageSizeA4.H}, text, gopdf.CellOption{Align: t.align})
	}

	if err := parseVars(vars, cb); err != nil {
		return errors.Wrap(err, "parseVars")
	}
	return nil
}

// Generate はPDFを出力します
func (gen *Generator) Generate(writer io.Writer) error {
	if err := gen.pdf.Write(writer); err != nil {
		gen.pdf.Close()
		return err
	}
	return gen.pdf.Close()
}

// Template はGeneratorに登録されたテンプレートです
type Template struct {
	id int
}
