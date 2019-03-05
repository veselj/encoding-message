package message

// Usage:
// 1. Define a struct for your message format. Create golang fields
// and add metadata using struct tags.
// 2. Use the Unmarshal method to parse a raw message into the defined struct
// by passing a pointer to that struct.

import (
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
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
	data         string
	offset       int // read offset in data
	errorContext struct {
		Struct string
		Field  string
	}
}

type fieldTags struct {
	len     int
	sep     rune
	padding string
}

func (t *fieldTags) decode(tag reflect.StructTag) error {
	var err error
	fSepRaw := tag.Get("sep")
	t.sep, _ = utf8.DecodeRuneInString(fSepRaw)
	if t.sep == utf8.RuneError {
		t.sep = 0
	}
	fLen := tag.Get("len")
	if len(fLen) > 0 {
		t.len, err = strconv.Atoi(fLen)
		if err != nil {
			return errors.New("Unable to parse len tag of field")
		}
	}
	return nil
}

// init initilises decodeState structure.
func (d *decodeState) init(data []byte) *decodeState {
	d.data = string(data)
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
	return d.decodeStruct(rv)
}

// getChunk gets the next chunk of data by length or by separator sep
// length can also restrict how far separator is looked
func (d *decodeState) getChunk(tags fieldTags) (string, error) {
	var rv string
	remaining := len(d.data) - d.offset
	if tags.sep == 0 {
		// by length
		if remaining >= tags.len {
			rv = d.data[d.offset:(d.offset + tags.len)]
			d.offset += tags.len
		} else {
			return "", errors.New("unsufficient remaining length in data")
		}
	} else {
		// by separator
		var found bool
		remData := d.data[d.offset:]
		for i, c := range remData {
			if c == tags.sep {
				rv = d.data[d.offset : d.offset+i]
				d.offset++
				found = true
				break
			}
			if tags.len > 0 && i == tags.len-1 {
				rv = d.data[d.offset : d.offset+tags.len]
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
	// read fieldTags
	var tags fieldTags
	err := tags.decode(st.Tag)
	if err != nil {
		return errors.Wrap(err, "unable to parse tags")
	}
	switch st.Type.Kind() {
	case reflect.Int:
		chunk, err := d.getChunk(tags)
		if err != nil {
			return errors.Wrap(err, "Unable to parse chunk of data")
		}
		iv, err := strconv.ParseInt(chunk, 10, 64)
		if err != nil {
			return errors.New("unable to parse int from " + chunk)
		}
		v.SetInt(iv)
	case reflect.String:
		chunk, err := d.getChunk(tags)
		if err != nil {
			return errors.Wrap(err, "Unable to parse chunk of data")
		}
		v.SetString(chunk)

	case reflect.Struct:
		d.decodeStruct(v)

	case reflect.Slice:
		elemType := st.Type.Elem()
		d.decodeSlice(tags, elemType, v)
	}
	return nil
}

// decodeSlice parses a slice of values
func (d *decodeState) decodeSlice(tags fieldTags, t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.String:
		result := make([]string, 0)
		for {
			chunk, err := d.getChunk(tags)
			if err != nil {
				return errors.Wrap(err, "unable to get chunk in decodeSlice")
			}
			if len(chunk) == 0 {
				break
			}
			result = append(result, chunk)
		}
		slice := reflect.MakeSlice(reflect.TypeOf(result), len(result), len(result))
		v.Set(slice)
		v.SetLen(len(result))
		for i, s := range result {
			v.Index(i).SetString(s)
		}

	case reflect.Int:
		result := make([]int, 0)
		for {
			chunk, err := d.getChunk(tags)
			if err != nil {
				return errors.Wrap(err, "unable to get chunk in decodeSlice")
			}
			if len(chunk) == 0 {
				break
			}
			iv, err := strconv.ParseInt(chunk, 10, 64)
			if err != nil {
				return errors.New("unable to parse int from " + chunk)
			}
			result = append(result, int(iv))
		}
		slice := reflect.MakeSlice(reflect.TypeOf(result), len(result), len(result))
		v.Set(slice)
		v.SetLen(len(result))
		for i, intV := range result {
			v.Index(i).SetInt(int64(intV))
		}
	}

	return nil
}
