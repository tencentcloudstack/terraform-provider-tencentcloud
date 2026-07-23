## 1. Service Layer Changes

- [x] 1.1 Add `qualifier` field to `scfFunctionInfo` struct in `service_tencentcloud_scf.go`
- [x] 1.2 Add `qualifier` parameter to `DescribeFunction` method and pass it to `GetFunction` API request
- [x] 1.3 Add `qualifier` parameter to `DeleteFunction` method and pass it to `DeleteFunction` API request
- [x] 1.4 Add `qualifier` parameter to `CreateTriggers` method and pass it to `CreateTrigger` API request
- [x] 1.5 Add `qualifier` parameter to `DeleteTriggers` method and pass it to `DeleteTrigger` API request

## 2. Resource Schema and CRUD Changes

- [x] 2.1 Add `qualifier` schema field (TypeString, Optional, Computed) to `ResourceTencentCloudScfFunction()` schema
- [x] 2.2 Update `resourceTencentCloudScfFunctionCreate` to populate `qualifier` in `scfFunctionInfo` and pass to `CreateTriggers`
- [x] 2.3 Update `resourceTencentCloudScfFunctionRead` to pass qualifier to `DescribeFunction`, read `response.Qualifier` and `response.Triggers[].Qualifier`
- [x] 2.4 Update `resourceTencentCloudScfFunctionUpdate` to pass qualifier to `CreateTriggers` and `DeleteTriggers`
- [x] 2.5 Update `resourceTencentCloudScfFunctionDelete` to pass qualifier to `DeleteFunction`

## 3. Documentation

- [x] 3.1 Update `resource_tc_scf_function.md` with qualifier parameter usage example

## 4. Testing

- [x] 4.1 Add test cases for qualifier parameter in `resource_tc_scf_function_test.go`
- [x] 4.2 Run `go test` with `-gcflags=all=-l` on the test file to verify unit tests pass

## 5. Validation

- [x] 5.1 Verify the code compiles correctly (no syntax errors, all imports resolved)
- [x] 5.2 Verify backward compatibility — existing configurations without qualifier continue to work