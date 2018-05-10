package luno

import (
	"net/url"
	"reflect"
	"strconv"
)

// makeURLValues converts a request struct into a url.Values map.
func makeURLValues(v interface{}) (url.Values, error) {
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

		var s string
		switch fieldValue.Interface().(type) {
		case int, int8, int16, int32, int64:
			s = strconv.FormatInt(fieldValue.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			s = strconv.FormatUint(fieldValue.Uint(), 10)
		case float32:
			s = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 32)
		case float64:
			s = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 64)
		case []byte:
			s = string(fieldValue.Bytes())
		case string:
			s = fieldValue.String()
		}
		values.Set(urlTag, s)
	}

	return values, nil
}
