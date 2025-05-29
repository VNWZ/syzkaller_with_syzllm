package syzllm_pkg

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSortedMapKeysDesc(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]string
		expected []string
	}{
		{
			name: "unordered",
			m: map[string]string{
				"7":  "c",
				"8":  "d",
				"2":  "f",
				"5":  "a",
				"6":  "b",
				"10": "g",
				"9":  "e",
			},
			expected: []string{"10", "9", "8", "7", "6", "5", "2"},
		},
		{
			name: "Ascending",
			m: map[string]string{
				"2":  "f",
				"5":  "a",
				"6":  "b",
				"7":  "c",
				"8":  "d",
				"9":  "e",
				"10": "g",
			},
			expected: []string{"10", "9", "8", "7", "6", "5", "2"},
		},
		{
			name: "Descending",
			m: map[string]string{
				"10": "g",
				"9":  "e",
				"8":  "f",
				"7":  "c",
				"6":  "b",
				"5":  "a",
				"2":  "f",
			},
			expected: []string{"10", "9", "8", "7", "6", "5", "2"},
		},
	}

	// Iterate over test cases.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SortedMapKeysDesc(tc.m)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("%s: got\n %v\n, want: \n%v", tc.name, result, tc.expected)
			}
		})
	}
}

func TestInsertItem(t *testing.T) {
	// Test cases for InsertItem with string slices
	stringTests := []struct {
		name     string
		slice    []string
		item     string
		pos      int
		expected []string
		err      error
	}{
		{
			name:     "InsertAtStart",
			slice:    []string{"a", "b", "c"},
			item:     "x",
			pos:      0,
			expected: []string{"x", "a", "b", "c"},
			err:      nil,
		},
		{
			name:     "InsertAtMiddle",
			slice:    []string{"a", "b", "c"},
			item:     "x",
			pos:      2,
			expected: []string{"a", "b", "x", "c"},
			err:      nil,
		},
		{
			name:     "InsertAtEnd",
			slice:    []string{"a", "b", "c"},
			item:     "x",
			pos:      3,
			expected: []string{"a", "b", "c", "x"},
			err:      nil,
		},
		{
			name:     "InsertInEmpty",
			slice:    []string{},
			item:     "x",
			pos:      0,
			expected: []string{"x"},
			err:      nil,
		},
		{
			name:     "InvalidPositionNegative",
			slice:    []string{"a", "b"},
			item:     "x",
			pos:      -1,
			expected: nil,
			err:      fmt.Errorf("invalid position -1; must be between 0 and 2"),
		},
		{
			name:     "InvalidPositionTooLarge",
			slice:    []string{"a", "b"},
			item:     "x",
			pos:      3,
			expected: nil,
			err:      fmt.Errorf("invalid position 3; must be between 0 and 2"),
		},
	}

	// Run string tests
	for _, tc := range stringTests {
		t.Run("String/"+tc.name, func(t *testing.T) {
			result, err := InsertItem(tc.slice, tc.item, tc.pos)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("%s: got \n%v\n, want: \n%v", tc.name, result, tc.expected)
			}
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("%s: error got %v, want %v", tc.name, err, tc.err)
			}
		})
	}
}

func TestEnlargeSlice(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		pos      int
		expected []string
		err      error
	}{
		{
			name:     "InsertAtStart",
			slice:    []string{"a", "b", "c"},
			pos:      0,
			expected: []string{"", "a", "b", "c"},
			err:      nil,
		},
		{
			name:     "InsertAtMiddle",
			slice:    []string{"a", "b", "c"},
			pos:      2,
			expected: []string{"a", "b", "", "c"},
			err:      nil,
		},
		{
			name:     "InsertAtEnd",
			slice:    []string{"a", "b", "c"},
			pos:      3,
			expected: []string{"a", "b", "c", ""},
			err:      nil,
		},
		{
			name:     "InsertInEmpty",
			slice:    []string{},
			pos:      0,
			expected: []string{""},
			err:      nil,
		},
		{
			name:     "InvalidPositionNegative",
			slice:    []string{"a", "b"},
			pos:      -1,
			expected: nil,
			err:      fmt.Errorf("invalid position -1"),
		},
		{
			name:     "InvalidPositionTooLarge",
			slice:    []string{"a", "b"},
			pos:      3,
			expected: nil,
			err:      fmt.Errorf("invalid position 3"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := enlargeSlice(tc.slice, tc.pos)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("%s: got %v, want %v", tc.name, result, tc.expected)
			}
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("%s: error got %v, want %v", tc.name, err, tc.err)
			}
		})
	}
}
