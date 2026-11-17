## ADDED Requirements

### Requirement: SyncWay on cluster creation
The `tencentcloud_cynosdb_cluster` resource SHALL support an optional `SyncWay` parameter (TypeString, valid values `async`/`semisync`/`sync`, immutable after creation) that is passed to the `CreateClusters` API as `SyncWay`.

#### Scenario: Create cluster with SyncWay
- **WHEN** a user specifies `SyncWay = "async"` (or `semisync`/`sync`) in the `tencentcloud_cynosdb_cluster` resource configuration
- **THEN** the provider SHALL set `SyncWay` in the `CreateClusters` API request with the specified value

#### Scenario: Create cluster without SyncWay
- **WHEN** a user does NOT specify `SyncWay` in the `tencentcloud_cynosdb_cluster` resource configuration
- **THEN** the provider SHALL NOT set `SyncWay` in the `CreateClusters` API request (API default applies)

#### Scenario: Validation of SyncWay values
- **WHEN** a user specifies an invalid `SyncWay` value other than `async`/`semisync`/`sync`
- **THEN** the provider SHALL return a validation error indicating valid values are `async`, `semisync`, and `sync`

#### Scenario: Update SyncWay triggers immutable error
- **WHEN** a user changes `SyncWay` after creation
- **THEN** the provider SHALL return an error indicating the argument cannot be modified

### Requirement: SemiSyncTimeout on cluster creation
The `tencentcloud_cynosdb_cluster` resource SHALL support an optional `SemiSyncTimeout` parameter (TypeInt, valid range `[1000, 4294967295]` ms, immutable after creation) that is passed to the `CreateClusters` API as `SemiSyncTimeout`.

#### Scenario: Create cluster with SemiSyncTimeout
- **WHEN** a user specifies `SemiSyncTimeout` within `[1000, 4294967295]` in the `tencentcloud_cynosdb_cluster` resource configuration
- **THEN** the provider SHALL set `SemiSyncTimeout` in the `CreateClusters` API request with the specified value

#### Scenario: Create cluster without SemiSyncTimeout
- **WHEN** a user does NOT specify `SemiSyncTimeout` in the `tencentcloud_cynosdb_cluster` resource configuration
- **THEN** the provider SHALL NOT set `SemiSyncTimeout` in the `CreateClusters` API request (API default `10000` applies)

#### Scenario: Validation of SemiSyncTimeout range
- **WHEN** a user specifies a `SemiSyncTimeout` value outside `[1000, 4294967295]`
- **THEN** the provider SHALL return a validation error indicating the value must be between `1000` and `4294967295`

#### Scenario: Update SemiSyncTimeout triggers immutable error
- **WHEN** a user changes `SemiSyncTimeout` after creation
- **THEN** the provider SHALL return an error indicating the argument cannot be modified

### Requirement: Read sync configuration from DescribeClusterDetail
The `tencentcloud_cynosdb_cluster` Read function SHALL refresh `BinlogSyncWay` (Computed, string) and `SemiSyncTimeout` (Computed) from the `DescribeClusterDetail` API response field `Detail.SlaveZoneAttr[0]`.

#### Scenario: Read existing cluster with slave zone
- **WHEN** the provider reads an existing `tencentcloud_cynosdb_cluster` resource whose `DescribeClusterDetail` response contains a non-empty `SlaveZoneAttr` list
- **THEN** `BinlogSyncWay` SHALL be set from `SlaveZoneAttr[0].BinlogSyncWay`
- **AND** `SemiSyncTimeout` SHALL be set from `SlaveZoneAttr[0].SemiSyncTimeout`

#### Scenario: Read cluster without slave zone
- **WHEN** the provider reads an existing `tencentcloud_cynosdb_cluster` resource whose `DescribeClusterDetail` response has an empty or nil `SlaveZoneAttr` list
- **THEN** the provider SHALL skip setting `BinlogSyncWay` and `SemiSyncTimeout` without error

#### Scenario: Import populates sync configuration
- **WHEN** a user imports an existing `tencentcloud_cynosdb_cluster` into Terraform state
- **THEN** the Read function SHALL populate `BinlogSyncWay` and `SemiSyncTimeout` from the `DescribeClusterDetail` API response
