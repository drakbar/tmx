package tmx

import (
	"errors"
	"reflect"
)

var (
	// data quantity errors
	missingData      = errors.New("base64 data is an empty string")
	dataSizeMismatch = errors.New("tile data and map size do not match")
)

var (
	// data loading errors
	templateNotLoaded = errors.New("the template wasn't loaded")
	nilDataPtr        = errors.New("data pointer is nil")
)

var (
	// data formatting errors
	unsupportedEncoding    = errors.New("the encoding type is unsupported")
	unsupportedCompression = errors.New("the compression type is unsupported")
)

var (
	// data type errors
	dataStringMismatch  = errors.New("the data is not of type string")
	csvDataMismatch     = errors.New("csv data structure incorrect")
	highBitDataMismatch = errors.New("tile data is not a byte array")
)

var (
	// reflection errors
	reflectionBothWrong = errors.New("both the src and dst are not a structures")
	reflectionSrcWrong  = errors.New("the src is not a structures")
	reflectionDstWrong  = errors.New("the dst is not a structures")
)

var (
	// invalid data errors
	badGlobalId       = errors.New("global id could not be found in any tileset")
	noMatchingTileset = errors.New("template does not match a valid tileset")
)

// copyFields copies the fields of one structure over to another. It does not
// copy slices or structure however.
func copyFields(src, dst *reflect.Value) (e error) {
	// verify that both are of type struct
	if e = checkStruct(*src, *dst); e != nil {
		return
	}
	for i := 0; i < src.NumField(); i++ {
		// get the src field name and value
		n, v := src.Type().Field(i).Name, src.Field(i)
		// get the dst field
		f := dst.FieldByName(n)
		// if the field exists and it can be assigned a value
		if f.IsValid() && f.CanSet() {
			// assign the field a value based on its type
			switch v.Type().Kind() {
			case reflect.String:
				f.SetString(v.String())
			case reflect.Int:
				// nothing has a gid of zero
				if n == "Gid" && v.Int() == 0 {
					continue
				}
				f.SetInt(v.Int())
			case reflect.Float64:
				f.SetFloat(v.Float())
			case reflect.Bool:
				f.SetBool(v.Bool())
			default:
				continue
			}
		}
	}
	return
}

// checkStruct verifies if both the source an destination reflect values are of
// the type struct.
func checkStruct(src, dst reflect.Value) (e error) {
	if src.Kind() != reflect.Struct && dst.Kind() != reflect.Struct {
		// src and dst values are not structs
		return reflectionBothWrong
	}
	if src.Kind() != reflect.Struct {
		// src value is not a struct
		return reflectionSrcWrong
	}
	if dst.Kind() != reflect.Struct {
		// dst value is not a struct
		return reflectionDstWrong
	}
	return
}
