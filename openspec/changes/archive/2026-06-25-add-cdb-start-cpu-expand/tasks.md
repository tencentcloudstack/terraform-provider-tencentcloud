## 1. Schema & Resource Definition

- [x] 1.1 Create `tencentcloud/services/cdb/resource_tc_cdb_start_cpu_expand_attachment.go` with the resource schema definition including all top-level fields (`instance_id`, `type`, `expand_cpu`, `auto_strategy`, `time_interval_strategy`, `period_strategy`, `async_request_id`) and nested block schemas for strategy parameters
- [x] 1.2 Add `Timeouts` block to the resource schema for async Create and Delete operations
- [x] 1.3 Add `Importer` with `schema.ImportStatePassthrough` to the resource schema

## 2. CRUD Implementation

- [x] 2.1 Implement `resourceTencentCloudCdbStartCpuExpandAttachmentCreate` function: call `StartCpuExpand` API with retry, validate response/AsyncRequestId, set resource ID to `instance_id`, and poll `DescribeCPUExpandStrategyInfo` until expansion strategy is confirmed active
- [x] 2.2 Implement `resourceTencentCloudCdbStartCpuExpandAttachmentRead` function: call `DescribeCPUExpandStrategyInfo` API with retry, handle nil response by logging and clearing ID, set all schema fields from response data including nested strategy blocks
- [x] 2.3 Implement `resourceTencentCloudCdbStartCpuExpandAttachmentUpdate` function: check `immutableArgs` (type, expand_cpu, auto_strategy, time_interval_strategy, period_strategy) and return error if any has changed — this is a CRD-only attachment resource
- [x] 2.4 Implement `resourceTencentCloudCdbStartCpuExpandAttachmentDelete` function: call `StopCpuExpand` API with retry, poll `DescribeCPUExpandStrategyInfo` until expansion strategy is confirmed removed

## 3. Service Layer

- [x] 3.1 Add helper methods in the CDB service layer for `DescribeCdbStartCpuExpandAttachmentById` to encapsulate the `DescribeCPUExpandStrategyInfo` API call and response parsing
- [x] 3.2 Add async polling helper in the CDB service layer to check async request status for Create/Delete operations

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_cdb_start_cpu_expand` resource in `tencentcloud/provider.go` with the function `cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()`
- [x] 4.2 Update `tencentcloud/provider.md` to include the `tencentcloud_cdb_start_cpu_expand` resource entry

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/cdb/resource_tc_cdb_start_cpu_expand_attachment_test.go` with gomonkey-based unit tests for Create, Read, and Delete operations
- [x] 5.2 Verify unit tests pass with `go test -gcflags=all=-l` (do NOT run go build/go vet)

## 6. Documentation

- [x] 6.1 Create `tencentcloud/services/cdb/resource_tc_cdb_start_cpu_expand_attachment.md` following gendoc format with description, example usage for all four expansion types (auto/manual/timeInterval/period), and import section
- [x] 6.2 Ensure `make doc` generates the corresponding `website/docs/r/cdb_start_cpu_expand_attachment.markdown` documentation (to be executed in the finalize phase)
