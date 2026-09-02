# mysql-clone-instance Specification

## Purpose
TBD - created by archiving change add-cdb-clone-instance (renamed to tencentcloud_mysql_clone_instance). Update Purpose after archive.
## Requirements
### Requirement: Resource schema defines clone instance parameters
The resource `tencentcloud_mysql_clone_instance` SHALL define a schema with the following top-level fields:
- `instance_id` (Required, ForceNew, TypeString): Source CDB instance ID to clone from
- `specified_rollback_time` (Optional, ForceNew, TypeString): Rollback time (yyyy-mm-dd hh:mm:ss); mutually exclusive with `specified_backup_id`
- `specified_backup_id` (Optional, ForceNew, TypeInt): Backup file ID to clone from; mutually exclusive with `specified_rollback_time`
- `uniq_vpc_id` (Optional, ForceNew, TypeString): VPC ID
- `uniq_subnet_id` (Optional, ForceNew, TypeString): Subnet ID
- `memory` (Optional, TypeInt): Instance memory in MB (updatable)
- `volume` (Optional, TypeInt): Instance disk in GB (updatable)
- `instance_name` (Optional, ForceNew, TypeString): Cloned instance name
- `security_group` (Optional, ForceNew, TypeList of TypeString): Security group IDs
- `resource_tags` (Optional, ForceNew, TypeList of tag block): Instance tags
- `cpu` (Optional, TypeInt): CPU cores (updatable)
- `protect_mode` (Optional, TypeInt): Data replication mode 0/1/2 (updatable)
- `deploy_mode` (Optional, TypeInt): Deploy mode 0/1 (updatable)
- `slave_zone` (Optional, TypeString): Slave 1 zone (updatable)
- `backup_zone` (Optional, TypeString): Slave 2 zone (updatable)
- `device_type` (Optional, TypeString): Instance type (updatable)
- `instance_nodes` (Optional, ForceNew, TypeInt): Instance node count
- `deploy_group_id` (Optional, ForceNew, TypeString): Placement group ID
- `dry_run` (Optional, ForceNew, TypeBool): Dry-run flag
- `cage_id` (Optional, ForceNew, TypeString): Financial cage ID
- `project_id` (Optional, ForceNew, TypeInt): Project ID
- `pay_type` (Optional, ForceNew, TypeString): Payment type PRE_PAID/USED_PAID
- `period` (Optional, ForceNew, TypeInt): Instance duration in months
- `cluster_topology` (Optional, TypeList, MaxItems:1): Cloud disk node topology (updatable)
- `src_region` (Optional, ForceNew, TypeString): Source instance region for cross-region clone
- `specified_sub_backup_id` (Optional, ForceNew, TypeInt): Cross-region backup ID
- `master_zone` (Optional, ForceNew, TypeString): Master zone
- `zone` (Optional, ForceNew, TypeString): Instance zone
- `fourth_zone` (Optional, TypeString): Slave 3 zone (updatable)
- `async_request_id` (Computed, TypeString): Async request ID returned by Create/Update APIs

The resource SHALL include a `Timeouts` block with Create (default 60 min), Update (default 60 min), and Delete (default 20 min) timeouts, and an `Importer` with `schema.ImportStatePassthrough`.

#### Scenario: Schema validation for clone with rollback time
- **WHEN** user sets `instance_id` and `specified_rollback_time` (without `specified_backup_id`)
- **THEN** the resource SHALL accept the configuration and call `CreateCloneInstance` with `SpecifiedRollbackTime`

#### Scenario: Schema validation for clone with backup id
- **WHEN** user sets `instance_id` and `specified_backup_id` (without `specified_rollback_time`)
- **THEN** the resource SHALL accept the configuration and call `CreateCloneInstance` with `SpecifiedBackupId`

#### Scenario: Updatable fields are not ForceNew
- **WHEN** user changes `memory`, `volume`, `cpu`, `protect_mode`, `deploy_mode`, `slave_zone`, `backup_zone`, `device_type`, `cluster_topology`, or `fourth_zone`
- **THEN** the resource SHALL trigger an in-place Update via `UpgradeDBInstance` rather than re-creation

### Requirement: Resource Create operation
The resource Create operation SHALL call the `CreateCloneInstance` cloud API with all configured parameters. Since `CreateCloneInstance` is an async interface returning `AsyncRequestId`, the Create operation SHALL:
1. Call `CreateCloneInstance` within a retry block using `tccommon.WriteRetryTimeout`
2. Validate that the response and `response.Response.AsyncRequestId` are non-nil and non-empty; if empty, return `NonRetryableError`
3. Store the `AsyncRequestId` and set the `async_request_id` schema field
4. Poll `DescribeAsyncRequestInfo` with the `AsyncRequestId` until the task status is `SUCCESS`; for `INITIAL`/`RUNNING`/`UNDEFINED`/`PAUSED` return `RetryableError`, for `FAILED`/`KILLED`/`REMOVED` return `NonRetryableError` with the task `Info` message
5. After async success, call `DescribeCloneList` with the source `InstanceId` to obtain the `DstInstanceId` (cloned instance ID) of the completed clone task
6. Set the resource ID to the `DstInstanceId`
7. Call the Read operation to populate state

#### Scenario: Successful creation with async polling
- **WHEN** `CreateCloneInstance` API succeeds and returns a valid `AsyncRequestId`, and `DescribeAsyncRequestInfo` eventually returns `SUCCESS`
- **THEN** the resource SHALL query `DescribeCloneList`, extract the `DstInstanceId`, set the resource ID to it, and run Read to populate state

#### Scenario: CreateCloneInstance returns nil response
- **WHEN** `CreateCloneInstance` API returns nil response or nil/empty `AsyncRequestId`
- **THEN** the resource SHALL return `NonRetryableError` to prevent writing an empty resource ID

#### Scenario: Async clone task fails
- **WHEN** `DescribeAsyncRequestInfo` returns status `FAILED` or `KILLED`
- **THEN** the resource SHALL return `NonRetryableError` with the task info message

### Requirement: Resource Read operation
The resource Read operation SHALL call `DescribeDBInstances` with `InstanceIds` set to `d.Id()` (the cloned instance ID) and `QueryClusterInfo` set to `true`. The Read operation SHALL:
1. Use retry with `tccommon.ReadRetryTimeout`
2. If the response items are empty (instance not found), log `log.Printf("[CRUD] mysql_clone_instance id=%s", d.Id())` before calling `d.SetId("")`
3. If the instance is found, set schema fields from the response only when the response field is non-nil: `memory`, `volume`, `cpu`, `protect_mode`, `deploy_mode`, `device_type`, `instance_name`, `zone`, `project_id`, `deploy_group_id`, `uniq_vpc_id` (from `UniqVpcId`), `uniq_subnet_id` (from `UniqSubnetId`)
4. Do NOT overwrite create-only/source fields (e.g. `instance_id`, `specified_rollback_time`, `specified_backup_id`) during Read — preserve the configured values

#### Scenario: Reading an existing cloned instance
- **WHEN** `DescribeDBInstances` returns a matching `InstanceInfo` item
- **THEN** the resource SHALL set all non-nil schema fields from the response and retain the resource ID

#### Scenario: Reading a deleted instance
- **WHEN** `DescribeDBInstances` returns an empty items list
- **THEN** the resource SHALL log `[CRUD] mysql_clone_instance id=<id>` and call `d.SetId("")`

### Requirement: Resource Update operation
The resource Update operation SHALL call the `UpgradeDBInstance` cloud API to adjust instance configuration. Since `UpgradeDBInstance` is an async interface returning `AsyncRequestId`, the Update operation SHALL:
1. Build the `UpgradeDBInstanceRequest` with `InstanceId` = `d.Id()`, and include updatable fields that changed: `memory`, `volume`, `cpu`, `protect_mode`, `deploy_mode`, `slave_zone`, `backup_zone`, `device_type`, `cluster_topology`, `fourth_zone`
2. Check `immutableArgs` for any changed field that is not updatable and not ForceNew; if found, return an error
3. Call `UpgradeDBInstance` within a retry block using `tccommon.WriteRetryTimeout`; validate response and `AsyncRequestId` are non-empty (return `NonRetryableError` if empty)
4. Poll `DescribeAsyncRequestInfo` until status is `SUCCESS`; return `RetryableError` for in-progress states and `NonRetryableError` for failure states
5. Call Read to refresh state

#### Scenario: Successful update of memory and volume
- **WHEN** user changes `memory` and `volume` and `UpgradeDBInstance` succeeds with async `SUCCESS`
- **THEN** the resource SHALL complete the Update and call Read to refresh state

#### Scenario: Changed field is not updatable
- **WHEN** user changes a field that is neither updatable via `UpgradeDBInstance` nor marked ForceNew
- **THEN** the resource SHALL return an error indicating the field is immutable

### Requirement: Resource Delete operation
The resource Delete operation SHALL destroy the cloned instance using the `OfflineIsolatedInstances` cloud API. Since `OfflineIsolatedInstances` requires the instance to be in isolated state, the Delete operation SHALL:
1. Call `IsolateDBInstance` within a retry block using `tccommon.WriteRetryTimeout` to isolate the instance
2. Poll `DescribeDBInstanceById` until the instance status reaches `MYSQL_STATUS_ISOLATED` (return `RetryableError` while status is `MYSQL_STATUS_ISOLATING`/`MYSQL_STATUS_RUNNING`)
3. Call `OfflineIsolatedInstances` with `InstanceIds` set to `d.Id()`
4. Poll `DescribeIsolatedDBInstanceById` until the instance disappears (returns nil)
5. If the instance is already gone at any polling step, return nil (already deleted)

#### Scenario: Successful deletion
- **WHEN** the instance is isolated, then offlined via `OfflineIsolatedInstances`, and the instance no longer appears in `DescribeIsolatedDBInstanceById`
- **THEN** the resource SHALL complete the Delete and remove the instance from state

#### Scenario: Instance already deleted
- **WHEN** `DescribeDBInstanceById` returns nil during isolation polling
- **THEN** the resource SHALL return nil without error (already deleted)

### Requirement: Provider registration
The provider SHALL register the `tencentcloud_mysql_clone_instance` resource in `tencentcloud/provider.go` with the factory function `cdb.ResourceTencentCloudMysqlCloneInstance()`, and SHALL add the corresponding entry to `tencentcloud/provider.md`.

#### Scenario: Resource is available in the provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_mysql_clone_instance` resource type SHALL be registered and usable in Terraform configurations

