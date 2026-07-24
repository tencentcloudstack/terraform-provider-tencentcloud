## Context

The `tencentcloud_ga2_forwarding_policy` resource uses the `DescribeForwardingPolicy` API (GA2 v20250115) to read the current state of a forwarding policy. The Read function (`resourceTencentCloudGa2ForwardingPolicyRead`) currently calls the service layer `DescribeGa2ForwardingPolicyById` and, if an error is returned, propagates it directly to Terraform without checking whether the error means the resource was deleted externally.

The GA2 SDK defines `ResourceNotFound` as a known error code (see `vendor/.../ga2/v20250115/errors.go`). Other resources in the provider (e.g., `resource_tc_ga2_forwarding_rule.go`) already handle this pattern correctly. This change brings the forwarding_policy resource in line with the established pattern.

## Goals / Non-Goals

**Goals:**
- When `DescribeGa2ForwardingPolicyById` returns a `ResourceNotFound` error, the Read function should log a warning, clear the resource ID, and return nil (instead of returning an error)
- The change must be minimal and not affect any other behavior

**Non-Goals:**
- No schema changes
- No changes to Create, Update, or Delete functions
- No changes to the service layer (`DescribeGa2ForwardingPolicyById`)

## Decisions

**Decision 1: Handle `ResourceNotFound` in the Read function, not in the service layer**

The service layer (`DescribeGa2ForwardingPolicyById`) already handles the case where the API returns successfully but the policy is not found in the result set (returns `nil, nil`). The `ResourceNotFound` SDK error is a different case — the API itself returns an error response. The Read function is the right place to translate this error into a "not found" state cleanup, because only the resource layer knows about Terraform concepts like `d.SetId("")`.

**Decision 2: Use `!d.IsNewResource()` guard**

Following the reference pattern from other resources, the `ResourceNotFound` check includes `!d.IsNewResource()` to prevent incorrectly clearing the ID during the initial resource creation cycle (when the Read is called from within Create). This is a safety guard.

**Decision 3: Follow the existing code pattern from the same package**

The `resourceTencentCloudGa2ForwardingRuleRead` function already implements this pattern. Consistency within the package reduces maintenance burden.

## Risks / Trade-offs

- **Risk**: If the API returns `ResourceNotFound` for a transient reason (e.g., API inconsistency) during a normal Read, the resource will be removed from state. → **Mitigation**: This is the same behavior as all other Terraform resources. The retry mechanism in the service layer already handles transient errors within the `ReadRetryTimeout` window. If the API returns `ResourceNotFound` consistently, the resource truly does not exist.