## Why

CDB (Cloud Database) instances currently support CPU elastic expansion through the cloud API, but there is no Terraform resource to manage this capability. Users need to manage CPU elastic expansion (auto/manual/timeInterval/period strategies) through Terraform to ensure consistent infrastructure configuration and avoid manual operations that can drift from the desired state.

## What Changes

- Add a new Terraform RESOURCE_KIND_ATTACHMENT resource `tencentcloud_cdb_start_cpu_expand` to manage the binding/unbinding of CPU elastic expansion strategies to CDB instances
- Create resource file: `resource_tc_cdb_start_cpu_expand_attachment.go`
- Create test file: `resource_tc_cdb_start_cpu_expand_attachment_test.go`
- Create documentation file: `resource_tc_cdb_start_cpu_expand_attachment.md`
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- The resource will support CRUD operations:
  - **Create**: Call `StartCpuExpand` API to enable CPU elastic expansion on a CDB instance (async, requires polling)
  - **Read**: Call `DescribeCPUExpandStrategyInfo` API to query current expansion strategy
  - **Delete**: Call `StopCpuExpand` API to disable CPU elastic expansion (async, requires polling)
  - **No Update**: This is an attachment resource — all top-level fields besides `instance_id` are immutable; changes require re-creation

## Capabilities

### New Capabilities
- `cdb-start-cpu-expand-attachment`: Manages the CPU elastic expansion strategy binding for CDB instances, supporting auto/manual/timeInterval/period expansion types with corresponding strategy parameters

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- New resource registration in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- New resource implementation in `tencentcloud/services/cdb/`
- New test file in `tencentcloud/services/cdb/`
- New documentation in `tencentcloud/services/cdb/` and `website/docs/r/`
- Cloud API dependencies: `StartCpuExpand`, `DescribeCPUExpandStrategyInfo`, `StopCpuExpand` from `cdb/v20170320` SDK package
