## Why

The `tencentcloud_teo_web_security_template` resource currently lacks support for configuring default deny security action parameters (`DefaultDenySecurityActionParameters`) in the `security_policy` block. The cloud API's `SecurityPolicy` struct already supports this field, which allows users to configure default interception handling behavior for managed rules and other security modules. Without this parameter, users cannot manage default deny actions through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `default_deny_security_action_parameters` field to the `security_policy` schema in `tencentcloud_teo_web_security_template` resource
  - This is an Optional + Computed block with MaxItems: 1
  - Contains `managed_rules` sub-block (DenyActionParameters) for managed rules default deny action configuration
  - Contains `other_modules` sub-block (DenyActionParameters) for other security modules default deny action configuration
  - Each DenyActionParameters sub-block contains: `block_ip`, `block_ip_duration`, `return_custom_page`, `response_code`, `error_page_id`, `stall`
- Update Create/Modify/Read operations to handle the new parameter
- Update unit tests to cover the new parameter
- Update the .md documentation file

## Capabilities

### New Capabilities
- `default-deny-action-params`: Adds default deny security action parameters configuration to the teo_web_security_template resource, enabling users to configure default interception handling for managed rules and other security modules

### Modified Capabilities

## Impact

- Affected files:
  - `tencentcloud/services/teo/resource_tc_teo_web_security_template.go` - Schema, CRUD logic
  - `tencentcloud/services/teo/resource_tc_teo_web_security_template_test.go` - Unit tests
  - `tencentcloud/services/teo/resource_tc_teo_web_security_template.md` - Documentation
- Cloud API: Uses existing `DefaultDenySecurityActionParameters` field in `SecurityPolicy` struct (teo v20220901)
- No breaking changes: The new field is optional, existing configurations remain compatible
