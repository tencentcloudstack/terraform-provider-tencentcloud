## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/cynosdb/resource_tc_cynosdb_libra_db_instance_attachment.go` with schema definition (all input parameters from AddLibraDBInstances, computed attributes from DescribeLibraDBInstanceDetail, delete parameters from IsolateLibraDBCluster) and CRUD functions (Create, Read, Update with immutableArgs check, Delete)
- [x] 1.2 Implement Create function: call AddLibraDBInstances API with retry, validate response ResourceIds is not empty, set composite ID (cluster_id#instance_id), poll DescribeLibraDBInstanceDetail until instance is ready
- [x] 1.3 Implement Read function: parse composite ID, call DescribeLibraDBInstanceDetail API with retry, set all non-nil computed attributes in state, handle instance-not-found by removing from state
- [x] 1.4 Implement Update function: check all top-level parameters (except cluster_id) against immutableArgs list, return error if any immutable parameter is changed
- [x] 1.5 Implement Delete function: parse composite ID, call IsolateLibraDBCluster API with retry using cluster_id and optional isolate_reason_types/isolate_reason parameters

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_cynosdb_libra_db_instance` resource in `tencentcloud/provider.go`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create resource documentation file `tencentcloud/services/cynosdb/resource_tc_cynosdb_libra_db_instance_attachment.md` with description, example usage, and import section using composite ID format (cluster_id#instance_id)

## 4. Testing

- [x] 4.1 Create unit test file `tencentcloud/services/cynosdb/resource_tc_cynosdb_libra_db_instance_attachment_test.go` using gomonkey to mock cloud API calls, test Create/Read/Delete logic, and verify with `go test -gcflags=all=-l`
