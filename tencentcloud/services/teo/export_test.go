package teo

// export_test.go exposes internal build/flatten helpers for unit testing in
// the teo_test package. These aliases are only compiled during `go test`.

var BuildBotManagementActionOverrideFromMap = buildBotManagementActionOverrideFromMap
var FlattenBotManagementActionOverride = flattenBotManagementActionOverride
