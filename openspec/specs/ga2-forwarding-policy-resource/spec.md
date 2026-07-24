# ga2-forwarding-policy-resource Specification

## Purpose
TBD - created by archiving change add-ga2-forwarding-policy-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource schema definition
The system SHALL provide a `tencentcloud_ga2_forwarding_policy` resource with the following schema fields:
- `global_accelerator_id` (TypeString, Required, ForceNew) — The global accelerator instance ID.
- `listener_id` (TypeString, Required, ForceNew) — The listener ID.
- `host` (TypeString, Required) — The domain/host for the forwarding policy.
- `forwarding_policy_id` (TypeString, Computed) — The forwarding policy ID returned by the cloud API.
- `default_host_flag` (TypeBool, Computed) — Whether this is the default host policy for the listener.

The resource SHALL declare a Timeouts block with Create/Update/Delete defaults of 5 minutes.

#### Scenario: Schema validation on create
- **WHEN** a user creates a `tencentcloud_ga2_forwarding_policy` resource without `global_accelerator_id`
- **THEN** Terraform SHALL report a validation error indicating the field is required.

#### Scenario: Schema validation on create without host
- **WHEN** a user creates a `tencentcloud_ga2_forwarding_policy` resource without `host`
- **THEN** Terraform SHALL report a validation error indicating the field is required.

### Requirement: Resource creation
The system SHALL create a forwarding policy via `CreateForwardingPolicy` API call, passing `GlobalAcceleratorId`, `ListenerId`, and `Host` from the resource configuration.

After a successful API call, the system SHALL:
1. Extract `ForwardingPolicyId` from the response and set it as a computed field.
2. Set `d.SetId()` with the 3-segment composite ID `<gaId>#<listenerId>#<policyId>`.
3. Poll `DescribeTaskResult` using the returned `TaskId` via `WaitForGa2TaskFinish` until status is `SUCCESS`.
4. Call Read to refresh the state.

#### Scenario: Successful creation
- **WHEN** a user applies a `tencentcloud_ga2_forwarding_policy` resource with valid `global_accelerator_id`, `listener_id`, and `host`
- **THEN** the system SHALL call `CreateForwardingPolicy`, set the ID, wait for task completion, and the resource SHALL be marked as created.

#### Scenario: Creation with duplicate host
- **WHEN** a user attempts to create a forwarding policy with a `host` that already exists under the same listener
- **THEN** the system SHALL return an error indicating `InvalidParameterValue.ForwardingPolicyHostConflict`.

#### Scenario: Creation on transport layer listener
- **WHEN** a user attempts to create a forwarding policy under a TCP/UDP (non-HTTP/HTTPS) listener
- **THEN** the system SHALL return an error indicating `UnsupportedOperation.TransportLayerUnsupportedOperateForwardingPolicy`.

#### Scenario: Creation response nil check
- **WHEN** `CreateForwardingPolicy` returns a response where `Response.ForwardingPolicyId` is nil or empty
- **THEN** the system SHALL return a `NonRetryableError` to avoid writing an empty resource ID into state.

### Requirement: Resource read
The system SHALL read a forwarding policy by calling `DescribeForwardingPolicy` with `GlobalAcceleratorId` and `ListenererId` from the parsed composite ID, then paginating (Limit=100) to find the entry whose `ForwardingPolicyId` matches the 3rd segment of `d.Id()`.

On successful match, the system SHALL:
1. Set `global_accelerator_id`, `listener_id`, `host`, `forwarding_policy_id` from the matched `ForwardingPolicySet` entry.
2. Set `default_host_flag` from the `DefaultHostFlag` field (if non-nil).

If the policy is not found, the system SHALL log the current id and call `d.SetId("")` to signal the resource has been deleted out of band.

#### Scenario: Successful read
- **WHEN** a `tencentcloud_ga2_forwarding_policy` resource exists and `DescribeForwardingPolicy` returns a matching entry
- **THEN** the system SHALL populate all schema fields from the response.

#### Scenario: Read not found
- **WHEN** a `tencentcloud_ga2_forwarding_policy` resource exists in state but the cloud API returns no matching entry
- **THEN** the system SHALL log the missing resource and call `d.SetId("")` to remove it from state.

### Requirement: Read handles ResourceNotFound gracefully
The `tencentcloud_ga2_forwarding_policy` resource Read function SHALL check whether the error returned by the `DescribeForwardingPolicy` API is a `ResourceNotFound` error. When the resource does not exist, the Read function SHALL log a warning, clear the Terraform resource ID, and return nil (instead of returning an error).

#### Scenario: Forwarding policy deleted outside Terraform
- **WHEN** the Read function calls `DescribeGa2ForwardingPolicyById` and receives a `TencentCloudSDKError` with code `"ResourceNotFound"` for an existing resource (not a new resource being created)
- **THEN** the Read function logs a WARN-level message containing the resource ID, clears the resource ID via `d.SetId("")`, and returns nil

#### Scenario: Forwarding policy exists normally
- **WHEN** the Read function calls `DescribeGa2ForwardingPolicyById` and receives a valid response without error
- **THEN** the Read function sets all resource attributes from the response and returns nil (existing behavior unchanged)

#### Scenario: Non-ResourceNotFound error
- **WHEN** the Read function calls `DescribeGa2ForwardingPolicyById` and receives an error that is NOT a `ResourceNotFound` error
- **THEN** the Read function returns the error to Terraform (existing behavior unchanged)

### Requirement: Resource update
The system SHALL update a forwarding policy via `ModifyForwardingPolicy` API call, passing `GlobalAcceleratorId`, `ListenererId`, `ForwardingPolicyId`, and the updated `Host`.

After a successful API call, the system SHALL:
1. Poll `DescribeTaskResult` using the returned `TaskId` via `WaitForGa2TaskFinish` until status is `SUCCESS`.
2. Call Read to refresh the state.

#### Scenario: Successful host update
- **WHEN** a user changes the `host` field on an existing `tencentcloud_ga2_forwarding_policy` resource
- **THEN** the system SHALL call `ModifyForwardingPolicy`, wait for task completion, and the resource SHALL reflect the new host.

#### Scenario: Attempt to modify default forwarding policy
- **WHEN** a user attempts to modify or delete a forwarding policy that is the default host for a listener
- **THEN** the system SHALL return an error indicating `UnsupportedOperation.DefaultForwardingPolicyOperate`.

### Requirement: Resource deletion
The system SHALL delete a forwarding policy via `DeleteForwardingPolicy` API call, passing `GlobalAcceleratorId`, `ListenererId`, and `ForwardingPolicyId` from the parsed composite ID.

After a successful API call, the system SHALL poll `DescribeTaskResult` using the returned `TaskId` via `WaitForGa2TaskFinish` until status is `SUCCESS`.

#### Scenario: Successful deletion
- **WHEN** a user destroys a `tencentcloud_ga2_forwarding_policy` resource
- **THEN** the system SHALL call `DeleteForwardingPolicy`, wait for task completion, and remove the resource from state.

### Requirement: Resource import
The system SHALL support `terraform import` using the 3-segment composite ID `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>`.

#### Scenario: Successful import
- **WHEN** a user runs `terraform import tencentcloud_ga2_forwarding_policy.example ga-xxx#lsr-xxx#fp-xxx`
- **THEN** the system SHALL parse the composite ID, call `DescribeForwardingPolicy` to find the matching policy, and populate the state.

### Requirement: Service layer helper
The `Ga2Service` SHALL provide a method `DescribeGa2ForwardingPolicyById(ctx, gaId, listenerId, policyId)` that:
1. Builds a `DescribeForwardingPolicyRequest` with `GlobalAcceleratorId`, `ListenererId`, `Offset=0`, `Limit=100`.
2. Paginates through all results.
3. Returns the first `*ForwardingPolicySet` whose `ForwardingPolicyId` equals `policyId`.
4. Returns `(nil, nil)` if no matching policy is found.

#### Scenario: Helper finds matching policy
- **WHEN** `DescribeGa2ForwardingPolicyById` is called with valid IDs and a matching policy exists
- **THEN** it SHALL return the matching `*ForwardingPolicySet` and nil error.

#### Scenario: Helper finds no matching policy
- **WHEN** `DescribeGa2ForwardingPolicyById` is called with valid IDs but no matching policy exists
- **THEN** it SHALL return `(nil, nil)`.

### Requirement: Provider registration
The system SHALL register `tencentcloud_ga2_forwarding_policy` in `tencentcloud/provider.go` under the `ga2` namespace block, adjacent to the existing GA2 resources.

The system SHALL append `tencentcloud_ga2_forwarding_policy` to the `Global Accelerator(GA2)` Resources section in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_ga2_forwarding_policy` resource SHALL be available for use.

