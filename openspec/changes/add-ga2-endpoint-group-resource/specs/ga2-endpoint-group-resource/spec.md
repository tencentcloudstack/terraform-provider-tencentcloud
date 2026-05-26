## ADDED Requirements

### Requirement: Resource manages GA2 endpoint group lifecycle
The provider SHALL expose a resource `tencentcloud_ga2_endpoint_group` that manages a Tencent Cloud Global Accelerator 2 endpoint group via `CreateEndpointGroup`, `DescribeEndpointGroups`, `ModifyEndpointGroup`, and `DeleteEndpointGroups` APIs in the `ga2/v20250115` SDK namespace.

#### Scenario: Create provisions and waits for task success
- **WHEN** a user defines the resource with `global_accelerator_id`, `listener_id`, `endpoint_group_type`, and `endpoint_group_configuration`
- **THEN** the provider SHALL call `CreateEndpointGroup`, capture the returned `EndpointGroupId` and `TaskId`, poll `DescribeTaskResult` until `Status == "SUCCESS"`, and set the resource ID to `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`

#### Scenario: Read reflects current cloud state
- **WHEN** the provider reads the resource
- **THEN** it SHALL split the resource ID into three parts, call `DescribeEndpointGroups` filtered by `endpoint-group-id`, and map every non-nil response field into Terraform state

#### Scenario: Read clears state when not found
- **WHEN** `DescribeEndpointGroups` returns no record matching the endpoint group ID
- **THEN** the provider SHALL call `d.SetId("")` and return nil

#### Scenario: Update modifies mutable fields and waits
- **WHEN** any mutable attribute under `endpoint_group_configuration` (or top-level `endpoint_configurations`/`port_overrides`/health-check fields) changes
- **THEN** the provider SHALL build a `ModifyEndpointGroup` request from the new state, call the API, and poll `DescribeTaskResult` until `Status == "SUCCESS"`

#### Scenario: Delete removes the group and waits
- **WHEN** the resource is destroyed
- **THEN** the provider SHALL call `DeleteEndpointGroups` with the single-element `EndpointGroupIds` slice and poll `DescribeTaskResult` until `Status == "SUCCESS"`

### Requirement: Schema mirrors CreateEndpointGroup parameters
The resource schema SHALL contain fields corresponding to every parameter in `CreateEndpointGroupRequestParams`.

#### Scenario: Top-level required fields enforced
- **WHEN** a user omits `global_accelerator_id`, `listener_id`, `endpoint_group_type`, or `endpoint_group_configuration`
- **THEN** Terraform SHALL return a validation error before any API call

#### Scenario: Nested endpoint_group_configuration mirrors SDK struct
- **WHEN** the user provides an `endpoint_group_configuration` block
- **THEN** every field inside (`name`, `endpoint_group_region`, `endpoint_configurations`, `check_type`, `description`, `check_port`, `context_type`, `check_send_context`, `check_recv_context`, `enable_health_check`, `connect_timeout`, `health_check_interval`, `unhealthy_threshold`, `healthy_threshold`, `forward_protocol`, `check_domain`, `check_path`, `check_method`, `status_mask`, `port_overrides`, `isp_type`, `cipher_policy_id`) SHALL map 1:1 onto the SDK request

#### Scenario: ForceNew on identifying and immutable fields
- **WHEN** the user changes `global_accelerator_id`, `listener_id`, or `endpoint_group_type`
- **THEN** Terraform SHALL plan a destroy-and-recreate, since `ModifyEndpointGroup` does not accept these fields

### Requirement: All SDK invocations are wrapped with retry logic
Every SDK API call (Create, Read, Update, Delete, DescribeTaskResult) SHALL be invoked inside `resource.Retry(...)` using project timeout constants.

#### Scenario: Read uses ReadRetryTimeout
- **WHEN** `DescribeEndpointGroups` is invoked from `Read` or the helper
- **THEN** the call SHALL be wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)` and use `tccommon.RetryError` for retryable errors

#### Scenario: Mutating calls use WriteRetryTimeout
- **WHEN** `CreateEndpointGroup`, `ModifyEndpointGroup`, or `DeleteEndpointGroups` is invoked
- **THEN** the call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`

#### Scenario: Task polling tolerates transient failures
- **WHEN** `DescribeTaskResult` returns a transient error or non-`SUCCESS` status
- **THEN** the helper SHALL return `resource.RetryableError` so the polling loop continues until `WriteRetryTimeout*2` elapses

### Requirement: Async task polling via DescribeTaskResult
After every mutating API call, the provider SHALL block until the task associated with the returned `TaskId` reaches terminal `SUCCESS` state via `DescribeTaskResult`.

#### Scenario: Task succeeds
- **WHEN** `DescribeTaskResult` returns `Status == "SUCCESS"`
- **THEN** polling SHALL stop and the operation SHALL be considered complete

#### Scenario: Task remains in progress
- **WHEN** `DescribeTaskResult` returns any non-`SUCCESS` status
- **THEN** polling SHALL continue until the timeout elapses; on timeout the operation SHALL fail with the last observed status

### Requirement: Pagination uses API maximum page size
List/query API calls SHALL use the documented maximum `Limit` to minimize round trips.

#### Scenario: DescribeEndpointGroups uses Limit=100
- **WHEN** the service helper queries endpoint groups
- **THEN** it SHALL pass `Limit=100` (the API maximum) and iterate with `Offset` until results are exhausted

### Requirement: Defensive nil-pointer handling
The provider SHALL guard against nil pointer dereferences when reading SDK response fields.

#### Scenario: Nil response is non-retryable
- **WHEN** any SDK call returns `result == nil` or `result.Response == nil`
- **THEN** the call SHALL return `resource.NonRetryableError` with a clear message

#### Scenario: Nil field is skipped, not dereferenced
- **WHEN** the read function processes a response field that is nil
- **THEN** it SHALL skip the corresponding `d.Set(...)` call rather than dereferencing the pointer

### Requirement: Composite resource ID
The resource ID SHALL be the string `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`.

#### Scenario: ID is set on successful create
- **WHEN** `CreateEndpointGroup` succeeds and the task reaches `SUCCESS`
- **THEN** `d.SetId` SHALL be called with the three IDs joined by `#`

#### Scenario: ID is parsed on subsequent operations
- **WHEN** Read, Update, or Delete is invoked
- **THEN** the provider SHALL split `d.Id()` on `#`, validate exactly 3 segments, and use the segments for API parameters; an invalid ID format SHALL return an error
