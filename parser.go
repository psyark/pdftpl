package pdftpl

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// pdftplタグ付きフィールドに対するコールバック
type parseCallback func(string, tag)

// パラメータ構造体をパースしてpdftplタグ付きフィールドに対してコールバックを呼び出します
func parseVars(vars interface{}, cb parseCallback) error {
	return parseVarsInternal(reflect.ValueOf(vars), cb, 0, 0)
}

func parseVarsInternal(v reflect.Value, cb parseCallback, x float64, y float64) error {
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("vars は構造体である必要があります")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tagStr, ok := f.Tag.Lookup("pdftpl")
		if !ok {
			continue
		}

		switch f.Type.Kind() {
		case reflect.String:
			tag, err := parseTag(tagStr)
			if err != nil {
				return errors.Wrap(err, "parseTag")
			}

			text := v.Field(i).String()
			cb(text, tag.move(x, y))

		case reflect.Array, reflect.Slice:
			rx, ry, err := parseRelTag(tagStr)
			if err != nil {
				return err
			}

			for j := 0; j < v.Field(i).Len(); j++ {
				fj := float64(j)
				// 座標を変えつつ再帰呼び出し
				if err := parseVarsInternal(v.Field(i).Index(j), cb, x+rx*fj, y+ry*fj); err != nil {
					return err
				}
			}

		case reflect.Struct:
			rx, ry, err := parseRelTag(tagStr)
			if err != nil {
				return err
			}

			// 座標を変えつつ再帰呼び出し
			if err := parseVarsInternal(v.Field(i), cb, x+rx, y+ry); err != nil {
				return err
			}

		default:
			return fmt.Errorf("vars の %vフィールドが未対応の型 (%v)です", f.Name, f.Type)
		}
	}

	return nil
}
