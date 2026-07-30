## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/bh/resource_tc_bh_bind_device_resource.go` with schema definition and CRUD functions (Create/Read/Update/Delete), following the pattern of `tencentcloud_igtm_strategy` resource
- [x] 1.2 Register the resource `tencentcloud_bh_bind_device_resource` in `tencentcloud/provider.go`

## 2. Documentation

- [x] 2.1 Create resource example file `tencentcloud/services/bh/resource_tc_bh_bind_device_resource.md` with Example Usage and Import sections
- [x] 2.2 Add resource entry to `tencentcloud/provider.md`

## 3. Testing

- [x] 3.1 Create unit test file `tencentcloud/services/bh/resource_tc_bh_bind_device_resource_test.go` using gomonkey to mock cloud API calls, covering Create/Read/Update/Delete scenarios
- [x] 3.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
