## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config.go` with schema definition (zone_id, function_id, function_component_bindings) and CRUD functions (Create, Read, Update, Delete) following the patterns from `resource_tc_teo_function_runtime_environment.go`
- [x] 1.2 Implement Read function: call `DescribeFunctionComponentBindings` with pagination (Limit=1000), loop through pages, map response to state
- [x] 1.3 Implement Create function: call `ModifyFunctionComponentBindings` with Operation="rebind", set composite ID (zone_id#function_id), call Read
- [x] 1.4 Implement Update function: detect changes in `function_component_bindings`, call `ModifyFunctionComponentBindings` with Operation="rebind", call Read
- [x] 1.5 Implement Delete function: call `ModifyFunctionComponentBindings` with Operation="rebind" and empty binding list

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_function_component_binding` resource in `tencentcloud/provider.go`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create resource documentation file `tencentcloud/services/teo/resource_tc_teo_function_component_binding.md` with one-line description, Example Usage, and Import section

## 4. Unit Tests

- [x] 4.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config_test.go` using gomonkey to mock `ModifyFunctionComponentBindings` and `DescribeFunctionComponentBindings` APIs
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify tests pass
