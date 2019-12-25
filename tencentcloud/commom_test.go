package tencentcloud

import (
	"testing"
)

func TestIsContains(t *testing.T) {
	var s *string
	var i interface{} = s
	tests := []struct {
		x interface{}
		y interface{}
		z bool
	}{
		{nil, nil, false},
		{s, s, false},
		{&i, &i, true},

		{[]int{0, 1, 2}, 1, true},
		{[]int{0, 1, 2}, 3, false},
		{[]int{0, 1, 2}, int64(1), false},
		{[]int{0, 1, 2}, "1", false},
		{[]int{0, 1, 2}, true, false},

		{[]int64{0, 1, 2}, int64(1), true},
		{[]int64{0, 1, 2}, int64(3), false},
		{[]int64{0, 1, 2}, 1, false},
		{[]int64{0, 1, 2}, "1", false},
		{[]int64{0, 1, 2}, true, false},

		{[]float64{0.0, 1.0, 2.0}, 1.0, true},
		{[]float64{0.0, 1.0, 2.0}, float64(1), true},
		{[]float64{0.0, 1.0, 2.0}, 3.0, false},
		{[]float64{0.0, 1.0, 2.0}, 1, false},

		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, 1, false},
		{[]string{"a", "b", "c"}, true, false},

		{[]interface{}{0, "1", 2}, "1", true},
		{[]interface{}{0, "1", 2}, 1, false},
		{[]interface{}{0, 1, 2}, true, false},
		{[]interface{}{0, true, 2}, true, true},
		{[]interface{}{0, false, 2}, true, false},

		{[]interface{}{[]int{0, 1}, []int{1, 2}}, []int{1, 2}, true},
		{[]interface{}{[]int{0, 1}, []int{1, 2, 3}}, []int{1, 2}, false},

		{map[string]int{"a": 1}, "a", true},
		{map[string]int{"a": 1}, "d", false},
		{map[string]int{"a": 1}, 1, false},
		{map[string]int{"a": 1}, true, false},

		{"abc", "a", true},
		{"abc", "d", false},
		{"abc", 1, false},
		{"abc", true, false},

		{"a", "a", true},
		{1, 1, true},
		{-1, -1, true},
		{1.0, 1.0, true},
		{true, true, true},
		{false, false, true},
	}

	for _, v := range tests {
		if IsContains(v.x, v.y) != v.z {
			t.Errorf("Failed: %#v", v)
		}
	}
}
