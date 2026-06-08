## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_edge_k_v_put_attachment.go` with schema definition (zone_id, namespace, key, value, expiration, expiration_ttl) and CRUD functions (Create calls EdgeKVPut, Read calls EdgeKVGet, Update enforces immutableArgs, Delete calls EdgeKVDelete)
- [x] 1.2 Implement composite ID using `tccommon.FILED_SP` separator (zone_id#namespace#key), with ID parsing in Read and Delete methods

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_edge_k_v_put` resource in `tencentcloud/provider.go`
- [x] 2.2 Add `tencentcloud_teo_edge_k_v_put` entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create resource documentation file `tencentcloud/services/teo/resource_tc_teo_edge_k_v_put_attachment.md` with Example Usage and Import sections

## 4. Unit Tests

- [x] 4.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_edge_k_v_put_attachment_test.go` with gomonkey mock tests covering Create, Read, and Delete operations
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass
