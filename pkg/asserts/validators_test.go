package asserts_test

import (
	"fmt"
	"testing"

	"github.com/bergzeit/goepcqrcode/pkg/asserts"
)

func TestAssertInBetweenInt(t *testing.T) {
	tests := []struct {
		left, right, value int
		expected           bool
	}{
		{1, 10, 5, true},
		{1, 10, 0, false},
		{1, 10, 11, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v in between %v and %v", tt.value, tt.left, tt.right), func(t *testing.T) {
			result := asserts.AssertInBetween(tt.left, tt.right, tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAssertInBetweenFloat(t *testing.T) {
	tests := []struct {
		left, right, value int
		expected           bool
	}{
		{1.0, 10.0, 5.0, true},
		{1.0, 10.0, 0.0, false},
		{1.0, 10.0, 11.0, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v in between %v and %v", tt.value, tt.left, tt.right), func(t *testing.T) {
			result := asserts.AssertInBetween(tt.left, tt.right, tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAssertContains(t *testing.T) {
	tests := []struct {
		slice    interface{}
		value    interface{}
		expected bool
	}{
		{[]int{1, 2, 3, 4, 5}, 3, true},
		{[]int{1, 2, 3, 4, 5}, 6, false},
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]float64{1.1, 2.2, 3.3}, 2.2, true},
		{[]float64{1.1, 2.2, 3.3}, 4.4, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			switch slice := tt.slice.(type) {
			case []int:
				result := asserts.AssertContains(slice, tt.value.(int))
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			case []string:
				result := asserts.AssertContains(slice, tt.value.(string))
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			case []float64:
				result := asserts.AssertContains(slice, tt.value.(float64))
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			default:
				t.Errorf("unsupported slice type")
			}
		})
	}
}

func TestAssertSize(t *testing.T) {
	tests := []struct {
		slice    []byte
		maxSize  int
		expected int
	}{
		{[]byte{1, 2, 3}, 5, -1},         // smaller
		{[]byte{1, 2, 3, 4, 5}, 5, 0},    // equal
		{[]byte{1, 2, 3, 4, 5, 6}, 5, 1}, // bigger
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := asserts.AssertSize(tt.slice, tt.maxSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
