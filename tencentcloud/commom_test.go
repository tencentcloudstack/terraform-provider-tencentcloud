package tencentcloud

import (
	"reflect"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/stretchr/testify/assert"
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

func TestGetListIncrement(t *testing.T) {
	var (
		old1      = []int{1, 2, 2, 3, 5}
		new1      = []int{1, 2, 3, 2, 4, 5, 6, 3}
		expected1 = []int{4, 6, 3}
	)
	actual1, _ := GetListIncrement(old1, new1)
	assert.Equalf(t, expected1, actual1, "incr1 should equal, got %v %v", expected1, actual1)

	var (
		old2      = []int{1, 2, 4, 5}
		new2      = []int{1, 2, 3, 2, 4, 5, 6, 3}
		expected2 = []int{3, 2, 6, 3}
	)
	actual2, _ := GetListIncrement(old2, new2)
	assert.Equalf(t, expected2, actual2, "incr1 should equal, got %v %v", expected2, actual2)

	var (
		old3      = []int{1, 2, 4, 5, 3, 6, 3, 2}
		new3      = []int{1, 2, 3, 2, 4, 5, 6, 3}
		expected3 = make([]int, 0)
	)
	actual3, _ := GetListIncrement(old3, new3)
	assert.Equalf(t, expected3, actual3, "incr1 should equal, got %v %v", expected3, actual3)

	var (
		old4 = []int{1}
		new4 = []int{2}
	)

	_, err := GetListIncrement(old4, new4)
	assert.EqualError(t, err, "elem 1 not exist")

}

func TestIsExpectError(t *testing.T) {

	err := sdkErrors.NewTencentCloudSDKError("ClientError.NetworkError", "", "")

	// Expected
	expectedFull := []string{"ClientError.NetworkError"}
	expectedShort := []string{"ClientError"}
	assert.Equalf(t, isExpectError(err, expectedFull), true, "")
	assert.Equalf(t, isExpectError(err, expectedShort), true, "")

	// Unexpected
	unExpectedMatchHead := []string{"ClientError.HttpStatusCodeError"}
	unExpectedShort := []string{"SystemError"}
	assert.Equalf(t, isExpectError(err, unExpectedMatchHead), false, "")
	assert.Equalf(t, isExpectError(err, unExpectedShort), false, "")
}

func TestYamlParser(t *testing.T) {
	yamlStr := `test`
	yaml, err := YamlParser(yamlStr)
	assert.Equalf(t, err != nil, true, "")
	assert.Equalf(t, yaml == nil, true, "")

	yamlStr1 := `version: v1
name: test-name
desc: this is a description`
	yaml1, err1 := YamlParser(yamlStr1)
	assert.Equalf(t, err1, nil, "")
	assert.Equalf(t, reflect.TypeOf(yaml1).String(), "map[interface {}]interface {}", "")
	assert.Equalf(t, yaml1["name"], "test-name", "")
}
