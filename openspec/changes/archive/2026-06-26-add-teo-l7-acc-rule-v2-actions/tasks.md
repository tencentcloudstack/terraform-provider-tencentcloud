## 1. Schema Definition

- [x] 1.1 Add `advanced_origin_routing_parameters` optional parameter to `branches.actions` Schema map in `resource_tc_teo_l7_acc_rule_extension.go`, with `direction` field
- [x] 1.2 Add `shield_parameters` optional parameter to `branches.actions` Schema map, with `shield_space_id` field
- [x] 1.3 Add `site_failover_parameters` optional parameter to `branches.actions` Schema map, with `site_failover_status_codes` and `site_failover_params` fields
- [x] 1.4 Update `actions.name` description to include `AdvancedOriginRouting`, `Shield`, `SiteFailover` enum values

## 2. Create Method (GetBranchs - Flatten)

- [x] 2.1 In `resourceTencentCloudTeoL7AccRuleGetBranchs`, add logic to read `advanced_origin_routing_parameters` from resource data and populate `RuleEngineAction.AdvancedOriginRoutingParameters`
- [x] 2.2 Add logic to read `shield_parameters` and populate `RuleEngineAction.ShieldParameters`
- [x] 2.3 Add logic to read `site_failover_parameters` and populate `RuleEngineAction.SiteFailoverParameters`, including nested `SiteFailover` structure

## 3. Read Method (SetBranchs - Set)

- [x] 3.1 In `resourceTencentCloudTeoL7AccRuleSetBranchs`, add logic to extract `AdvancedOriginRoutingParameters` from API response and set to `advanced_origin_routing_parameters`
- [x] 3.2 Add logic to extract `ShieldParameters` and set to `shield_parameters`
- [x] 3.3 Add logic to extract `SiteFailoverParameters` and set to `site_failover_parameters`, including nested `SiteFailover` structure

## 4. Rollback Top-Level Actions Parameter

- [x] 4.1 Remove top-level `actions` parameter from resource schema in `resource_tc_teo_l7_acc_rule_v2.go`
- [x] 4.2 Remove top-level `actions` processing logic from Create method
- [x] 4.3 Remove top-level `actions` read logic from Read method
- [x] 4.4 Remove top-level `actions` processing logic from Update method
- [x] 4.5 Update `.changelog/4258.txt` to reflect the actual changes

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_l7_acc_rule_v2.md` to replace top-level actions example with AdvancedOriginRouting/Shield/SiteFailover example
- [x] 5.2 Update website documentation to add new parameter descriptions

## 6. Unit Tests

- [x] 6.1 Add schema validation tests for `advanced_origin_routing_parameters`
- [x] 6.2 Add schema validation tests for `shield_parameters`
- [x] 6.3 Add schema validation tests for `site_failover_parameters`
- [x] 6.4 Add flatten/set roundtrip tests for AdvancedOriginRoutingParameters
- [x] 6.5 Add flatten/set roundtrip tests for ShieldParameters
- [x] 6.6 Add flatten/set roundtrip tests for SiteFailoverParameters (basic, nested, private params, redirect URL)
- [x] 6.7 Fix unused variable compilation errors in test file
