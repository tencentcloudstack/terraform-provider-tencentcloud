## 1. Schema Definition

- [x] 1.1 Add `schema_mode` (`Optional`, `Type: schema.TypeString`) field to the `databases` block schema inside the `objects` block of `ResourceTencentCloudDtsSyncConfig()` in `tencentcloud/services/dts/resource_tc_dts_sync_config.go`, placed alongside the existing `schema_name`/`new_schema_name`/`db_mode` fields, with a description matching the SDK doc (schema selection mode: `All`/`Partial`, used by PG and SQLServer).

## 2. CRUD Implementation

- [x] 2.1 In `resourceTencentCloudDtsSyncConfigUpdate` (Update path, `ConfigureSyncJob`), within the `objects.databases` build loop, add `if v, ok := databasesMap["schema_mode"]; ok { database.SchemaMode = helper.String(v.(string)) }` so the configured value is sent to the cloud API.
- [x] 2.2 In `resourceTencentCloudDtsSyncConfigRead` (Read path, `DescribeSyncJobs` via `DescribeDtsSyncConfigById`), within the `objects.databases` flatten loop, add `if databases.SchemaMode != nil { databasesMap["schema_mode"] = databases.SchemaMode }` so the value is read back into state only when non-nil.

## 3. Tests

- [x] 3.1 Add a mock-based (gomonkey) unit test case in `tencentcloud/services/dts/resource_tc_dts_sync_config_test.go` that verifies `schema_mode` is correctly mapped into the `ConfigureSyncJob` request `Objects.Databases[].SchemaMode` (write path) and read back from `DescribeSyncJobs` response `JobList[].Objects.Databases[].SchemaMode` (read path).

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/dts/resource_tc_dts_sync_config.md` to document the new `schema_mode` field within the `objects.databases` block (description + example usage snippet), following the existing sibling-field documentation style. (The `website/docs/r/dts_sync_config.html.markdown` is generated via `make doc` during finalization.)

## 5. Verification

- [x] 5.1 Verify the generated Go code compiles (via the dedicated build/lint flow run outside this skill) and that no existing fields were altered.
- [x] 5.2 Run `make doc` during the finalization phase to regenerate `website/docs/r/dts_sync_config.html.markdown` with the new argument included.
