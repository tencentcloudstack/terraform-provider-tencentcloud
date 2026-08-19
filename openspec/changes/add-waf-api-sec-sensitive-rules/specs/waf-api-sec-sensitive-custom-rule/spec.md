## ADDED Requirements

### Requirement: Custom Sensitive Rule Resource
The provider SHALL offer a `tencentcloud_waf_api_sec_sensitive_custom_rule` resource that manages the `CustomRule` (`ApiSecCustomSensitiveRule`) structure of the WAF `ModifyApiSecSensitiveRule` API. Its schema SHALL contain only `domain`, `rule_name`, `status`, and the fields of `ApiSecCustomSensitiveRule` (`position`, `match_key`, `match_value`, `level`, `match_cond`, `is_pan`), and no other fields.

#### Scenario: Create custom sensitive rule
- **WHEN** a user applies a `tencentcloud_waf_api_sec_sensitive_custom_rule` with `domain`, `rule_name`, `status` and `CustomRule` fields
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the top-level `Domain`, `RuleName`, `Status`, and the populated `CustomRule`, and sets the resource ID to `Domain#RuleName`

#### Scenario: Read reflects remote state
- **WHEN** the resource is read
- **THEN** the provider calls `DescribeApiSecSensitiveRuleList` for the `Domain`, matches the entry by `RuleName`, and sets the schema attributes nil-safely
- **AND** if no matching rule exists the provider clears the ID and returns without error

#### Scenario: Update applies changes
- **WHEN** any settable attribute changes
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the updated `CustomRule` and the same `Domain`/`RuleName`

#### Scenario: Delete uses internal status 3
- **WHEN** the resource is destroyed
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with `Status = 3` for the `Domain`/`RuleName`

### Requirement: Status Exposure
The `status` attribute SHALL accept only `0` (off) and `1` (on); the delete value `3` SHALL NOT be user-settable and SHALL be applied only internally during destroy.

#### Scenario: Reject invalid status
- **WHEN** a user sets `status` to a value other than `0` or `1`
- **THEN** the provider rejects the configuration during validation
