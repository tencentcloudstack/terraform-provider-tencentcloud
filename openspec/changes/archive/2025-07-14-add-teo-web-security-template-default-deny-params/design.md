## Context

The `tencentcloud_teo_web_security_template` resource manages TEO (TencentCloud EdgeOne) web security policy templates. Currently, the `security_policy` schema supports: `custom_rules`, `managed_rules`, `http_ddos_protection`, `rate_limiting_rules`, `exception_rules`, and `bot_management`.

The cloud API's `SecurityPolicy` struct has a `DefaultDenySecurityActionParameters` field that allows configuring default deny behavior for managed rules and other security modules. This field is not yet exposed in the Terraform resource schema.

The cloud API struct reference:
- `SecurityPolicy.DefaultDenySecurityActionParameters` (*DefaultDenySecurityActionParameters)
  - `ManagedRules` (*DenyActionParameters) - default deny action for managed rules
  - `OtherModules` (*DenyActionParameters) - default deny action for other modules (custom rules, rate limiting, bot management)
- `DenyActionParameters` struct contains: BlockIp, BlockIpDuration, ReturnCustomPage, ResponseCode, ErrorPageId, Stall

The CRUD API interfaces:
- Create: `CreateWebSecurityTemplate` - accepts `SecurityPolicy` with `DefaultDenySecurityActionParameters`
- Modify: `ModifyWebSecurityTemplate` - accepts `SecurityPolicy` with `DefaultDenySecurityActionParameters`
- Describe: `DescribeWebSecurityTemplate` - returns `SecurityPolicy` with `DefaultDenySecurityActionParameters`
- Delete: `DeleteWebSecurityTemplate` - no relevant changes

## Goals / Non-Goals

**Goals:**
- Add `default_deny_security_action_parameters` field to the `security_policy` schema
- Support Create/Read/Update operations for the new field
- Maintain backward compatibility with existing Terraform configurations
- Add unit tests covering the new field

**Non-Goals:**
- No changes to other existing fields in the resource schema
- No changes to the data source (if any)
- No changes to the delete operation (delete only uses ZoneId + TemplateId)

## Decisions

1. **Schema design**: The `default_deny_security_action_parameters` field will be `Optional + Computed` with `MaxItems: 1`, following the same pattern as other sub-blocks in `security_policy`. Each sub-block (`managed_rules`, `other_modules`) will be `Optional + Computed` with `MaxItems: 1`.

2. **DenyActionParameters sub-schema reuse**: The `managed_rules` and `other_modules` sub-blocks both use the same `DenyActionParameters` struct. We will define the schema inline (as done with other similar patterns in this resource) rather than creating a shared schema helper, since the existing resource code already uses inline definitions for similar repeating structures.

3. **Field mapping**: All fields in `DenyActionParameters` (block_ip, block_ip_duration, return_custom_page, response_code, error_page_id, stall) are `Optional` strings, matching the existing `deny_action_parameters` pattern already used in `custom_rules` action blocks.

4. **Read handling**: In the Read operation, check if `DefaultDenySecurityActionParameters` is nil before flattening. For each sub-block, check nil before accessing fields.

5. **Create/Update handling**: In Create and Update operations, if `default_deny_security_action_parameters` is specified in the Terraform configuration, expand it into the `SecurityPolicy.DefaultDenySecurityActionParameters` struct before making the API call.

## Risks / Trade-offs

- [Backward compatibility] The new field is purely additive (Optional) → No risk of breaking existing configurations. Existing state files without this field will continue to work.
- [API nil handling] The cloud API may return nil for `DefaultDenySecurityActionParameters` and its sub-fields → Mitigate by checking nil before accessing nested fields in Read/flatten operations.
