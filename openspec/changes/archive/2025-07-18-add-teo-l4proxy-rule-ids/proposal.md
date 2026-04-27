## Why

The `tencentcloud_teo_l4_proxy_rule` resource currently does not expose the `l4proxy_rule_ids` field from the `CreateL4ProxyRules` API response. The API returns `L4ProxyRuleIds` upon successful creation, but this value is not persisted in the Terraform state, making it unavailable for downstream references or debugging. Adding this computed parameter will allow users to reference the created rule IDs in other resources or outputs.

## What Changes

- Add a new computed parameter `l4proxy_rule_ids` (TypeList of TypeString) to the `tencentcloud_teo_l4_proxy_rule` resource schema
- Populate `l4proxy_rule_ids` from `response.Response.L4ProxyRuleIds` in the create function
- Persist `l4proxy_rule_ids` in the read function by reading from the resource ID (the rule ID is already part of the composite ID)

## Capabilities

### New Capabilities
- `teo-l4proxy-rule-ids`: Adds the `l4proxy_rule_ids` computed attribute to the `tencentcloud_teo_l4_proxy_rule` resource, exposing the rule IDs returned by the `CreateL4ProxyRules` API

### Modified Capabilities

## Impact

- **Affected code**: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_rule.go` (schema definition, create and read functions)
- **API**: Uses existing `CreateL4ProxyRules` API response field `L4ProxyRuleIds`, no new API calls needed
- **Backward compatibility**: Fully backward compatible — adding a computed attribute does not break existing configurations
