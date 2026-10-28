## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_bh_bind_device_resource` resource SHALL define the following schema fields:
- `device_id_set` (Required, ForceNew, List of Int): 资产ID集合
- `resource_id` (Required, String): 堡垒机服务ID
- `domain_id` (Optional, String): 网络域ID
- `manage_dimension` (Optional, Int): K8S集群托管账号维度，1-集群，2-命名空间，3-工作负载
- `manage_account_id` (Optional, Int): K8S集群托管账号id
- `manage_account` (Optional, Sensitive, String): K8S集群托管账号名称
- `manage_kubeconfig` (Optional, Sensitive, String): K8S集群托管账号凭证
- `namespace` (Optional, String): K8S集群托管的namespace
- `workload` (Optional, String): K8S集群托管的workload

#### Scenario: Schema validation
- **WHEN** user defines a `tencentcloud_bh_bind_device_resource` resource with `device_id_set` and `resource_id`
- **THEN** the resource SHALL be accepted by Terraform plan

#### Scenario: Missing required fields
- **WHEN** user defines a `tencentcloud_bh_bind_device_resource` resource without `device_id_set` or `resource_id`
- **THEN** Terraform SHALL return a validation error

### Requirement: Create binds devices to bastion host service
The resource SHALL call `BindDeviceResource` API with all specified parameters to bind devices to the bastion host service instance.

#### Scenario: Successful creation
- **WHEN** user applies a `tencentcloud_bh_bind_device_resource` with valid `device_id_set` and `resource_id`
- **THEN** the system SHALL call `BindDeviceResource` API with the provided parameters
- **THEN** the resource ID SHALL be set as `<device_ids_comma_separated>#<resource_id>`

#### Scenario: API returns error on create
- **WHEN** the `BindDeviceResource` API returns an error during creation
- **THEN** the system SHALL retry with `tccommon.ReadRetryTimeout` and return the error wrapped with `tccommon.RetryError()`

### Requirement: Read queries device binding status
The resource SHALL call `DescribeDevices` API with the first device ID from `device_id_set` to read the current binding state.

#### Scenario: Device is bound to service
- **WHEN** `DescribeDevices` returns a device with non-nil `Resource` field
- **THEN** the system SHALL set `resource_id` from `Device.Resource.ResourceId`
- **THEN** the system SHALL set `domain_id`, `manage_dimension`, `manage_account_id`, `namespace`, `workload` from the Device struct fields if they are non-nil

#### Scenario: Device is not bound (Resource is nil)
- **WHEN** `DescribeDevices` returns a device with nil `Resource` field
- **THEN** the system SHALL call `d.SetId("")` to remove the resource from state

#### Scenario: Device not found
- **WHEN** `DescribeDevices` returns empty `DeviceSet`
- **THEN** the system SHALL call `d.SetId("")` to remove the resource from state

### Requirement: Update modifies binding parameters
The resource SHALL call `BindDeviceResource` API with updated parameters when non-ForceNew fields change.

#### Scenario: Update resource_id
- **WHEN** user changes `resource_id` in the configuration
- **THEN** the system SHALL call `BindDeviceResource` with the new `resource_id` and existing `device_id_set`

#### Scenario: Update K8S parameters
- **WHEN** user changes `manage_dimension`, `namespace`, or `workload`
- **THEN** the system SHALL call `BindDeviceResource` with all current parameters including the changes

### Requirement: Delete unbinds devices from service
The resource SHALL call `BindDeviceResource` API with empty `ResourceId` to unbind devices.

#### Scenario: Successful deletion
- **WHEN** user destroys the `tencentcloud_bh_bind_device_resource` resource
- **THEN** the system SHALL call `BindDeviceResource` with `ResourceId` set to empty string and the original `device_id_set`

### Requirement: Import support with composite ID
The resource SHALL support import using the composite ID format `<device_ids_comma_separated>#<resource_id>`.

#### Scenario: Import existing binding
- **WHEN** user runs `terraform import tencentcloud_bh_bind_device_resource.example 123,456#bh-abc123`
- **THEN** the system SHALL parse the composite ID and read the binding state via `DescribeDevices`

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** user references `tencentcloud_bh_bind_device_resource` in their Terraform configuration
- **THEN** the provider SHALL recognize and handle the resource type
