## Why

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policies. The `bot_management.client_attestation_rules` block currently exposes the core fields of the cloud API's `ClientAttestationRule` struct (id, name, enabled, priority, condition, attester_id, invalid_attestation_action) but does NOT expose the `DeviceProfiles` field. The cloud API supports per-rule client device configuration (`DeviceProfiles`, a list of `DeviceProfile`) through both `DescribeSecurityPolicy` (response: `SecurityPolicy.BotManagement.ClientAttestationRules.Rules.DeviceProfiles`) and `ModifySecurityPolicy` (request: `SecurityPolicy.BotManagement.ClientAttestationRules.Rules.DeviceProfiles`). Without this parameter, Terraform users cannot configure device-specific attestation scoring thresholds (risk levels and the corresponding request actions for iOS / Android / WebView / WeChatMiniProgram clients), which is required to fully manage client attestation rules through Terraform.

## What Changes

- Add a new `device_profiles` parameter (TypeList, Optional) under each rule of the `client_attestation_rules` block in the `bot_management` section of the `tencentcloud_teo_security_policy_config` resource
- The `device_profiles` list will contain `DeviceProfile` elements, each exposing the following sub-fields:
  - `client_type` (TypeString, Required): Client device type. Valid values: `iOS`, `Android`, `WebView`, `WeChatMiniProgram`.
  - `high_risk_min_score` (TypeInt, Optional, Computed): Minimum score (1-99) for judging a request as high risk. Default 50.
  - `high_risk_request_action` (TypeList, MaxItems: 1, Optional): Action for high-risk requests, using the existing `SecurityAction` schema pattern (Name supports `Deny`, `Monitor`, `Redirect`, `Challenge`; default `Monitor`).
  - `medium_risk_min_score` (TypeInt, Optional, Computed): Minimum score (1-99) for judging a request as medium risk. Default 15.
  - `medium_risk_request_action` (TypeList, MaxItems: 1, Optional): Action for medium-risk requests, using the existing `SecurityAction` schema pattern (default `Monitor`).
- Implement Read (flatten) logic to populate `device_profiles` from the `DescribeSecurityPolicy` response's `BotManagement.ClientAttestationRules.Rules[].DeviceProfiles`, with nil checks at every nesting level
- Implement Create/Update (expand) logic to expand `device_profiles` into the `ModifySecurityPolicy` request's `SecurityPolicy.BotManagement.ClientAttestationRules.Rules[].DeviceProfiles` field
- Add unit tests using the gomonkey mock approach for the new parameter
- Update the resource `.md` documentation with example usage

## Capabilities

### New Capabilities
- `teo-security-policy-client-attestation-device-profiles`: Adds the `device_profiles` parameter to the `client_attestation_rules` block of the `tencentcloud_teo_security_policy_config` resource, enabling Terraform users to configure per-rule client device attestation scoring thresholds (risk levels and request actions) for iOS / Android / WebView / WeChatMiniProgram clients through the SecurityPolicy API.

### Modified Capabilities

## Impact

- **Affected files**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`, `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`, `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md`
- **Cloud API**: Uses existing `DescribeSecurityPolicy` and `ModifySecurityPolicy` APIs (no new API dependencies). The `DeviceProfiles` field and `DeviceProfile` struct already exist in the vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).
- **Backward compatibility**: Adding a new Optional parameter is backward compatible; existing Terraform configurations and state will continue to work without changes. The cloud API preserves existing device configuration when `DeviceProfiles` is not specified.
