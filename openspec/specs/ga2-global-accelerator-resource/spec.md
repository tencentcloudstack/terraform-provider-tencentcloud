## ADDED Requirements

### Requirement: Create global accelerator instance

The system SHALL provide a Terraform resource `tencentcloud_ga2_global_accelerator` that creates a GA2 global accelerator instance by calling the `CreateGlobalAccelerator` API with the following parameters:
- `name` (Required, string): Instance name, max 60 bytes
- `instance_charge_type` (Optional, string, ForceNew): Billing mode, valid values: `PREPAID`, `POSTPAID`. Default: `POSTPAID`
- `description` (Optional, string): Description, max 100 bytes
- `cross_border_type` (Optional, string): Cross-border type, valid values: `HighQuality`, `Unicom`
- `cross_border_promise_flag` (Optional, bool): Cross-border service promise flag
- `tags` (Optional, map of string, ForceNew): Tag information

After successful creation, the system SHALL set the resource ID to the returned `GlobalAcceleratorId` and wait for the async task to complete by polling `DescribeTaskResult` until status is SUCCESS.

#### Scenario: Successful creation with required parameters only

- **WHEN** user provides a valid `name` in the Terraform configuration
- **THEN** the system calls `CreateGlobalAccelerator` API, sets the resource ID to the returned `GlobalAcceleratorId`, waits for the task to complete, and reads back the resource state

#### Scenario: Successful creation with all parameters

- **WHEN** user provides `name`, `instance_charge_type`, `description`, `cross_border_type`, `cross_border_promise_flag`, and `tags`
- **THEN** the system calls `CreateGlobalAccelerator` API with all provided parameters, sets the resource ID, waits for the task to complete, and reads back the resource state

#### Scenario: Creation fails with nil response

- **WHEN** the `CreateGlobalAccelerator` API returns a nil response or nil `GlobalAcceleratorId`
- **THEN** the system returns a NonRetryableError

### Requirement: Read global accelerator instance

The system SHALL read the global accelerator instance state by calling `DescribeGlobalAccelerators` API with a `global-accelerator-id` filter. The system SHALL set the following computed attributes from the response:
- `name` (string)
- `description` (string)
- `instance_charge_type` (string)
- `cross_border_type` (string)
- `create_time` (Computed, string)
- `state` (Computed, string)
- `status` (Computed, string)
- `ddos_id` (Computed, string)
- `cname` (Computed, string)

#### Scenario: Successful read

- **WHEN** the resource exists and `DescribeGlobalAccelerators` returns a matching instance
- **THEN** the system sets all computed and configured attributes from the `GlobalAcceleratorSet` response

#### Scenario: Resource not found

- **WHEN** `DescribeGlobalAccelerators` returns an empty `GlobalAcceleratorSet` for the given ID
- **THEN** the system removes the resource from state (calls `d.SetId("")`)

### Requirement: Update global accelerator instance

The system SHALL update the global accelerator instance by calling `ModifyGlobalAccelerator` API when any of the following fields change: `name`, `description`, `cross_border_type`, `cross_border_promise_flag`. After the API call, the system SHALL wait for the async task to complete.

The fields `instance_charge_type` and `tags` are ForceNew and SHALL NOT be included in the update request.

#### Scenario: Update name and description

- **WHEN** user changes `name` or `description` in the Terraform configuration
- **THEN** the system calls `ModifyGlobalAccelerator` with the new values and the `GlobalAcceleratorId`, waits for the task to complete, and reads back the resource state

#### Scenario: Update cross-border settings

- **WHEN** user changes `cross_border_type` or `cross_border_promise_flag`
- **THEN** the system calls `ModifyGlobalAccelerator` with the updated values, waits for the task to complete, and reads back the resource state

### Requirement: Delete global accelerator instance

The system SHALL delete the global accelerator instance by calling `DeleteGlobalAccelerator` API with the `GlobalAcceleratorId`. After the API call, the system SHALL wait for the async task to complete.

#### Scenario: Successful deletion

- **WHEN** user destroys the Terraform resource
- **THEN** the system calls `DeleteGlobalAccelerator` with the resource ID, waits for the task to complete

#### Scenario: Resource already deleted

- **WHEN** `DeleteGlobalAccelerator` returns a resource not found error
- **THEN** the system treats the deletion as successful (no error)

### Requirement: Import global accelerator instance

The system SHALL support importing an existing global accelerator instance using its `GlobalAcceleratorId`.

#### Scenario: Successful import

- **WHEN** user runs `terraform import tencentcloud_ga2_global_accelerator.example <global_accelerator_id>`
- **THEN** the system reads the instance state using `DescribeGlobalAccelerators` and populates all attributes

### Requirement: Register resource in provider

The system SHALL register `tencentcloud_ga2_global_accelerator` in `tencentcloud/provider.go` resource map and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource is available after registration

- **WHEN** the provider is initialized
- **THEN** `tencentcloud_ga2_global_accelerator` is available as a valid resource type
