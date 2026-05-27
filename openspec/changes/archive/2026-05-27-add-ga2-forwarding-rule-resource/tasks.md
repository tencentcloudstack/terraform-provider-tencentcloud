## 1. Service Layer

- [x] 1.1 Add `DescribeGa2ForwardingRuleById` method to `tencentcloud/services/ga2/service_tencentcloud_ga2.go` that calls DescribeForwardingRule API with pagination (Limit=100) and filters by forwarding_rule_id to return a single `*ga2v20250115.ForwardingRuleSet`

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` with schema definition including: global_accelerator_id (Required, ForceNew), listener_id (Required, ForceNew), forwarding_policy_id (Required, ForceNew), rule_conditions (Required, TypeList), rule_actions (Required, TypeList), origin_headers (Optional, TypeList), enable_origin_sni (Optional, Bool), origin_sni (Optional, String), origin_host (Optional, String), forwarding_rule_id (Computed), task_id (Computed). Include Timeouts block (Create/Update/Delete: 20 minutes) and Importer
- [x] 2.2 Implement `resourceTencentCloudGa2ForwardingRuleCreate` function: call CreateForwardingRule API with retry, validate response and ForwardingRuleId not nil, wait for task via WaitForGa2TaskFinish, set composite ID, call Read
- [x] 2.3 Implement `resourceTencentCloudGa2ForwardingRuleRead` function: parse composite ID, call DescribeGa2ForwardingRuleById, handle not-found by clearing ID, set all attributes from ForwardingRuleSet response (check nil before setting)
- [x] 2.4 Implement `resourceTencentCloudGa2ForwardingRuleUpdate` function: detect changes with HasChange, call ModifyForwardingRule API with retry, wait for task via WaitForGa2TaskFinish, call Read
- [x] 2.5 Implement `resourceTencentCloudGa2ForwardingRuleDelete` function: parse composite ID, call DeleteForwardingRule API with retry, wait for task via WaitForGa2TaskFinish
- [x] 2.6 Implement helper functions: `parseGa2ForwardingRuleId` (split composite ID into 4 parts), `buildRuleConditions`/`flattenRuleConditions`, `buildRuleActions`/`flattenRuleActions`, `buildOriginHeaders`/`flattenOriginHeaders`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_ga2_forwarding_rule` resource in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add `tencentcloud_ga2_forwarding_rule` entry to `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.md` with Example Usage section showing a complete forwarding rule configuration, and Import section explaining the composite ID format (global_accelerator_id#listener_id#forwarding_policy_id#forwarding_rule_id)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule_test.go` with gomonkey-based unit tests covering: Create (success, nil response error), Read (success, not found), Update (success), Delete (success). Mock the GA2 SDK client methods and verify correct API calls and state handling. Run tests with `go test -gcflags=all=-l`
