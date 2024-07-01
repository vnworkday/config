package config

import (
	"fmt"
	"reflect"
	"testing"
)

func failTestMessage(funcName string, expected, got any) string {
	return fmt.Sprintf(`%s()
expected:	%v
actual:		%v`, funcName, expected, got)
}

func TestMergeMaps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dst  map[string]string
		src  map[string]string
		want map[string]string
	}{
		{
			name: "When source map has new keys then they should be added to the destination map",
			dst:  map[string]string{"key1": "value1", "key2": "value2"},
			src:  map[string]string{"key2": "new value2", "key3": "value3"},
			want: map[string]string{"key1": "value1", "key2": "new value2", "key3": "value3"},
		},
		{
			name: "When source map is empty then the destination map should remain the same",
			dst:  map[string]string{"key1": "value1", "key2": "value2"},
			src:  map[string]string{},
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "When destination map is empty then the destination map should be the source map",
			dst:  map[string]string{},
			src:  map[string]string{"key1": "value1", "key2": "value2"},
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "When both maps are empty then the destination map should be empty",
			dst:  map[string]string{},
			src:  map[string]string{},
			want: map[string]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mergeMaps(test.dst, test.src)

			if !reflect.DeepEqual(test.dst, test.want) {
				t.Errorf(failTestMessage("mergeMaps", test.want, test.dst))
			}
		})
	}
}

func TestKeyValsToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ss   []string
		want map[string]string
	}{
		{
			name: "When slice contains valid key-value pairs then they should be added to the map",
			ss:   []string{"key1=value1", "key2=value2"},
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "When slice is empty then the map should be empty",
			ss:   []string{},
			want: map[string]string{},
		},
		{
			name: "When slice contains strings without '=' then they should be ignored",
			ss:   []string{"key1", "key2=value2"},
			want: map[string]string{"key2": "value2"},
		},
		{
			name: "When slice contains strings with multiple '=' then only the first '=' should be considered",
			ss:   []string{"key1=value1=extra", "key2=value2"},
			want: map[string]string{"key1": "value1=extra", "key2": "value2"},
		},
		{
			name: "When slice contains strings with leading/trailing spaces then they should be trimmed",
			ss:   []string{" key1=value1", "key2=value2   "},
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "When slice contains strings with different cases then the keys should be converted to lowercase",
			ss:   []string{"Key1=value1", "KEY2=value2"},
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := keyValsToMap(test.ss); !reflect.DeepEqual(got, test.want) {
				t.Errorf(failTestMessage("keyValsToMap", test.want, got))
			}
		})
	}
}

type TestStruct struct {
	Field1       string `config:"field_1"`
	Field2       int
	NestedStruct struct {
		Field3 bool `config:"field_3"`
	}
}

func TestMapKeysToFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		structPtr any
		want      map[string]reflect.Value
	}{
		{
			name: "When struct has fields then they should be added to the map",
			structPtr: &TestStruct{
				Field1: "value1",
				Field2: 2,
				NestedStruct: struct {
					Field3 bool `config:"field_3"`
				}(struct{ Field3 bool }{Field3: true}),
			},
			want: map[string]reflect.Value{
				"app_field_1":              reflect.ValueOf("value1"), // Field1 is tagged with config
				"app_Field2":               reflect.ValueOf(2),        // Field2 is not tagged with config
				"app_NestedStruct_field_3": reflect.ValueOf(true),     // NestedStruct.Field3 is tagged with config
			},
		},
		{
			name:      "When struct has no fields then the map should be empty",
			structPtr: &struct{}{},
			want:      map[string]reflect.Value{},
		},
		{
			name: "When struct has nested structs then their fields should be added to the map",
			structPtr: &struct {
				NestedStruct struct {
					Field1 string `config:"field_1"`
				}
			}{NestedStruct: struct {
				Field1 string `config:"field_1"`
			}{Field1: "value1"}},
			want: map[string]reflect.Value{"app_NestedStruct_field_1": reflect.ValueOf("value1")},
		},
		{
			name: "When struct has nested structs with no fields then the map should be empty",
			structPtr: &struct {
				NestedStruct struct{}
			}{NestedStruct: struct{}{}},
			want: map[string]reflect.Value{},
		},
		{
			name: "When struct has nested structs with nested structs then their fields should be added to the map",
			structPtr: &struct {
				NestedStruct struct {
					NestedStruct struct {
						Field1 string `config:"field_1"`
					}
				}
			}{NestedStruct: struct {
				NestedStruct struct {
					Field1 string `config:"field_1"`
				}
			}{NestedStruct: struct {
				Field1 string `config:"field_1"`
			}{Field1: "value1"}}},
			want: map[string]reflect.Value{"app_NestedStruct_NestedStruct_field_1": reflect.ValueOf("value1")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			valMap := make(map[string]reflect.Value)

			mapKeysToFields(reflect.ValueOf(test.structPtr), valMap, "app_", "_")

			for key, val := range valMap {
				if !reflect.DeepEqual(val.Interface(), test.want[key].Interface()) {
					t.Errorf(failTestMessage("mapKeysToFields", test.want[key], val))
				}
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		field  reflect.StructField
		prefix string
		want   string
	}{
		{
			name:   "When field has a tag then it should be used as the key",
			field:  reflect.StructField{Name: "Field1", Tag: `config:"tag1"`},
			prefix: "app_",
			want:   "app_tag1",
		},
		{
			name:   "When field does not have a tag then the field name should be used as the key",
			field:  reflect.StructField{Name: "Field1"},
			prefix: "app_",
			want:   "app_Field1",
		},
		{
			name:   "When field has an empty tag then the field name should be used as the key",
			field:  reflect.StructField{Name: "Field1", Tag: `config:""`},
			prefix: "app_",
			want:   "app_Field1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := getKey(test.field, test.prefix); got != test.want {
				t.Errorf(failTestMessage("getKey", test.want, got))
			}
		})
	}
}

func TestStringToSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		str       string
		delimiter string
		want      []string
	}{
		{
			name:      "When string contains multiple elements then they should be split into a slice",
			str:       "a, b, c",
			delimiter: ",",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "When string is empty then the slice should be empty",
			str:       "",
			delimiter: ",",
			want:      []string{},
		},
		{
			name:      "When string contains leading/trailing spaces then they should be trimmed",
			str:       " a , b , c ",
			delimiter: ",",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "When string contains multiple delimiters then they should be ignored",
			str:       "a,,b,c",
			delimiter: ",",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "When string contains only delimiters then the slice should be empty",
			str:       ",,",
			delimiter: ",",
			want:      []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := stringToSlice(test.str, test.delimiter); !reflect.DeepEqual(got, test.want) {
				t.Errorf(failTestMessage("stringToSlice", test.want, got))
			}
		})
	}
}
