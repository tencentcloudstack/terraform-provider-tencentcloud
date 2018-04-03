package goset

import "reflect"

func isAvailableSlice(v reflect.Value) bool {
	if v.Kind() != reflect.Slice {
		return false
	}

	var kind string
	for i := 0; i < v.Len(); i++ {
		eleKind := reflect.TypeOf(v.Index(i)).Kind().String()
		if i == 0 {
			kind = eleKind
		} else {
			if kind != eleKind {
				return false
			}
		}
	}

	return true
}

func areAvailableSlices(v1, v2 reflect.Value) bool {
	if !isAvailableSlice(v1) {
		return false
	}
	if !isAvailableSlice(v2) {
		return false
	}
	if v1.Len() == 0 && v2.Len() == 0 {
		return true
	}
	return reflect.TypeOf(v1).Kind().String() == reflect.TypeOf(v2).Kind().String()
}
