## Why

The `tencentcloud_dts_sync_jobs` data source reads sync job details from the DTS cloud API `DescribeSyncJobs`. The cloud API already returns `Database.SchemaMode` (the schema selection mode, `All` or `Partial`, needed by PostgreSQL and SQL Server links) inside `Objects.Databases[]`, but the Terraform data source schema does not expose it. Users querying sync jobs via this data source therefore cannot read the schema selection mode of each database object. This change adds the missing `schema_mode` computed field so the value is available to Terraform configurations.

## What Changes

- Add a new computed field `schema_mode` (string) inside the `list.objects.databases[]` nested block of the `tencentcloud_dts_sync_jobs` data source schema.
- In the data source Read handler, populate `schema_mode` from `Database.SchemaMode` when the cloud API returns a non-nil value (skipping the set when nil, consistent with the existing nil-guard pattern used for sibling fields).
- Update the data source documentation (`data_source_tc_dts_sync_jobs.md`) to document the new `schema_mode` field.

## Capabilities

### New Capabilities
- `dts-sync-jobs-schema-mode`: Defines that the `tencentcloud_dts_sync_jobs` data source SHALL expose a computed `schema_mode` field inside the `list.objects.databases[]` block, mapping to the cloud API `Database.SchemaMode` returned by `DescribeSyncJobs`.

### Modified Capabilities
<!-- (no existing capability spec to modify; the existing `dts-sync-config-schema-mode` capability covers the `tencentcloud_dts_sync_config` resource, which is a distinct resource from this data source) -->

## Impact

- **Affected code**: `tencentcloud/services/dts/data_source_tc_dts_sync_jobs.go` (schema definition + Read handler flattening logic).
- **Affected docs**: `tencentcloud/services/dts/data_source_tc_dts_sync_jobs.md` (regenerated via `make doc` during finalization).
- **Tests**: `tencentcloud/services/dts/data_source_tc_dts_sync_jobs_test.go` (add a mock-based unit test verifying the `schema_mode` field is read from the cloud API response).
- **Cloud API**: No SDK changes required; the `Database.SchemaMode` field already exists in `vendor/.../dts/v20211206/models.go` (confirmed at the `Database` struct).
- **Compatibility**: Backwards compatible — the new field is computed-only and additive; existing configurations and state are not affected.
