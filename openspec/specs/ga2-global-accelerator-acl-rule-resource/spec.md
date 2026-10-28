## ADDED Requirements

### Requirement: Resource schema defines all ACL rule fields
The `tencentcloud_ga2_global_accelerator_acl_rule` resource SHALL expose the following schema fields:
- `global_accelerator_id` (Required, ForceNew): The GA2 instance ID.
- `global_accelerator_acl_policy_id` (Required, ForceNew): The ACL policy ID.
- `protocol` (Required): The protocol type (TCP, UDP, or ALL).
- `port` (Required): The port number or range.
- `source_cidr_block` (Required): The source CIDR block.
- `policy` (Required): The action policy (ACCEPT or DROP).
- `description` (Optional): Description of the ACL rule.
- `global_accelerator_acl_rule_id` (Computed): The server-assigned ACL rule ID.
- `task_id` (Computed): The async task ID from the latest write operation.

The resource SHALL include a `Timeouts` block with Create/Update/Delete defaults of 5 minutes.

#### Scenario: Schema validation
- **WHEN** a user defines a `tencentcloud_ga2_global_accelerator_acl_rule` resource without `global_accelerator_id`
- **THEN** Terraform plan validation fails with a missing required field error.

### Requirement: Create ACL rule
The Create operation SHALL call `CreateGlobalAcceleratorAclRule` with the user-provided fields assembled into a single `AclEntries` element. After the API call, the operation SHALL poll the returned `TaskId` via `WaitForGa2TaskFinish`. On success, the resource ID SHALL be set to the composite `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId#GlobalAcceleratorAclRuleId`.

#### Scenario: Successful create
- **WHEN** a user applies a `tencentcloud_ga2_global_accelerator_acl_rule` resource with valid inputs
- **THEN** the ACL rule is created on the cloud side and the resource ID is set to the composite ID.

#### Scenario: Create with nil response
- **WHEN** the `CreateGlobalAcceleratorAclRule` API returns a nil response or nil `GlobalAcceleratorAclRuleIds`
- **THEN** the Create operation returns a non-retryable error.

### Requirement: Read ACL rule
The Read operation SHALL parse the composite ID to extract `GlobalAcceleratorAclPolicyId` and `GlobalAcceleratorAclRuleId`, call `DescribeGlobalAcceleratorAclRules` with pagination (Limit=200), and match the result by `GlobalAcceleratorAclRuleId`. If no matching rule is found, the resource SHALL be removed from state (`d.SetId("")`). If found, all schema fields SHALL be populated from the API response.

#### Scenario: Read existing rule
- **WHEN** Terraform refreshes state for an existing ACL rule resource
- **THEN** all schema fields are populated from the cloud API response.

#### Scenario: Read deleted rule
- **WHEN** the ACL rule has been deleted outside of Terraform
- **THEN** the resource is removed from state with a warning log.

### Requirement: Update ACL rule
The Update operation SHALL call `ModifyGlobalAcceleratorAclRule` with the current field values when any non-computed field changes. After the API call, the operation SHALL poll the returned `TaskId` via `WaitForGa2TaskFinish`.

#### Scenario: Update description
- **WHEN** a user changes the `description` field of an existing ACL rule resource
- **THEN** `ModifyGlobalAcceleratorAclRule` is called with the new description and the task is polled to completion.

### Requirement: Delete ACL rule
The Delete operation SHALL call `DeleteGlobalAcceleratorAclRule` with `GlobalAcceleratorId`, `GlobalAcceleratorAclPolicyId`, and `GlobalAcceleratorAclRuleIds` (as a single-element slice). After the API call, the operation SHALL poll the returned `TaskId` via `WaitForGa2TaskFinish`.

#### Scenario: Delete existing rule
- **WHEN** a user destroys an ACL rule resource
- **THEN** the ACL rule is deleted from the cloud side and the resource is removed from state.

### Requirement: Import ACL rule
The resource SHALL support import via `terraform import` using the composite ID format `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId#GlobalAcceleratorAclRuleId`. The Importer SHALL use `schema.ImportStatePassthrough`.

#### Scenario: Import existing rule
- **WHEN** a user runs `terraform import tencentcloud_ga2_global_accelerator_acl_rule.example ga-xxx#aclpol-xxx#aclrule-xxx`
- **THEN** the resource is imported into state and subsequent plan shows no changes.