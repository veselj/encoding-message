package message

// Usage:
// 1. Define a struct for your message format. Create golang fields
// and add metadata using struct tags.
// 2. Use the Unmarshal method to parse a raw message into the defined struct
// by passing a pointer to that struct.

import (
	"errors"
	"reflect"
	"strconv"
)

// Unmarshal will read data and parse it in order to
// populate the struct referenced by v using struct tags (metadata).
// v should be a pointer to the struct.
func Unmarshal(data []byte, v interface{}) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}

// decodeState keeps internal state of the decoder.
type decodeState struct {
	data         []byte
	offset       int // read offset in data
	errorContext struct {
		Struct string
		Field  string
	}
}

// init initializes decodeState
func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.offset = 0
	d.errorContext.Struct = ""
	d.errorContext.Field = ""
	return d
}

// unmarshal is the internal unmarshal implementation over decodeState
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

// value receives the struct as reflect.Value, iterates
// through its fields and populates them using the defined metadata.
func (d *decodeState) value(rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		rft := rt.Field(i)
		fType := rft.Tag.Get("type")
		fLen := rft.Tag.Get("len")
		fl, err := strconv.Atoi(fLen)
		if err != nil {
			return errors.New("Unable to parse len tag of field")
		}
		//fPadding := rft.Tag.Get("padding")
		return d.valueByLen(rv, i, fType, fl)
	}
	return nil
}

func (d *decodeState) valueByLen(rv reflect.Value, index int, fType string, fLen int) error {
	switch fType {
	case "int":
		chunk := d.data[d.offset:(d.offset + fLen)]
		iv, err := strconv.Atoi(string(chunk))
		if err != nil {
			return errors.New("unable to parse int from " + string(chunk))
		}
		rv.Field(index).SetInt(int64(iv))
	}
	return nil
}
