package pdftpl

import (
	"bytes"
	"fmt"
	"io"

	"github.com/signintech/gopdf"
)

// NewGenerator は新しいGeneratorを返します
func NewGenerator(pageSize *gopdf.Rect) *Generator {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *pageSize})
	return &Generator{gopdf: pdf, pageSize: pageSize}
}

type Generator struct {
	gopdf    *gopdf.GoPdf
	pageSize *gopdf.Rect
}

func (gen *Generator) RegisterFont(name string, ttfData []byte) error {
	return gen.gopdf.AddTTFFontData(name, ttfData)
}

// RegisterTemplate はテンプレートとなるPDFをページ単位で登録します
func (gen *Generator) RegisterTemplate(data []byte, pageNumber int) (tmpl *Template, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("panic: %v", err)
		}
	}()

	var rs io.ReadSeeker = bytes.NewReader(data)
	id := gen.gopdf.ImportPageStream(&rs, pageNumber, "/MediaBox")
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
	gen.gopdf.AddPage()
	if tpl != nil {
		gen.gopdf.UseImportedTemplate(tpl.id, 0, 0, gen.pageSize.W, gen.pageSize.H)
	}

	elements, err := parseVars(vars)
	if err != nil {
		return fmt.Errorf("parseVars: %w", err)
	}

	for _, element := range elements {
		if err := element.draw(gen.gopdf, debug); err != nil {
			return err
		}
	}

	return nil
}

// Generate はPDFのバイト列を返却します
func (gen *Generator) Generate() ([]byte, error) {
	return gen.gopdf.GetBytesPdfReturnErr()
}

// Template はGeneratorに登録されたテンプレートです
type Template struct {
	id int
}
