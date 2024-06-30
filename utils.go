package config

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func mergeMaps(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

// stringsToMap converts a slice of strings to a map.
// Each string must be in the format "key=value".
// If a string does not contain "=", it will be ignored.
// The key will be converted to lowercase.
func stringsToMap(ss []string) map[string]string {
	retMap := make(map[string]string)

	numParts := 2

	for _, s := range ss {
		if !strings.Contains(s, "=") {
			continue
		}

		parts := strings.SplitN(s, "=", numParts)

		key, value := strings.ToLower(parts[0]), parts[1]

		if key != "" && value != "" {
			retMap[key] = value
		}
	}

	return retMap
}

// mapKeysToFields recursively maps keys to fields in a struct.
//
// @param ptr: pointer to a struct.
func mapKeysToFields(structPtr reflect.Value, valMap map[string]reflect.Value, prefix string, structDelimiter string) {
	structVal := structPtr.Elem()

	for i := range structPtr.NumField() {
		field := structPtr.Type().Field(i)
		fieldPtr := structVal.Field(i).Addr()

		key := getKey(field, prefix)

		if !strings.Contains(prefix, structDelimiter) {
			prefix += structDelimiter
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			mapKeysToFields(fieldPtr, valMap, prefix, structDelimiter)
		default:
			valMap[key] = fieldPtr
		}
	}
}

// getKey returns the key for a field, based on its tag or name.
// If a tag is present, it will be used as the key.
// Otherwise, the field name will be used.
func getKey(field reflect.StructField, prefix string) string {
	name := field.Name

	if tag, exists := field.Tag.Lookup("config"); exists {
		if tag = strings.TrimSpace(tag); tag != "" {
			name = tag
		}
	}

	return strings.ToLower(prefix + name)
}

// stringToSlice converts a string to a slice.
// The string will be split by the delimiter.
// Each element will be trimmed of whitespace.
// Empty elements will be ignored.
// If the delimiter is empty, it will default to c.sliceDelimiter.
func stringToSlice(str string, delimiter string) []string {
	if delimiter == "" {
		panic("delimiter cannot be empty")
	}

	str = strings.TrimSpace(str)

	if str == "" {
		return nil
	}

	splits := strings.Split(str, delimiter)
	filtered := splits[:0] // zero-length slice of the same underlying array

	for _, split := range splits {
		if split = strings.TrimSpace(split); split != "" {
			filtered = append(filtered, split)
		}
	}

	return filtered
}

// convertAndSetSlice converts a slice of strings to a slice of values, and sets it on a settable.
// It returns a slice of indices that failed to convert.
// Supported types:
//   - int, uint, float variants
//   - bool, string
//   - time.Duration
//   - *url.URL
func convertAndSetSlice(slicePtr reflect.Value, values []string) []int {
	sliceVal := slicePtr.Elem()
	elemType := sliceVal.Type().Elem()

	var failedIndices []int

	for i, s := range values {
		elemPtr := reflect.New(elemType)
		if !convertAndSetValue(elemPtr, s) {
			failedIndices = append(failedIndices, i)
		} else {
			sliceVal.Set(reflect.Append(sliceVal, elemPtr.Elem()))
		}
	}

	return failedIndices
}

// convertAndSetValue converts a string to a value, and sets it on a settable.
// It returns true if the conversion was successful.
// Supported types:
//   - int, uint, float variants
//   - bool, string
//   - time.Duration
//   - *url.URL
func convertAndSetValue(settable reflect.Value, str string) bool {
	settableValue := settable.Elem()

	switch settableValue.Kind() {
	case reflect.Pointer:
		return convertAndSetPointer(settableValue, str)
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

func convertAndSetPointer(settableValue reflect.Value, str string) bool {
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
