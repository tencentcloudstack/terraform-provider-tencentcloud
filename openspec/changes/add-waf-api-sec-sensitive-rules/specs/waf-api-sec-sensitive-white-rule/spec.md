## ADDED Requirements

### Requirement: Sensitive White Rule Resource
The provider SHALL offer a `tencentcloud_waf_api_sec_sensitive_white_rule` resource that manages the `ApiSecSensitiveWhiteRuleRule` (`ApiSecSensitiveWhiteRule`) structure of the WAF `ModifyApiSecSensitiveRule` API. Its schema SHALL contain only `domain`, `rule_name`, `status`, and the input fields of `ApiSecSensitiveWhiteRule` (`white_mode`, `description`, the nested `api_name_op` block, and the `white_fields` block of `ApiSecSensitiveWhiteField` with `field_name`, `field_type`, `sensitive_types`), plus computed `update_time`, and no other fields. The sub-struct `RuleName`/`Status` SHALL be unified with the top-level `rule_name`/`status`.

#### Scenario: Create white rule
- **WHEN** a user applies the resource with the white-rule fields
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the populated `ApiSecSensitiveWhiteRuleRule` and sets the ID to `Domain#RuleName`

#### Scenario: Read via sensitive white rule list
- **WHEN** the resource is read
- **THEN** the provider calls `DescribeApiSecSensitiveRuleList` with `IsQueryApiSensitiveWhiteRule = true`, matches `ApiSecSensitiveWhiteRule` by `RuleName`, and sets attributes nil-safely
- **AND** if not found the provider clears the ID and returns without error

#### Scenario: Update applies changes
- **WHEN** any settable attribute changes
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the updated `ApiSecSensitiveWhiteRuleRule`

#### Scenario: Delete uses internal status 3
- **WHEN** the resource is destroyed
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with `Status = 3` for the `Domain`/`RuleName`

### Requirement: Status Exposure
The `status` attribute SHALL accept only `0` (off) and `1` (on); the delete value `3` SHALL be applied only internally during destroy.

#### Scenario: Reject invalid status
- **WHEN** a user sets `status` to a value other than `0` or `1`
- **THEN** the provider rejects the configuration during validation
