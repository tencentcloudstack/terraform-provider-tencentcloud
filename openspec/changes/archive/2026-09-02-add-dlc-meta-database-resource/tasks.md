## 1. Resource implementation

- [x] 1.1 Create `tencentcloud/services/dlc/resource_tc_dlc_meta_database.go` (single file). Top-level layout: package + imports → `ResourceTencentCloudDlcMetaDatabase()` schema → `resourceTencentCloudDlcMetaDatabaseCreate/Read/Update/Delete` → build/flatten helpers. Do NOT add a file-level comment at the top.
- [x] 1.2 Define the schema per the spec. Top-level fields in this order: `database_name` (Required, ForceNew), `datasource_connection_name` (Optional), `comment` (Optional), `govern_policy` (Optional, TypeList, MaxItems 1), `smart_policy` (Optional, TypeList, MaxItems 1), then computed fields: `batch_id`, `task_id_set`, `properties`, `create_time`, `modified_time`, `location`, `user_alias`, `user_sub_uin`, `database_id`, `catalog_name`, `catalog_type`, `is_information_schema`. Add a `Timeouts` block with `Create/Read/Update/Delete` defaults.
- [x] 1.3 Implement the nested `govern_policy` block: `rule_type` (Optional string), `govern_engine` (Optional string).
- [x] 1.4 Implement the nested `smart_policy` block: `base_info` (TypeList, MaxItems 1) with `uin` (Required), `policy_type`, `catalog`, `database`, `table`, `app_id` (Optional strings); `policy` (TypeList, MaxItems 1) with `inherit` (Optional string), `resources` (TypeList), `written` (TypeList, MaxItems 1), `lifecycle` (TypeList, MaxItems 1), `index` (TypeList, MaxItems 1), `change_table` (TypeList, MaxItems 1), `table_expiration` (TypeList, MaxItems 1).
- [x] 1.5 Implement the `resources` sub-block: `attribution_type`, `resource_type`, `name`, `instance` (Optional strings), `favor` (TypeList), `status` (Optional int), `resource_group_name` (Optional string), `resource_conf` (TypeList, MaxItems 1 with `parallelism` Optional int). Implement `favor` sub-block: `priority` (Optional int), `catalog`, `data_base`, `table` (Optional strings).
- [x] 1.6 Implement the `written` sub-block: `written_enable` (Optional string), `advance_policy` (TypeList, MaxItems 1). Implement `advance_policy` sub-block: `compact_enable`, `delete_enable`, `cow_compact_enable`, `compact_strategy` (Optional strings), `min_input_files`, `target_file_size_bytes`, `retain_last`, `before_days`, `expired_snapshots_interval_min`, `remove_orphan_interval_min` (Optional ints), `sort_orders` (TypeList). Implement `sort_orders` sub-block: `column`, `sort_direction`, `null_order` (Optional strings).
- [x] 1.7 Implement the `lifecycle` sub-block: `lifecycle_enable`, `expired_field`, `expired_field_format` (Optional strings), `expiration` (Optional int), `drop_table` (Optional bool). Implement `index` sub-block: `index_enable` (Optional string). Implement `change_table` sub-block: `data_retention_time` (Optional int). Implement `table_expiration` sub-block: `enabled` (Required bool), `expiration` (Required int).
- [x] 1.8 Implement `Create`: parse schema into `CreateMetaDatabaseRequest` (including all nested blocks); wrap SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil `Response`; log `logId` and `d.Id()` before ID check; set `d.SetId(...)` outside retry; poll `DescribeDatabase` via `resource.Retry(tccommon.ReadRetryTimeout, ...)` until the database is found; return Read.
- [x] 1.9 Implement `Read`: parse `d.Id()` by `tccommon.FILED_SP` to extract `database_name` and optionally `datasource_connection_name`; call `DescribeDatabase` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`; on nil/empty response log `[CRUD] dlc_meta_database id=<id>` then `d.SetId("")`; flatten all computed fields with nil checks before each `d.Set(...)`.
- [x] 1.10 Implement `Update`: define `immutableArgs := []string{"datasource_connection_name", "comment", "govern_policy", "smart_policy"}`; loop and return `fmt.Errorf("argument `%s` cannot be changed", v)` on any change; if no changes, return Read.
- [x] 1.11 Implement `Delete`: parse `d.Id()`; call `DeleteMetaDatabase` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil `Response`; poll `DescribeDatabase` until the database is no longer found (or timeout).
- [x] 1.12 Add `defer tccommon.LogElapsed("resource.tencentcloud_dlc_meta_database.<op>")()` and `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function. Use `tccommon.GetLogId(tccommon.ContextNil)` for logId.
- [x] 1.13 Add `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}` to the resource definition.

## 2. Provider registration

- [x] 2.1 In `tencentcloud/provider.go`, register `"tencentcloud_dlc_meta_database": dlc.ResourceTencentCloudDlcMetaDatabase()` in `ResourcesMap`. Place it adjacent to existing `dlc` entries to keep the namespace contiguous.
- [x] 2.2 In `tencentcloud/provider.md`, add the documentation entry for `tencentcloud_dlc_meta_database` following the existing `dlc` resource entries.

## 3. Resource documentation

- [x] 3.1 Create `tencentcloud/services/dlc/resource_tc_dlc_meta_database.md`. Content sections: one-line description mentioning DLC product; Example Usage HCL showing `database_name`, `comment`, and a `govern_policy` block; Import section explaining the composite ID format (`<datasource_connection_name>#<database_name>` or bare `<database_name>`). Do NOT add `Argument Reference` or `Attribute Reference` sections (auto-generated by `make doc`).

## 4. Unit tests

- [x] 4.1 Create `tencentcloud/services/dlc/resource_tc_dlc_meta_database_test.go` using gomonkey mock (NOT the terraform test suite). Stub `CreateMetaDatabaseWithContext`, `DescribeDatabaseWithContext`, and `DeleteMetaDatabaseWithContext` on the DLC client.
- [x] 4.2 Write `TestResourceTencentCloudDlcMetaDatabaseCreate` covering: build request from schema, mock create response with BatchId/TaskIdSet, mock describe response with DatabaseResponseInfo fields, verify `d.SetId()` and state population.
- [x] 4.3 Write `TestResourceTencentCloudDlcMetaDatabaseRead` covering: ID parsing for both composite and bare formats, nil response guard (`d.SetId("")` with `[CRUD]` log), flattening of all computed fields.
- [x] 4.4 Write `TestResourceTencentCloudDlcMetaDatabaseDelete` covering: ID parsing, mock delete response, mock describe returning not-found, verify no error.
- [x] 4.5 Write `TestResourceTencentCloudDlcMetaDatabaseUpdate` covering: immutable field change returns error, no-change returns nil.

## 5. Build & verification

- [x] 5.1 Verify the resource compiles (visual check; do NOT run `go build`).
- [x] 5.2 Verify all cloud API parameter mappings are correct: CreateMetaDatabase input fields exist in the SDK request struct, DescribeDatabase input fields exist in its request struct, DeleteMetaDatabase input fields exist in its request struct.
- [x] 5.3 Verify `make doc` will be run in the finalize phase to regenerate `website/docs/r/dlc_meta_database.html.markdown` from the resource Schema/Description and the `.md` example file. Do NOT hand-edit the generated file.
