package config

import (
	"reflect"
	"strings"
)

const (
	keyValueDelimiter = "="
	keyValueNumParts  = 2
)

// mergeMaps merges the source map into the destination map.
// If a key exists in both maps, the value in the source map will overwrite the value in the destination map.
func mergeMaps(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

// keyValsToMap converts a slice of strings to a map.
// Each string must be in the format "key=value".
// If a string does not contain "=", it will be ignored.
// The key will be converted to lowercase.
func keyValsToMap(ss []string) map[string]string {
	retMap := make(map[string]string)

	for _, s := range ss {
		if !strings.Contains(s, keyValueDelimiter) {
			continue
		}

		parts := strings.SplitN(s, keyValueDelimiter, keyValueNumParts)

		if len(parts) == keyValueNumParts {
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			val := strings.TrimSpace(parts[1])
			retMap[key] = val
		}
	}

	return retMap
}

// mapKeysToFields recursively maps keys to fields in a struct.
//
// Params:
//   - structPtr: A pointer to the struct to map keys to.
//   - valMap: A map of keys to values.
//   - prefix: The prefix to prepend to the keys.
//   - structDelimiter: The delimiter to use when joining the prefix and field names.
//
// Example:
//
//	type Config struct {
//	  Host string `config:"server_host"`
//	}
//
//	config := Config{}
//	structPtr := reflect.ValueOf(&config)
//	valMap := make(map[string]reflect.Value)
//	mapKeysToFields(structPtr, valMap, "app_", "_")
//
//	fmt.Println(valMap) // Output: map[app_server_host:<value>]
func mapKeysToFields(structPtr reflect.Value, valMap map[string]reflect.Value, prefix string, structDelimiter string) {
	structVal := structPtr.Elem()

	for i := range structVal.NumField() {
		field := structVal.Type().Field(i)
		fieldPtr := structVal.Field(i).Addr()

		key := getKey(field, prefix)

		switch field.Type.Kind() {
		case reflect.Struct:
			mapKeysToFields(fieldPtr, valMap, key+structDelimiter, structDelimiter)
		case reflect.Pointer, reflect.Interface:
			valMap[key] = fieldPtr
		default:
			valMap[key] = fieldPtr.Elem()
		}
	}
}

// getKey returns the key for a field, based on its tag or name.
// If a tag is present, it will be used as the key.
// Otherwise, the field name will be used.
//
// Params:
//   - field: The field to get the key for.
//   - prefix: The prefix to prepend to the key.
//
// Result: The key for the field.
//
// Example:
//
//	type Config struct {
//	  Host string `config:"server_host"`
//	}
//
//	field := reflect.TypeOf(Config{}).Field(0)
//	key := getKey(field, "app_")
//	fmt.Println(key) // Output: app_server_host
func getKey(field reflect.StructField, prefix string) string {
	name := field.Name

	tag, exists := field.Tag.Lookup("config")

	if exists {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			name = tag
		}
	}

	return prefix + name
}

// stringToSlice converts a string to a slice using the given delimiter.
//
// Params:
//   - str: The string to convert.
//   - delimiter: The delimiter to split the string by.
//
// Result: The resulting slice will have all elements trimmed of whitespace.
//
// Example:
//
//	stringToSlice("a, b, c, ", ",") // Output: []string{"a", "b", "c"}
func stringToSlice(str string, delimiter string) []string {
	splits := strings.Split(str, delimiter)
	filtered := splits[:0] // zero-length slice of the same underlying array

	for _, split := range splits {
		split = strings.TrimSpace(split)
		if split != "" {
			filtered = append(filtered, split)
		}
	}

	return filtered
}
