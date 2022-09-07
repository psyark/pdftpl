package pdftpl

import (
	"fmt"
	"io"

	"github.com/signintech/gopdf"
)

// Generator はPDFを出力するためのインターフェースです
// 名前を変えるかもしれません
type Generator interface {
	RegisterFont(name string, ttfData []byte) error
	RegisterTemplate(io.ReadSeeker, int) (*Template, error)
	Generate(io.Writer) error
}

// NewGenerator は新しいGeneratorを返します
func NewGenerator() (Generator, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	return &generatorImpl{pdf: pdf}, nil
}

type generatorImpl struct {
	pdf *gopdf.GoPdf
}

func (gen *generatorImpl) RegisterFont(name string, ttfData []byte) error {
	return gen.pdf.AddTTFFontData(name, ttfData)
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
