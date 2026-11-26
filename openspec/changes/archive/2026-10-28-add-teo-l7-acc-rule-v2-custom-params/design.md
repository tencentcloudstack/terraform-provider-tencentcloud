## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource (file `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`) manages EdgeOne L7 acceleration rules via three cloud APIs in package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`:
- `CreateL7AccRules` — request body `Rules []*RuleEngineItem`
- `ModifyL7AccRule` — request body `Rule *RuleEngineItem`
- `DescribeL7AccRules` — response body `Rules []*RuleEngineItem`

The resource schema is split across two files:
- `resource_tc_teo_l7_acc_rule_v2.go` — top-level schema (`zone_id`, `status`, `rule_name`, `description`, `branches`, `rule_id`, `rule_priority`) plus the CRUD functions.
- `resource_tc_teo_l7_acc_rule_extension.go` — the recursive `TencentTeoL7RuleBranchBasicInfo(depth)` helper that defines the `branches`/`sub_rules` schema and the `actions` block (40+ action parameters), plus the matching expand (`resourceTencentCloudTeoL7AccRuleGetBranchs`) and flatten (`resourceTencentCloudTeoL7AccRuleSetBranchs`) functions.

The existing `RuleEngineItem` schema already surfaces `Status`, `RuleId`, `RuleName`, `Description`, `Branches`, and `RulePriority`. The `actions` block already covers `advanced_origin_routing_parameters`, `shield_parameters`, and `site_failover_parameters` (per the prior `teo-l7-acc-rule-v2-actions` capability). This change adds the remaining `RuleEngineItem` fields — `CustomVariables` and `CustomVariableOperations` — plus the missing `CustomActionParameters` action.

## Goals / Non-Goals

**Goals:**
- Surface `CustomVariables`, `CustomVariableOperations` (top-level) and `CustomActionParameters` (under `branches.actions`) as Optional Terraform parameters, wired through Create/Read/Update.
- Preserve backward compatibility: all new parameters are Optional; omitted parameters produce no change in API behavior.
- Follow the existing extension file patterns (expand-to-struct / flatten-to-map) so the new code integrates cleanly with the recursive `sub_rules` handling.

**Non-Goals:**
- Do NOT change `rule_id` or `rule_priority` from Computed to settable (the cloud API marks `RulePriority` as output-only and `RuleId` is server-generated on create).
- Do NOT alter any existing action parameter schema.
- Do NOT add new cloud API dependencies — all structs already exist in the vendored SDK.

## Decisions

### Decision 1: Add `custom_variables` and `custom_variable_operations` as top-level schema entries in `resource_tc_teo_l7_acc_rule_v2.go`

These map directly to `RuleEngineItem.CustomVariables` and `RuleEngineItem.CustomVariableOperations`, which are siblings of `Branches` (not nested inside it). Placing them at the top level matches the SDK structure and the existing `branches`/`description` placement.

**Alternatives considered**: Putting them inside the extension helper. Rejected because the extension helper defines the `branches`/`actions`/`sub_rules` recursion — these fields live one level above, on `RuleEngineItem` itself.

### Decision 2: Handle `custom_variable_operations` expand/flatten in dedicated helpers (not the generic branch helpers)

`CustomVariableOperations` uses distinct SDK types (`CustomVariableOperation`, `CustomVariableOperationRuleBranch`, `CustomVariableOperationRuleAction`, `CustomVariableOperationSubRule`) that differ from the regular rule branch types (`RuleBranch`, `RuleEngineAction`, `RuleEngineSubRule`). Although the JSON shape mirrors `branches`, the Go types are separate, so dedicated helpers avoid type-mismatch compile errors.

### Decision 3: Add `custom_action_parameters` to the `actions` block inside `TencentTeoL7RuleBranchBasicInfo`

This follows the existing pattern where each `RuleEngineAction` field gets a matching `*_parameters` schema entry and matching expand/flatten blocks alongside the others (`shield_parameters`, `site_failover_parameters`, etc.).

### Decision 4: `CustomActionParameter` value handling — single `value` string field

The SDK `CustomActionParameter` struct has 7 typed value fields driven by `ValueType` (`StringValue`, `IntegerValue`, `FloatValue`, `BooleanValue`, `StringArrayValue`, `IntegerArrayValue`, `FloatArrayValue`). Exposing all 7 would bloat the schema and conflict with the requirement's single `Value` SchemaName. Instead, expose one `value` (TypeString) field per parameter element. The **expand** logic sets the SDK field corresponding to `ValueType`; the **flatten** logic serializes whichever typed field is populated back to the string. This keeps the schema simple while remaining faithful to the API.

### Decision 5: Wire Update through the existing single `ModifyL7AccRule` call

The existing `ResourceTencentCloudTeoL7AccRuleV2Update` already sends a single `ModifyL7AccRule` request whenever `status`/`rule_name`/`description`/`branches` change (overwrite semantics). Extend the `HasChange` guard and the rule-building block to also expand `custom_variables` and `custom_variable_operations` into `request.Rule`. The cloud API uses覆盖式 (overwrite) update, so the full rule content is rebuilt on each update — consistent with current behavior.

## Risks / Trade-offs

- **[Risk] `CustomActionParameter.ValueType` not exposed in schema** → The expand logic derives the SDK value field from `ValueType`. Because `ValueType` is itself a field of `CustomActionParameter` but not requested in the mapping, the implementation will set `StringValue` by default and rely on the cloud API to coerce; alternatively, expose `ValueType` as an optional discriminator if the API rejects generic strings. **Mitigation**: Implement the value field selection by inspecting whichever typed field the flatten populated; if the API requires `ValueType`, add it as an Optional schema field during implementation (covered by the `name`/`value` mapping — the implementation step will verify against the SDK).
- **[Risk] Deeply nested `sub_rules` recursion in `custom_variable_operations`** → The SDK only supports one layer of `SubRules` inside `CustomVariableOperations.Branches`. The schema mirrors this with a finite nesting (no depth-parameterized recursion needed). **Mitigation**: Define `sub_rules` with a fixed single-layer `branches` (no further recursion) matching the SDK constraint.
- **[Risk] Large diff in `resource_tc_teo_l7_acc_rule_extension.go`** → The extension file is already large (~3200 lines). Adding `custom_action_parameters` expand/flatten increases it further. **Mitigation**: Keep the new code colocated with the existing action handlers; no refactor of existing code.

## Reference: Vendored SDK struct definitions

The following struct definitions (from `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`) are the source of truth for the new parameters. Implementation should reference these rather than re-reading the vendor file.

### `RuleEngineItem` (models.go:25625) — the top-level rule struct
```go
type RuleEngineItem struct {
	Status                   *string                   `json:"Status,omitnil,omitempty"`
	RuleId                   *string                   `json:"RuleId,omitnil,omitempty"`
	RuleName                 *string                   `json:"RuleName,omitnil,omitempty"`
	Description              []*string                 `json:"Description,omitnil,omitempty"`
	CustomVariables          []*CustomVariable         `json:"CustomVariables,omitnil,omitempty"`          // NEW
	CustomVariableOperations []*CustomVariableOperation `json:"CustomVariableOperations,omitnil,omitempty"` // NEW
	Branches                 []*RuleBranch             `json:"Branches,omitnil,omitempty"`
	RulePriority             *int64                    `json:"RulePriority,omitnil,omitempty"`             // output-only
}
```

### `CustomVariable` (models.go:6211) — element of `CustomVariables`
```go
type CustomVariable struct {
	Name          *string `json:"Name,omitnil,omitempty"`          // user.rule.* prefix
	InitialValue  *string `json:"InitialValue,omitnil,omitempty"`
	Description   *string `json:"Description,omitnil,omitempty"`
}
```

### `CustomVariableOperation` (models.go:6222) — element of `CustomVariableOperations`
```go
type CustomVariableOperation struct {
	Branches    []*CustomVariableOperationRuleBranch `json:"Branches,omitnil,omitempty"`
	Description []*string                            `json:"Description,omitnil,omitempty"`
}
```

### `CustomVariableOperationRuleBranch` (models.go:6240)
```go
type CustomVariableOperationRuleBranch struct {
	Condition *string                              `json:"Condition,omitnil,omitempty"`
	Actions   []*CustomVariableOperationRuleAction `json:"Actions,omitnil,omitempty"`
	SubRules  []*CustomVariableOperationSubRule    `json:"SubRules,omitnil,omitempty"`
}
```

### `CustomVariableOperationRuleAction` (models.go:6231)
```go
type CustomVariableOperationRuleAction struct {
	Name         *string        `json:"Name,omitnil,omitempty"` // only "Set"
	SetParameters *SetParameters `json:"SetParameters,omitnil,omitempty"`
}
```

### `CustomVariableOperationSubRule` (models.go:6253)
```go
type CustomVariableOperationSubRule struct {
	Branches    []*CustomVariableOperationRuleBranch `json:"Branches,omitnil,omitempty"`
	Description []*string                            `json:"Description,omitnil,omitempty"`
}
```

### `SetParameters` (models.go:26037) — used by custom variable operation actions
```go
type SetParameters struct {
	Name  *string `json:"Name,omitnil,omitempty"`
	Value *string `json:"Value,omitnil,omitempty"`
}
```

### `RuleEngineAction` (models.go:25389) — the action struct inside `RuleBranch.Actions`
```go
type RuleEngineAction struct {
	Name                     *string                     `json:"Name,omitnil,omitempty"` // incl. "CustomAction"
	// ... 40 existing *_parameters fields ...
	CustomActionParameters   *CustomActionParameters     `json:"CustomActionParameters,omitnil,omitempty"`   // NEW
}
```

### `CustomActionParameters` (models.go:6090)
```go
type CustomActionParameters struct {
	CustomActions []*CustomAction `json:"CustomActions,omitnil,omitempty"`
}
```

### `CustomAction` (models.go:6053)
```go
type CustomAction struct {
	ActionId   *string                 `json:"ActionId,omitnil,omitempty"`
	Parameters []*CustomActionParameter `json:"Parameters,omitnil,omitempty"`
}
```

### `CustomActionParameter` (models.go:6061) — typed-value union
```go
type CustomActionParameter struct {
	Name               *string   `json:"Name,omitnil,omitempty"`
	ValueType          *string   `json:"ValueType,omitnil,omitempty"` // String|Integer|Float|Boolean|ArrayOfString|ArrayOfInteger|ArrayOfFloat
	StringValue        *string   `json:"StringValue,omitnil,omitempty"`
	IntegerValue       *int64    `json:"IntegerValue,omitnil,omitempty"`
	FloatValue         *float64  `json:"FloatValue,omitnil,omitempty"`
	BooleanValue       *bool     `json:"BooleanValue,omitnil,omitempty"`
	StringArrayValue   []*string `json:"StringArrayValue,omitnil,omitempty"`
	IntegerArrayValue  []*int64  `json:"IntegerArrayValue,omitnil,omitempty"`
	FloatArrayValue    []*float64 `json:"FloatArrayValue,omitnil,omitempty"`
}
```

## Reference: Existing schema structure (relevant excerpts)

### Top-level schema (`resource_tc_teo_l7_acc_rule_v2.go:25-69`)
The `branches` element uses `TencentTeoL7RuleBranchBasicInfo(1)`. Create builds `request.Rules = []*RuleEngineItem{rule}`; Update builds `request.Rule = rule` with `rule.RuleId` set; Read calls `service.DescribeTeoL7AccRuleById` and flattens `respData.Rules[0]`.

### `TencentTeoL7RuleBranchBasicInfo` actions block (`resource_tc_teo_l7_acc_rule_extension.go:20-1416`)
Defines `condition`, `actions` (with `name` + ~40 `*_parameters` entries), and `sub_rules` (recursive when `depth < 8`). The `name` description currently lists `AdvancedOriginRouting`, `Shield`, `SiteFailover` but NOT `CustomAction`.

### Expand function `resourceTencentCloudTeoL7AccRuleGetBranchs` (`:1450+`)
Iterates `branches`, builds `RuleBranch` with `Condition`, `Actions []*RuleEngineAction`, `SubRules`. Each action parameter is expanded from `actionsMap["<param>"]` into the corresponding SDK struct pointer. New `custom_action_parameters` expand block goes alongside `shield_parameters`/`site_failover_parameters`.

### Flatten function `resourceTencentCloudTeoL7AccRuleSetBranchs` (`:2900+`)
Iterates SDK `RuleBranch`, builds `actionsMap` from each `RuleEngineAction` field with nil checks. New `custom_action_parameters` flatten block goes alongside `shield_parameters`/`site_failover_parameters`.

## Migration Plan

No migration required. All new parameters are Optional and additive. Existing Terraform configurations and state files remain valid. The cloud API ignores omitted optional fields.

## Open Questions

- **`CustomActionParameter.ValueType`**: The requirement mapping lists only `Name` and `Value` as the schema fields for each `parameters` element. The SDK requires `ValueType` to know which typed value field to read. The implementation will set `ValueType` based on the value content (defaulting to `String`) OR expose `value_type` as an additional optional discriminator if the API rejects the default. This will be verified during implementation by checking the cloud API behavior; the spec's `name`/`value` mapping is the baseline.
