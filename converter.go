package config

import (
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// convertAndSetSlice converts a slice of strings to a value of the type that the slice holds,
// and appends it to the slice. It returns a slice of indices that failed to convert.
// Supported types:
//   - int, uint, float variants
//   - bool, string
//   - time.Duration
//   - *url.URL
//
// Parameters:
//   - slicePtr - A reflect.Value that points to the slice that will be set.
//   - values - A slice of strings that will be converted and set on the slice.
//
// Returns:
// A slice of indices that failed to convert.
func convertAndSetSlice(slicePtr reflect.Value, values []string) []int {
	sliceVal := slicePtr.Elem()
	elemType := sliceVal.Type().Elem()

	var failedIndices []int

	for i, s := range values {
		elemPtr := reflect.New(elemType)

		if success := convertAndSetValue(elemPtr, s); !success {
			failedIndices = append(failedIndices, i)
		} else {
			sliceVal.Set(reflect.Append(sliceVal, elemPtr.Elem()))
		}
	}

	// If the slice is nil, create a new slice.
	if sliceVal.IsNil() {
		sliceVal.Set(reflect.MakeSlice(sliceVal.Type(), 0, 0))
	}

	return failedIndices
}

// convertAndSetValue converts a string to a value of the type that the reflect.Value holds,
// and sets it on the reflect.Value. It returns true if the conversion was successful.
// Supported types:
//   - int, uint, float variants
//   - bool, string
//   - time.Duration
//   - *url.URL
//
// Parameters:
//   - settable - A reflect.Value that will be set with the converted value.
//   - str - The string that will be converted and set on the settable.
//
// Returns:
// A boolean indicating if the conversion was successful.
func convertAndSetValue(settable reflect.Value, str string) bool {
	var settableValue reflect.Value
	if settable.Kind() == reflect.Ptr || settable.Kind() == reflect.Interface {
		settableValue = settable.Elem()
	} else {
		settableValue = settable
	}

	switch settableValue.Kind() {
	case reflect.Pointer:
		// Only URL is supported as a pointer type.
		return convertAndSetURL(settableValue, str)
	case reflect.String:
		return convertAndSetString(settableValue, str)
	case reflect.Bool:
		return convertAndSetBool(settableValue, str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convertAndSetInt(settableValue, str)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convertAndSetUint(settableValue, str)
	case reflect.Float32, reflect.Float64:
		return convertAndSetFloat(settableValue, str)
	default:
		return false
	}
}

// The following functions are helper functions used by convertAndSetValue to convert and set specific types.

func convertAndSetURL(settableValue reflect.Value, str string) bool {
	urlVal, err := url.Parse(str)

	if err == nil {
		settableValue.Set(reflect.ValueOf(urlVal))
	}

	return err == nil
}

func convertAndSetString(settableValue reflect.Value, str string) bool {
	settableValue.SetString(str)

	return true
}

func convertAndSetBool(settableValue reflect.Value, str string) bool {
	boolVal, err := strconv.ParseBool(str)

	if err == nil {
		settableValue.SetBool(boolVal)
	}

	return err == nil
}

func convertAndSetInt(settableValue reflect.Value, str string) bool {
	if settableValue.Type().PkgPath() == "time" && settableValue.Type().Name() == "Duration" {
		return convertAndSetDuration(settableValue, str)
	}

	intVal, err := strconv.ParseInt(str, 10, settableValue.Type().Bits())

	if err == nil {
		settableValue.SetInt(intVal)
	}

	return err == nil
}

func convertAndSetDuration(settableValue reflect.Value, str string) bool {
	d, err := time.ParseDuration(str)

	if err == nil {
		settableValue.SetInt(int64(d))
	}

	return err == nil
}

func convertAndSetUint(settableValue reflect.Value, str string) bool {
	uintVal, err := strconv.ParseUint(str, 10, settableValue.Type().Bits())

	if err == nil {
		settableValue.SetUint(uintVal)
	}

	return err == nil
}

func convertAndSetFloat(settableValue reflect.Value, str string) bool {
	floatVal, err := strconv.ParseFloat(str, settableValue.Type().Bits())

	if err == nil {
		settableValue.SetFloat(floatVal)
	}

	return err == nil
}
