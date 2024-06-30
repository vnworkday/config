package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

const (
	defaultStructDelimiter = "."
	defaultSliceDelimiter  = " "
)

type Builder struct {
	structDelimiter string
	sliceDelimiter  string
	configMap       map[string]string
	failedFields    []string
}

// FromEnv reads environment variables and adds them to the config map.
func (b *Builder) FromEnv() *Builder {
	mergeMaps(b.configMap, stringsToMap(os.Environ()))

	return b
}

func (b *Builder) FromFile(file string) *Builder {
	return b.appendFile(file, true)
}

func (b *Builder) MapTo(target any) error {
	return b.decode(target, "")
}

// Sub accepts a struct pointer and a prefix, and populates it with the current config state.
// The prefix is prepended to all keys when looking up values.
func (b *Builder) Sub(target any, prefix string) error {
	return b.decode(target, prefix+b.structDelimiter)
}

func (b *Builder) decode(target any, prefix string) error {
	structPtr := reflect.ValueOf(target)

	if structPtr.Kind() != reflect.Ptr || structPtr.Elem().Kind() != reflect.Struct {
		panic("target must be a struct pointer")
	}

	m := make(map[string]reflect.Value)
	mapKeysToFields(structPtr, m, prefix, b.structDelimiter)

	for key, fieldPtr := range m {
		stringValue, ok := b.configMap[key]

		if !ok {
			continue
		}

		switch fieldPtr.Elem().Type().Kind() {
		case reflect.Slice:
			for _, i := range convertAndSetSlice(fieldPtr, stringToSlice(stringValue, b.sliceDelimiter)) {
				b.failedFields = append(b.failedFields, fmt.Sprintf("%s[%d]", key, i))
			}
		default:
			if !convertAndSetValue(fieldPtr, stringValue) {
				b.failedFields = append(b.failedFields, key)
			}
		}
	}

	sort.Strings(b.failedFields) // sort for deterministic output

	if len(b.failedFields) > 0 {
		return errors.Errorf("failed to convert fields: %s", strings.Join(b.failedFields, ", "))
	}

	return nil
}

// appendFile reads a file and adds its contents to the config map.
// If includeErr is true, it will also add any errors to the failedFields slice.
func (b *Builder) appendFile(file string, includeErr bool) *Builder {
	content, err := os.ReadFile(file)

	if includeErr && err != nil {
		b.failedFields = append(b.failedFields, fmt.Sprintf("file: %s, error: %s", file, err))
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	var scannedStrings []string
	for scanner.Scan() {
		scannedStrings = append(scannedStrings, scanner.Text())
	}

	if includeErr && scanner.Err() != nil {
		b.failedFields = append(b.failedFields, fmt.Sprintf("file: %s, error: %s", file, scanner.Err()))
	}

	mergeMaps(b.configMap, stringsToMap(scannedStrings))

	return b
}

func newBuilder(opts ...Option) *Builder {
	builder := &Builder{
		structDelimiter: defaultStructDelimiter,
		sliceDelimiter:  defaultSliceDelimiter,
		configMap:       make(map[string]string),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

type Option func(*Builder)

func WithStructDelimiter(d string) Option {
	return func(b *Builder) {
		b.structDelimiter = d
	}
}

func WithSliceDelimiter(d string) Option {
	return func(b *Builder) {
		b.sliceDelimiter = d
	}
}
