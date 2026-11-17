## ADDED Requirements

### Requirement: DTS sync config supports table column mode

The `tencentcloud_dts_sync_config` resource SHALL expose an optional `column_mode` string field inside the `objects.databases.tables[]` block. The field SHALL map to the cloud API `Table.ColumnMode` parameter on both the `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs. The field SHALL be optional and backwards compatible with existing configurations and state.

#### Scenario: User configures column_mode on a table object

- **WHEN** a user sets `objects.0.databases.0.tables.0.column_mode` to `"All"` (or `"Partial"`) in a `tencentcloud_dts_sync_config` resource configuration
- **THEN** the provider SHALL include `ColumnMode` with that value in the `ConfigureSyncJob` request `Objects.Databases[].Tables[]` entry sent to the DTS cloud API

#### Scenario: column_mode is omitted from configuration

- **WHEN** a user does not set `objects.0.databases.0.tables.0.column_mode`
- **THEN** the provider SHALL NOT set the `ColumnMode` field on the corresponding `Table` entry in the `ConfigureSyncJob` request, preserving existing behavior

#### Scenario: Reading back column_mode from the cloud API

- **WHEN** the provider reads the `tencentcloud_dts_sync_config` resource via `DescribeSyncJobs` and the returned `Table.ColumnMode` is non-nil
- **THEN** the provider SHALL set `objects.0.databases[].tables[].column_mode` in Terraform state to the returned value

#### Scenario: Reading back when column_mode is absent

- **WHEN** the provider reads the resource and the returned `Table.ColumnMode` is nil
- **THEN** the provider SHALL NOT set `column_mode` in state for that table entry, leaving existing state untouched

### Requirement: DTS sync config supports table columns for partial column sync

The `tencentcloud_dts_sync_config` resource SHALL expose an optional `columns` list field inside the `objects.databases.tables[]` block, mapping to the cloud API `Table.Columns`. Each list entry SHALL be a nested object with:
- `column_name` → `Column.ColumnName`
- `new_column_name` → `Column.NewColumnName`

Both fields SHALL be optional strings. The `columns` field SHALL map to the `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs and SHALL be optional and backwards compatible.

#### Scenario: User configures columns with column name and new column name

- **WHEN** a user sets `objects.0.databases.0.tables.0.columns` to a list containing an object with `column_name = "c1"` and `new_column_name = "c1_new"` in a `tencentcloud_dts_sync_config` resource configuration
- **THEN** the provider SHALL include a `Columns` array in the `ConfigureSyncJob` request `Objects.Databases[].Tables[]` entry, where each element has `ColumnName` and `NewColumnName` set to the configured values

#### Scenario: columns is omitted from configuration

- **WHEN** a user does not set `objects.0.databases.0.tables.0.columns`
- **THEN** the provider SHALL NOT set the `Columns` field on the corresponding `Table` entry in the `ConfigureSyncJob` request, preserving existing behavior

#### Scenario: Reading back columns from the cloud API

- **WHEN** the provider reads the `tencentcloud_dts_sync_config` resource via `DescribeSyncJobs` and the returned `Table.Columns` is non-nil and non-empty
- **THEN** the provider SHALL set `objects.0.databases[].tables[].columns` in Terraform state to a list whose entries preserve `column_name` and `new_column_name` from each returned `Column` element, including only non-nil fields

#### Scenario: Reading back when columns is absent

- **WHEN** the provider reads the resource and the returned `Table.Columns` is nil or empty
- **THEN** the provider SHALL NOT set `columns` in state for that table entry, leaving existing state untouched

### Requirement: DTS sync config supports table tmp_tables

The `tencentcloud_dts_sync_config` resource SHALL expose an optional `tmp_tables` set-of-strings field inside the `objects.databases.tables[]` block. The field SHALL map to the cloud API `Table.TmpTables` parameter on both the `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs. The field SHALL be optional and backwards compatible with existing configurations and state.

#### Scenario: User configures tmp_tables on a table object

- **WHEN** a user sets `objects.0.databases.0.tables.0.tmp_tables` to `["_t1_new", "_t1_old"]` in a `tencentcloud_dts_sync_config` resource configuration
- **THEN** the provider SHALL include a `TmpTables` array with those values in the `ConfigureSyncJob` request `Objects.Databases[].Tables[]` entry sent to the DTS cloud API

#### Scenario: tmp_tables is omitted from configuration

- **WHEN** a user does not set `objects.0.databases.0.tables.0.tmp_tables`
- **THEN** the provider SHALL NOT set the `TmpTables` field on the corresponding `Table` entry in the `ConfigureSyncJob` request, preserving existing behavior

#### Scenario: Reading back tmp_tables from the cloud API

- **WHEN** the provider reads the `tencentcloud_dts_sync_config` resource via `DescribeSyncJobs` and the returned `Table.TmpTables` is non-nil and non-empty
- **THEN** the provider SHALL set `objects.0.databases[].tables[].tmp_tables` in Terraform state to the returned list of strings

#### Scenario: Reading back when tmp_tables is absent

- **WHEN** the provider reads the resource and the returned `Table.TmpTables` is nil or empty
- **THEN** the provider SHALL NOT set `tmp_tables` in state for that table entry, leaving existing state untouched

### Requirement: DTS sync config supports table_edit_mode

The `tencentcloud_dts_sync_config` resource SHALL expose an optional `table_edit_mode` string field inside the `objects.databases.tables[]` block. The field SHALL map to the cloud API `Table.TableEditMode` parameter on both the `ConfigureSyncJob` (write) and `DescribeSyncJobs` (read) APIs. The field SHALL be optional and backwards compatible with existing configurations and state.

#### Scenario: User configures table_edit_mode on a table object

- **WHEN** a user sets `objects.0.databases.0.tables.0.table_edit_mode` to `"rename"` (or `"pt"`) in a `tencentcloud_dts_sync_config` resource configuration
- **THEN** the provider SHALL include `TableEditMode` with that value in the `ConfigureSyncJob` request `Objects.Databases[].Tables[]` entry sent to the DTS cloud API

#### Scenario: table_edit_mode is omitted from configuration

- **WHEN** a user does not set `objects.0.databases.0.tables.0.table_edit_mode`
- **THEN** the provider SHALL NOT set the `TableEditMode` field on the corresponding `Table` entry in the `ConfigureSyncJob` request, preserving existing behavior

#### Scenario: Reading back table_edit_mode from the cloud API

- **WHEN** the provider reads the `tencentcloud_dts_sync_config` resource via `DescribeSyncJobs` and the returned `Table.TableEditMode` is non-nil
- **THEN** the provider SHALL set `objects.0.databases[].tables[].table_edit_mode` in Terraform state to the returned value

#### Scenario: Reading back when table_edit_mode is absent

- **WHEN** the provider reads the resource and the returned `Table.TableEditMode` is nil
- **THEN** the provider SHALL NOT set `table_edit_mode` in state for that table entry, leaving existing state untouched
