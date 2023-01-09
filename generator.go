package pdftpl

import (
	"fmt"
	"io"

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

	elements, err := parseVars(vars)
	if err != nil {
		return fmt.Errorf("parseVars: %w", err)
	}

	for _, element := range elements {
		if err := element.draw(gen.pdf, debug); err != nil {
			return err
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
