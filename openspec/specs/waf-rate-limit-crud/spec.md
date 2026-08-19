## ADDED Requirements

### Requirement: Create WAF rate limit rule
The system SHALL create a WAF rate limit rule by calling `CreateRateLimitV2` API with the user-specified parameters. The resource SHALL store the composite ID (`domain#limit_rule_id`) in Terraform state after successful creation. The system MUST validate that the API response contains a non-nil, non-zero `LimitRuleID` before setting the resource ID.

#### Scenario: Successful creation with required parameters
- **WHEN** user provides `domain`, `name`, `priority`, `status`, `limit_window`, `limit_object`, and `limit_strategy`
- **THEN** the system calls `CreateRateLimitV2` API and stores the returned `LimitRuleID` combined with `domain` as the resource ID

#### Scenario: Creation returns empty response
- **WHEN** the `CreateRateLimitV2` API returns a nil response or nil `LimitRuleID`
- **THEN** the system returns a `NonRetryableError` indicating the creation failed with empty response

#### Scenario: Creation with all optional parameters
- **WHEN** user provides all optional parameters including `limit_method`, `limit_paths`, `limit_headers`, `limit_header_name`, `get_params_name`, `get_params_value`, `post_params_name`, `post_params_value`, `ip_location`, `redirect_info`, `block_page`, `object_src`, `quota_share`, `paths_option`, and `order`
- **THEN** the system passes all parameters to the `CreateRateLimitV2` API and creates the rule successfully

### Requirement: Read WAF rate limit rule
The system SHALL read a WAF rate limit rule by calling `DescribeRateLimitsV2` API with the `Domain` and `Id` parameters extracted from the composite resource ID. The system SHALL set all readable fields from the API response into Terraform state.

#### Scenario: Successful read
- **WHEN** the resource exists and `DescribeRateLimitsV2` returns a non-empty `RateLimits` array
- **THEN** the system sets all fields from the first matching rule into Terraform state

#### Scenario: Resource not found
- **WHEN** `DescribeRateLimitsV2` returns an empty `RateLimits` array or nil response
- **THEN** the system logs the resource ID, calls `d.SetId("")` to remove from state, and returns nil

#### Scenario: Read with retry on transient errors
- **WHEN** the `DescribeRateLimitsV2` API call fails with a transient error
- **THEN** the system retries the call within `tccommon.ReadRetryTimeout`

### Requirement: Update WAF rate limit rule
The system SHALL update a WAF rate limit rule by calling `UpdateRateLimitV2` API with the `Domain`, `LimitRuleId`, and all changed parameters. The system SHALL detect which parameters have changed and include them in the update request.

#### Scenario: Successful update of mutable fields
- **WHEN** user modifies any mutable field (name, priority, status, limit_object, limit_strategy, limit_header_name, limit_method, limit_paths, limit_headers, limit_window, get_params_name, get_params_value, post_params_name, post_params_value, ip_location, redirect_info, block_page, object_src, quota_share, paths_option, order)
- **THEN** the system calls `UpdateRateLimitV2` API with the updated values and refreshes state via Read

#### Scenario: Update with retry on transient errors
- **WHEN** the `UpdateRateLimitV2` API call fails with a transient error
- **THEN** the system retries the call within `tccommon.WriteRetryTimeout`

### Requirement: Delete WAF rate limit rule
The system SHALL delete a WAF rate limit rule by calling `DeleteRateLimitsV2` API with the `Domain` and a single-element `LimitRuleIds` array containing the rule ID extracted from the composite resource ID.

#### Scenario: Successful deletion
- **WHEN** the resource exists and `DeleteRateLimitsV2` is called
- **THEN** the system removes the rule and clears the resource ID from state

#### Scenario: Delete with retry on transient errors
- **WHEN** the `DeleteRateLimitsV2` API call fails with a transient error
- **THEN** the system retries the call within `tccommon.WriteRetryTimeout`

### Requirement: Import WAF rate limit rule
The system SHALL support importing an existing WAF rate limit rule using the composite ID format `domain#limit_rule_id`. The import SHALL trigger a Read operation to populate all state fields.

#### Scenario: Successful import
- **WHEN** user runs `terraform import tencentcloud_waf_rate_limit.example domain#limit_rule_id`
- **THEN** the system parses the composite ID, calls Read to populate state, and the resource is managed by Terraform

#### Scenario: Invalid import ID format
- **WHEN** user provides an import ID that does not contain exactly one `#` separator
- **THEN** the system returns an error indicating the expected format

### Requirement: Provider registration
The system SHALL register `tencentcloud_waf_rate_limit` in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`.

#### Scenario: Resource available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_waf_rate_limit` is available as a valid resource type

### Requirement: Resource documentation
The system SHALL provide a documentation file at `tencentcloud/services/waf/resource_tc_waf_rate_limit.md` with a one-line description, Example Usage section, and Import section showing the composite ID format.

#### Scenario: Documentation includes example
- **WHEN** user views the resource documentation
- **THEN** the documentation shows a complete HCL example with required parameters and an import section with the composite ID format
