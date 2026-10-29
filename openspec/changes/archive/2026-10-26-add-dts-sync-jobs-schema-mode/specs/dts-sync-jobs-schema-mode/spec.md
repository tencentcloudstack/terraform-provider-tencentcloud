## ADDED Requirements

### Requirement: DTS sync jobs data source exposes per-database schema mode

The `tencentcloud_dts_sync_jobs` data source SHALL expose a computed `schema_mode` string field inside the `list[].objects[].databases[]` nested block. The field SHALL map to the cloud API `Database.SchemaMode` value returned by the `DescribeSyncJobs` API. The field SHALL be computed and backwards compatible with existing configurations and state.

#### Scenario: Reading schema_mode from the cloud API

- **WHEN** the provider reads the `tencentcloud_dts_sync_jobs` data source via `DescribeSyncJobs` and a returned `Database.SchemaMode` value is non-nil
- **THEN** the provider SHALL set `list[].objects[].databases[].schema_mode` in Terraform state to the returned value

#### Scenario: schema_mode is absent in the cloud API response

- **WHEN** the provider reads the data source and the returned `Database.SchemaMode` is nil
- **THEN** the provider SHALL NOT set `schema_mode` in state for that database entry, leaving existing state untouched
