## Context

The DLC (Data Lake Compute) product exposes `GenerateInternalTable` to create a managed/internal table (Iceberg, Hive, etc.) and `DeleteTable` / `AlterTableComment` / `DescribeTable` to manage it. There is no single "table id" returned by the create API; a table is uniquely identified within a datasource by `DatabaseName` + `TableName` (and optionally `DatasourceConnectionName`). The provider already hosts DLC resources under `tencentcloud/services/dlc/` (e.g. `tencentcloud_dlc_data_engine`, `tencentcloud_dlc_user`), so the new resource follows the same package and the established `tencentcloud_igtm_strategy` code style.

All four APIs are synchronous: `GenerateInternalTable` returns the generated SQL (`Execution.SQL`) and an `IsTIcebergSql` flag, but no task/job id. No async polling is required.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_dlc_internal_table` resource with full CRUD: create via `GenerateInternalTable`, read via `DescribeTable`, update via `AlterTableComment`, delete via `DeleteTable`.
- Surface every cloud-API parameter mapped in the proposal as a Terraform schema field, including the deeply nested smart-governance policy tree (`smart_policy` → `base_info` / `policy` → `resources` / `written` / `lifecycle` / `index` / `change_table` / `table_expiration`).
- Support import using the composite `database_name#table_name` id.
- Keep the implementation consistent with provider conventions: `tccommon` retry on reads, nil-checks before `d.Set`, `NonRetryableError` on empty create results, composite id with `tccommon.FILED_SP`.

**Non-Goals:**
- No support for altering columns/partitions/properties after creation. `AlterTableComment` only updates `TableBaseInfo` (comment/type/format/etc.); structural changes (columns, partitions, properties, smart policy) are therefore `ForceNew` to avoid drift, because no cloud API exists to update them in place. (This matches the code-gen requirement: only fields present in the update API are mutable; the rest are immutable.)
- No async-job polling (no async APIs involved).
- No data-source variant in this change.

## Decisions

### 1. Composite resource ID = `database_name#table_name`
**Rationale**: `GenerateInternalTable` does not return an opaque id; `DescribeTable` queries by `DatabaseName` + `TableName`. Joining them with `tccommon.FILED_SP` (`#`) gives a stable, import-friendly id.
**Alternative considered**: Including `datasource_connection_name` in the id. Rejected because it is optional and frequently empty; keeping the id to the two required keys is simpler and matches how `DescribeTable` is typically called. `datasource_connection_name` remains a regular (non-id) schema field passed to each API.

### 2. Update scope limited to `AlterTableComment`-editable fields
**Rationale**: The only update API is `AlterTableComment`, which accepts a `TableBaseInfo`. Only the scalar `TableBaseInfo` fields that the cloud API actually honors for an alter are mutable in Update; everything else (columns, partitions, properties, upsert_keys, smart_policy, primary_keys) is `ForceNew`. Per the code-gen rules, fields not supported by the update API must not appear in the update path.
**Alternative considered**: Calling `GenerateInternalTable` again to "recreate". Rejected — that would be delete+create, which Terraform already does via `ForceNew`.

### 3. Read populates read-only computed fields from `DescribeTable`
`DescribeTable` returns extra computed-only fields (`location`, `modified_time`, `create_time`, `input_format`, `storage_size`, `record_count`, `map_materialized_view_name`, `heat_value`, `input_format_short`, per-column `nullable`/`create_time`/`modified_time`/`data_mask_strategy_info`/`type_text`, per-partition `create_time`, and the create-response `execution`/`is_t_iceberg_sql`). These are exposed as `Computed` schema attributes so users can reference them, but they are never sent on create/update.

### 4. Schema structure follows the JsonPath mapping (flattened, not re-nested)
The `table_base_info` object is flattened to top-level scalars (`database_name`, `table_name`, `datasource_connection_name`, `table_comment`, `type`, `table_format`, `user_alias`, `user_sub_uin`, `govern_policy.*`, `db_govern_policy_is_disable`, `smart_policy.*`, `primary_keys`). `columns`, `partitions`, `properties`, `upsert_keys` are top-level lists. This matches the provider convention of flattening nested request objects and avoids redundant wrapper blocks. The deprecated `govern_policy` / `db_govern_policy_is_disable` fields are still exposed (optional) because the cloud API accepts them.

### 5. Smart-governance policy tree modeled with `TypeList`/`MaxItems: 1` blocks
`smart_policy`, `base_info`, `policy`, `resources` (list, since the API uses `[]*ResourceInfo`), `favor` (list), `resource_conf`, `written`, `advance_policy`, `sort_orders` (list), `lifecycle`, `index`, `change_table`, `table_expiration` are nested `TypeList` blocks with `MaxItems: 1` where the API field is a single struct, and plain `TypeList` where the API field is a slice. Required sub-fields (`base_info.uin`, `table_expiration.enabled`, `table_expiration.expiration`) are `Required` inside their blocks.

## Risks / Trade-offs

- **[Risk] Structural changes require recreation** → Mitigation: documented as `ForceNew` for columns/partitions/properties/smart_policy/primary_keys/upsert_keys; the `AlterTableComment` path handles only the mutable `TableBaseInfo` scalars. Acceptable because the cloud API offers no in-place column/partition update.
- **[Risk] `datasource_connection_name` not in id may collide across datasources** → Mitigation: rare in practice; documented in the `.md` import section that the id is `database_name#table_name` and that cross-datasource tables should rely on distinct database names. If needed later, the id scheme can be extended without breaking existing state (the two-part id remains a valid prefix).
- **[Risk] Deprecated fields (`govern_policy`, `db_govern_policy_is_disable`, `drop_table`) may be removed by the cloud API** → Mitigation: kept optional and clearly documented as deprecated; their absence from a response is handled by nil-checks before `d.Set`.
- **[Trade-off] Large nested schema** → The smart-policy tree is deep, but flattening would lose the API's structure and make the HCL harder to write. Nested blocks are the lesser evil and match the JsonPath mapping exactly.
