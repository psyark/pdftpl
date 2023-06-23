package pdftpl_test

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/aquilax/go-perlin"
	"github.com/psyark/pdftpl"
	"github.com/signintech/gopdf"
)

var (
	//go:embed "testdata/fonts/ipaexm.ttf"
	ipaexmBytes []byte
	//go:embed "testdata/fonts/ipaexg.ttf"
	ipaexgBytes []byte
	//go:embed "testdata/pdf-templates/meter-clinic-building.pdf"
	templateBytes []byte
)

func TestNewGenerator(t *testing.T) {
	gen := pdftpl.NewGenerator()

	if err := gen.RegisterFont("", ipaexgBytes); err != nil {
		t.Fatal(err)
	}
	if err := gen.RegisterFont("ipaexm", ipaexmBytes); err != nil {
		t.Fatal(err)
	}

	tpl, err := gen.RegisterPageTemplate(gopdf.PageSizeA4, templateBytes, 1)
	if err != nil {
		t.Fatal(err)
	}

	vars := PDFVars{
		Title:     "電気料金請求書",
		Date:      "2022/09/07",
		Recipient: "あいうえおかきくけこさしすせそたちつてと\nなにぬねの",
		Image:     createMockImage(),
	}

	for i := range vars.Meters {
		vars.Meters[i] = PDFMeterVars{
			Label:   fmt.Sprintf("ラベル%d", i),
			Prev:    fmt.Sprintf("%d", i),
			Current: fmt.Sprintf("%d", i*10),
			Delta:   fmt.Sprintf("%d", i*9),
		}
	}

	if err := gen.AddPageWithTemplate(vars, tpl, pdftpl.Debug(color.RGBA{R: 0xFF, A: 0xFF})); err != nil {
		t.Fatal(err)
	}
	if err := gen.AddPage(vars, gopdf.PageSizeA4Landscape, pdftpl.Debug(color.RGBA{B: 0xFF, A: 0xFF})); err != nil {
		t.Fatal(err)
	}

	pdfData, err := gen.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("testdata/out.pdf", pdfData, 0666); err != nil {
		t.Fatal(err)
	}
}

// PDFVars は請求書のpdftpl用パラメータです
type PDFVars struct {
	Image      image.Image         `pdftpl:"x=50,y=50,w=120,h=120,f=contain"`
	Date       string              `pdftpl:"x=425,y=50,w=120,s=12,a=r"`
	Title      string              `pdftpl:"x=50,y=100,w=495,s=26,a=c,f=ipaexm"`
	Recipient  string              `pdftpl:"x=50,y=190,w=300,s=18,a=l"`
	Amount     string              `pdftpl:"x=132,y=252,w=95,s=12,a=r"`
	Room       string              `pdftpl:"x=132,y=300,w=145,s=10,a=c"`
	Period     string              `pdftpl:"x=132,y=316,w=145,s=10,a=c"`
	Meters     [4]PDFMeterVars     `pdftpl:"dy=16"`
	TotalDelta string              `pdftpl:"x=494.5,y=347,w=48,s=10,a=r"`
	Statements [2]PDFStatementVars `pdftpl:"dy=16"`
	Amount2    string              `pdftpl:"x=397,y=629,w=138,s=10,a=r"`
	Prompt     string              `pdftpl:"x=220,y=662,w=80,s=10,a=c"`
}

// PDFMeterVars は請求書のうちメーター欄のpdftpl用パラメータです
type PDFMeterVars struct {
	Label   string `pdftpl:"x=337,y=283,w=48,s=10,a=c"`
	Prev    string `pdftpl:"x=389.5,y=283,w=48,s=10,a=r"`
	Current string `pdftpl:"x=442,y=283,w=48,s=10,a=r"`
	Delta   string `pdftpl:"x=494.5,y=283,w=48,s=10,a=r"`
}

// PDFStatementVars は請求書のうち詳細欄のpdftpl用パラメータです
type PDFStatementVars struct {
	Title  string `pdftpl:"x=59,y=393,w=320,s=10,a=l"`
	Amount string `pdftpl:"x=397,y=393,w=138,s=10,a=r"`
}

func createMockImage() image.Image {
	noise := perlin.NewPerlin(2, 2, 3, 0)
	mock := image.NewGray(image.Rect(0, 0, 256, 256))
	for y := 0; y < mock.Rect.Dy(); y++ {
		for x := 0; x < mock.Rect.Dx(); x++ {
			fx := float64(x) / float64(mock.Rect.Dx())
			fy := float64(y) / float64(mock.Rect.Dy())
			mock.SetGray(x, y, color.Gray{Y: uint8(noise.Noise2D(fx, fy)*127 + 127)})
		}
	}
	return mock
}
