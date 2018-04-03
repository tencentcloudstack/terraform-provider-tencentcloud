package goset

import (
	"fmt"
	"reflect"
)

// Uniq the slice of objects, the objects must be the same type, both builtin and custom types are supported.
func Uniq(elements interface{}) interface{} {
	v := reflect.ValueOf(elements)
	if !isAvailableSlice(v) {
		return elements
	}
	if v.Len() <= 1 {
		return elements
	}

	slim := reflect.MakeSlice(reflect.TypeOf(elements), 0, v.Cap())
	for i := 0; i < v.Len(); i++ {
		found := false
		for j := 0; j < slim.Len(); j++ {
			if reflect.DeepEqual(
				v.Index(i).Interface(),
				slim.Index(j).Interface(),
			) {
				found = true
			}
		}
		if found {
			continue
		}
		slim = reflect.Append(slim, v.Index(i))
	}

	return slim.Interface()
}

// Return the intersection of aSet and bSet, aSet and bSet must be a slice type, and the two slice must have the same type of elements.
// Empty slice is allowed, but nil is not allowed to represent the null set concept.
func Intersect(aSet interface{}, bSet interface{}) interface{} {
	_, iSet, _, _ := Difference(aSet, bSet)
	return iSet
}

// Return the union of aSet and bSet, aSet and bSet must be a slice type, and the two slice must have the same type of elements.
// Empty slice is allowed, but nil is not allowed to represent the null set concept.
func Union(aSet interface{}, bSet interface{}) interface{} {
	uSet, _, _, _ := Difference(aSet, bSet)
	return uSet
}

// Return the difference of aSet and bSet, aSet and bSet must be a slice type, and the two slice must have the same type of elements.
// Empty slice is allowed, but nil is not allowed to represent the null set concept.
func Difference(aSet interface{}, bSet interface{}) (iUnion, iIntersection, iADifferenceSet, iBDifferenceSet interface{}) {
	av := reflect.ValueOf(aSet)
	bv := reflect.ValueOf(bSet)
	if !areAvailableSlices(av, bv) {
		panic("A set and B set should be slices and have the same type of elements")
	}

	var union = reflect.MakeSlice(reflect.TypeOf(aSet), 0, av.Cap()+bv.Cap())
	var intersection = reflect.MakeSlice(reflect.TypeOf(aSet), 0, av.Cap()+bv.Cap())
	var aDifferenceSet = reflect.MakeSlice(reflect.TypeOf(aSet), 0, av.Cap())
	var bDifferenceSet = reflect.MakeSlice(reflect.TypeOf(aSet), 0, bv.Cap())

	for i := 0; i < av.Len(); i++ {
		skip := false
		for j := 0; j < bv.Len(); j++ {
			if reflect.DeepEqual(
				bv.Index(j).Interface(),
				av.Index(i).Interface(),
			) {
				intersection = reflect.Append(intersection, bv.Index(j))
				skip = true
			}
		}
		if !skip {
			aDifferenceSet = reflect.Append(aDifferenceSet, av.Index(i))
		}
	}

	for i := 0; i < bv.Len(); i++ {
		skip := false
		for j := 0; j < intersection.Len(); j++ {
			if reflect.DeepEqual(
				intersection.Index(j).Interface(),
				bv.Index(i).Interface(),
			) {
				skip = true
			}
		}
		if !skip {
			bDifferenceSet = reflect.Append(bDifferenceSet, bv.Index(i))
		}
	}

	union = reflect.AppendSlice(aDifferenceSet, intersection)
	union = reflect.AppendSlice(union, bDifferenceSet)

	iUnion = union.Interface()
	iIntersection = intersection.Interface()
	iADifferenceSet = aDifferenceSet.Interface()
	iBDifferenceSet = bDifferenceSet.Interface()

	return iUnion, iIntersection, iADifferenceSet, iBDifferenceSet
}

// Add an element to a set, the element's type must be the same as the type of existed set's element
// If the element's value already exists on the set, the return set keep the same as the old set.
func AddElement(set interface{}, e interface{}) interface{} {
	v := reflect.ValueOf(set)
	if v.Type().Elem() != reflect.TypeOf(e) {
		panic("Set and element are not the same type")
	}

	if !isAvailableSlice(v) {
		panic("Invalid Slice")
	}

	if !IsUniq(set) {
		panic("Set should be uniq")
	}

	ev := reflect.ValueOf(e)

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(
			e,
			v.Index(i).Interface(),
		) {
			return set
		}
	}

	v = reflect.Append(v, ev)
	return v.Interface()
}

// Merge two sets, aSet and bSet must be the same type
// Same effect as the union aSet and bSet.
func AddElements(aSet interface{}, bSet interface{}) interface{} {
	av := reflect.ValueOf(aSet)
	bv := reflect.ValueOf(bSet)
	if !areAvailableSlices(av, bv) {
		panic("Invalid Slices")
	}

	for i := 0; i < bv.Len(); i++ {
		aSet = AddElement(aSet, bv.Index(i).Interface())
	}
	return Uniq(aSet)
}

// Remove an element to a set, the element's type must be the same as the type of existed set's element
// If the element's value did not exists on the set, the return set keep the same as the old set.
func RemoveElement(set interface{}, e interface{}) interface{} {
	v := reflect.ValueOf(set)
	if !isAvailableSlice(v) {
		panic("Invalid Slice")
	}

	if v.Len() == 0 {
		return set
	}

	ev := reflect.ValueOf(e)
	if !ev.IsValid() {
		return set
	}

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(
			e,
			v.Index(i).Interface(),
		) {
			v = reflect.AppendSlice(
				v.Slice(0, i),
				v.Slice(i+1, v.Len()),
			)
			return v.Interface()
		}
	}

	return set
}

// Reduce aSet's elements by looking up the bSet, aSet and bSet must be the same type
// Same effect as the difference aSet and bSet.
func RemoveElements(aSet interface{}, bSet interface{}) interface{} {
	av := reflect.ValueOf(aSet)
	bv := reflect.ValueOf(bSet)
	if !areAvailableSlices(av, bv) {
		panic("Invalid Slices")
	}

	newSetSlice := reflect.MakeSlice(reflect.TypeOf(aSet), av.Len(), av.Cap())
	for i := 0; i < av.Len(); i++ {
		newSetSlice.Index(i).Set(av.Index(i))
	}
	var newSet = newSetSlice.Interface()
	for i := 0; i < bv.Len(); i++ {
		newSet = RemoveElement(newSet, bv.Index(i).Interface())
	}
	return Uniq(newSet)
}

// To tell if the aSet is uniq, aSet must be a slices.
func IsUniq(aSet interface{}) bool {
	v := reflect.ValueOf(aSet)
	if !isAvailableSlice(v) {
		return false
	}
	if v.Len() <= 1 {
		return true
	}
	ele := v.Index(0).Interface()
	others := reflect.MakeSlice(reflect.TypeOf(aSet), v.Len()-1, v.Cap())
	for i := 1; i < v.Len(); i++ {
		if reflect.DeepEqual(
			ele,
			v.Index(i).Interface(),
		) {
			return false
		}
		others = v.Slice(1, v.Len())
	}
	return IsUniq(others.Interface())
}

// Tell if the two set is equal, both aSet and bSet must be slices with same types.
// Sequence of elements is ignored to be the factor of equal detection.
func IsEqual(aSet interface{}, bSet interface{}) bool {
	av := reflect.ValueOf(aSet)
	bv := reflect.ValueOf(bSet)
	if av.Len() != bv.Len() {
		return false
	}
	if av.Len() == 0 && bv.Len() == 0 {
		return true
	}
	if !areAvailableSlices(av, bv) {
		return false
	}

	aMap := make(map[int]bool)
	bMap := make(map[int]bool)

	hits := 0
	for i := 0; i < av.Len(); i++ {
		if aMap[i] {
			continue
		}
		found := false
		for j := 0; j < bv.Len(); j++ {
			if bMap[j] {
				continue
			}
			if reflect.DeepEqual(
				av.Index(i).Interface(),
				bv.Index(j).Interface(),
			) {
				aMap[i] = true
				bMap[j] = true
				hits += 1
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return hits == av.Len() && hits == bv.Len()
}

// Return true if set contains the ele element, set must be a slice with the same type of ele.
func IsIncluded(set interface{}, ele interface{}) bool {
	ev := reflect.ValueOf(ele)
	if !ev.IsValid() {
		return true
	}
	v := reflect.ValueOf(set)
	if !isAvailableSlice(v) {
		return false
	}
	if v.Len() == 0 {
		return false
	}
	if reflect.TypeOf(ev).String() != reflect.TypeOf(v.Index(0)).String() {
		return false
	}

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(
			v.Index(i).Interface(),
			ev.Interface(),
		) {
			return true
		}
	}

	return false
}

// Return true if aSet is the subset of bSet, bothe aSet and bSet must be slices and have the same type of elements.
func IsSubset(aSet interface{}, bSet interface{}) bool {
	_, _, aSubSet, _ := Difference(aSet, bSet)
	return reflect.ValueOf(aSubSet).Len() == 0
}

// Return true if aSet is the superset of bSet, bothe aSet and bSet must be slices and their elements's type must be the same too.
func IsSuperset(aSet interface{}, bSet interface{}) bool {
	_, _, _, bSubSet := Difference(aSet, bSet)
	return reflect.ValueOf(bSubSet).Len() == 0
}

// mapFunc accepts one parameter and output one parameter
// like func(string)int
// defaultReturnSlice is exacly mapped when the set is an empty slice
func Map(set interface{}, mapFunc interface{}, defaultReturnSlice interface{}) (mapped interface{}) {
	vSet := reflect.ValueOf(set)
	if vSet.Kind() != reflect.Slice {
		panic("source set must be slice")
	}

	vFunc := reflect.ValueOf(mapFunc)
	if vFunc.Kind() != reflect.Func {
		panic("mapper must be a func")
	}

	if vSet.Len() == 0 {
		return defaultReturnSlice
	}

	slim := reflect.MakeSlice(reflect.TypeOf(defaultReturnSlice), 0, vSet.Cap())
	for i := 0; i < vSet.Len(); i++ {
		ele := vFunc.Call([]reflect.Value{vSet.Index(i)})
		slim = reflect.Append(slim, ele[0])
	}
	return slim.Interface()
}

// ids is a slice of buildin types
// objs is a slice of structs
// fieldName is the name of ids
// reorderObjects is the same set of objs with order changed according to the ids
func Reorder(ids interface{}, objs interface{}, fieldName string) (reordedObjects interface{}) {
	objsType := reflect.TypeOf(objs)
	if objsType.Kind() != reflect.Slice {
		panic("objs should be slice")
	}

	idsType := reflect.TypeOf(ids)
	if idsType.Kind() != reflect.Slice {
		panic("ids should be slice")
	}

	idsValue := reflect.ValueOf(ids)
	objsValue := reflect.ValueOf(objs)

	// fast return if objs is empty or ids is abnormal
	if objsValue.Len() == 0 || objsValue.Len() != idsValue.Len() {
		return objs
	}

	idSlice := reflect.MakeSlice(idsType, 0, idsValue.Cap())
	for i := 0; i < idsValue.Len(); i++ {
		ev := idsValue.Index(i)
		idSlice = reflect.Append(idSlice, ev)
	}

	m := make(map[string]reflect.Value)

	for i := 0; i < objsValue.Len(); i++ {
		var ev reflect.Value
		old_ev := objsValue.Index(i)
		if old_ev.Kind() == reflect.Ptr {
			ev = old_ev.Elem()
		} else {
			ev = old_ev
		}

		id := ev.FieldByName(fieldName)
		s := fmt.Sprintf("%v", id)
		m[s] = old_ev
	}

	newObjs := reflect.MakeSlice(objsType, 0, objsValue.Cap())
	for i := 0; i < idSlice.Len(); i++ {
		id := idSlice.Index(i)
		k := fmt.Sprintf("%v", id)
		if v, ok := m[k]; ok {
			newObjs = reflect.Append(newObjs, v)
		}
	}

	return newObjs.Interface()

}
