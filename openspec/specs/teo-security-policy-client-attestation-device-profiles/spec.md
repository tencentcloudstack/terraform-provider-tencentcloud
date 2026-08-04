## ADDED Requirements

### Requirement: device_profiles schema definition
The `tencentcloud_teo_security_policy_config` resource SHALL include an Optional `device_profiles` parameter (TypeList, no MaxItems) under each rule of the `bot_management.client_attestation_rules` block. Each element SHALL be a `DeviceProfile` exposing:
- `client_type` (TypeString, Required): Client device type. Valid values: `iOS`, `Android`, `WebView`, `WeChatMiniProgram`.
- `high_risk_min_score` (TypeInt, Optional, Computed): Minimum score (1-99) for judging a request as high risk. Default 50.
- `high_risk_request_action` (TypeList, MaxItems: 1, Optional): Action for high-risk requests, using the shared `SecurityAction` schema (Name supports `Deny`, `Monitor`, `Redirect`, `Challenge`; default `Monitor`).
- `medium_risk_min_score` (TypeInt, Optional, Computed): Minimum score (1-99) for judging a request as medium risk. Default 15.
- `medium_risk_request_action` (TypeList, MaxItems: 1, Optional): Action for medium-risk requests, using the shared `SecurityAction` schema (default `Monitor`).

#### Scenario: Resource accepts device_profiles configuration
- **WHEN** a user provides a `device_profiles` block inside a `client_attestation_rules` rule of `tencentcloud_teo_security_policy_config`
- **THEN** the resource SHALL accept and process the configuration without errors

#### Scenario: device_profiles is optional
- **WHEN** a user creates a `tencentcloud_teo_security_policy_config` resource with a `client_attestation_rules` rule that omits `device_profiles`
- **THEN** the resource SHALL be created successfully and the existing behavior SHALL be preserved

### Requirement: client_type field
The `client_type` field SHALL be a Required TypeString that identifies the client device type. Valid values SHALL be `iOS`, `Android`, `WebView`, `WeChatMiniProgram`.

#### Scenario: client_type is required
- **WHEN** a user provides a `device_profiles` entry without a `client_type`
- **THEN** the resource SHALL reject the configuration at plan/validate time

#### Scenario: client_type set to Android
- **WHEN** a user configures `device_profiles` with `client_type = "Android"`
- **THEN** the resource SHALL set `DeviceProfile.ClientType` to "Android" in the ModifySecurityPolicy API request

### Requirement: Score threshold fields
The `high_risk_min_score` and `medium_risk_min_score` fields SHALL be TypeInt (Optional, Computed) mapping to the SDK's `*uint64` fields, with API defaults of 50 and 15 respectively.

#### Scenario: User-specified high risk score
- **WHEN** a user configures `device_profiles` with `high_risk_min_score = 60`
- **THEN** the resource SHALL set `DeviceProfile.HighRiskMinScore` to 60 in the ModifySecurityPolicy API request

#### Scenario: Omitted scores default from API
- **WHEN** a user omits `high_risk_min_score` and `medium_risk_min_score`, and the DescribeSecurityPolicy response returns the API defaults (50 and 15)
- **THEN** the resource SHALL populate `high_risk_min_score = 50` and `medium_risk_min_score = 15` in state during Read

### Requirement: Request action fields
The `high_risk_request_action` and `medium_risk_request_action` fields SHALL be TypeList (MaxItems: 1, Optional) reusing the shared `SecurityAction` schema, identical in shape to the existing `invalid_attestation_action` field.

#### Scenario: High-risk action with Deny
- **WHEN** a user configures `high_risk_request_action` with `name = "Deny"` and `deny_action_parameters`
- **THEN** the resource SHALL set `DeviceProfile.HighRiskRequestAction.Name` to "Deny" and populate `DeviceProfile.HighRiskRequestAction.DenyActionParameters` in the API request

#### Scenario: Medium-risk action with Monitor
- **WHEN** a user configures `medium_risk_request_action` with `name = "Monitor"`
- **THEN** the resource SHALL set `DeviceProfile.MediumRiskRequestAction.Name` to "Monitor" with no additional action parameters

### Requirement: Read operation for device_profiles
The resource Read operation SHALL populate `device_profiles` from the `DescribeSecurityPolicy` API response, extracting `SecurityPolicy.BotManagement.ClientAttestationRules.Rules[].DeviceProfiles` and flattening each `DeviceProfile` (including its nested `SecurityAction` fields) into the Terraform state.

#### Scenario: Read with DeviceProfiles present
- **WHEN** the DescribeSecurityPolicy response contains a non-nil `DeviceProfiles` field on a `ClientAttestationRule`
- **THEN** the resource SHALL flatten each `DeviceProfile` into the `device_profiles` state attribute of that rule

#### Scenario: Read with DeviceProfiles absent
- **WHEN** the DescribeSecurityPolicy response has a nil `DeviceProfiles` field on a `ClientAttestationRule`
- **THEN** the resource SHALL not set the `device_profiles` attribute for that rule (it remains empty)

### Requirement: Nil checks in Read operation
The Read operation for `device_profiles` SHALL check for nil at each nesting level before accessing sub-fields, consistent with existing resource patterns. Specifically:
- Check `rule.DeviceProfiles != nil` before iterating the list
- For each `DeviceProfile`, check each scalar field (`ClientType`, `HighRiskMinScore`, `MediumRiskMinScore`) for nil before setting it
- Check `DeviceProfile.HighRiskRequestAction` / `DeviceProfile.MediumRiskRequestAction` for nil before flattening via `flattenSecurityAction()`

#### Scenario: Partial DeviceProfile response
- **WHEN** the DescribeSecurityPolicy response returns a `DeviceProfile` with only `ClientType` and `HighRiskMinScore` set (and the action fields are nil)
- **THEN** the resource SHALL flatten `client_type` and `high_risk_min_score` only, and SHALL NOT set `high_risk_request_action` or `medium_risk_request_action`

### Requirement: Create and Update operations for device_profiles
The resource Create and Update operations SHALL expand the `device_profiles` Terraform configuration into the `SecurityPolicy.BotManagement.ClientAttestationRules.Rules[].DeviceProfiles` field of the `ModifySecurityPolicy` API request, including all nested structures.

#### Scenario: Create with device_profiles
- **WHEN** a user creates a resource with a `client_attestation_rules` rule containing `device_profiles`
- **THEN** the Create operation SHALL expand the configuration and include `DeviceProfiles` on the corresponding `ClientAttestationRule` in the ModifySecurityPolicy request

#### Scenario: Update with device_profiles changes
- **WHEN** a user updates the resource's `device_profiles` configuration on an existing `client_attestation_rules` rule
- **THEN** the Update operation SHALL expand the new configuration and include `DeviceProfiles` on the corresponding `ClientAttestationRule` in the ModifySecurityPolicy request

#### Scenario: Expand omits device_profiles when absent
- **WHEN** a `client_attestation_rules` rule in Terraform state does not specify `device_profiles`
- **THEN** the expand logic SHALL leave `ClientAttestationRule.DeviceProfiles` unset (nil), so the cloud API preserves any existing device configuration
