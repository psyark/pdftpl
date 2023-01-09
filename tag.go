package pdftpl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	"github.com/signintech/gopdf"
)

// `pdftpl:"x=425,y=50,w=120,s=12,f=gothic,a=r"`
type textTag struct {
	X           float64 `schema:"x,required"`
	Y           float64 `schema:"y,required"`
	W           float64 `schema:"w,required"`
	FontSize    float64 `schema:"s,required"`
	FontFace    string  `schema:"f"`
	LineHeight  float64 `schema:"lh"`
	AlignString string  `schema:"a"`
	Align       int
}

func (t textTag) fromOrigin(x float64, y float64) textTag {
	t.X += x
	t.Y += y
	return t
}

func (t *textTag) parse(tagStr string) error {
	vars, err := url.ParseQuery(strings.ReplaceAll(tagStr, ",", "&"))
	if err != nil {
		return fmt.Errorf("url.ParseQuery: %w", err)
	}

	if err := schema.NewDecoder().Decode(t, vars); err != nil {
		return err
	}

	switch t.AlignString {
	case "l", "":
		t.Align = gopdf.Left
	case "r":
		t.Align = gopdf.Right
	case "c":
		t.Align = gopdf.Center
	default:
		return fmt.Errorf("unsupported align: %v", t.AlignString)
	}

	return nil
}

// `pdftpl:"x=425,y=50,w=120,h=120,f=contain"`
type imageTag struct {
	X   float64 `schema:"x,required"`
	Y   float64 `schema:"y,required"`
	W   float64 `schema:"w,required"`
	H   float64 `schema:"h,required"`
	Fit string  `schema:"f,required"`
}

func (t imageTag) fromOrigin(x float64, y float64) imageTag {
	t.X += x
	t.Y += y
	return t
}

func (t *imageTag) parse(tagStr string) error {
	vars, err := url.ParseQuery(strings.ReplaceAll(tagStr, ",", "&"))
	if err != nil {
		return fmt.Errorf("url.ParseQuery: %w", err)
	}

	if err := schema.NewDecoder().Decode(t, vars); err != nil {
		return err
	}

	return nil
}

type transTag struct {
	Dx float64 `schema:"dx"`
	Dy float64 `schema:"dy"`
}

func (t *transTag) parse(tagStr string) error {
	vars, err := url.ParseQuery(strings.ReplaceAll(tagStr, ",", "&"))
	if err != nil {
		return fmt.Errorf("url.ParseQuery: %w", err)
	}

	return schema.NewDecoder().Decode(t, vars)
}
