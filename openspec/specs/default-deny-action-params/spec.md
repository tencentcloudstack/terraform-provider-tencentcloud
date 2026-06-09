## ADDED Requirements

### Requirement: Default deny security action parameters in security_policy
The `tencentcloud_teo_web_security_template` resource SHALL support a `default_deny_security_action_parameters` block within `security_policy`, allowing users to configure default interception handling behavior for managed rules and other security modules.

#### Scenario: Create resource with default_deny_security_action_parameters
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` resource with `default_deny_security_action_parameters` specified in the `security_policy` block
- **THEN** the resource SHALL be created with the specified default deny action parameters, and the state SHALL reflect the configured values

#### Scenario: Read resource with default_deny_security_action_parameters
- **WHEN** a user reads a `tencentcloud_teo_web_security_template` resource that has `DefaultDenySecurityActionParameters` configured in the cloud
- **THEN** the resource state SHALL contain the `default_deny_security_action_parameters` block with `managed_rules` and `other_modules` sub-blocks populated from the cloud API response

#### Scenario: Update resource default_deny_security_action_parameters
- **WHEN** a user updates the `default_deny_security_action_parameters` block in an existing `tencentcloud_teo_web_security_template` resource
- **THEN** the resource SHALL be updated with the new default deny action parameters via the ModifyWebSecurityTemplate API

#### Scenario: Create resource without default_deny_security_action_parameters
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` resource without specifying `default_deny_security_action_parameters`
- **THEN** the resource SHALL be created successfully, and the field SHALL not be included in the API request

### Requirement: managed_rules sub-block in default_deny_security_action_parameters
The `default_deny_security_action_parameters` block SHALL contain a `managed_rules` sub-block of type DenyActionParameters for configuring default deny actions for managed rules.

#### Scenario: Configure managed_rules deny action parameters
- **WHEN** a user specifies `managed_rules` within `default_deny_security_action_parameters` with block_ip, block_ip_duration, return_custom_page, response_code, error_page_id, and stall values
- **THEN** the values SHALL be sent to the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API under `SecurityPolicy.DefaultDenySecurityActionParameters.ManagedRules`

### Requirement: other_modules sub-block in default_deny_security_action_parameters
The `default_deny_security_action_parameters` block SHALL contain an `other_modules` sub-block of type DenyActionParameters for configuring default deny actions for other security modules (custom rules, rate limiting, bot management).

#### Scenario: Configure other_modules deny action parameters
- **WHEN** a user specifies `other_modules` within `default_deny_security_action_parameters` with block_ip, block_ip_duration, return_custom_page, response_code, error_page_id, and stall values
- **THEN** the values SHALL be sent to the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API under `SecurityPolicy.DefaultDenySecurityActionParameters.OtherModules`

### Requirement: DenyActionParameters fields
Both `managed_rules` and `other_modules` sub-blocks SHALL support the following fields: `block_ip` (Optional string), `block_ip_duration` (Optional string), `return_custom_page` (Optional string), `response_code` (Optional string), `error_page_id` (Optional string), `stall` (Optional string).

#### Scenario: All DenyActionParameters fields are optional
- **WHEN** a user specifies only some fields within `managed_rules` or `other_modules`
- **THEN** only the specified fields SHALL be sent to the API, and unspecified fields SHALL be omitted from the request

#### Scenario: Read returns nil DenyActionParameters
- **WHEN** the cloud API returns nil for `ManagedRules` or `OtherModules` within `DefaultDenySecurityActionParameters`
- **THEN** the corresponding sub-block SHALL not be set in the Terraform state
