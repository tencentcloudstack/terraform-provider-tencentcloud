# teo-l4proxy-rule-ids Specification

## Purpose
TBD - created by archiving change add-teo-l4proxy-rule-ids. Update Purpose after archive.
## Requirements
### Requirement: l4proxy_rule_ids computed attribute
The `tencentcloud_teo_l4_proxy_rule` resource SHALL include a computed attribute `l4proxy_rule_ids` of type `TypeList` with element type `TypeString`. This attribute SHALL be populated from the `L4ProxyRuleIds` field in the `CreateL4ProxyRules` API response during resource creation and persisted in the Terraform state during read operations.

#### Scenario: l4proxy_rule_ids populated after resource creation
- **WHEN** a `tencentcloud_teo_l4_proxy_rule` resource is created successfully via the `CreateL4ProxyRules` API
- **THEN** the `l4proxy_rule_ids` attribute SHALL be set to the list of rule IDs returned in `response.Response.L4ProxyRuleIds`

#### Scenario: l4proxy_rule_ids populated during resource read
- **WHEN** the `tencentcloud_teo_l4_proxy_rule` resource is read (refresh state)
- **THEN** the `l4proxy_rule_ids` attribute SHALL be set to a list containing the rule ID extracted from the composite resource ID

#### Scenario: l4proxy_rule_ids is computed and read-only
- **WHEN** a user attempts to set `l4proxy_rule_ids` in their Terraform configuration
- **THEN** the provider SHALL reject the configuration since the attribute is computed-only and not user-configurable
