## 1. Resource Implementation

- [x] 1.1 Update resource schema in `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` to ensure all parameters are correctly defined with proper types, constraints (Required/Optional/Computed/ForceNew), and descriptions
- [x] 1.2 Update Create function to properly validate API response (check nil response, empty RuleIds) and return NonRetryableError on empty ID
- [x] 1.3 Update Read function to log resource ID before clearing state when API returns empty response
- [x] 1.4 Verify Update function correctly handles HasChange detection and passes all modified fields to ModifyL7AccRule API
- [x] 1.5 Verify Delete function correctly passes zone_id and rule_id to DeleteL7AccRules API

## 2. Provider Registration

- [x] 2.1 Verify resource is registered in `tencentcloud/provider.go`
- [x] 2.2 Verify resource is listed in `tencentcloud/provider.md`

## 3. Testing

- [x] 3.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` using gomonkey to mock cloud API calls, covering Create/Read/Update/Delete operations

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` with example usage showing all parameters including zone_id, status, rule_name, description, and branches
