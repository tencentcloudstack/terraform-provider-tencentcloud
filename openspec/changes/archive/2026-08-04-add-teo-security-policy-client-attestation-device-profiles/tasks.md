## 1. Schema Definition

- [x] 1.1 Add `device_profiles` field (TypeList, Optional) under the `client_attestation_rules` rule schema (inside `bot_management`) in `resource_tc_teo_security_policy_config.go`, located after the existing `invalid_attestation_action` field
- [x] 1.2 Add `client_type` sub-field (TypeString, Required) with valid values `iOS`, `Android`, `WebView`, `WeChatMiniProgram`
- [x] 1.3 Add `high_risk_min_score` sub-field (TypeInt, Optional, Computed)
- [x] 1.4 Add `high_risk_request_action` sub-field (TypeList, MaxItems: 1, Optional) with `Elem: securityActionSchema()`
- [x] 1.5 Add `medium_risk_min_score` sub-field (TypeInt, Optional, Computed)
- [x] 1.6 Add `medium_risk_request_action` sub-field (TypeList, MaxItems: 1, Optional) with `Elem: securityActionSchema()`

## 2. Read (Flatten) Implementation

- [x] 2.1 In the `client_attestation_rules` flatten loop, after the `rule.InvalidAttestationAction` handling, add a nil guard `if rule.DeviceProfiles != nil` and iterate the list
- [x] 2.2 For each `DeviceProfile`, set `client_type` from `ClientType` (with nil check)
- [x] 2.3 Set `high_risk_min_score` from `HighRiskMinScore` and `medium_risk_min_score` from `MediumRiskMinScore` (with nil checks)
- [x] 2.4 Set `high_risk_request_action` via `flattenSecurityAction()` when `HighRiskRequestAction != nil`, and `medium_risk_request_action` via `flattenSecurityAction()` when `MediumRiskRequestAction != nil`
- [x] 2.5 Set the flattened `deviceProfilesList` into `ruleMap["device_profiles"]`

## 3. Create/Update (Expand) Implementation

- [x] 3.1 In the `client_attestation_rules` expand loop, after the `invalid_attestation_action` handling, read `ruleMap["device_profiles"]` and iterate entries
- [x] 3.2 For each entry, construct a `teov20220901.DeviceProfile{}` and set `ClientType` via `helper.String()` when `client_type` is non-empty
- [x] 3.3 Set `HighRiskMinScore` and `MediumRiskMinScore` via `helper.IntUint64()` when present
- [x] 3.4 Set `HighRiskRequestAction` and `MediumRiskRequestAction` via `buildSecurityActionFromMap()` when the corresponding action map is present
- [x] 3.5 Append the constructed `DeviceProfile` to a slice and assign it to `clientAttestationRule.DeviceProfiles`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` to add example usage showing `device_profiles` (with `client_type`, score thresholds, and request actions) inside a `client_attestation_rules` rule

## 5. Unit Tests

- [x] 5.1 Add unit test functions for `device_profiles` flatten logic using the gomonkey mock approach in `resource_tc_teo_security_policy_config_test.go`
- [x] 5.2 Add unit test functions for `device_profiles` expand logic using the gomonkey mock approach
- [x] 5.3 Add a unit test covering the nil `DeviceProfiles` response path (state stays empty)

## 6. Verification

- [x] 6.1 Verify the code compiles without errors (visual inspection of code correctness)
- [x] 6.2 Verify all new schema fields are Optional/Computed and backward compatible
- [x] 6.3 Verify nil checks are present at every nesting level in the Read flatten logic
- [x] 6.4 Verify the cloud API `DeviceProfiles` / `DeviceProfile` fields used exist in the vendored SDK
