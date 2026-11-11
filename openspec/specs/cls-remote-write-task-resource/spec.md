# cls-remote-write-task-resource Specification

## Purpose
TBD - created by archiving change add-cls-remote-write-task-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource Schema Definition

The system SHALL provide a Terraform resource `tencentcloud_cls_remote_write_task` with a schema that maps to the CLS RemoteWrite Task cloud API parameters. The schema MUST include required fields for task creation and optional fields for advanced configuration.

#### Scenario: Schema contains all required creation parameters

- **WHEN** the resource schema is defined
- **THEN** the following fields MUST be present and marked as Required: `topic_id` (string), `name` (string), `target` (string), `remote_write_url` (string), `auth_type` (int), `net_type` (int)

#### Scenario: Schema contains optional configuration parameters

- **WHEN** the resource schema is defined
- **THEN** the following fields MUST be present and marked as Optional: `vpc_id` (string), `virtual_gateway_type` (int), `enable` (int)

#### Scenario: Schema contains auth_info nested block

- **WHEN** the resource schema is defined
- **THEN** an `auth_info` nested block (TypeList, MaxItems: 1) MUST be present with optional fields: `username` (string), `password` (string), `token` (string)

#### Scenario: Schema contains read-only computed fields

- **WHEN** the resource schema is defined
- **THEN** the following fields MUST be present and marked as Computed (read-only): `task_id` (string), `status` (int), `create_time` (string), `update_time` (string), `logset_id` (string)

#### Scenario: Resource supports import

- **WHEN** the resource is imported using `terraform import`
- **THEN** the import MUST accept a composite ID in the format `task_id#topic_id`

### Requirement: Create Operation

The system SHALL implement a Create operation that calls `CreateRemoteWriteTask` cloud API to create a RemoteWrite task and sets the resource ID to a composite of `task_id` and `topic_id`.

#### Scenario: Successful task creation

- **WHEN** `terraform apply` is run with valid configuration
- **THEN** the provider MUST call `CreateRemoteWriteTask` with all required and optional parameters
- **AND** the provider MUST check that the response contains a non-empty `TaskId`
- **AND** the provider MUST set `d.SetId()` to `task_id#topic_id` (using `tccommon.FILED_SP` separator)
- **AND** the provider MUST call Read to populate the full state

#### Scenario: Create returns empty TaskId

- **WHEN** `CreateRemoteWriteTask` returns an empty `TaskId`
- **THEN** the provider MUST return a `NonRetryableError` with a descriptive message
- **AND** the provider MUST NOT set an empty ID in state

#### Scenario: Create with auth_info block

- **WHEN** the configuration includes an `auth_info` block with username, password, and token
- **THEN** the provider MUST construct a `RemoteWriteAuthInfo` struct and pass it to `CreateRemoteWriteTask` via the `AuthInfo` parameter

### Requirement: Read Operation

The system SHALL implement a Read operation that calls `DescribeRemoteWriteTasks` cloud API with a `taskId` filter to retrieve the current state of the RemoteWrite task.

#### Scenario: Successful read of existing task

- **WHEN** the Read operation is invoked for an existing task
- **THEN** the provider MUST call `DescribeRemoteWriteTasks` with `Filters` containing `Key=taskId` and `Values=[task_id]`
- **AND** the provider MUST extract the task info from the first element of `Response.Infos`
- **AND** the provider MUST set all schema fields from the API response with nil checks before each `d.Set()`

#### Scenario: Read returns empty result

- **WHEN** `DescribeRemoteWriteTasks` returns an empty `Infos` list or nil response
- **THEN** the provider MUST log `[CRUD] tencentcloud_cls_remote_write_task id=<id>` before clearing the ID
- **AND** the provider MUST call `d.SetId("")` to remove the resource from state

#### Scenario: Read handles nil AuthInfo

- **WHEN** the API response contains a nil `AuthInfo` field
- **THEN** the provider MUST skip setting the `auth_info` block without causing a panic

#### Scenario: Read uses retry with ReadRetryTimeout

- **WHEN** the `DescribeRemoteWriteTasks` API call is made
- **THEN** the call MUST be wrapped in a retry loop using `tccommon.ReadRetryTimeout` as the timeout
- **AND** API errors MUST be wrapped with `tccommon.RetryError()` for retry handling

### Requirement: Update Operation

The system SHALL implement an Update operation that calls `ModifyRemoteWriteTask` cloud API to modify the RemoteWrite task configuration when any non-ForceNew field changes.

#### Scenario: Successful task modification

- **WHEN** a non-ForceNew field changes and `terraform apply` is run
- **THEN** the provider MUST call `ModifyRemoteWriteTask` with `TaskId` and `TopicId` parsed from the composite ID
- **AND** the provider MUST include all changed optional parameters in the request
- **AND** after a successful Modify, the provider MUST call Read to refresh state

#### Scenario: Update enable field

- **WHEN** the `enable` field is changed
- **THEN** the provider MUST pass the new `Enable` value to `ModifyRemoteWriteTask`

#### Scenario: Update auth_info block

- **WHEN** the `auth_info` block is changed
- **THEN** the provider MUST construct an updated `RemoteWriteAuthInfo` struct and pass it to `ModifyRemoteWriteTask`

### Requirement: Delete Operation

The system SHALL implement a Delete operation that calls `DeleteRemoteWriteTask` cloud API to remove the RemoteWrite task.

#### Scenario: Successful task deletion

- **WHEN** `terraform destroy` is run for an existing task
- **THEN** the provider MUST call `DeleteRemoteWriteTask` with `TaskId` and `TopicId` parsed from the composite ID
- **AND** after successful deletion, the provider MUST NOT return an error

#### Scenario: Delete uses retry with ReadRetryTimeout

- **WHEN** the `DeleteRemoteWriteTask` API call is made
- **THEN** the call MUST be wrapped in a retry loop using `tccommon.ReadRetryTimeout` as the timeout
- **AND** API errors MUST be wrapped with `tccommon.RetryError()` for retry handling

### Requirement: Provider Registration

The system SHALL register the `tencentcloud_cls_remote_write_task` resource in the Terraform provider.

#### Scenario: Resource registered in provider

- **WHEN** the provider is initialized
- **THEN** `tencentcloud/provider.go` MUST contain a registration entry for `tencentcloud_cls_remote_write_task`
- **AND** `tencentcloud/provider.md` MUST contain a documentation entry for the resource

### Requirement: Resource Documentation

The system SHALL provide a markdown documentation file for the `tencentcloud_cls_remote_write_task` resource.

#### Scenario: Documentation file exists with example usage

- **WHEN** the resource is implemented
- **THEN** a file `tencentcloud/services/cls/resource_tc_cls_remote_write_task.md` MUST exist
- **AND** the file MUST contain a one-line description mentioning CLS
- **AND** the file MUST contain an Example Usage section with HCL configuration
- **AND** the file MUST contain an Import section explaining the composite ID format

### Requirement: Unit Tests with Mock

The system SHALL provide unit tests for the `tencentcloud_cls_remote_write_task` resource using gomonkey to mock cloud API calls.

#### Scenario: Test file exists and covers CRUD operations

- **WHEN** the resource is implemented
- **THEN** a file `tencentcloud/services/cls/resource_tc_cls_remote_write_task_test.go` MUST exist
- **AND** the tests MUST use gomonkey to mock `CreateRemoteWriteTask`, `DescribeRemoteWriteTasks`, `ModifyRemoteWriteTask`, and `DeleteRemoteWriteTask` API calls
- **AND** the tests MUST NOT use Terraform test suite (TF_ACC)

