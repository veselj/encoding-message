package message

import (
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// Unmarshal will read data and populate a struct referenced by v
// using struct tags (metadata)
func Unmarshal(data []byte, v interface{}) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}

// decodeState stores the state of the unmarshalling
type decodeState struct {
	data         string
	offset       int // read offset in data
	errorContext struct {
		Struct string
		Field  string
	}
}

// init initilises decodeState structure.
func (d *decodeState) init(data []byte) *decodeState {
	d.data = string(data)
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
	return d.decodeStruct(rv)
}

// getChunk gets the next chunk of data by length or by separator sep
// length can also restrict how far separator is looked
func (d *decodeState) getChunk(length int, sep rune) (string, error) {
	var rv string
	remaining := len(d.data) - d.offset
	if sep == 0 {
		// by length
		if remaining >= length {
			rv = d.data[d.offset:(d.offset + length)]
			d.offset += length
		} else {
			return "", errors.New("unsufficient remaining length in data")
		}
	} else {
		// by separator
		var found bool
		remData := d.data[d.offset:]
		for i, c := range remData {
			if c == sep {
				rv = d.data[d.offset : d.offset+i]
				d.offset++
				found = true
				break
			}
			if length > 0 && i == length-1 {
				rv = d.data[d.offset : d.offset+length]
				found = true
				break
			}
		}
		if !found {
			rv = remData
		}
		d.offset += len(rv)
	}
	return rv, nil
}

func (d *decodeState) decodeStruct(rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		d.decodeField(rt.Field(i), rv.Field(i))
	}
	return nil
}

// decodeField processes one field of a struct and stores it
// in the reflect.Value param.
func (d *decodeState) decodeField(st reflect.StructField, v reflect.Value) error {
	var fl int
	var err error
	fSepRaw := st.Tag.Get("sep")
	fSep, _ := utf8.DecodeRuneInString(fSepRaw)
	if fSep == utf8.RuneError {
		fSep = 0
	}
	fLen := st.Tag.Get("len")
	if len(fLen) > 0 {
		fl, err = strconv.Atoi(fLen)
		if err != nil {
			return errors.New("Unable to parse len tag of field")
		}
	}
	chunk, err := d.getChunk(fl, fSep)
	if err != nil {
		return errors.Wrap(err, "Unable to parse chunk of data")
	}
	switch st.Type.Kind() {
	case reflect.Int:
		iv, err := strconv.ParseInt(chunk, 10, 64)
		if err != nil {
			return errors.New("unable to parse int from " + chunk)
		}
		v.SetInt(iv)
	case reflect.String:
		v.SetString(chunk)

	case reflect.Struct:
		d.decodeStruct(v)

	case reflect.Slice:
		elemType := st.Type.Elem()
		d.decodeSlice(elemType, v)
	}
	return nil
}

func (d *decodeState) decodeSlice(t reflect.Type, v reflect.Value) {
	return
}
