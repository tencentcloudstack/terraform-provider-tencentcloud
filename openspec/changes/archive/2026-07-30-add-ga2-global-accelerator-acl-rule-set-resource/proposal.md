## Why

Tencent Cloud Global Accelerator V2 (`ga2`) already ships a single-rule resource `tencentcloud_ga2_global_accelerator_acl_rule` that manages exactly one ACL rule per Terraform resource instance. When an ACL policy contains many rules, users must declare one resource block per rule, which is verbose and makes it hard to keep the rule set in sync with a source-of-truth list. The cloud API `CreateGlobalAcceleratorAclRule` natively accepts an `AclEntries` array (batch create) and `DeleteGlobalAcceleratorAclRule` accepts a `GlobalAcceleratorAclRuleIds` array (batch delete), so a set-style resource that owns the full rule collection under one ACL policy is both idiomatic to the API and far more ergonomic for users.

## What Changes

- Add a new resource `tencentcloud_ga2_global_accelerator_acl_rule_set` backed by the `ga2` v20250115 SDK that manages the **complete collection** of ACL rules under one `GlobalAcceleratorAclPolicyId`.
- Schema fields:
  - `global_accelerator_id` (Required, ForceNew): the GA2 instance ID.
  - `global_accelerator_acl_policy_id` (Required, ForceNew): the ACL policy ID that owns the rule set.
  - `acl_entries` (Required, `TypeList`): the desired full set of ACL rules. Each element exposes `protocol`, `port`, `source_cidr_block`, `policy`, `description` (user-editable) and `global_accelerator_acl_rule_id` (computed, server-assigned).
  - `task_id` (Computed): the async task ID from the latest write operation.
- Implement async-aware CRUD on top of the four SDK calls:
  - **Create**: assemble all `acl_entries` into a single `AclEntries` array and call `CreateGlobalAcceleratorAclRule` once; poll the returned `TaskId` via `Ga2Service.WaitForGa2TaskFinish`; map the returned `GlobalAcceleratorAclRuleIds` back onto the `acl_entries` list items.
  - **Read**: call `DescribeGlobalAcceleratorAclRules` with pagination (`Limit=200`, the documented maximum) filtered by `GlobalAcceleratorAclPolicyId`; flatten the returned `GlobalAcceleratorAclRuleSet` into the `acl_entries` list.
  - **Update**: diff old vs. new `acl_entries` — create new entries (batch via `CreateGlobalAcceleratorAclRule`), modify changed entries one-by-one via `ModifyGlobalAcceleratorAclRule`, and delete removed entries (batch via `DeleteGlobalAcceleratorAclRule`); poll every returned `TaskId`.
  - **Delete**: collect all `global_accelerator_acl_rule_id` values from `acl_entries` and call `DeleteGlobalAcceleratorAclRule` once with the full `GlobalAcceleratorAclRuleIds` array; poll the returned `TaskId`.
- Resource ID is the composite `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId` (using `tccommon.FILED_SP` as separator), so exactly one rule-set resource exists per ACL policy.
- All SDK calls are wrapped with `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`).
- Resource supports import via the composite ID `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId`.
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace and author the resource markdown doc `resource_tc_ga2_global_accelerator_acl_rule_set.md`.

## Capabilities

### New Capabilities
- `ga2-global-accelerator-acl-rule-set-resource`: Lifecycle management (create / read / update / delete / import) of the full set of ACL rules under one GA2 ACL policy as a single Terraform resource, including batch create/delete, per-rule modify, async task polling, and composite-ID import.

### Modified Capabilities
<!-- None: this change only introduces a new set-style resource; it does not alter the behavior of the existing single-rule resource `tencentcloud_ga2_global_accelerator_acl_rule`. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_rule_set.go` (schema + CRUD + build/flatten/diff helpers, single file).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_rule_set.md` (resource doc + import syntax).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_rule_set_test.go` (acceptance test skeleton).
- **Modified code**:
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_global_accelerator_acl_rule_set` in `ResourcesMap`.
- **APIs consumed**: `CreateGlobalAcceleratorAclRule`, `DescribeGlobalAcceleratorAclRules`, `ModifyGlobalAcceleratorAclRule`, `DeleteGlobalAcceleratorAclRule` (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified. Coexists with the existing single-rule resource.
- **No SDK upgrade required**: all required APIs are already present in the vendored SDK.
