## ADDED Requirements

### Requirement: API Privilege Rule Resource
The provider SHALL offer a `tencentcloud_waf_api_sec_sensitive_privilege_rule` resource that manages the `ApiSecPrivilegeRule` structure of the WAF `ModifyApiSecSensitiveRule` API. Its schema SHALL contain only `domain`, `rule_name`, `status`, and the input fields of `ApiSecPrivilegeRule` (`api_name`, `position`, `parameter_list`, `source`, `option`, and the nested `api_name_op` block with `value`, `op`, and nested `api_name_method`), plus computed `update_time`, and no other fields. The sub-struct `RuleName`/`Status` SHALL be unified with the top-level `rule_name`/`status`.

#### Scenario: Create privilege rule
- **WHEN** a user applies the resource with the privilege fields
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the populated `ApiSecPrivilegeRule` and sets the ID to `Domain#RuleName`

#### Scenario: Read via privilege rule list
- **WHEN** the resource is read
- **THEN** the provider calls `DescribeApiSecSensitiveRuleList` with `IsQueryApiPrivilegeRule = true`, matches `ApiSecPrivilegeRule` by `RuleName`, and sets attributes nil-safely
- **AND** if not found the provider clears the ID and returns without error

#### Scenario: Update applies changes
- **WHEN** any settable attribute changes
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the updated `ApiSecPrivilegeRule`

#### Scenario: Delete uses internal status 3
- **WHEN** the resource is destroyed
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with `Status = 3` for the `Domain`/`RuleName`

### Requirement: Status Exposure
The `status` attribute SHALL accept only `0` (off) and `1` (on); the delete value `3` SHALL be applied only internally during destroy.

#### Scenario: Reject invalid status
- **WHEN** a user sets `status` to a value other than `0` or `1`
- **THEN** the provider rejects the configuration during validation
