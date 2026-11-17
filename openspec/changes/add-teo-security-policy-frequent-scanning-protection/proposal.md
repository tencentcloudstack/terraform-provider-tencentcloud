## Why

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policies. The `security_policy.managed_rules` block currently exposes `enabled`, `detection_only`, `semantic_analysis`, `auto_update`, and `managed_rule_groups`, but does NOT expose the `FrequentScanningProtection` field. The cloud API supports high-frequency scanning protection configuration through both `DescribeSecurityPolicy` (response: `SecurityPolicy.ManagedRules.FrequentScanningProtection`) and `ModifySecurityPolicy` (request: `SecurityPolicy.ManagedRules.FrequentScanningProtection`). Without this parameter, Terraform users cannot configure automatic banning of visitors who frequently trigger "block" managed rules within a time window, which is required to fully manage the high-frequency scanning protection feature through Terraform.

## What Changes

- Add a new `frequent_scanning_protection` parameter (TypeList, Optional, MaxItems: 1, Computed) under the `managed_rules` block (inside `security_policy`) of the `tencentcloud_teo_security_policy_config` resource
- The `frequent_scanning_protection` block will expose the following sub-fields:
  - `enabled` (TypeString, Optional, Computed): Whether the high-frequency scanning protection rule is enabled. Valid values: `on`, `off`.
  - `id` (TypeString, Computed): Rule ID, output-only.
  - `action` (TypeList, MaxItems: 1, Optional, Computed): The disposal action, reusing the existing `securityActionSchema()` helper (Name supports `Deny`, `Monitor`, `JSChallenge`).
  - `count_by` (TypeString, Optional, Computed): Request statistical matching method. Valid values: `http.request.xff_header_ip`, `http.request.ip`.
  - `block_threshold` (TypeInt, Optional, Computed): Threshold for blocking, cumulative block count within the counting period. Range 1-4294967294.
  - `counting_period` (TypeString, Optional, Computed): Time window for statistics. Range 5-1800, unit seconds (s), e.g. `5s`.
  - `action_duration` (TypeString, Optional, Computed): Duration of the action. Range 60-86400, unit seconds (s), e.g. `60s`.
- Implement Read (flatten) logic to populate `frequent_scanning_protection` from the `DescribeSecurityPolicy` response's `ManagedRules.FrequentScanningProtection`, with nil checks at every nesting level; reuse `flattenSecurityAction()` for the nested `Action`.
- Implement Create/Update (expand) logic to expand `frequent_scanning_protection` into the `ModifySecurityPolicy` request's `SecurityPolicy.ManagedRules.FrequentScanningProtection` field; reuse `buildSecurityActionFromMap()` for the nested `Action`.
- Add unit tests using the gomonkey mock approach for the new parameter.
- Update the resource `.md` documentation with example usage (generated via `make doc` in the finalize stage).

## Capabilities

### New Capabilities
- `teo-security-policy-frequent-scanning-protection`: Adds the `frequent_scanning_protection` parameter to the `managed_rules` block of the `tencentcloud_teo_security_policy_config` resource, enabling Terraform users to configure high-frequency scanning protection (automatic banning of visitors who frequently trigger block-configured managed rules within a time window) through the SecurityPolicy API.

### Modified Capabilities

## Impact

- **Affected files**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`, `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`, `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md`
- **Cloud API**: Uses existing `DescribeSecurityPolicy` and `ModifySecurityPolicy` APIs (no new API dependencies). The `FrequentScanningProtection` field and struct already exist in the vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).
- **Backward compatibility**: Adding a new Optional parameter is backward compatible; existing Terraform configurations and state will continue to work without changes. The cloud API preserves existing configuration when `FrequentScanningProtection` is not specified.
