package pdftpl

import (
	"fmt"
	"image"
	"reflect"
)

type stackEntry struct {
	value   reflect.Value
	originX float64
	originY float64
}

// パラメータ構造体をパースしてelementのスライスを返します
func getElementsFromStyledVars(styledVars any) ([]element, error) {
	elements := []element{}

	stack := []stackEntry{{value: reflect.ValueOf(styledVars)}}
	for len(stack) != 0 {
		entry := stack[0]
		stack = stack[1:]

		v := entry.value
		t := entry.value.Type()
		if t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("styledVars must be Struct")
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
					return nil, fmt.Errorf("parseTag: %v, %w", tagStr, err)
				}

				element := &textElement{textTag: tag.fromOrigin(entry.originX, entry.originY), text: v.Field(i).String()}
				elements = append(elements, element)
			case reflect.Array, reflect.Slice:
				tag := &transTag{}
				if err := tag.parse(tagStr); err != nil {
					return nil, fmt.Errorf("parseTag: %v, %w", tagStr, err)
				}

				for j := 0; j < v.Field(i).Len(); j++ {
					fj := float64(j)
					stack = append(stack, stackEntry{v.Field(i).Index(j), entry.originX + tag.Dx*fj, entry.originY + tag.Dy*fj})
				}

			case reflect.Struct:
				tag := &transTag{}
				if err := tag.parse(tagStr); err != nil {
					return nil, fmt.Errorf("parseTag: %v, %w", tagStr, err)
				}

				stack = append(stack, stackEntry{v.Field(i), entry.originX + tag.Dx, entry.originY + tag.Dy})

			default:
				if f.Type == reflect.TypeOf((*image.Image)(nil)).Elem() {
					tag := &imageTag{}
					if err := tag.parse(tagStr); err != nil {
						return nil, fmt.Errorf("parseTag: %v, %w", tagStr, err)
					}

					element := &imageElement{imageTag: tag.fromOrigin(entry.originX, entry.originY)}
					if img, ok := v.Field(i).Interface().(image.Image); ok {
						element.image = img
					}
					elements = append(elements, element)
				} else {
					return nil, fmt.Errorf("vars の %vフィールドが未対応の型 (%v)です", f.Name, f.Type)
				}
			}
		}

	}

	return elements, nil
}
