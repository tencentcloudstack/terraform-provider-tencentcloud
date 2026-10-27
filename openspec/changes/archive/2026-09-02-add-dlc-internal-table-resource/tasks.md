## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/dlc/resource_tc_dlc_internal_table.go` with `ResourceTencentCloudDlcInternalTable()` schema: flattened `table_base_info` scalars (`database_name` Required+ForceNew, `table_name` Required+ForceNew, `datasource_connection_name` Optional, `table_comment` Optional, `type` Optional, `table_format` Optional, `user_alias` Optional, `user_sub_uin` Optional, `db_govern_policy_is_disable` Optional), `govern_policy` block (Optional, `rule_type`, `govern_engine`), `smart_policy` nested block tree (Optional; `base_info` with Required `uin`, `policy_type`, `catalog`, `database`, `table`, `app_id`; `policy` with `inherit`, `resources` list, `written`/`advance_policy`/`sort_orders`, `lifecycle`, `index`, `change_table`, `table_expiration` with Required `enabled`/`expiration`), `primary_keys` (Optional+ForceNew), `columns` (Required+ForceNew; `name`/`type` Required, `comment`/`default`/`not_null`/`precision`/`scale`/`position`/`is_partition` Optional), `partitions` (Optional+ForceNew), `properties` (Optional+ForceNew; `key`/`value` Required), `upsert_keys` (Optional+ForceNew), plus computed fields (`location`, `modified_time`, `create_time`, `input_format`, `storage_size`, `record_count`, `map_materialized_view_name`, `heat_value`, `input_format_short`, per-column computed `nullable`/`create_time`/`modified_time`/`data_mask_strategy_info`/`type_text`, per-partition `create_time`, `execution`, `is_t_iceberg_sql`)
- [x] 1.2 Implement `resourceTencentCloudDlcInternalTableCreate`: build `GenerateInternalTableRequest` from schema (`TableBaseInfo`, `Columns` as `[]*TColumn`, `Partitions` as `[]*TPartition`, `Properties`, `UpsertKeys`), call API with retry, check empty result → `NonRetryableError` (log logId + d.Id() first), set `d.SetId(database_name#table_name)` outside retry, then call Read
- [x] 1.3 Implement `resourceTencentCloudDlcInternalTableRead`: split composite id into `database_name`/`table_name`, call `DescribeTable` (with `DatasourceConnectionName` if set) inside `tccommon.ReadRetryTimeout` retry; on empty response log `[CRUD] dlc_internal_table id=%s` then `d.SetId("")`; otherwise nil-check each field before `d.Set`; populate computed fields from `TableResponseInfo`; set state outside retry
- [x] 1.4 Implement `resourceTencentCloudDlcInternalTableUpdate`: diff mutable `TableBaseInfo` scalars + `govern_policy`/`smart_policy` subtree; call `AlterTableComment` with updated `TableBaseInfo`; then call Read. `database_name`/`table_name`/`columns`/`partitions`/`properties`/`upsert_keys`/`primary_keys` are `ForceNew` so they never reach Update
- [x] 1.5 Implement `resourceTencentCloudDlcInternalTableDelete`: build `TableBaseInfo` from id + `datasource_connection_name`, call `DeleteTable` with retry
- [x] 1.6 Add `Importer: schema.ImportStatePassthrough` to support `terraform import` via composite id

## 2. Provider Registration

- [x] 2.1 Register `"tencentcloud_dlc_internal_table": dlc.ResourceTencentCloudDlcInternalTable()` in `tencentcloud/provider.go` `resourcesMap`
- [ ] 2.2 Add the resource entry to `tencentcloud/provider.md` (via `make doc` in the finalize phase; do NOT hand-edit `website/`)

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/dlc/resource_tc_dlc_internal_table.md` with one-line description (mentioning DLC), Example Usage (use `jsonencode()` for any JSON-string fields), and an Import section documenting the composite `database_name#table_name` id

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/dlc/resource_tc_dlc_internal_table_test.go` using gomonkey mocks for the DLC SDK client (no Terraform acc suite): mock `GenerateInternalTable`, `DescribeTable`, `AlterTableComment`, `DeleteTable`; cover create-then-read, update mutable fields, delete, empty-create-result error, and table-not-found read-clears-id paths

## 5. Finalize (deferred to tfpacer-finalize)

- [ ] 5.1 Run `gofmt` on changed Go files
- [ ] 5.2 Run `make doc` to generate `website/docs/` docs and update `provider.md`
- [ ] 5.3 Generate `.changelog/` entry
