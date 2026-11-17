## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_tdmysql_db_instance` resource SHALL define a schema that maps to the tdmysql cloud API `CreateDBInstances`, `DescribeDBInstanceDetail`, `ModifyInstanceName`, and `IsolateDBInstance` interfaces. The schema SHALL include:
- Create input fields (Required/Optional): `zone` (Required, ForceNew), `vpc_id` (Required, ForceNew), `subnet_id` (Required, ForceNew), `spec_code` (Required, ForceNew), `disk` (Required, ForceNew), `storage_node_num` (Required, ForceNew), `replications` (Required, ForceNew), `instance_name` (Required, not ForceNew), `instance_count` (Optional, ForceNew, default 1), `full_replications` (Optional, ForceNew), `create_version` (Optional, ForceNew), `resource_tags` (Optional, ForceNew, TypeList), `init_params` (Optional, ForceNew, TypeList), `time_unit` (Optional, ForceNew), `time_span` (Optional, ForceNew), `storage_node_cpu` (Optional, ForceNew), `storage_node_mem` (Optional, ForceNew), `pay_mode` (Optional, ForceNew), `mc_num` (Optional, ForceNew), `vport` (Optional, ForceNew), `zones` (Optional, ForceNew, TypeList), `auto_voucher` (Optional, ForceNew), `voucher_ids` (Optional, ForceNew, TypeList), `instance_type` (Optional, ForceNew), `storage_type` (Optional, ForceNew), `az_mode` (Optional, ForceNew), `instance_mode` (Optional, ForceNew), `template_id` (Optional, ForceNew), `sql_mode` (Optional, ForceNew), `auto_scale_config` (Optional, ForceNew, TypeList), `security_group_ids` (Optional, ForceNew, TypeList), `user_name` (Optional, ForceNew), `password` (Optional, ForceNew, Sensitive), `encryption_enable` (Optional, ForceNew)
- Create output fields (Computed): `instance_ids` (TypeList of TypeString), `flow_id`
- Read output fields (Computed): `instance_id`, `vip`, `vport`, `status`, `create_time`, `update_time`, `char_set`, `node` (TypeList), `region`, `status_desc`, `renew_flag`, `expire_at`, `isolated_at`, `zones`, `disk_usage`, `binlog_status`, `standby_flag`, `binlog_type`, `timing_modify_instance_flag`, `columnar_node_cpu`, `columnar_node_mem`, `columnar_node_num`, `columnar_node_disk`, `columnar_node_storage_type`, `columnar_node_spec_code`, `columnar_vip`, `columnar_vport`, `is_support_columnar`, `instance_category`, `is_switch_full_replications_enable`, `dumper_vip`, `dumper_vport`, `template_name`, `analysis_mode`, `analysis_relation_infos` (TypeList), `analysis_instance_info` (TypeList), `maintenance_window` (TypeList), `encryption_kms_region`

#### Scenario: Schema validation on create
- **WHEN** a user creates a `tencentcloud_tdmysql_db_instance` resource with required fields `zone`, `vpc_id`, `subnet_id`, `spec_code`, `disk`, `storage_node_num`, `replications`, and `instance_name`
- **THEN** the resource SHALL accept the configuration and proceed with creation

#### Scenario: Missing required fields
- **WHEN** a user creates the resource without any of the required fields (`zone`, `vpc_id`, `subnet_id`, `spec_code`, `disk`, `storage_node_num`, `replications`, `instance_name`)
- **THEN** Terraform SHALL report a validation error

### Requirement: Resource Create Operation
The resource SHALL call `CreateDBInstances` API to create tdmysql instances. The request SHALL include all configured Create input fields mapped to the corresponding API parameters (`Zone`, `VpcId`, `SubnetId`, `SpecCode`, `Disk`, `StorageNodeNum`, `Replications`, `InstanceCount`, `FullReplications`, `CreateVersion`, `InstanceName`, `ResourceTags`, `InitParams`, `TimeUnit`, `TimeSpan`, `StorageNodeCpu`, `StorageNodeMem`, `PayMode`, `MCNum`, `Vport`, `Zones`, `AutoVoucher`, `VoucherIds`, `InstanceType`, `StorageType`, `AZMode`, `InstanceMode`, `TemplateId`, `SQLMode`, `AutoScaleConfig`, `SecurityGroupIds`, `UserName`, `Password`, `EncryptionEnable`). Since `CreateDBInstances` is an asynchronous interface returning `FlowId`, after a successful API call the resource SHALL poll `DescribeFlow` with the returned `FlowId` until `Status` equals `success` before proceeding. The resource ID SHALL be set to the first element of the returned `InstanceIds` array.

#### Scenario: Successful creation with async flow polling
- **WHEN** the create operation is called with valid parameters
- **THEN** the `CreateDBInstances` API SHALL be called with the corresponding parameters
- **AND** the returned `FlowId` SHALL be polled via `DescribeFlow` until `Status` is `success`
- **AND** the resource ID SHALL be set to the first element of `InstanceIds`
- **AND** `instance_ids` and `flow_id` computed fields SHALL be populated

#### Scenario: Async flow polling failure
- **WHEN** the `DescribeFlow` polling returns `Status` of `failed` or `paused`
- **THEN** the create operation SHALL return an error

#### Scenario: Create with empty response
- **WHEN** the `CreateDBInstances` API returns `nil` response or empty `InstanceIds`
- **THEN** the create operation SHALL return a `NonRetryableError`

#### Scenario: Create with retry on transient error
- **WHEN** the `CreateDBInstances` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.WriteRetryTimeout`

### Requirement: Resource Read Operation
The resource SHALL call `DescribeDBInstanceDetail` API to query the tdmysql instance detail. The request SHALL include `InstanceId` set to `d.Id()`. The response fields SHALL be mapped to the corresponding schema fields. Before calling `setXX()` for each field, the resource SHALL check if the corresponding response field is `nil` and skip setting if nil. If the response or `Response` is empty, the resource SHALL first log `log.Printf("[CRUD] tdmysql_db_instance id=%s", d.Id())` and then call `d.SetId("")`.

#### Scenario: Successful read
- **WHEN** the read operation is called and the resource exists
- **THEN** the `DescribeDBInstanceDetail` API SHALL be called with `InstanceId`
- **AND** all non-nil response fields SHALL be populated to the schema

#### Scenario: Resource not found
- **WHEN** the read operation is called and the API returns empty response
- **THEN** the resource SHALL log `[CRUD] tdmysql_db_instance id=<id>` first
- **AND** then the resource ID SHALL be cleared (`d.SetId("")`)

#### Scenario: Read with retry on transient error
- **WHEN** the `DescribeDBInstanceDetail` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.ReadRetryTimeout`

### Requirement: Resource Update Operation
The resource SHALL call `ModifyInstanceName` API to update the instance name when `instance_name` changes. The request SHALL include `InstanceId` (from `d.Id()`) and `InstanceName`. All other top-level Create input fields SHALL be treated as immutable; the Update method SHALL maintain an `immutableArgs` array containing these fields, and if any of them has changed, the operation SHALL return an error.

#### Scenario: Update instance_name
- **WHEN** the `instance_name` field changes
- **THEN** the `ModifyInstanceName` API SHALL be called with `InstanceId` and `InstanceName`

#### Scenario: Immutable field change rejected
- **WHEN** any field in the immutableArgs array (e.g., `zone`, `vpc_id`, `spec_code`, `disk`, etc.) changes
- **THEN** the update operation SHALL return an error

#### Scenario: Update with retry on transient error
- **WHEN** the `ModifyInstanceName` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.WriteRetryTimeout`

### Requirement: Resource Delete Operation
The resource SHALL call `IsolateDBInstance` API to isolate the tdmysql instance. The request SHALL include `InstanceIds` set to a single-element array containing `d.Id()`. The operation SHALL verify that the returned `SuccessInstanceIds` contains the target instance ID.

#### Scenario: Successful deletion
- **WHEN** the delete operation is called
- **THEN** the `IsolateDBInstance` API SHALL be called with `InstanceIds` containing the instance ID
- **AND** the `SuccessInstanceIds` SHALL contain the target instance ID

#### Scenario: Delete with retry on transient error
- **WHEN** the `IsolateDBInstance` API call fails with a transient error
- **THEN** the operation SHALL retry with `tccommon.WriteRetryTimeout`

### Requirement: Resource Import
The resource SHALL support Terraform import via `schema.ImportStatePassthrough`. The imported ID SHALL be the `instance_id` of the tdmysql instance.

#### Scenario: Import existing resource
- **WHEN** a user imports an existing `tencentcloud_tdmysql_db_instance` resource by its instance ID
- **THEN** Terraform SHALL call the Read operation to populate the resource state

### Requirement: SDK Client Helper
The `tencentcloud/connectivity/client.go` file SHALL add the tdmysql SDK import (`tdmysqlv20211122`), a `tdmysqlv20211122Conn` struct field on `TencentCloudClient`, and a `UseTdmysqlV20211122Client()` method that returns a `*tdmysqlv20211122.Client` instance following the existing client helper pattern (e.g., `UseTdcpgClient`).

#### Scenario: Client helper available
- **WHEN** a resource CRUD function calls `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client()`
- **THEN** a non-nil `*tdmysqlv20211122.Client` SHALL be returned

### Requirement: Provider Registration
The resource SHALL be registered in `tencentcloud/provider.go` under `ResourcesMap` with key `tencentcloud_tdmysql_db_instance` and value `tdmysql.ResourceTencentCloudTdmysqlDbInstance()`. The `tencentcloud/provider.md` file SHALL be updated to include the resource in the Tdmysql section. The `tencentcloud/services/tdmysql` package SHALL be imported in `provider.go`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_tdmysql_db_instance` SHALL be available as a resource type

### Requirement: Unit Tests
The resource SHALL have unit tests using gomonkey mock approach. Tests SHALL cover create (including async flow polling), read, update, and delete operations. The mock SHALL stub the `CreateDBInstances`, `DescribeFlow`, `DescribeDBInstanceDetail`, `ModifyInstanceName`, and `IsolateDBInstance` API calls. Tests SHALL be run with `go test -gcflags=all=-l`.

#### Scenario: Unit test for create operation
- **WHEN** the create function is tested
- **THEN** it SHALL mock the `CreateDBInstances` and `DescribeFlow` API calls and verify the resource state (including `instance_ids` and `flow_id`)

#### Scenario: Unit test for read operation
- **WHEN** the read function is tested
- **THEN** it SHALL mock the `DescribeDBInstanceDetail` API call and verify the resource state

#### Scenario: Unit test for update operation
- **WHEN** the update function is tested with `instance_name` change
- **THEN** it SHALL mock the `ModifyInstanceName` API call and verify the resource state

#### Scenario: Unit test for delete operation
- **WHEN** the delete function is tested
- **THEN** it SHALL mock the `IsolateDBInstance` API call and verify the resource is removed

### Requirement: Documentation
The resource SHALL have a `.md` documentation file at `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.md` following the gendoc/README.md format, including:
- A one-sentence description mentioning the TDSQL-C for MySQL (tdmysql) product name (format: "Provides a resource to ...")
- Example Usage section with HCL configuration (using `jsonencode()` for any JSON string field values if applicable)
- Import section (as RESOURCE_KIND_GENERAL resource), noting that the instance ID is used for import

The documentation SHALL NOT include `Argument Reference` or `Attribute Reference` sections (these are auto-generated).

#### Scenario: Documentation file exists
- **WHEN** the resource is created
- **THEN** a `resource_tc_tdmysql_db_instance.md` file SHALL exist in the `tencentcloud/services/tdmysql/` directory with the required sections
