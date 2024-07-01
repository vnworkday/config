package config

import (
	"net/url"
	"reflect"
	"slices"
	"testing"
	"time"
)

func toURL(s string) *url.URL {
	u, _ := url.Parse(s)

	return u
}

func TestConvertAndSetSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		slicePtr     reflect.Value
		values       []string
		want         any
		wantFailures []int
	}{
		{
			name:         "WhenValuesAreURLs",
			slicePtr:     reflect.ValueOf(new([]*url.URL)),
			values:       []string{"https://example.com", "://malformed", "https://example.org"},
			want:         []*url.URL{toURL("https://example.com"), toURL("https://example.org")},
			wantFailures: []int{1},
		},
		{
			name:         "WhenValuesAreStrings",
			slicePtr:     reflect.ValueOf(new([]string)),
			values:       []string{"test1", "test2", "test3"},
			want:         []string{"test1", "test2", "test3"},
			wantFailures: []int{},
		},
		{
			name:         "WhenValuesAreBools",
			slicePtr:     reflect.ValueOf(new([]bool)),
			values:       []string{"true", "false", "notabool"},
			want:         []bool{true, false},
			wantFailures: []int{2},
		},
		{
			name:         "WhenValuesAreInts",
			slicePtr:     reflect.ValueOf(new([]int)),
			values:       []string{"123", "456", "notanint"},
			want:         []int{123, 456},
			wantFailures: []int{2},
		},
		{
			name:         "WhenValuesAreUnsupportedType",
			slicePtr:     reflect.ValueOf(new([]complex128)),
			values:       []string{"1+2i", "3+4i"},
			want:         []complex128{},
			wantFailures: []int{0, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotFailures := convertAndSetSlice(test.slicePtr, test.values)

			if !slices.Equal(test.wantFailures, gotFailures) {
				t.Errorf("convertAndSetSlice() = %v, want %v", gotFailures, test.wantFailures)
			}

			if !reflect.DeepEqual(test.want, test.slicePtr.Elem().Interface()) {
				t.Errorf("convertAndSetSlice() = %v, want %v", test.slicePtr.Elem().Interface(), test.want)
			}
		})
	}
}

func TestConvertAndSetValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		settable reflect.Value
		str      string // value to convert
		want     any    // expected value after conversion
		wantOk   bool   // true if it is expected that the conversion was successful
	}{
		{
			name:     "When value is URL",
			settable: reflect.ValueOf(new(*url.URL)),
			str:      "https://example.com",
			want:     toURL("https://example.com"),
			wantOk:   true,
		},
		{
			name:     "When value is string",
			settable: reflect.ValueOf(new(string)),
			str:      "test",
			want:     "test",
			wantOk:   true,
		},
		{
			name:     "When value is bool",
			settable: reflect.ValueOf(new(bool)),
			str:      "true",
			want:     true,
			wantOk:   true,
		},
		{
			name:     "When value is int",
			settable: reflect.ValueOf(new(int)),
			str:      "123",
			want:     123,
			wantOk:   true,
		},
		{
			name:     "When value is uint",
			settable: reflect.ValueOf(new(uint)),
			str:      "123",
			want:     uint(123),
			wantOk:   true,
		},
		{
			name:     "When value is float",
			settable: reflect.ValueOf(new(float64)),
			str:      "123.45",
			want:     123.45,
			wantOk:   true,
		},
		{
			name:     "When value is time.Duration",
			settable: reflect.ValueOf(new(time.Duration)),
			str:      "1h",
			want:     time.Hour,
			wantOk:   true,
		},
		{
			name:     "When value is complex",
			settable: reflect.ValueOf(new(complex128)),
			str:      "1+2i",
			want:     nil,
			wantOk:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ok := convertAndSetValue(test.settable, test.str)
			if ok != test.wantOk {
				t.Errorf("convertAndSetValue() ok = %v, wantOk %v", ok, test.wantOk)
			}

			if ok && !reflect.DeepEqual(test.settable.Elem().Interface(), test.want) {
				t.Errorf("convertAndSetValue() = %v, want %v", test.settable.Elem().Interface(), test.want)
			}
		})
	}
}
