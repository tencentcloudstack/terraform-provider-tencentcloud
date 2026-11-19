## ADDED Requirements

### Requirement: frequent_scanning_protection schema definition
The `tencentcloud_teo_security_policy_config` resource SHALL include an Optional `frequent_scanning_protection` parameter (TypeList, MaxItems: 1, Computed) under the `managed_rules` block (inside `security_policy`). Each element SHALL expose:
- `enabled` (TypeString, Optional, Computed): Whether the high-frequency scanning protection rule is enabled. Valid values: `on`, `off`.
- `id` (TypeString, Computed): Rule ID, output-only.
- `action` (TypeList, MaxItems: 1, Optional, Computed): The disposal action, reusing the shared `SecurityAction` schema (Name supports `Deny`, `Monitor`, `JSChallenge`).
- `count_by` (TypeString, Optional, Computed): Request statistical matching method. Valid values: `http.request.xff_header_ip`, `http.request.ip`.
- `block_threshold` (TypeInt, Optional, Computed): Cumulative block count threshold within the counting period. Range 1-4294967294.
- `counting_period` (TypeString, Optional, Computed): Statistics time window. Range 5-1800, unit seconds (s).
- `action_duration` (TypeString, Optional, Computed): Action duration. Range 60-86400, unit seconds (s).

#### Scenario: Resource accepts frequent_scanning_protection configuration
- **WHEN** a user provides a `frequent_scanning_protection` block inside the `managed_rules` of `tencentcloud_teo_security_policy_config`
- **THEN** the resource SHALL accept and process the configuration without errors

#### Scenario: frequent_scanning_protection is optional
- **WHEN** a user creates a `tencentcloud_teo_security_policy_config` resource with a `managed_rules` block that omits `frequent_scanning_protection`
- **THEN** the resource SHALL be created successfully and the existing behavior SHALL be preserved

### Requirement: enabled field
The `enabled` field SHALL be an Optional TypeString that controls whether the high-frequency scanning protection rule is enabled. Valid values SHALL be `on`, `off`.

#### Scenario: enabled set to on
- **WHEN** a user configures `frequent_scanning_protection` with `enabled = "on"`
- **THEN** the resource SHALL set `FrequentScanningProtection.Enabled` to "on" in the ModifySecurityPolicy API request

### Requirement: id field is output-only
The `id` field SHALL be a Computed-only TypeString. The resource SHALL populate it from the `DescribeSecurityPolicy` response and SHALL NOT send it in the `ModifySecurityPolicy` request.

#### Scenario: Read populates id
- **WHEN** the DescribeSecurityPolicy response returns a non-nil `FrequentScanningProtection.Id`
- **THEN** the resource SHALL flatten `id` into state during Read

#### Scenario: Expand omits id
- **WHEN** a user creates or updates the resource with `frequent_scanning_protection`
- **THEN** the expand logic SHALL NOT set `FrequentScanningProtection.Id` in the API request

### Requirement: action field reuses SecurityAction schema
The `action` field SHALL be a TypeList (MaxItems: 1, Optional, Computed) reusing the shared `SecurityAction` schema, identical in shape to other `*SecurityAction` fields in the resource (e.g. `invalid_attestation_action`).

#### Scenario: Action with Deny
- **WHEN** a user configures `action` with `name = "Deny"` and `deny_action_parameters`
- **THEN** the resource SHALL set `FrequentScanningProtection.Action.Name` to "Deny" and populate `FrequentScanningProtection.Action.DenyActionParameters` in the API request

#### Scenario: Action with Monitor
- **WHEN** a user configures `action` with `name = "Monitor"`
- **THEN** the resource SHALL set `FrequentScamingProtection.Action.Name` to "Monitor" with no additional action parameters

### Requirement: Read operation for frequent_scanning_protection
The resource Read operation SHALL populate `frequent_scanning_protection` from the `DescribeSecurityPolicy` API response, extracting `SecurityPolicy.ManagedRules.FrequentScanningProtection` and flattening its fields (including the nested `Action` via `flattenSecurityAction()`) into the Terraform state.

#### Scenario: Read with FrequentScanningProtection present
- **WHEN** the DescribeSecurityPolicy response contains a non-nil `FrequentScanningProtection` field on `ManagedRules`
- **THEN** the resource SHALL flatten its fields into the `frequent_scanning_protection` state attribute

#### Scenario: Read with FrequentScanningProtection absent
- **WHEN** the DescribeSecurityPolicy response has a nil `FrequentScanningProtection` field on `ManagedRules`
- **THEN** the resource SHALL not set the `frequent_scanning_protection` attribute (it remains empty)

### Requirement: Nil checks in Read operation
The Read operation for `frequent_scanning_protection` SHALL check for nil at each nesting level before accessing sub-fields, consistent with existing resource patterns. Specifically:
- Check `ManagedRules.FrequentScanningProtection != nil` before building the map
- For each scalar field (`Enabled`, `Id`, `CountBy`, `CountingPeriod`, `ActionDuration`, `BlockThreshold`), check the pointer for nil before setting it
- Check `FrequentScanningProtection.Action != nil` before flattening via `flattenSecurityAction()`

#### Scenario: Partial FrequentScanningProtection response
- **WHEN** the DescribeSecurityPolicy response returns a `FrequentScanningProtection` with only `Enabled` and `Action` set (and the scalar threshold fields are nil)
- **THEN** the resource SHALL flatten `enabled` and `action` only, and SHALL NOT set `block_threshold`, `count_by`, `counting_period`, or `action_duration`

### Requirement: Create and Update operations for frequent_scanning_protection
The resource Create and Update operations SHALL expand the `frequent_scanning_protection` Terraform configuration into the `SecurityPolicy.ManagedRules.FrequentScanningProtection` field of the `ModifySecurityPolicy` API request, including all nested structures (except the output-only `Id`).

#### Scenario: Create with frequent_scanning_protection
- **WHEN** a user creates a resource with a `managed_rules` block containing `frequent_scanning_protection`
- **THEN** the Create operation SHALL expand the configuration and include `FrequentScanningProtection` on `ManagedRules` in the ModifySecurityPolicy request

#### Scenario: Update with frequent_scanning_protection changes
- **WHEN** a user updates the resource's `frequent_scanning_protection` configuration on an existing `managed_rules` block
- **THEN** the Update operation SHALL expand the new configuration and include `FrequentScanningProtection` on `ManagedRules` in the ModifySecurityPolicy request

#### Scenario: Expand omits frequent_scanning_protection when absent
- **WHEN** a `managed_rules` block in Terraform state does not specify `frequent_scanning_protection`
- **THEN** the expand logic SHALL leave `ManagedRules.FrequentScanningProtection` unset (nil), so the cloud API preserves any existing configuration