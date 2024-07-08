package config

import "testing"

func Test_shouldPanic(t *testing.T) {
	t.Parallel()

	var emptyStruct struct{}
	var someNumber int

	tests := []struct {
		name      string
		target    any
		wantPanic bool
	}{
		{
			name:      "struct",
			target:    emptyStruct,
			wantPanic: true,
		},
		{
			name:      "*int",
			target:    &someNumber,
			wantPanic: true,
		},
		{
			name:      "int",
			target:    someNumber,
			wantPanic: true,
		},
		{
			name:      "nil",
			target:    nil,
			wantPanic: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if r := recover(); (r == nil) == test.wantPanic {
					t.Errorf("should have caused a panic")
				}
			}()

			_, _ = LoadConfig[any](&test.target)
		})
	}
}
