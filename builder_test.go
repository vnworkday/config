package config

import (
	"os"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  map[string]string
		target  any
		prefix  string
		wantOut any
		wantErr bool
	}{
		{
			name: "When config map matches struct fields then it should populate the struct",
			config: map[string]string{
				"Field1": "value1",
				"Field2": "value2",
			},
			target: &struct {
				Field1 string
				Field2 string
			}{},
			prefix: "",
			wantOut: &struct {
				Field1 string
				Field2 string
			}{
				Field1: "value1",
				Field2: "value2",
			},
			wantErr: false,
		},
		{
			name: "When config map has extra fields then it should ignore them",
			config: map[string]string{
				"Field1": "value1",
				"Field2": "value2",
				"Field3": "value3",
			},
			target: &struct {
				Field1 string
				Field2 string
			}{},
			prefix: "",
			wantOut: &struct {
				Field1 string
				Field2 string
			}{
				Field1: "value1",
				Field2: "value2",
			},
			wantErr: false,
		},
		{
			name: "When config map is missing fields then it should leave them as zero values",
			config: map[string]string{
				"Field1": "value1",
			},
			target: &struct {
				Field1 string
				Field2 string
			}{},
			prefix: "",
			wantOut: &struct {
				Field1 string
				Field2 string
			}{
				Field1: "value1",
			},
			wantErr: false,
		},
		{
			name: "When config map has invalid values then it should return an error",
			config: map[string]string{
				"Field1": "value1",
				"Field2": "not a number",
			},
			target: &struct {
				Field1 string
				Field2 int
			}{},
			prefix: "",
			wantOut: &struct {
				Field1 string
				Field2 int
			}{
				Field1: "value1",
				Field2: 0,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			b := &Builder{
				configMap: test.config,
			}

			err := b.decode(test.target, test.prefix)
			if (err != nil) != test.wantErr {
				t.Errorf(failTestMessage("decode", test.wantErr, err))
			}

			if !test.wantErr && !reflect.DeepEqual(test.target, test.wantOut) {
				t.Errorf(failTestMessage("decode", test.wantOut, test.target))
			}
		})
	}
}

func TestAppendFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		file       string
		includeErr bool
		wantErr    bool
	}{
		{
			name:       "When file contains valid content then it should be added to the config map",
			file:       "valid.txt",
			includeErr: true,
			wantErr:    false,
		},
		{
			name:       "When file does not exist then it should return an error",
			file:       "nonexistent.txt",
			includeErr: true,
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			builder := &Builder{
				configMap: make(map[string]string),
			}

			// Create a temporary file with the specified content
			if test.file != "nonexistent.txt" {
				tmpFile, err := os.CreateTemp("", "example.*.txt")
				if err != nil {
					t.Fatal(err)
				}
				defer func(name string) {
					_ = os.Remove(name)
				}(tmpFile.Name())

				if _, err := tmpFile.WriteString("key=value"); err != nil {
					t.Fatal(err)
				}

				if err := tmpFile.Close(); err != nil {
					t.Fatal(err)
				}

				test.file = tmpFile.Name()
			}

			builder.appendFile(test.file, test.includeErr)

			if (len(builder.failedFields) > 0) != test.wantErr {
				t.Errorf(failTestMessage("appendFile", test.wantErr, len(builder.failedFields) > 0))
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		structDelimiter string
		sliceDelimiter  string
		wantStructDelim string
		wantSliceDelim  string
	}{
		{
			name:            "When no options are provided then it should use default delimiters",
			structDelimiter: "",
			sliceDelimiter:  "",
			wantStructDelim: ".",
			wantSliceDelim:  " ",
		},
		{
			name:            "When struct delimiter option is provided then it should be set",
			structDelimiter: "_",
			sliceDelimiter:  "",
			wantStructDelim: "_",
			wantSliceDelim:  " ",
		},
		{
			name:            "When slice delimiter option is provided then it should be set",
			structDelimiter: "",
			sliceDelimiter:  ",",
			wantStructDelim: ".",
			wantSliceDelim:  ",",
		},
		{
			name:            "When both delimiters options are provided then they should be set",
			structDelimiter: "_",
			sliceDelimiter:  ",",
			wantStructDelim: "_",
			wantSliceDelim:  ",",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var opts []Option

			if test.structDelimiter != "" {
				opts = append(opts, WithStructDelimiter(test.structDelimiter))
			}

			if test.sliceDelimiter != "" {
				opts = append(opts, WithSliceDelimiter(test.sliceDelimiter))
			}

			builder := newBuilder(opts...)

			if builder.structDelimiter != test.wantStructDelim {
				t.Errorf("newBuilder() structDelimiter = %v, want %v", builder.structDelimiter, test.wantStructDelim)
			}

			if builder.sliceDelimiter != test.wantSliceDelim {
				t.Errorf("newBuilder() sliceDelimiter = %v, want %v", builder.sliceDelimiter, test.wantSliceDelim)
			}
		})
	}
}
