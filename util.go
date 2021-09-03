package luno

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// makeURLValues converts a request struct into a url.Values map.
func makeURLValues(v interface{}) url.Values {
	values := make(url.Values)

	valElem := reflect.ValueOf(v).Elem()
	typElem := reflect.TypeOf(v).Elem()

	for i := 0; i < typElem.NumField(); i++ {
		field := typElem.Field(i)
		urlTag := field.Tag.Get("url")
		if urlTag == "" || urlTag == "-" {
			continue
		}

		fieldValue := valElem.Field(i)

		stringer, ok := fieldValue.Interface().(QueryValuer)
		if ok {
			values.Set(urlTag, stringer.QueryValue())
			continue
		}

		k := fieldValue.Kind()
		var s string
		ss := make([]string, 0)
		switch k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:
			s = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			s = strconv.FormatUint(fieldValue.Uint(), 10)
		case reflect.Float32:
			s = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 32)
		case reflect.Float64:
			s = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 64)
		case reflect.Slice:
			 if field.Type.Elem().Kind() == reflect.String {
				for i := 0; i < fieldValue.Len(); i++ {
					ss = append(ss, fieldValue.Index(i).String())
				}
			}
		case reflect.String:
			s = fieldValue.String()
		case reflect.Bool:
			s = fmt.Sprintf("%v", fieldValue.Bool())
		}
		if len(ss) > 0 {
			for _, str := range ss {
				values.Add(urlTag, str)
			}
		} else {
			values.Set(urlTag, s)
		}
	}

	return values
}

type QueryValuer interface {
	QueryValue() string
}
