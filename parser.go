package pdftpl

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type taggedText struct {
	textTag
	Text string
}

// パラメータ構造体をパースしてpdftplタグ付きフィールドに対してコールバックを呼び出します
func parseVars(vars interface{}) ([]taggedText, error) {
	return parseVarsInternal(reflect.ValueOf(vars), 0, 0)
}

func parseVarsInternal(v reflect.Value, ox float64, oy float64) ([]taggedText, error) {
	texts := []taggedText{}
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("vars は構造体である必要があります")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tagStr, ok := f.Tag.Lookup("pdftpl")
		if !ok {
			continue
		}

		switch f.Type.Kind() {
		case reflect.String:
			tag := &textTag{}
			if err := tag.parse(tagStr); err != nil {
				return nil, errors.Wrap(err, "parseTag")
			}

			tt := taggedText{textTag: tag.fromOrigin(ox, oy), Text: v.Field(i).String()}
			texts = append(texts, tt)
		case reflect.Array, reflect.Slice:
			tag := &transTag{}
			if err := tag.parse(tagStr); err != nil {
				return nil, errors.Wrap(err, "parseTag")
			}

			for j := 0; j < v.Field(i).Len(); j++ {
				fj := float64(j)
				// 座標を変えつつ再帰呼び出し
				texts2, err := parseVarsInternal(v.Field(i).Index(j), ox+tag.Dx*fj, oy+tag.Dy*fj)
				if err != nil {
					return nil, err
				}
				texts = append(texts, texts2...)
			}

		case reflect.Struct:
			tag := &transTag{}
			if err := tag.parse(tagStr); err != nil {
				return nil, errors.Wrap(err, "parseTag")
			}

			// 座標を変えつつ再帰呼び出し
			texts2, err := parseVarsInternal(v.Field(i), ox+tag.Dx, oy+tag.Dy)
			if err != nil {
				return nil, err
			}
			texts = append(texts, texts2...)

		default:
			return nil, fmt.Errorf("vars の %vフィールドが未対応の型 (%v)です", f.Name, f.Type)
		}
	}

	return texts, nil
}
