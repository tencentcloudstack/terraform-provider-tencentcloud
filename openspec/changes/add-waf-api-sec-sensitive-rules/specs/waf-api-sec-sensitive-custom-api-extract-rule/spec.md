## ADDED Requirements

### Requirement: Custom API Extract Rule Resource
The provider SHALL offer a `tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule` resource that manages the `CustomApiExtractRule` (`ApiSecExtractRule`) structure of the WAF `ModifyApiSecSensitiveRule` API. Its schema SHALL contain only `domain`, `rule_name`, `status`, and the input fields of `ApiSecExtractRule` (`api_name`, `methods`, `regex`), plus the computed `update_time`, and no other fields. The sub-struct `RuleName` SHALL be unified with the top-level `rule_name`, and the sub-struct `Status` SHALL be unified with the top-level `status`.

#### Scenario: Create API extract rule
- **WHEN** a user applies the resource with `domain`, `rule_name`, `status`, `api_name`, `methods`, `regex`
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the populated `CustomApiExtractRule` and sets the ID to `Domain#RuleName`

#### Scenario: Read via extract rule list
- **WHEN** the resource is read
- **THEN** the provider calls `DescribeApiSecSensitiveRuleList` with `IsQueryApiExtractRule = true`, matches `ApiExtractRule` by `RuleName`, and sets attributes nil-safely
- **AND** if not found the provider clears the ID and returns without error

#### Scenario: Update applies changes
- **WHEN** any settable attribute changes
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the updated `CustomApiExtractRule`

#### Scenario: Delete uses internal status 3
- **WHEN** the resource is destroyed
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with `Status = 3` for the `Domain`/`RuleName`

### Requirement: Status Exposure
The `status` attribute SHALL accept only `0` (off) and `1` (on); the delete value `3` SHALL be applied only internally during destroy.

#### Scenario: Reject invalid status
- **WHEN** a user sets `status` to a value other than `0` or `1`
- **THEN** the provider rejects the configuration during validation
