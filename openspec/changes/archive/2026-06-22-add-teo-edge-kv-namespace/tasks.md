## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_edge_k_v_namespace.go` with schema definition (zone_id, namespace, remark as input; capacity, capacity_used, created_on, modified_on as computed) and CRUD functions (Create, Read, Update, Delete) following the igtm_strategy resource pattern. Use zone_id#namespace as composite ID, retry with tccommon.ReadRetryTimeout, and handle nil responses properly.
- [x] 1.2 Register `tencentcloud_teo_edge_k_v_namespace` resource in `tencentcloud/provider.go` resource map
- [x] 1.3 Add `tencentcloud_teo_edge_k_v_namespace` entry in `tencentcloud/provider.md`

## 2. Documentation

- [x] 2.1 Create resource example file `tencentcloud/services/teo/resource_tc_teo_edge_k_v_namespace.md` with Example Usage and Import sections

## 3. Unit Tests

- [x] 3.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_edge_k_v_namespace_test.go` using gomonkey to mock TEO cloud API calls, covering Create, Read, Update, and Delete operations
- [x] 3.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass
