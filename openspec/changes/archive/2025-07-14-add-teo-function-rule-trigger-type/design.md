## Context

The `tencentcloud_teo_function_rule` resource manages EdgeOne function trigger rules. Currently, the resource only supports the `direct` trigger type where a single function is specified by `function_id`. The TencentCloud TEO API (`CreateFunctionRule`/`ModifyFunctionRule`) supports a `TriggerType` parameter that controls how functions are selected for execution:

- `direct`: Directly specify a single function (current behavior, default)
- `weight`: Select functions based on weight ratios
- `region`: Select functions based on client IP country/region

The `TriggerType` field is present in `CreateFunctionRuleRequest`, `ModifyFunctionRuleRequest`, and `FunctionRule` (read response) in the cloud SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).

## Goals / Non-Goals

**Goals:**
- Add `trigger_type` as an Optional string parameter to `tencentcloud_teo_function_rule` resource
- Support `trigger_type` in Create, Read, and Update operations
- Maintain full backward compatibility — existing TF configurations without `trigger_type` continue to work (API defaults to `direct` when not specified)

**Non-Goals:**
- Adding `RegionMappingSelections` and `WeightedSelections` parameters (these are complex nested objects that depend on `trigger_type` and will be added in a separate change)
- Changing the composite ID format or existing schema fields
- Modifying the Delete operation (it does not use `TriggerType`)

## Decisions

1. **Schema: `trigger_type` as Optional string without ForceNew**
   - The `TriggerType` field is present in both Create and Modify API requests, so it is mutable (not ForceNew).
   - Default value is `direct` per API behavior — we do NOT set `Default:` in the schema to avoid spurious diffs; instead we let the API default handle it.
   - Validate values using `ValidateFunc` with `validation.StringInSlice([]string{"direct", "weight", "region"}, false)`.

2. **CRUD handler modifications**
   - **Create**: Read `trigger_type` from schema and set `request.TriggerType` if specified.
   - **Read**: Read `respData.TriggerType` from `FunctionRule` response and set `trigger_type` in state (with nil check).
   - **Update**: Add `"trigger_type"` to `mutableArgs` list, and set `request.TriggerType` if specified in the Modify request.

3. **Extension file not needed**
   - The change is simple enough to be made directly in the main resource file without creating an `_extension.go` file.

## Risks / Trade-offs

- **[Risk] `trigger_type` set to `weight` or `region` without corresponding selection configs** → The API may reject the request if `RegionMappingSelections` or `WeightedSelections` are required but not provided. However, since this change only adds `trigger_type` and not the selection configs, users who set `trigger_type` to `weight` or `region` may encounter API errors. This is acceptable since the selection configs will be added in a follow-up change. The parameter description should note this dependency.
- **[Risk] API returns empty/null `TriggerType` for existing rules** → The Read handler checks for nil before setting the field, so existing state will not be affected.
