package pdftpl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

// `pdftpl:"x=425,y=50,w=120,s=12,f=gothic,a=r"`

type tag struct {
	x, y, w  float64
	fontSize int
	fontFace string
	align    int
}

func (t tag) move(x float64, y float64) tag {
	t.x += x
	t.y += y
	return t
}

func parseTag(tagStr string) (*tag, error) {
	vars, err := url.ParseQuery(strings.ReplaceAll(tagStr, ",", "&"))
	if err != nil {
		return nil, errors.Wrap(err, "url.ParseQuery")
	}

	t := tag{}
	if _, err := fmt.Sscanf(vars.Get("x"), "%f", &t.x); err != nil {
		return nil, fmt.Errorf("xパラメータが不正です: %v", tagStr)
	}
	if _, err := fmt.Sscanf(vars.Get("y"), "%f", &t.y); err != nil {
		return nil, fmt.Errorf("yパラメータが不正です: %v", tagStr)
	}
	if _, err := fmt.Sscanf(vars.Get("w"), "%f", &t.w); err != nil {
		return nil, fmt.Errorf("wパラメータが不正です: %v", tagStr)
	}
	if _, err := fmt.Sscanf(vars.Get("s"), "%d", &t.fontSize); err != nil {
		return nil, fmt.Errorf("sパラメータが不正です: %v", tagStr)
	}
	if vars.Has("f") {
		if _, err := fmt.Sscanf(vars.Get("f"), "%s", &t.fontFace); err != nil {
			return nil, fmt.Errorf("fパラメータが不正です: %v", tagStr)
		}
	}
	if vars.Has("a") {
		switch vars.Get("a") {
		case "c":
			t.align = gopdf.Center
		case "l":
			t.align = gopdf.Left
		case "r":
			t.align = gopdf.Right
		default:
			return nil, fmt.Errorf("unsupported align: %v", vars.Get("a"))
		}
	}
	return &t, nil
}

func parseRelTag(tagStr string) (float64, float64, error) {
	var dx, dy float64

	vars, err := url.ParseQuery(strings.ReplaceAll(tagStr, ",", "&"))
	if err != nil {
		return 0, 0, errors.Wrap(err, "url.ParseQuery")
	}

	if vars.Has("dx") {
		if _, err := fmt.Sscanf(vars.Get("dx"), "%f", &dx); err != nil {
			return 0, 0, fmt.Errorf("dxパラメータが不正です: %v", tagStr)
		}
	}
	if vars.Has("dy") {
		if _, err := fmt.Sscanf(vars.Get("dy"), "%f", &dy); err != nil {
			return 0, 0, fmt.Errorf("dyパラメータが不正です: %v", tagStr)
		}
	}

	return dx, dy, nil
}
