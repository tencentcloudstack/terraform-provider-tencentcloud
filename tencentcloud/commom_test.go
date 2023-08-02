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

func TestMatchAny(t *testing.T) {
	assert.False(t, MatchAny(1))
	assert.True(t, MatchAny(1, 1, 2, 3))
	assert.False(t, MatchAny(1, 4, 5, 6))
	assert.True(t, MatchAny("a", "b", "c", "a", "aa"))

	one := 1
	two := 2
	var ptrOne *int
	var nilVal *string

	ptrOne = &one
	assert.True(t, MatchAny(ptrOne, &one, two))

	assert.False(t, MatchAny(ptrOne, 5, 6, 7))
	assert.False(t, MatchAny(nilVal, nil))

	var oneI64 int64 = 1
	var oneUI64 uint64 = 1

	assert.False(t, MatchAny(oneI64, 2, 1, 3))
	assert.False(t, MatchAny(oneUI64, 1))
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

func TestGetListDiffs(t *testing.T) {
	var (
		o1           = []int{1, 2, 3, 4, 5}
		n1           = []int{1, 2, 3, 5, 6}
		expectAdds1  = []int{6}
		expectLacks1 = []int{4}
	)

	adds1, lacks1 := GetListDiffs(o1, n1)
	assert.Contains(t, expectAdds1, adds1[0])
	assert.Contains(t, expectLacks1, lacks1[0])

	var (
		o2           = []int{1, 1, 3, 4, 5}
		n2           = []int{4, 1, 1, 5, 3}
		expectAdds2  = make([]int, 0)
		expectLacks2 = make([]int, 0)
	)

	adds2, lacks2 := GetListDiffs(o2, n2)
	assert.Equalf(t, len(expectAdds2), len(adds2), "adds2 should be %v, got %v", expectAdds2, adds2)
	assert.Equalf(t, len(expectLacks2), len(lacks2), "lacks2 should be %v, got %v", expectLacks2, lacks2)

	// TODO
	var (
		o3           = []int{1, 3, 3, 4, 4}
		n3           = []int{4, 3, 1, 7, 3, 6}
		expectAdds3  = []int{6, 7}
		expectLacks3 = []int{4}
	)

	adds3, lacks3 := GetListDiffs(o3, n3)
	assert.Contains(t, expectAdds3, adds3[0])
	assert.Contains(t, expectAdds3, adds3[1])
	assert.Equalf(t, expectLacks3, lacks3, "lacks3 should be %v, got %v", expectLacks3, lacks3)

	var (
		o4           = []int{1, 2, 3, 4, 5}
		n4           = []int{4}
		expectAdds4  = make([]int, 0)
		expectLacks4 = []int{1, 2, 3, 5}
	)

	adds4, lacks4 := GetListDiffs(o4, n4)
	assert.Equalf(t, len(expectAdds4), len(adds4), "adds4 should be %v, got %v", expectAdds4, adds4)
	assert.Contains(t, expectLacks4, lacks4[0])
	assert.Contains(t, expectLacks4, lacks4[1])
	assert.Contains(t, expectLacks4, lacks4[2])
	assert.Contains(t, expectLacks4, lacks4[3])

	adds5, lacks5 := GetListDiffs([]int{100003, 100004}, []int{100003, 100003, 100004})
	assert.Equalf(t, 1, len(adds5), "")
	assert.Equalf(t, 0, len(lacks5), "")
	assert.Contains(t, []int{100003, 100003, 100004}, adds5[0])

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
