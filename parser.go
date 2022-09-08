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

type stackEntry struct {
	value   reflect.Value
	originX float64
	originY float64
}

// パラメータ構造体をパースしてtaggedTextのスライスを返します
func parseVars(vars interface{}) ([]taggedText, error) {
	texts := []taggedText{}

	stack := []stackEntry{{value: reflect.ValueOf(vars)}}
	for len(stack) != 0 {
		entry := stack[0]
		stack = stack[1:]

		v := entry.value
		t := entry.value.Type()
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

				tt := taggedText{textTag: tag.fromOrigin(entry.originX, entry.originY), Text: v.Field(i).String()}
				texts = append(texts, tt)
			case reflect.Array, reflect.Slice:
				tag := &transTag{}
				if err := tag.parse(tagStr); err != nil {
					return nil, errors.Wrap(err, "parseTag")
				}

				for j := 0; j < v.Field(i).Len(); j++ {
					fj := float64(j)
					stack = append(stack, stackEntry{v.Field(i).Index(j), entry.originX + tag.Dx*fj, entry.originY + tag.Dy*fj})
				}

			case reflect.Struct:
				tag := &transTag{}
				if err := tag.parse(tagStr); err != nil {
					return nil, errors.Wrap(err, "parseTag")
				}

				stack = append(stack, stackEntry{v.Field(i), entry.originX + tag.Dx, entry.originY + tag.Dy})
			default:
				return nil, fmt.Errorf("vars の %vフィールドが未対応の型 (%v)です", f.Name, f.Type)
			}
		}

	}

	return texts, nil
}
