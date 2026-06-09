## 1. Schema Definition

- [x] 1.1 Add `trigger_type` field to the `tencentcloud_teo_function_rule` resource schema in `tencentcloud/services/teo/resource_tc_teo_function_rule.go`: Type `schema.TypeString`, Optional, with `ValidateFunc: validation.StringInSlice([]string{"direct", "weight", "region"}, false)`, and Description documenting valid values and that it defaults to `direct`

## 2. CRUD Handler Updates

- [x] 2.1 Update `resourceTencentCloudTeoFunctionRuleCreate`: Read `trigger_type` from schema and set `request.TriggerType` if the value is specified
- [x] 2.2 Update `resourceTencentCloudTeoFunctionRuleRead`: After calling `DescribeTeoFunctionRuleById`, check if `respData.TriggerType` is not nil and set `trigger_type` in state via `d.Set("trigger_type", respData.TriggerType)`
- [x] 2.3 Update `resourceTencentCloudTeoFunctionRuleUpdate`: Add `"trigger_type"` to the `mutableArgs` list, and in the Modify request block, set `request.TriggerType` if the value is specified

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_function_rule.md` to add `trigger_type` parameter description and usage example

## 4. Unit Tests

- [x] 4.1 Add unit test cases in `tencentcloud/services/teo/resource_tc_teo_function_rule_test.go` to verify Create/Read/Update handlers correctly handle the `trigger_type` parameter, including: setting trigger_type in create request, reading trigger_type from response, and updating trigger_type in modify request
