package pdftpl

import (
	"bytes"
	"fmt"
	"image/color"
	"io"

	"github.com/signintech/gopdf"
)

// NewGenerator は新しいGeneratorを返します
func NewGenerator() *Generator {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})
	return &Generator{gopdf: pdf}
}

type Generator struct {
	gopdf *gopdf.GoPdf
}

func (gen *Generator) RegisterFont(name string, ttfData []byte) error {
	return gen.gopdf.AddTTFFontData(name, ttfData)
}

// RegisterPageTemplate はテンプレートとなるPDFをページ単位で登録します
func (gen *Generator) RegisterPageTemplate(pageSize *gopdf.Rect, pdfBytes []byte, pageNumber int) (basePDF *PageTemplate, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("panic: %v", err)
		}
	}()

	var readSeeker io.ReadSeeker = bytes.NewReader(pdfBytes)
	basePDF = &PageTemplate{
		templateID: gen.gopdf.ImportPageStream(&readSeeker, pageNumber, "/MediaBox"),
		pageSize:   pageSize,
	}
	return
}

// AddPageWithTemplate はページテンプレートを指定してページを追加します
func (gen *Generator) AddPageWithTemplate(vars any, tpl *PageTemplate, options ...AddPageOption) error {
	return gen.addPage(vars, tpl.pageSize, tpl, options...)
}

// AddPage はページサイズを指定してページを追加します
func (gen *Generator) AddPage(vars any, pageSize *gopdf.Rect, options ...AddPageOption) error {
	return gen.addPage(vars, pageSize, nil, options...)
}

func (gen *Generator) addPage(vars any, pageSize *gopdf.Rect, tpl *PageTemplate, options ...AddPageOption) error {
	opts := &addPageOptions{
		DebugBorderColor: color.Transparent,
	}
	for _, o := range options {
		o(opts)
	}

	// ページ追加
	gen.gopdf.AddPageWithOption(gopdf.PageOption{PageSize: pageSize})
	if tpl != nil {
		gen.gopdf.UseImportedTemplate(tpl.templateID, 0, 0, pageSize.W, pageSize.H)
	}

	elements, err := parseVars(vars)
	if err != nil {
		return fmt.Errorf("parseVars: %w", err)
	}

	for _, element := range elements {
		if err := element.draw(gen.gopdf, opts); err != nil {
			return err
		}
	}

	return nil
}

// Generate はPDFのバイト列を返却します
func (gen *Generator) Generate() ([]byte, error) {
	return gen.gopdf.GetBytesPdfReturnErr()
}

// PageTemplate は登録されたページテンプレートです
type PageTemplate struct {
	templateID int
	pageSize   *gopdf.Rect
}
