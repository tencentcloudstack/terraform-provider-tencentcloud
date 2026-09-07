## Why

Tencent Cloud DLC (Data Lake Compute) supports creating internal tables (e.g. Iceberg/Hive managed tables) via the `GenerateInternalTable` API, but there is currently no Terraform resource to manage the full lifecycle of an internal table. Users who manage their lakehouse schema through Terraform have no way to declare DLC internal tables as code, forcing them to use the console or ad-hoc SQL. Adding a `tencentcloud_dlc_internal_table` resource closes this gap and lets DLC tables be version-controlled, drift-detected, and destroyed alongside other infrastructure.

## What Changes

- Add a new RESOURCE_KIND_GENERAL resource `tencentcloud_dlc_internal_table` under `tencentcloud/services/dlc/`.
- Implement Create (`GenerateInternalTable`), Read (`DescribeTable`), Update (`AlterTableComment`), and Delete (`DeleteTable`) handlers.
- Expose the table base info, columns, partitions, properties, upsert keys, and the (optional) smart-governance policy tree as schema fields, matching the cloud API parameter mapping.
- Register the resource in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md` + a `resource_tc_dlc_internal_table.md` example.
- Add unit tests (`resource_tc_dlc_internal_table_test.go`) using gomonkey mocks (no Terraform acc suite).

## Capabilities

### New Capabilities
- `dlc-internal-table`: Manages the full CRUD lifecycle of a DLC internal table (base info, columns, partitions, properties, smart governance policy).

### Modified Capabilities
<!-- None. Existing DLC specs (data-engine, user/work-group policy attachments) are unaffected. -->

## Impact

- **New files**:
  - `tencentcloud/services/dlc/resource_tc_dlc_internal_table.go`
  - `tencentcloud/services/dlc/resource_tc_dlc_internal_table_test.go`
  - `tencentcloud/services/dlc/resource_tc_dlc_internal_table.md`
- **Modified files**:
  - `tencentcloud/provider.go` (register resource in `resourcesMap`)
  - `tencentcloud/provider.md` (auto-generated doc entry, produced via `make doc`)
- **Cloud APIs** (dlc v20210125): `GenerateInternalTable`, `DescribeTable`, `AlterTableComment`, `DeleteTable`. All four are synchronous (no async polling required).
- **Resource ID**: composite `database_name#table_name` (joined with `tccommon.FILED_SP`), since `GenerateInternalTable` does not return a single unique ID and `DescribeTable` queries by `DatabaseName` + `TableName` (+ optional `DatasourceConnectionName`). Supports import via the composite ID.
- **Backward compatibility**: Pure addition; no existing schema or state is affected.
