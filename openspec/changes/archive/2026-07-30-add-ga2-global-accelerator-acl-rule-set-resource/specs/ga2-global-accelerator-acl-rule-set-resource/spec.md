## ADDED Requirements

### Requirement: Resource schema defines the ACL rule set fields
The `tencentcloud_ga2_global_accelerator_acl_rule_set` resource SHALL expose the following schema fields:
- `global_accelerator_id` (Required, ForceNew): The GA2 instance ID.
- `global_accelerator_acl_policy_id` (Required, ForceNew): The ACL policy ID that owns the rule set.
- `acl_entries` (Required, `TypeList`): The desired full set of ACL rules under the policy. Each element SHALL expose:
  - `protocol` (Required): The protocol (`TCP`, `UDP`, or `ALL`).
  - `port` (Required): The port number or range.
  - `source_cidr_block` (Required): The source CIDR block.
  - `policy` (Required): The action policy (`ACCEPT` or `DROP`).
  - `description` (Optional): Description of the ACL rule (max 100 bytes).
  - `global_accelerator_acl_rule_id` (Computed): The server-assigned ACL rule ID.
- `task_id` (Computed): The async task ID from the latest write operation.

The resource SHALL include a `Timeouts` block with Create/Update/Delete defaults of 5 minutes. The resource SHALL NOT introduce an extra wrapping schema level around the list elements (rule fields SHALL be placed directly on each `acl_entries` element).

#### Scenario: Schema validation missing required field
- **WHEN** a user defines a `tencentcloud_ga2_global_accelerator_acl_rule_set` resource without `global_accelerator_acl_policy_id`
- **THEN** Terraform plan validation fails with a missing required field error.

#### Scenario: Empty rule set is valid
- **WHEN** a user defines a `tencentcloud_ga2_global_accelerator_acl_rule_set` resource with `acl_entries = []`
- **THEN** Terraform plan validation succeeds and the resource represents a policy with zero rules.

### Requirement: Create ACL rule set
The Create operation SHALL assemble all `acl_entries` elements into a single `AclEntries` array and call `CreateGlobalAcceleratorAclRule` once. After the API call, the operation SHALL poll the returned `TaskId` via `Ga2Service.WaitForGa2TaskFinish`. On success, the returned `GlobalAcceleratorAclRuleIds` SHALL be mapped back onto the corresponding `acl_entries` elements in input order, and the resource ID SHALL be set to the composite `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId`.

#### Scenario: Successful create with multiple rules
- **WHEN** a user applies a `tencentcloud_ga2_global_accelerator_acl_rule_set` resource with three `acl_entries`
- **THEN** `CreateGlobalAcceleratorAclRule` is called once with a three-element `AclEntries` array, the task is polled to completion, each `acl_entries` element is populated with its server-assigned `global_accelerator_acl_rule_id`, and the resource ID is set to the composite ID.

#### Scenario: Create with empty acl_entries
- **WHEN** a user applies a `tencentcloud_ga2_global_accelerator_acl_rule_set` resource with `acl_entries = []`
- **THEN** the Create operation skips the `CreateGlobalAcceleratorAclRule` API call and sets the resource ID to the composite `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId`.

#### Scenario: Create with nil response or empty rule IDs
- **WHEN** the `CreateGlobalAcceleratorAclRule` API returns a nil `Response`, a nil/empty `GlobalAcceleratorAclRuleIds`, or a number of IDs that does not match the input `AclEntries` length
- **THEN** the Create operation returns a non-retryable error and does not write an ID into state.

### Requirement: Read ACL rule set
The Read operation SHALL parse the composite ID to extract `GlobalAcceleratorId` and `GlobalAcceleratorAclPolicyId`, call `DescribeGlobalAcceleratorAclRules` with pagination (`Limit=200`) filtered by `GlobalAcceleratorAclPolicyId`, and flatten the returned `GlobalAcceleratorAclRuleSet` into the `acl_entries` list. If the API returns an empty set, the resource SHALL NOT be removed from state; instead `acl_entries` SHALL be set to an empty list. If found, all schema fields SHALL be populated from the API response.

#### Scenario: Read existing rule set
- **WHEN** Terraform refreshes state for an existing ACL rule set resource whose policy contains rules
- **THEN** all `acl_entries` elements are populated from the cloud API response, including each rule's `global_accelerator_acl_rule_id`, `protocol`, `port`, `source_cidr_block`, `policy`, and `description`.

#### Scenario: Read policy with no rules
- **WHEN** Terraform refreshes state for an ACL rule set resource whose policy contains zero rules
- **THEN** `acl_entries` is set to an empty list and the resource remains in state (it is not removed).

#### Scenario: Read after out-of-band deletion of all rules
- **WHEN** all ACL rules under the policy have been deleted outside of Terraform
- **THEN** the Read operation sets `acl_entries` to an empty list and logs `[CRUD] ga2 global_accelerator_acl_rule_set id=<id>` to preserve context, then keeps the resource in state.

### Requirement: Update ACL rule set
The Update operation SHALL, on `d.HasChange("acl_entries")`, diff the old and new `acl_entries` lists keyed by `global_accelerator_acl_rule_id`: batch-create new entries via one `CreateGlobalAcceleratorAclRule` call, batch-delete removed entries via one `DeleteGlobalAcceleratorAclRule` call, and call `ModifyGlobalAcceleratorAclRule` once per changed entry. Each write SHALL poll its returned `TaskId` via `Ga2Service.WaitForGa2TaskFinish`. Changing `global_accelerator_id` or `global_accelerator_acl_policy_id` SHALL force resource recreation rather than entering Update.

#### Scenario: Add a rule to the set
- **WHEN** a user adds a new element to `acl_entries` of an existing rule set resource
- **THEN** the new entries are batch-created via `CreateGlobalAcceleratorAclRule`, the task is polled, and the new `global_accelerator_acl_rule_id` values are mapped back onto the new elements.

#### Scenario: Remove a rule from the set
- **WHEN** a user removes an element from `acl_entries` of an existing rule set resource
- **THEN** the removed `global_accelerator_acl_rule_id` values are passed to `DeleteGlobalAcceleratorAclRule` in a single batch call and the task is polled.

#### Scenario: Modify an existing rule
- **WHEN** a user changes the `description` of an existing `acl_entries` element (keeping its `global_accelerator_acl_rule_id`)
- **THEN** `ModifyGlobalAcceleratorAclRule` is called with the new field values for that rule and the task is polled to completion.

#### Scenario: ForceNew on identity change
- **WHEN** a user changes `global_accelerator_acl_policy_id` on an existing rule set resource
- **THEN** Terraform destroys and recreates the resource rather than updating in place.

### Requirement: Delete ACL rule set
The Delete operation SHALL collect every `global_accelerator_acl_rule_id` from `acl_entries` into a single `GlobalAcceleratorAclRuleIds` array and call `DeleteGlobalAcceleratorAclRule` once. After the API call, the operation SHALL poll the returned `TaskId` via `Ga2Service.WaitForGa2TaskFinish`. If `acl_entries` is empty at delete time, the Delete operation SHALL skip the API call and return success.

#### Scenario: Delete rule set with rules
- **WHEN** a user destroys an ACL rule set resource whose `acl_entries` contains rules
- **THEN** all rules are deleted from the cloud side in a single batch `DeleteGlobalAcceleratorAclRule` call, the task is polled, and the resource is removed from state.

#### Scenario: Delete rule set with no rules
- **WHEN** a user destroys an ACL rule set resource whose `acl_entries` is empty
- **THEN** the Delete operation skips the `DeleteGlobalAcceleratorAclRule` API call and the resource is removed from state.

### Requirement: Import ACL rule set
The resource SHALL support import via `terraform import` using the composite ID format `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId`. The Importer SHALL use `schema.ImportStatePassthrough`.

#### Scenario: Import existing rule set
- **WHEN** a user runs `terraform import tencentcloud_ga2_global_accelerator_acl_rule_set.example ga-xxx#sp-xxx`
- **THEN** the resource is imported into state, all `acl_entries` are populated from the cloud API, and a subsequent plan shows no changes.
