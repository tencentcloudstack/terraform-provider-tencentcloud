## Context

The Tencent Cloud GA2 (Global Accelerator V2) service already has several Terraform resources in this provider: `tencentcloud_ga2_global_accelerator`, `tencentcloud_ga2_listener`, `tencentcloud_ga2_endpoint_group`, `tencentcloud_ga2_forwarding_rule`, `tencentcloud_ga2_forwarding_policy`, and `tencentcloud_ga2_accelerate_area`. However, ACL rules — which control access at the accelerator level via protocol/port/cidr/policy tuples — are not yet managed. The GA2 SDK (`ga2/v20250115`) already contains the four CRUD APIs (`CreateGlobalAcceleratorAclRule`, `DescribeGlobalAcceleratorAclRules`, `ModifyGlobalAcceleratorAclRule`, `DeleteGlobalAcceleratorAclRule`) and the `GlobalAcceleratorAclRuleSet` response model, as well as the `AclEntries` input model.

The existing `Ga2Service` already provides `WaitForGa2TaskFinish` for async task polling and follows a consistent pagination pattern for describe-by-ID helpers. This resource will follow the same patterns.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_ga2_global_accelerator_acl_rule` that manages a single ACL rule under a GA2 ACL policy.
- Support all CRUD operations with async task polling via `DescribeTaskResult`.
- Support import via composite ID.
- Follow existing ga2 resource conventions (schema layout, retry wrappers, deferred consistency checks, async task polling).

**Non-Goals:**
- Bulk ACL rule management (this resource manages exactly one rule).
- The `AclEntries` array from Create is flattened into individual fields; the resource does not expose the raw `AclEntries` list structure.

## Decisions

### 1. Composite ID: `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId#GlobalAcceleratorAclRuleId`

**Why**: The `DescribeGlobalAcceleratorAclRules` API is keyed by `GlobalAcceleratorAclPolicyId` only — it does not accept `GlobalAcceleratorId` or `GlobalAcceleratorAclRuleId` as a filter. The `ModifyGlobalAcceleratorAclRule` and `DeleteGlobalAcceleratorAclRule` APIs require `GlobalAcceleratorId`, `GlobalAcceleratorAclPolicyId`, and `GlobalAcceleratorAclRuleId` (or `GlobalAcceleratorAclRuleIds`). Therefore, all three identifiers must be stored in the Terraform resource ID so they are available for Read, Update, and Delete operations.

**Alternatives considered**:
- Store only `GlobalAcceleratorAclRuleId` as the ID: rejected because the Describe API requires `GlobalAcceleratorAclPolicyId`, and Modify/Delete require `GlobalAcceleratorId`.
- Use a flat ID with additional schema fields storing the parent IDs: rejected because it breaks the `terraform import` flow (import only receives the ID string).

### 2. Schema: Flatten AclEntries into individual fields

**Why**: The `CreateGlobalAcceleratorAclRule` API accepts `AclEntries` as a list of `AclEntries` structs (each with `Protocol`, `Port`, `SourceCidrBlock`, `Policy`, `Description`), while `ModifyGlobalAcceleratorAclRule` operates on a single rule with these same fields as top-level parameters. Since this resource manages exactly one rule, flattening the `AclEntries` fields into top-level schema attributes (matching the Modify API shape) is the natural approach. On Create, the single `AclEntries` element is assembled from these fields.

**Alternatives considered**:
- Expose `acl_entries` as a TypeList with TypeList elements: rejected because Modify uses top-level fields, creating a mismatch between Create and Update parameter shapes.

### 3. Service helper: `DescribeGa2GlobalAcceleratorAclRuleById(ctx, policyId, ruleId)`

**Why**: The `DescribeGlobalAcceleratorAclRules` API paginates results and returns `[]*GlobalAcceleratorAclRuleSet`. Since there is no per-rule filter, we paginate (Limit=200, the documented maximum) and match client-side by `GlobalAcceleratorAclRuleId`.

### 4. Async task polling on all writes

**Why**: All three write APIs (`CreateGlobalAcceleratorAclRule`, `ModifyGlobalAcceleratorAclRule`, `DeleteGlobalAcceleratorAclRule`) return a `TaskId` in their response. The existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper is reused for all three, with per-operation timeouts from the resource's `Timeouts` block.

### 5. Resource schema field ordering

Following the convention of existing ga2 resources, fields are ordered as: parent identifiers first (`global_accelerator_id`, `global_accelerator_acl_policy_id`), then the rule's own fields (`protocol`, `port`, `source_cidr_block`, `policy`, `description`), and finally computed/read-only fields (`global_accelerator_acl_rule_id`, `task_id`).

## Risks / Trade-offs

- **[Risk] Describe API has no `GlobalAcceleratorId` filter**: If the policy ID is reused across different accelerators (unlikely but possible), the client-side match by `GlobalAcceleratorAclRuleId` alone could return the wrong rule. → **Mitigation**: The `GlobalAcceleratorAclRuleId` is unique within the GA2 system, so matching by rule ID alone is sufficient. The `GlobalAcceleratorId` is stored in the composite ID for use in Modify/Delete calls only.
- **[Risk] Async task may fail after initial API call**: The `WaitForGa2TaskFinish` helper polls until success or timeout; if the task fails, the error is propagated to the user. → **Mitigation**: This is the same pattern used by all other ga2 resources; no additional handling needed.
- **[Trade-off] Single-rule-per-resource**: Users managing many ACL rules will need multiple Terraform resource blocks. → **Acceptable**: This matches the Terraform resource model and allows per-rule lifecycle management.