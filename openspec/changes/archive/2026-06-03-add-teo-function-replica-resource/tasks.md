## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_function_replica.go` with Schema definition (zone_id, function_id, replica_name as Required+ForceNew; content as Required; remark as Optional) and CRUD functions (Create, Read, Update, Delete) following tencentcloud_igtm_strategy patterns
- [x] 1.2 Implement Create function: call CreateFunctionReplica API with retry, set composite ID as `zone_id#function_id#replica_name`
- [x] 1.3 Implement Read function: call DescribeFunctionReplicas API with Filters (replica-name), Limit=200, match exact replica_name from response, set content and remark fields; if not found, remove from state
- [x] 1.4 Implement Update function: call ModifyFunctionReplica API with retry when content or remark changes
- [x] 1.5 Implement Delete function: call DeleteFunctionReplica API with retry, passing replica_name as single-element list in ReplicaNames

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_function_replica` resource in `tencentcloud/provider.go`
- [x] 2.2 Add `tencentcloud_teo_function_replica` entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create resource example file `tencentcloud/services/teo/resource_tc_teo_function_replica.md` with Example Usage and Import sections

## 4. Unit Tests

- [x] 4.1 Create test file `tencentcloud/services/teo/resource_tc_teo_function_replica_test.go` with gomonkey-based unit tests covering Create, Read, Update, Delete operations
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify tests pass
