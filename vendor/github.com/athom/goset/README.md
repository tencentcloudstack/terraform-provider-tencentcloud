# Go Set

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/athom/goset?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Go Set Package is a simple library for set operations with generic type supported.

[![Build Status](https://api.travis-ci.org/athom/goset.png?branch=master)](https://travis-ci.org/athom/goset)
[![GoDoc](https://godoc.org/github.com/athom/goset?status.png)](http://godoc.org/github.com/athom/goset)


## Installation

```bash
	go get "github.com/athom/goset"
```

## Features

- **Generic**

  All Go builtin types and custom defined types are supported.
  Even slice of pointers!

```go
a := goset.Uniq([]int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}).([]int)
b := goset.Uniq([]string{"1", "2", "2", "3", "3", "3", "4"}).([]string)

type Avatar struct {
        Age  int
        Name string
}

avatars := []Avatar{
        Avatar{112, "Angg"},
        Avatar{70, "Roku"},
        Avatar{230, "Kyoshi"},
        Avatar{230, "Kyoshi"},
        Avatar{33, "Kuruk"},
        Avatar{33, "Kuruk"},
        Avatar{33, "Kuruk"},
}
filteredAvatars := goset.Uniq(avatars).([]Avatar)
```

- **Handy**

  One Line Style calling design, aims to be **developer friendly**.   
  But not enough shit are given on the performance and mathmatical rigour.

Think about how many times you want to tell if a element is in a slice, you have to wrote like this:

```go
found := false
for _, e := range records {
        if e == theRecord {
                found = true
                break
        }
}
if found {
        //your code
}
```

Now you can do this in just one line:

```go
if goset.IsIncluded(records, theRecord) {
        //your code
}
```


## Useage

#### Detections

###### 1. IsUniq

```go
a := []int{1, 2, 3, 4, 4, 2, 3, 3, 4, 4}
ua := goset.Uniq(a).([]int)
// result: ok = false
```

###### 2. IsEqual

```go
a := []int{1, 2, 3}
b := []int{2, 1, 3}
ok := goset.IsEqual(a, b)
// result: ok = true
```

###### 3. IsIncluded

```go
a := []int{1, 2, 3, 4}
ok := goset.IsIncluded(a, 1)
// result: ok = true
```

###### 4. IsSubset

```go
a := []int{1, 2, 3, 4}
a1 := []int{1, 2, 3}
ok := goset.IsSubset(a1, a)
// result: ok = true
```

###### 5. IsSuperset

```go
a := []int{1, 2, 3, 4}
a1 := []int{1, 2, 3}
ok := goset.IsSuperset(a, a1)
// result: ok = true
```


#### Operations

###### 1. Uniq

```go
a := []int{1, 2, 3, 4, 4, 2, 3, 3, 4, 4}
ua := goset.Uniq(a).([]int)
// result: ua = []int{1, 2, 3, 4}
```

###### 2. Intersect 

```go
a1 := []int{1, 2, 3, 4}
b1 := []int{3, 4, 5, 6}
c1 := goset.Intersect(a1, b1).([]int)
// result: c1 = []int{3, 4}
```

###### 3. Union

```go
a1 := []int{1, 2, 3, 4}
b1 := []int{3, 4, 5, 6}
c1 := goset.Union(a1, b1).([]int)
// result: c1 = []int{1, 2, 3, 4, 5, 6}
```

###### 4. Difference

```go
a1 := []int{1, 2, 3, 4}
b1 := []int{3, 4, 5, 6}
_, _, c1, d1 := goset.Difference(a1, b1)
// result: c1 = []int{1, 2}
//         d1 = []int{5, 6}
```

###### 5. AddElement

```go
a := []int{1, 2, 3, 4}
a = goset.AddElement(a, 5).([]int)
// result: a = []int{1, 2, 3, 4, 5}
```

###### 6. AddElements

```go
a := []int{1, 2, 3, 4}
a = goset.AddElements(a, []int{5, 6}).([]int)
// result: a = []int{1, 2, 3, 4, 5, 6}
```

###### 7. RemoveElement

```go
a := []int{1, 2, 3, 4}
a = goset.RemoveElement(a, 4).([]int)
// result: a = []int{1, 2, 3}
```

###### 8. RemoveElements

```go
a := []int{1, 2, 3, 4}
a = goset.RemoveElements(a, []int{3, 4}).([]int)
// result: a = []int{1, 2}
```

###### 9. Map

```go
x := []int{1, 2, 3, 4}
y := goset.Map(a, func(i int) string {
       return "(" + strconv.Itoa(i) + ")"
}, []string{}).([]string)
// result: y = []string{"(1)", "(2)", "(3)", "(4)"}
```

###### 10. Reorder

```go
type Cat struct {
        Id   string
        Name string
}

cats := []*Cat{
        &Cat{
                1,
                "Tom",
        },
        &Cat{
                2,
                "Jerry",
        },
        &Cat{
                3,
                "HeiMao",
        },
        &Cat{
                4,
                "Coffee",
        },
}
order := []int{3, 1, 4, 2}
cats = goset.Reorder(order, cats, "Id").([]*Cat)
// result: 
cats = []*Cat{
        &Cat{
                3,
                "HeiMao",
        },
        &Cat{
                1,
                "Tom",
        },
        &Cat{
                4,
                "Coffee",
        },
        &Cat{
                2,
                "Jerry",
        },
}
```

## License

Go Set is released under the [WTFPL License](http://www.wtfpl.net/txt/copying).


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/athom/goset/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

