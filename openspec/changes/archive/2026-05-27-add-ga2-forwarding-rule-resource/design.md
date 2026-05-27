## Context

TencentCloud GA2 (Global Accelerator 2.0) provides Layer-7 forwarding rules that route traffic based on conditions (host, path) and perform actions (forward to endpoint groups). The GA2 SDK package (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`) is already vendored and provides four APIs for forwarding rule management: CreateForwardingRule, DescribeForwardingRule, ModifyForwardingRule, and DeleteForwardingRule.

An existing resource `tencentcloud_ga2_endpoint_group` in `tencentcloud/services/ga2/` demonstrates the established patterns for this service, including async task polling via `WaitForGa2TaskFinish`.

## Goals / Non-Goals

**Goals:**
- Implement a fully functional `tencentcloud_ga2_forwarding_rule` resource with Create, Read, Update, Delete operations.
- Handle async operations (Create/Modify/Delete all return TaskId) by polling DescribeTaskResult until success.
- Use a composite ID (`global_accelerator_id#listener_id#forwarding_policy_id#forwarding_rule_id`) for resource identification and import support.
- Provide unit tests using gomonkey mocks for all CRUD paths.
- Register the resource in provider.go and provider.md.

**Non-Goals:**
- Data source for listing forwarding rules (out of scope for this change).
- Support for tags on forwarding rules (the API does not expose tag fields).
- Acceptance tests requiring real cloud credentials.

## Decisions

### 1. Composite ID Structure
**Decision**: Use 4-part composite ID: `global_accelerator_id#listener_id#forwarding_policy_id#forwarding_rule_id` with `tccommon.FILED_SP` as separator.
**Rationale**: The Delete and Modify APIs require all four IDs. The Read API (DescribeForwardingRule) requires the first three plus pagination to find the specific rule by `forwarding_rule_id`. This matches the pattern used by `tencentcloud_ga2_endpoint_group` (3-part ID).

### 2. ForceNew Fields
**Decision**: Mark `global_accelerator_id`, `listener_id`, and `forwarding_policy_id` as ForceNew. The `forwarding_rule_id` is computed (returned by Create).
**Rationale**: These three fields identify the parent context (accelerator → listener → policy). Changing them means a different forwarding policy context, requiring resource recreation. The Modify API accepts these as identifiers, not mutable fields.

### 3. Async Operation Handling
**Decision**: Reuse the existing `Ga2Service.WaitForGa2TaskFinish` method after Create, Modify, and Delete operations.
**Rationale**: All three mutation APIs return a TaskId indicating async processing. The existing helper already polls DescribeTaskResult with proper retry logic and timeout support.

### 4. Read Implementation with Pagination
**Decision**: In the Read function, call DescribeForwardingRule with pagination (Limit=100) and iterate through results to find the matching `forwarding_rule_id`.
**Rationale**: The DescribeForwardingRule API returns all rules under a forwarding policy, not a single rule. We must paginate and filter client-side by `forwarding_rule_id`, similar to how `DescribeGa2EndpointGroupById` works.

### 5. Schema Design for Nested Structures
**Decision**: Use `TypeList` for `rule_conditions`, `rule_actions`, and `origin_headers` with nested `Elem` schemas.
**Rationale**: These are arrays of structured objects in the API. TypeList preserves ordering which is important for rule evaluation order.

## Risks / Trade-offs

- [Risk] DescribeForwardingRule returns rules for the entire policy, not a single rule → Mitigation: Client-side filtering by `forwarding_rule_id` with pagination ensures correctness.
- [Risk] Async task polling may time out for large configurations → Mitigation: Use configurable Timeouts block (default 20 minutes) consistent with endpoint_group resource.
- [Trade-off] Using TypeList for rule_conditions/rule_actions means order matters in state → This is intentional as rule evaluation order is significant for Layer-7 routing.
