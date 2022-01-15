package pointer

import (
	"reflect"
	"testing"
)

const (
	maxInt = int(^uint(0) >> 1)
	minInt = -maxInt - 1
)

func TestIntP(t *testing.T) {
	tests := []struct {
		name string
		in   int
	}{
		{
			name: "Negative",
			in:   -1,
		},
		{
			name: "Zero",
			in:   0,
		},
		{
			name: "Positive",
			in:   1,
		},
		{
			name: "Maximum",
			in:   maxInt,
		},
		{
			name: "Minimum",
			in:   minInt,
		},
		{
			name: "ImplicitZero",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IntP(test.in)

			if actual == nil {
				t.Error("expected pointer, got nil")
			} else if reflect.TypeOf(actual).Kind() != reflect.Ptr {
				t.Errorf("expected pointer, got %v", actual)
			} else if *actual != test.in {
				t.Errorf("expected %d, got %d", test.in, *actual)
			}
		})
	}
}

func TestStringP(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{
			name: "Blank",
			in:   "",
		},
		{
			name: "ImplicitBlank",
		},
		{
			name: "FooBar",
			in:   "foobar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := StringP(test.in)

			if actual == nil {
				t.Error("expected pointer, got nil")
			} else if reflect.TypeOf(actual).Kind() != reflect.Ptr {
				t.Errorf("expected pointer, got %v", actual)
			} else if *actual != test.in {
				t.Errorf("expected %s, got %s", test.in, *actual)
			}
		})
	}
}
