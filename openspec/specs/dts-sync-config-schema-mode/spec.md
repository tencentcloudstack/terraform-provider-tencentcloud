# dts-sync-config-schema-mode Specification

## Purpose
Defines that the `tencentcloud_dts_sync_config` resource SHALL support an optional `schema_mode` field inside the `objects.databases[]` block, mapping to the cloud API `Database.SchemaMode` parameter on the `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs.

## Requirements
### Requirement: DTS sync config supports per-database schema mode

The `tencentcloud_dts_sync_config` resource SHALL expose an optional `schema_mode` string field inside the `objects.databases[]` block. The field SHALL map to the cloud API `Database.SchemaMode` parameter on both the `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs. The field SHALL be optional and backwards compatible with existing configurations and state.

#### Scenario: User configures schema_mode on a database object

- **WHEN** a user sets `objects.0.databases.0.schema_mode` to `"Partial"` (or `"All"`) in a `tencentcloud_dts_sync_config` resource configuration
- **THEN** the provider SHALL include `SchemaMode` with that value in the `ConfigureSyncJob` request `Objects.Databases[]` entry sent to the DTS cloud API

#### Scenario: schema_mode is omitted from configuration

- **WHEN** a user does not set `objects.0.databases.0.schema_mode`
- **THEN** the provider SHALL NOT set the `SchemaMode` field on the corresponding `Database` entry in the `ConfigureSyncJob` request, preserving existing behavior

#### Scenario: Reading back schema_mode from the cloud API

- **WHEN** the provider reads the `tencentcloud_dts_sync_config` resource via `DescribeSyncJobs` and the returned `Database.SchemaMode` is non-nil
- **THEN** the provider SHALL set `objects.0.databases[].schema_mode` in Terraform state to the returned value

#### Scenario: Reading back when schema_mode is absent

- **WHEN** the provider reads the resource and the returned `Database.SchemaMode` is nil
- **THEN** the provider SHALL NOT set `schema_mode` in state for that database entry, leaving existing state untouched
