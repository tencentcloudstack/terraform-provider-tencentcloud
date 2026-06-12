## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation.go` with schema definition (cluster_id, port, security_group_ids, flow_id) and CRUD functions (Create calls OpenClusterReadOnlyInstanceGroupAccess + polls DescribeFlow; Read/Delete return nil)
- [x] 1.2 Register the resource `tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation` in `tencentcloud/provider.go`
- [x] 1.3 Add the resource entry to `tencentcloud/provider.md`

## 2. Documentation

- [x] 2.1 Create resource documentation file `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation.md` with one-line description and Example Usage section

## 3. Unit Tests

- [x] 3.1 Create unit test file `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation_test.go` using gomonkey to mock SDK client methods and verify Create logic with `go test -gcflags=all=-l`
