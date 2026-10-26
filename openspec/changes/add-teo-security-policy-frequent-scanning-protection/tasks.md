## 1. Schema Definition

- [x] 1.1 Add `frequent_scanning_protection` field (TypeList, Optional, Computed, MaxItems: 1) under the `managed_rules` schema block (inside `security_policy`) in `resource_tc_teo_security_policy_config.go`, placed as a sibling of `managed_rule_groups`
- [x] 1.2 Add `enabled` sub-field (TypeString, Optional, Computed) with valid values `on`, `off`
- [x] 1.3 Add `id` sub-field (TypeString, Computed only) — output-only rule ID, not sent in requests
- [x] 1.4 Add `action` sub-field (TypeList, MaxItems: 1, Optional, Computed) with `Elem: securityActionSchema()`
- [x] 1.5 Add `count_by` sub-field (TypeString, Optional, Computed) with valid values `http.request.xff_header_ip`, `http.request.ip`
- [x] 1.6 Add `block_threshold` sub-field (TypeInt, Optional, Computed) — range 1-4294967294
- [x] 1.7 Add `counting_period` sub-field (TypeString, Optional, Computed) — range 5-1800, unit seconds
- [x] 1.8 Add `action_duration` sub-field (TypeString, Optional, Computed) — range 60-86400, unit seconds

## 2. Read (Flatten) Implementation

- [x] 2.1 In the `ManagedRules` flatten block (after the `ManagedRuleGroups` handling, before `securityPolicyMap["managed_rules"] = ...`), add a nil guard `if respData.ManagedRules.FrequentScanningProtection != nil` and build `frequentScanningProtectionMap`
- [x] 2.2 Set `enabled` from `FrequentScanningProtection.Enabled` (with nil check)
- [x] 2.3 Set `id` from `FrequentScanningProtection.Id` (with nil check)
- [x] 2.4 Set `action` via `flattenSecurityAction()` when `FrequentScanningProtection.Action != nil`, wrapped in `[]interface{}` when non-empty
- [x] 2.5 Set `count_by` from `FrequentScanningProtection.CountBy` (with nil check)
- [x] 2.6 Set `block_threshold` from `FrequentScanningProtection.BlockThreshold` (with nil check)
- [x] 2.7 Set `counting_period` from `FrequentScanningProtection.CountingPeriod` (with nil check)
- [x] 2.8 Set `action_duration` from `FrequentScanningProtection.ActionDuration` (with nil check)
- [x] 2.9 Set `frequentScanningProtectionMap` into `managedRulesMap["frequent_scanning_protection"] = []interface{}{frequentScanningProtectionMap}`

## 3. Create/Update (Expand) Implementation

- [x] 3.1 In the `ManagedRules` expand block (after the `managed_rule_groups` handling, before `securityPolicy.ManagedRules = &managedRules`), read `managedRulesMap["frequent_scanning_protection"]` and convert to head map
- [x] 3.2 When present, construct a `teov20220901.FrequentScanningProtection{}` and set `Enabled` via `helper.String()` when `enabled` is non-empty
- [x] 3.3 Set `CountBy`, `CountingPeriod`, `ActionDuration` via `helper.String()` when non-empty
- [x] 3.4 Set `BlockThreshold` via `helper.IntInt64()` when present (non-zero)
- [x] 3.5 Set `Action` via `buildSecurityActionFromMap()` when the `action` map is present and non-empty
- [x] 3.6 Do NOT set `Id` (output-only) in expand
- [x] 3.7 Assign the constructed struct to `managedRules.FrequentScanningProtection`

## 4. Documentation

- [ ] 4.1 Run `make doc` in the finalize stage to regenerate `website/docs/` and the resource `.md` documentation, including example usage showing `frequent_scanning_protection` (with `enabled`, `action`, `count_by`, `block_threshold`, `counting_period`, `action_duration`) inside `managed_rules`

  > Deferred to the tfpacer-finalize skill stage (`make doc` and `website/` edits are only permitted there per project rules). Schema/code/flatten/expand implementation is complete.

## 5. Unit Tests

- [x] 5.1 Add unit test functions for `frequent_scanning_protection` flatten logic using the gomonkey mock approach in `resource_tc_teo_security_policy_config_test.go`
- [x] 5.2 Add unit test functions for `frequent_scanning_protection` expand logic using the gomonkey mock approach
- [x] 5.3 Add a unit test covering the nil `FrequentScanningProtection` response path (state stays empty)

## 6. Verification

- [x] 6.1 Verify the code compiles without errors (visual inspection of code correctness)
- [x] 6.2 Verify all new schema fields are Optional/Computed and backward compatible
- [x] 6.3 Verify nil checks are present at every nesting level in the Read flatten logic
- [x] 6.4 Verify the cloud API `FrequentScanningProtection` field/struct used exist in the vendored SDK
- [x] 6.5 Verify the expand logic does NOT send the output-only `Id` field