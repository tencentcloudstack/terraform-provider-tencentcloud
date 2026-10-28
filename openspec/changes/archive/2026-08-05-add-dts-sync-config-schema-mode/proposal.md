## Why

The `tencentcloud_dts_sync_config` resource manages DTS (Data Transmission Service) sync job configuration. The underlying cloud API `ConfigureSyncJob` and `DescribeSyncJobs` already support a `SchemaMode` field on each database object (`Objects.Databases[]`) that controls the schema selection mode (`All`/`Partial`), which is required for PostgreSQL and SQL Server sync scenarios. The current Terraform resource schema does not expose this field, so users cannot configure or read back the schema selection mode per database. Adding this optional field keeps the resource feature-complete with the cloud API and unblocks PG/SQLServer sync object configuration.

## What Changes

- Add a new optional top-level-nested field `schema_mode` to the `objects.databases` block of the `tencentcloud_dts_sync_config` resource (`Optional`, `TypeString`), mapping to the cloud API `Database.SchemaMode`.
- Populate `schema_mode` in the resource Update (`ConfigureSyncJob`) path by setting `database.SchemaMode` from the Terraform state value.
- Read `schema_mode` back in the resource Read (`DescribeSyncJobs`) path by setting `databasesMap["schema_mode"]` from `databases.SchemaMode` when non-nil.
- No breaking changes: the new field is `Optional` and backwards compatible with existing configurations and state.

## Capabilities

### New Capabilities
<!-- Capabilities being introduced. Replace <name> with kebab-case identifier (e.g., user-auth, data-export, api-rate-limiting). Each creates specs/<name>/spec.md -->
- `dts-sync-config-schema-mode`: Adds the ability to configure and read the per-database schema selection mode (`SchemaMode`) on the `tencentcloud_dts_sync_config` resource.

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation).
     Only list here if spec-level behavior changes. Each needs a delta spec file.
     Use existing spec names from openspec/specs/. Leave empty if no requirement changes. -->

## Impact

- Affected code:
  - `tencentcloud/services/dts/resource_tc_dts_sync_config.go` (schema definition, Read, Update)
  - `tencentcloud/services/dts/resource_tc_dts_sync_config_test.go` (test coverage for the new field)
  - `tencentcloud/services/dts/resource_tc_dts_sync_config.md` (documentation, auto-generated argument reference will include the new field)
- Cloud APIs: `ConfigureSyncJob` (write path), `DescribeSyncJobs` (read path) — both already expose `Database.SchemaMode` in the vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206`).
- Dependencies: None new; uses the vendored SDK only.
- Backwards compatible: existing configurations and state are unaffected.
