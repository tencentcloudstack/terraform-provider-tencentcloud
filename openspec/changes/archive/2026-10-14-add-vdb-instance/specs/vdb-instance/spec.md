## ADDED Requirements

### Requirement: Resource Registration
The system SHALL register a new Terraform resource named `tencentcloud_vdb_instance` in the provider, implemented in `tencentcloud/services/vdb/resource_tc_vdb_instance.go`.

#### Scenario: Resource is registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_vdb_instance` SHALL be available as a managed resource backed by `vdb.ResourceTencentCloudVdbInstance()`

#### Scenario: Resource supports import
- **WHEN** a user runs `terraform import tencentcloud_vdb_instance.example <instance_id>`
- **THEN** the resource SHALL support import via `schema.ImportStatePassthrough`

### Requirement: Create VDB Instance
The system SHALL create a VDB instance by calling the `CreateInstance` cloud API with retry, and SHALL set `d.SetId()` before polling for the instance to reach `online` status.

#### Scenario: Successful creation
- **WHEN** a user creates a `tencentcloud_vdb_instance` with required fields `vpc_id`, `subnet_id`, `pay_mode`, `security_group_ids`
- **THEN** the system SHALL call `CreateInstanceWithContext` inside `resource.Retry(WriteRetryTimeout)`
- **AND** SHALL check that `response.Response.InstanceIds` is non-nil and non-empty, returning `NonRetryableError` if empty
- **AND** SHALL set `d.SetId(instanceIds[0])` before polling
- **AND** SHALL poll `DescribeInstances` until instance status is `online` or timeout

#### Scenario: API returns empty InstanceIds
- **WHEN** `CreateInstance` returns nil Response or empty InstanceIds
- **THEN** the system SHALL return a `NonRetryableError` with a descriptive message

### Requirement: Read VDB Instance
The system SHALL read a VDB instance by calling the `DescribeInstances` cloud API, and SHALL handle the case where the instance is not found by printing a `[CRUD]` log before clearing the ID.

#### Scenario: Instance exists
- **WHEN** reading an existing instance
- **THEN** the system SHALL populate all computed fields from `InstanceInfo`, performing nil checks before each `d.Set()`
- **AND** SHALL populate `networks` sub-structure from `InstanceInfo.Networks`
- **AND** SHALL populate `resource_tags` sub-structure from `InstanceInfo.ResourceTags`
- **AND** SHALL call `DescribeInstanceNodes` to populate `nodes` sub-structure (failures are warned, not fatal)

#### Scenario: Instance not found
- **WHEN** `DescribeInstances` returns an empty list
- **THEN** the system SHALL print `log.Printf("[CRUD] vdb_instance id=%s", d.Id())` before calling `d.SetId("")`

### Requirement: Update VDB Instance
The system SHALL support updating instance specifications via `ScaleUpInstance` (cpu/memory/disk_size) and `ScaleOutInstance` (worker_node_num), updating security groups via `ModifyDBInstanceSecurityGroups`, and SHALL reject changes to immutable fields.

#### Scenario: Scale up instance
- **WHEN** `cpu`, `memory`, or `disk_size` changes
- **THEN** the system SHALL call `ScaleUpInstanceWithContext` inside `resource.Retry(WriteRetryTimeout)` with `RunNow=true`
- **AND** SHALL poll until CPU, memory, and disk match target values

#### Scenario: Scale out instance
- **WHEN** `worker_node_num` changes
- **THEN** the system SHALL call `ScaleOutInstanceWithContext` inside `resource.Retry(WriteRetryTimeout)` with `RunNow=true`
- **AND** SHALL poll until ReplicaNum matches target value

#### Scenario: Update security groups
- **WHEN** `security_group_ids` changes
- **THEN** the system SHALL call `ModifyDBInstanceSecurityGroupsWithContext` with the full replacement list
- **AND** SHALL poll until security groups match target

#### Scenario: Immutable field change rejected
- **WHEN** any of `pay_mode`, `instance_name`, `pay_period`, `auto_renew`, `params`, `resource_tags`, `instance_type`, `mode`, `product_type`, `node_type` changes
- **THEN** the system SHALL return an error indicating the field cannot be changed

### Requirement: Delete VDB Instance
The system SHALL delete a VDB instance by first isolating it via `IsolateInstance`, and when `force_delete` is true, SHALL further destroy it via `DestroyInstances`.

#### Scenario: Isolate only (default)
- **WHEN** a user destroys the resource and `force_delete` is false (default)
- **THEN** the system SHALL call `IsolateInstanceWithContext` inside `resource.Retry(WriteRetryTimeout)`
- **AND** SHALL poll until instance status is `isolated`

#### Scenario: Force destroy
- **WHEN** a user destroys the resource and `force_delete` is true
- **THEN** the system SHALL isolate the instance first
- **AND** SHALL call `DestroyInstancesWithContext` inside `resource.Retry(WriteRetryTimeout)`
- **AND** SHALL poll until the instance is not found

### Requirement: Retry and Error Handling
All cloud API calls SHALL use retry mechanisms with appropriate timeouts, and SHALL handle nil pointer dereferences for response fields.

#### Scenario: API call with retry
- **WHEN** any cloud API is called
- **THEN** the call SHALL be wrapped in `resource.Retry()` with the appropriate timeout (WriteRetryTimeout for create/update/delete, ReadRetryTimeout for read)
- **AND** errors SHALL be wrapped with `tccommon.RetryError(e)`

#### Scenario: Nil response handling
- **WHEN** an API response is nil
- **THEN** the system SHALL return `NonRetryableError` with a descriptive message instead of panicking

### Requirement: Pagination Max Value
Query interfaces with pagination fields SHALL use the maximum value allowed by the API documentation.

#### Scenario: DescribeInstances pagination
- **WHEN** calling `DescribeInstances` in service layer
- **THEN** the `Limit` field SHALL be set to the maximum value (100)

#### Scenario: DescribeInstanceNodes pagination
- **WHEN** calling `DescribeInstanceNodes` in service layer
- **THEN** the `Limit` field SHALL be set to the maximum value (100)