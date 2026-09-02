# dlc-meta-database-resource Specification

## Purpose
TBD - created by archiving change add-dlc-meta-database-resource. Update Purpose after archive.
## Requirements
### Requirement: dlc_meta_database resource schema definition
The system SHALL define a Terraform resource named `tencentcloud_dlc_meta_database` with a schema that maps all `CreateMetaDatabase` input parameters and `DescribeDatabase` computed output fields.

#### Scenario: top-level input fields
- **WHEN** the `tencentcloud_dlc_meta_database` resource schema is defined
- **THEN** it includes `database_name` of type `schema.TypeString` with `Required: true` and `ForceNew: true`
- **AND** it includes `datasource_connection_name` of type `schema.TypeString` with `Optional: true`
- **AND** it includes `comment` of type `schema.TypeString` with `Optional: true`
- **AND** it includes `govern_policy` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`
- **AND** it includes `smart_policy` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`

#### Scenario: govern_policy block structure
- **WHEN** the `govern_policy` element schema is defined
- **THEN** it is a `schema.Resource` with `rule_type` of type `schema.TypeString` with `Optional: true`
- **AND** it has `govern_engine` of type `schema.TypeString` with `Optional: true`

#### Scenario: smart_policy block structure
- **WHEN** the `smart_policy` element schema is defined
- **THEN** it is a `schema.Resource` with `base_info` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`
- **AND** it has `policy` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`

#### Scenario: smart_policy base_info block
- **WHEN** the `base_info` element schema is defined
- **THEN** it is a `schema.Resource` with `uin` of type `schema.TypeString` with `Required: true`
- **AND** it has `policy_type`, `catalog`, `database`, `table`, `app_id` of type `schema.TypeString` with `Optional: true`

#### Scenario: smart_policy policy block
- **WHEN** the `policy` element schema is defined
- **THEN** it is a `schema.Resource` with `inherit` of type `schema.TypeString` with `Optional: true`
- **AND** it has `resources` of type `schema.TypeList` with `Optional: true`
- **AND** it has `written` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`
- **AND** it has `lifecycle` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`
- **AND** it has `index` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`
- **AND** it has `change_table` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`
- **AND** it has `table_expiration` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`

#### Scenario: smart_policy resources block
- **WHEN** the `resources` element schema is defined
- **THEN** it is a `schema.Resource` with `attribution_type`, `resource_type`, `name`, `instance` of type `schema.TypeString` with `Optional: true`
- **AND** it has `favor` of type `schema.TypeList` with `Optional: true`
- **AND** it has `status` of type `schema.TypeInt` with `Optional: true`
- **AND** it has `resource_group_name` of type `schema.TypeString` with `Optional: true`
- **AND** it has `resource_conf` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`

#### Scenario: smart_policy resources favor block
- **WHEN** the `favor` element schema is defined
- **THEN** it is a `schema.Resource` with `priority` of type `schema.TypeInt` with `Optional: true`
- **AND** it has `catalog`, `data_base`, `table` of type `schema.TypeString` with `Optional: true`

#### Scenario: smart_policy resources resource_conf block
- **WHEN** the `resource_conf` element schema is defined
- **THEN** it is a `schema.Resource` with `parallelism` of type `schema.TypeInt` with `Optional: true`

#### Scenario: smart_policy written block
- **WHEN** the `written` element schema is defined
- **THEN** it is a `schema.Resource` with `written_enable` of type `schema.TypeString` with `Optional: true`
- **AND** it has `advance_policy` of type `schema.TypeList` with `Optional: true` and `MaxItems: 1`

#### Scenario: smart_policy written advance_policy block
- **WHEN** the `advance_policy` element schema is defined
- **THEN** it is a `schema.Resource` with `compact_enable`, `delete_enable`, `cow_compact_enable`, `compact_strategy` of type `schema.TypeString` with `Optional: true`
- **AND** it has `min_input_files`, `target_file_size_bytes`, `retain_last`, `before_days`, `expired_snapshots_interval_min`, `remove_orphan_interval_min` of type `schema.TypeInt` with `Optional: true`
- **AND** it has `sort_orders` of type `schema.TypeList` with `Optional: true`

#### Scenario: smart_policy written advance_policy sort_orders block
- **WHEN** the `sort_orders` element schema is defined
- **THEN** it is a `schema.Resource` with `column`, `sort_direction`, `null_order` of type `schema.TypeString` with `Optional: true`

#### Scenario: smart_policy lifecycle block
- **WHEN** the `lifecycle` element schema is defined
- **THEN** it is a `schema.Resource` with `lifecycle_enable`, `expired_field`, `expired_field_format` of type `schema.TypeString` with `Optional: true`
- **AND** it has `expiration` of type `schema.TypeInt` with `Optional: true`
- **AND** it has `drop_table` of type `schema.TypeBool` with `Optional: true`

#### Scenario: smart_policy index block
- **WHEN** the `index` element schema is defined
- **THEN** it is a `schema.Resource` with `index_enable` of type `schema.TypeString` with `Optional: true`

#### Scenario: smart_policy change_table block
- **WHEN** the `change_table` element schema is defined
- **THEN** it is a `schema.Resource` with `data_retention_time` of type `schema.TypeInt` with `Optional: true`

#### Scenario: smart_policy table_expiration block
- **WHEN** the `table_expiration` element schema is defined
- **THEN** it is a `schema.Resource` with `enabled` of type `schema.TypeBool` with `Required: true`
- **AND** it has `expiration` of type `schema.TypeInt` with `Required: true`

#### Scenario: computed output fields
- **WHEN** the `tencentcloud_dlc_meta_database` resource schema is defined
- **THEN** it includes `batch_id` of type `schema.TypeString` with `Computed: true`
- **AND** it includes `task_id_set` of type `schema.TypeList` with `Computed: true` and string elements
- **AND** it includes `properties` of type `schema.TypeList` with `Computed: true`
- **AND** it includes `create_time`, `modified_time`, `location`, `user_alias`, `user_sub_uin`, `database_id`, `catalog_name`, `catalog_type` of type `schema.TypeString` with `Computed: true`
- **AND** it includes `is_information_schema` of type `schema.TypeBool` with `Computed: true`

### Requirement: dlc_meta_database resource ID format
The system SHALL use a composite resource ID that encodes `datasource_connection_name` and `database_name` so that Read and Delete can reconstruct cloud API requests.

#### Scenario: ID with datasource_connection_name
- **WHEN** the user provides `datasource_connection_name` during create
- **THEN** the resource ID is set to `datasource_connection_name` + `tccommon.FILED_SP` + `database_name`

#### Scenario: ID without datasource_connection_name
- **WHEN** the user does not provide `datasource_connection_name` during create
- **THEN** the resource ID is set to just `database_name`

#### Scenario: ID parsing on read/delete
- **WHEN** the Read or Delete handler parses `d.Id()`
- **THEN** it splits by `tccommon.FILED_SP`
- **AND** if the split yields two parts, the first is `datasource_connection_name` and the second is `database_name`
- **AND** if the split yields one part, it is `database_name` and `datasource_connection_name` is empty

### Requirement: dlc_meta_database create operation
The system SHALL create a DLC meta database by calling `CreateMetaDatabase` with all configured schema fields mapped to the request, then poll `DescribeDatabase` until the resource is confirmed present.

#### Scenario: build CreateMetaDatabase request from schema
- **WHEN** the Create handler is invoked
- **THEN** it builds a `CreateMetaDatabaseRequest` with `DatasourceConnectionName` from `datasource_connection_name` if set
- **AND** it sets `MetaDatabaseInfo.DatabaseName` from `database_name`
- **AND** it sets `MetaDatabaseInfo.Comment` from `comment` if set
- **AND** it sets `GovernPolicy.RuleType` and `GovernPolicy.GovernEngine` from the `govern_policy` block if set
- **AND** it sets `SmartPolicy.BaseInfo` and `SmartPolicy.Policy` from the `smart_policy` block if set

#### Scenario: async polling after create
- **WHEN** `CreateMetaDatabase` returns successfully with a non-nil `Response`
- **THEN** the system sets `d.SetId(...)` outside the retry block
- **AND** it polls `DescribeDatabase` using `resource.Retry` with `tccommon.ReadRetryTimeout` until the database is found
- **AND** if the database is not found within the timeout, the system returns a `NonRetryableError`

#### Scenario: nil response guard
- **WHEN** `CreateMetaDatabase` returns a nil `Response`
- **THEN** the system returns a `NonRetryableError` and does not set the resource ID

### Requirement: dlc_meta_database read operation
The system SHALL read the DLC meta database by calling `DescribeDatabase` and flatten the `DatabaseResponseInfo` response into the Terraform state.

#### Scenario: build DescribeDatabase request from ID
- **WHEN** the Read handler is invoked
- **THEN** it parses `d.Id()` to extract `database_name` and optionally `datasource_connection_name`
- **AND** it builds a `DescribeDatabaseRequest` with both fields

#### Scenario: flatten response into state
- **WHEN** `DescribeDatabase` returns a non-nil `DatabaseInfo` (of type `DatabaseResponseInfo`)
- **THEN** the system sets `database_name`, `comment`, `location`, `create_time`, `modified_time`, `user_alias`, `user_sub_uin`, `database_id`, `catalog_name`, `catalog_type`, `is_information_schema` from the response fields when they are non-nil
- **AND** it flattens `Properties` into the `properties` list
- **AND** it flattens `GovernPolicy` into the `govern_policy` block

#### Scenario: resource not found on read
- **WHEN** `DescribeDatabase` returns an empty or nil response
- **THEN** the system logs `[CRUD] dlc_meta_database id=<id>` for traceability
- **AND** it calls `d.SetId("")` to remove the resource from state

#### Scenario: retry on read
- **WHEN** the Read handler calls `DescribeDatabase`
- **THEN** the SDK call is wrapped in `resource.Retry` with `tccommon.ReadRetryTimeout`
- **AND** on failure the error is wrapped with `tccommon.RetryError()`

### Requirement: dlc_meta_database delete operation
The system SHALL delete the DLC meta database by calling `DeleteMetaDatabase`, then poll `DescribeDatabase` until the resource is confirmed gone.

#### Scenario: build DeleteMetaDatabase request from ID
- **WHEN** the Delete handler is invoked
- **THEN** it parses `d.Id()` to extract `database_name` and optionally `datasource_connection_name`
- **AND** it builds a `DeleteMetaDatabaseRequest` with both fields

#### Scenario: async polling after delete
- **WHEN** `DeleteMetaDatabase` returns successfully
- **THEN** the system polls `DescribeDatabase` using `resource.Retry` until the database is no longer found
- **AND** if the database is still present within the timeout, the system returns an error

### Requirement: dlc_meta_database update immutability enforcement
The system SHALL implement an Update handler that enforces immutability of all non-ID top-level schema fields, because no cloud API Update endpoint exists.

#### Scenario: immutable field changed
- **WHEN** the Update handler is invoked and any of `datasource_connection_name`, `comment`, `govern_policy`, `smart_policy` has changed
- **THEN** the system returns `fmt.Errorf("argument `%s` cannot be changed", v)` naming the changed field

#### Scenario: no changes
- **WHEN** the Update handler is invoked and no immutable fields have changed
- **THEN** the system returns the result of the Read handler without error

### Requirement: dlc_meta_database import support
The system SHALL support Terraform import via `schema.ImportStatePassthrough`.

#### Scenario: import with composite ID
- **WHEN** the user runs `terraform import tencentcloud_dlc_meta_database.example <datasource_connection_name>#<database_name>`
- **THEN** the import succeeds and the Read handler populates the state from `DescribeDatabase`

#### Scenario: import with bare database_name
- **WHEN** the user runs `terraform import tencentcloud_dlc_meta_database.example <database_name>`
- **THEN** the import succeeds and the Read handler populates the state from `DescribeDatabase` with default `datasource_connection_name`

### Requirement: dlc_meta_database provider registration
The system SHALL register `tencentcloud_dlc_meta_database` in `tencentcloud/provider.go` `ResourcesMap` and document it in `tencentcloud/provider.md`.

#### Scenario: provider.go registration
- **WHEN** the provider is initialised
- **THEN** `ResourcesMap` contains `"tencentcloud_dlc_meta_database": dlc.ResourceTencentCloudDlcMetaDatabase()`

#### Scenario: provider.md documentation
- **WHEN** `tencentcloud/provider.md` is generated
- **THEN** it includes an entry for `tencentcloud_dlc_meta_database`
