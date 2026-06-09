## Why

The `tencentcloud_teo_function_rule` resource currently only supports the `direct` trigger type (where a single function is directly specified by `function_id`). The cloud API (`CreateFunctionRule`/`ModifyFunctionRule`) also supports `weight` (weight-based function selection) and `region` (region-based function selection) trigger types via the `TriggerType` parameter. Without this parameter, users cannot configure weight-based or region-based function routing in their Terraform configurations.

## What Changes

- Add `trigger_type` parameter (Optional, string) to the `tencentcloud_teo_function_rule` resource schema, supporting values `direct`, `weight`, and `region`. Defaults to `direct` when not specified.
- Update the Create handler to pass `TriggerType` to the `CreateFunctionRule` API request.
- Update the Update handler to include `trigger_type` in the mutable arguments list and pass it to the `ModifyFunctionRule` API request.
- Update the Read handler to read `TriggerType` from the `DescribeFunctionRules` API response and set it in the Terraform state.
- Update the resource documentation (`.md` file) to describe the new `trigger_type` parameter.

## Capabilities

### New Capabilities
- `teo-function-rule-trigger-type`: Adds `trigger_type` parameter to the `tencentcloud_teo_function_rule` resource, enabling users to specify the function selection configuration type (direct, weight, or region).

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_function_rule.go` - Schema definition and CRUD handlers
- `tencentcloud/services/teo/resource_tc_teo_function_rule.md` - Resource documentation
- `tencentcloud/services/teo/resource_tc_teo_function_rule_test.go` - Unit tests
- Backward compatible: `trigger_type` is Optional with no ForceNew, existing configurations continue to work unchanged
