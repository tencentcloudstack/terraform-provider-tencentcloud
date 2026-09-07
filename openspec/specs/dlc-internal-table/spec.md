# dlc-internal-table Specification

## Purpose
TBD - created by archiving change add-dlc-internal-table-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource manages DLC internal table lifecycle
The provider SHALL provide a `tencentcloud_dlc_internal_table` resource that manages the full lifecycle (create, read, update, delete) of a Tencent Cloud DLC internal table via the `GenerateInternalTable`, `DescribeTable`, `AlterTableComment`, and `DeleteTable` cloud APIs (dlc v20210125).

#### Scenario: Create a new internal table
- **WHEN** a user applies a `tencentcloud_dlc_internal_table` configuration with required fields `database_name`, `table_name`, and `columns` (each column with `name` and `type`)
- **THEN** the provider SHALL call `GenerateInternalTable` with `TableBaseInfo`, `Columns`, and any provided `Partitions`/`Properties`/`UpsertKeys`, and SHALL set the Terraform resource id to `database_name#table_name` (joined with `tccommon.FILED_SP`)

#### Scenario: Read an existing internal table
- **WHEN** the provider reads the resource state
- **THEN** the provider SHALL call `DescribeTable` using `TableName`, `DatabaseName`, and (if set) `DatasourceConnectionName` derived from the composite id and schema, and SHALL populate all schema fields from the returned `Table` object after nil-checking each field

#### Scenario: Table not found during read
- **WHEN** `DescribeTable` returns an empty response (no `Table` / `TableBaseInfo`)
- **THEN** the provider SHALL log `log.Printf("[CRUD] dlc_internal_table id=%s", d.Id())` preserving the id, then call `d.SetId("")` to remove the resource from state

#### Scenario: Delete an internal table
- **WHEN** the user destroys the resource
- **THEN** the provider SHALL call `DeleteTable` with a `TableBaseInfo` populated from the composite id (`database_name`, `table_name`) and the `datasource_connection_name` if set

#### Scenario: Empty create result is rejected
- **WHEN** `GenerateInternalTable` succeeds but returns an empty response (nil `Response` or empty result)
- **THEN** the provider SHALL return a `NonRetryableError` rather than writing an empty id, after logging the current logId and `d.Id()`

### Requirement: Resource id is a composite of database and table name
The resource id SHALL be the composite `database_name#table_name` joined with `tccommon.FILED_SP`, because `GenerateInternalTable` returns no opaque id and `DescribeTable` is queried by `DatabaseName` + `TableName`.

#### Scenario: Import by composite id
- **WHEN** a user runs `terraform import tencentcloud_dlc_internal_table.foo <database_name>#<table_name>`
- **THEN** the provider SHALL split the id on `tccommon.FILED_SP` to recover `database_name` and `table_name` and SHALL populate the schema from a subsequent `DescribeTable` read

### Requirement: Update is limited to AlterTableComment-editable fields
The provider SHALL handle updates by calling `AlterTableComment` with the updated `TableBaseInfo`. Only the scalar `TableBaseInfo` fields accepted by `AlterTableComment` (`table_comment`, `type`, `table_format`, `user_alias`, `user_sub_uin`, `datasource_connection_name`, and the smart-governance policy tree rooted at `table_base_info`) SHALL be mutable in place.

#### Scenario: Update mutable table base info
- **WHEN** a user changes `table_comment`, `type`, `table_format`, `user_alias`, `user_sub_uin`, or the `govern_policy`/`smart_policy` subtree
- **THEN** the provider SHALL call `AlterTableComment` with the new `TableBaseInfo` and SHALL refresh state from `DescribeTable`

#### Scenario: Structural change forces recreation
- **WHEN** a user changes `database_name`, `table_name`, `columns`, `partitions`, `properties`, `upsert_keys`, or `primary_keys`
- **THEN** the provider SHALL treat the change as `ForceNew` (destroy + create), because no cloud API supports in-place structural updates

### Requirement: Schema exposes all mapped cloud-API parameters
The resource schema SHALL expose every cloud-API parameter mapped in the proposal: the flattened `table_base_info` scalars (`database_name`, `table_name`, `datasource_connection_name`, `table_comment`, `type`, `table_format`, `user_alias`, `user_sub_uin`, `govern_policy.*`, `db_govern_policy_is_disable`, `smart_policy.*`, `primary_keys`), the `columns` list (`name`, `type`, `comment`, `default`, `not_null`, `precision`, `scale`, `position`, `is_partition`), the `partitions` list, the `properties` list (`key`, `value`), and `upsert_keys`.

#### Scenario: Required fields enforced
- **WHEN** a user omits `database_name`, `table_name`, or `columns` (or a column's `name`/`type`), or `smart_policy.base_info.uin`, or `smart_policy.policy.table_expiration.enabled`/`expiration`
- **THEN** Terraform SHALL reject the configuration before apply because these fields are `Required`

#### Scenario: Optional fields default to unset
- **WHEN** a user omits optional fields such as `partitions`, `properties`, `upsert_keys`, `datasource_connection_name`, or any optional smart-policy sub-field
- **THEN** the provider SHALL not send those fields to the cloud API on create

### Requirement: Read populates computed-only fields from DescribeTable
The resource schema SHALL expose computed-only attributes returned by `DescribeTable` but never sent on create/update: `location`, `modified_time`, `create_time`, `input_format`, `storage_size`, `record_count`, `map_materialized_view_name`, `heat_value`, `input_format_short`, per-column `nullable`/`create_time`/`modified_time`/`data_mask_strategy_info`/`type_text`, per-partition `create_time`, and the create response `execution`/`is_t_iceberg_sql`.

#### Scenario: Computed fields populated after read
- **WHEN** `DescribeTable` returns a `Table` with `Location`, `StorageSize`, `RecordCount`, etc.
- **THEN** the provider SHALL set the corresponding computed schema fields, after nil-checking each, so they are available in state

### Requirement: Retry and error handling follow provider conventions
All cloud-API calls in Read SHALL be wrapped with `tccommon.ReadRetryTimeout` retry; failures SHALL be wrapped with `tccommon.RetryError()`. Setting ids/state SHALL happen outside the retry block. Create SHALL check the API result for emptiness and return `NonRetryableError` when empty.

#### Scenario: Transient read failure is retried
- **WHEN** `DescribeTable` fails with a transient error
- **THEN** the provider SHALL retry up to `tccommon.ReadRetryTimeout` before surfacing the error wrapped in `tccommon.RetryError()`

#### Scenario: Id/state set outside retry
- **WHEN** a Read succeeds inside the retry block
- **THEN** the provider SHALL perform `d.SetId()` and `d.Set()` calls outside the retry block, after the retry returns
