## ADDED Requirements

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