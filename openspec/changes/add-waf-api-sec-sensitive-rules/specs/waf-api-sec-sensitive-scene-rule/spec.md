## ADDED Requirements

### Requirement: API Scene Rule Resource
The provider SHALL offer a `tencentcloud_waf_api_sec_sensitive_scene_rule` resource that manages the `ApiSecSceneRule` structure of the WAF `ModifyApiSecSensitiveRule` API. Its schema SHALL contain only `domain`, `rule_name`, `status`, and the input fields of `ApiSecSceneRule` (`source` and the nested `rule_list` block of `ApiSecSceneRuleEntry` with `key`, `value`, `operate`, `name`), plus computed `update_time`, and no other fields. The sub-struct `RuleName`/`Status` SHALL be unified with the top-level `rule_name`/`status`.

#### Scenario: Create scene rule
- **WHEN** a user applies the resource with the scene fields
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the populated `ApiSecSceneRule` and sets the ID to `Domain#RuleName`

#### Scenario: Read via scene rule list
- **WHEN** the resource is read
- **THEN** the provider calls `DescribeApiSecSensitiveRuleList` with `IsQueryApiSceneRule = true`, matches `ApiSecSceneRule` by `RuleName`, and sets attributes nil-safely
- **AND** if not found the provider clears the ID and returns without error

#### Scenario: Update applies changes
- **WHEN** any settable attribute changes
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with the updated `ApiSecSceneRule`

#### Scenario: Delete uses internal status 3
- **WHEN** the resource is destroyed
- **THEN** the provider calls `ModifyApiSecSensitiveRule` with `Status = 3` for the `Domain`/`RuleName`

### Requirement: Status Exposure
The `status` attribute SHALL accept only `0` (off) and `1` (on); the delete value `3` SHALL be applied only internally during destroy.

#### Scenario: Reject invalid status
- **WHEN** a user sets `status` to a value other than `0` or `1`
- **THEN** the provider rejects the configuration during validation
