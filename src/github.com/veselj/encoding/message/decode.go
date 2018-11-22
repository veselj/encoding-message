package message

import (
	"errors"
	"reflect"
	"strconv"
)

// Unmarshal will read data and populate a struct referenced by v
// using struct tags (metadata)
func Unmarshal(data []byte, v interface{}) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}

type decodeState struct {
	data         []byte
	offset       int // read offset in data
	errorContext struct {
		Struct string
		Field  string
	}
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.offset = 0
	d.errorContext.Struct = ""
	d.errorContext.Field = ""
	return d
}

func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("input not a pointer or nil")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("input must point to a struct")
	}
	return d.value(rv)
}

func (d *decodeState) value(rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		rft := rt.Field(i)
		fType := rft.Tag.Get("type")
		fLen := rft.Tag.Get("len")
		//fPadding := rft.Tag.Get("padding")
		switch fType {
		case "int":
			fl, err := strconv.Atoi(fLen)
			if err != nil {
				return errors.New("Unable to parse len tag of field")
			}
			chunk := d.data[d.offset:(d.offset + fl)]
			iv, err := strconv.ParseInt(string(chunk), 10, 64)
			if err != nil {
				return errors.New("unable to parse int from " + string(chunk))
			}
			rv.Field(i).SetInt(iv)
		}
	}
	return nil
}
