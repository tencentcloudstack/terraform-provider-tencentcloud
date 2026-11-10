## Context

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policies via the cloud API's `DescribeSecurityPolicy` (read) and `ModifySecurityPolicy` (create/update) interfaces. Both interfaces share the same `SecurityPolicy` struct in the vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).

Within `SecurityPolicy.ManagedRules`, the resource already exposes `enabled`, `detection_only`, `semantic_analysis`, `auto_update`, `managed_rule_groups`. The SDK's `ManagedRules` struct additionally has a `FrequentScanningProtection *FrequentScanningProtection` field that is currently NOT exposed in the Terraform schema, flatten (read), or expand (create/update) logic.

### Vendored SDK struct definitions (reference for implementation)

**`FrequentScanningProtection`** (models.go:17319-17340):
```go
type FrequentScanningProtection struct {
    // 高频扫描防护规则是否开启。取值有：<li>on：开启；</li><li>off：关闭。</li>
    Enabled *string `json:"Enabled,omitnil,omitempty" name:"Enabled"`
    // 高频扫描防护的规则 ID，仅作为出参返回。
    Id *string `json:"Id,omitnil,omitempty" name:"Id"`
    // 高频扫描防护的处置动作。当 Enabled 为 on 时，此字段必填。
    Action *SecurityAction `json:"Action,omitnil,omitempty" name:"Action"`
    // 请求统计的匹配方式，当 Enabled 为 on 时，此字段必填。
    // 取值：http.request.xff_header_ip / http.request.ip
    CountBy *string `json:"CountBy,omitnil,omitempty" name:"CountBy"`
    // 阈值，累计拦截次数，取值范围 1 ~ 4294967294。
    BlockThreshold *int64 `json:"BlockThreshold,omitnil,omitempty" name:"BlockThreshold"`
    // 统计时间窗口，取值 5 ~ 1800，单位秒（s）。
    CountingPeriod *string `json:"CountingPeriod,omitnil,omitempty" name:"CountingPeriod"`
    // 处置动作持续时长，取值范围 60 ~ 86400，单位秒（s）。
    ActionDuration *string `json:"ActionDuration,omitnil,omitempty" name:"ActionDuration"`
}
```

**`ManagedRules`** (models.go:19167-19185) — the parent struct:
```go
type ManagedRules struct {
    Enabled *string `json:"Enabled,omitnil,omitempty" name:"Enabled"`
    DetectionOnly *string `json:"DetectionOnly,omitnil,omitempty" name:"DetectionOnly"`
    SemanticAnalysis *string `json:"SemanticAnalysis,omitnil,omitempty" name:"SemanticAnalysis"`
    AutoUpdate *ManagedRuleAutoUpdate `json:"AutoUpdate,omitnil,omitempty" name:"AutoUpdate"`
    ManagedRuleGroups []*ManagedRuleGroup `json:"ManagedRuleGroups,omitnil,omitempty" name:"ManagedRuleGroups"`
    // 高频扫描防护配置选项
    FrequentScanningProtection *FrequentScanningProtection `json:"FrequentScanningProtection,omitnil,omitempty" name:"FrequentScanningProtection"`
}
```

**`SecurityAction`** (models.go:25849-25870) — reused via existing `securityActionSchema()` / `flattenSecurityAction()` / `buildSecurityActionFromMap()` helpers:
```go
type SecurityAction struct {
    Name *string `json:"Name,omitnil,omitempty" name:"Name"`
    DenyActionParameters *DenyActionParameters `json:"DenyActionParameters,omitnil,omitempty"`
    RedirectActionParameters *RedirectActionParameters `json:"RedirectActionParameters,omitnil,omitempty"`
    AllowActionParameters *AllowActionParameters `json:"AllowActionParameters,omitnil,omitempty"`
    ChallengeActionParameters *ChallengeActionParameters `json:"ChallengeActionParameters,omitnil,omitempty"`
    BlockIPActionParameters *BlockIPActionParameters `json:"BlockIPActionParameters,omitnil,omitempty"`
    ReturnCustomPageActionParameters *ReturnCustomPageActionParameters `json:"ReturnCustomPageActionParameters,omitnil,omitempty"`
}
```

The `Action.Name` for `FrequentScanningProtection` supports: `Deny`, `Monitor`, `JSChallenge` (per the SDK documentation). The shared `securityActionSchema()` helper already covers all these action types (deny/redirect/allow/challenge action parameters), so it can be reused directly.

### Existing resource code locations (reference for implementation)

The resource file `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` is ~7258 lines. Key existing code locations:

- **Schema**: The `managed_rules` schema is defined inline at lines 365-609 (inside the `security_policy` top-level block). The block ends after the `meta_data` field (line ~609). The new `frequent_scanning_protection` field will be added as a sibling of `managed_rule_groups` within this schema.
- **Read (flatten)**: The `ManagedRules` flatten logic is at lines 3113-3297. It builds `managedRulesMap` and sets `securityPolicyMap["managed_rules"]`. The new `FrequentScanningProtection` flatten will be added here, before `securityPolicyMap["managed_rules"] = ...` at line 3296.
- **Create/Update (expand)**: The `ManagedRules` expand logic is at lines 4587-4719. It builds `managedRules := teov20220901.ManagedRules{}` and assigns `securityPolicy.ManagedRules = &managedRules` at line 4718. The new `FrequentScanningProtection` expand will be added here, before the assignment at line 4718.
- **Shared helpers**: `securityActionSchema()` at line 6928, `buildSecurityActionFromMap()` at line 7073, `flattenSecurityAction()` at line 7163.

## Goals / Non-Goals

**Goals:**
- Add `frequent_scanning_protection` as an Optional TypeList (MaxItems: 1, Computed) parameter under the `managed_rules` schema block (inside `security_policy`)
- Expose the 7 sub-fields of `FrequentScanningProtection`: `enabled`, `id`, `action`, `count_by`, `block_threshold`, `counting_period`, `action_duration`
- Implement Read (flatten) logic to populate `frequent_scanning_protection` from `DescribeSecurityPolicy` response, with nil checks at every nesting level (`ManagedRules`, `FrequentScanningProtection`, and each nested scalar / `Action` field); reuse `flattenSecurityAction()` for the nested `Action`
- Implement Create/Update (expand) logic to expand `frequent_scanning_protection` into the `ModifySecurityPolicy` request's `SecurityPolicy.ManagedRules.FrequentScanningProtection`; reuse `buildSecurityActionFromMap()` for the nested `Action`
- Add unit tests using the gomonkey mock approach
- Update the resource `.md` documentation (via `make doc` in the finalize stage)

**Non-Goals:**
- Modifying any existing schema fields or behavior of `managed_rules` or other blocks
- Adding `BlockIPActionParameters` / `ReturnCustomPageActionParameters` handling beyond what `securityActionSchema()` already provides (these are already covered by the shared helper)
- Changing the resource's ID format, import logic, or retry semantics
- Touching the `security_config` (write-only) top-level field

## Decisions

### 1. Schema structure for `frequent_scanning_protection`

**Decision**: Use `TypeList` (Optional, Computed, MaxItems: 1) for `frequent_scanning_protection`, with each element being a `schema.Resource` exposing the 7 sub-fields. Place it as a sibling of `managed_rule_groups` within the `managed_rules` schema block.

**Rationale**: The SDK `FrequentScanningProtection` is a single `*FrequentScanningProtection` (not a list), so MaxItems: 1 mirrors the cardinality. Using TypeList (not TypeMap / single object) matches the established pattern for nested config blocks throughout this resource (e.g. `auto_update`, `rate_limit_template`). Marking it Computed lets the API populate the `id` (output-only) and any API defaults on read.

**Alternatives considered**: A flat set of top-level `frequent_scanning_*` fields under `managed_rules` was rejected — grouping into a sub-block keeps the schema organized and matches the SDK struct boundary, which aids readability and future maintenance.

### 2. Reuse `securityActionSchema()` for the `action` sub-field

**Decision**: Define `action` as `TypeList` (MaxItems: 1, Optional, Computed) with `Elem: securityActionSchema()`, exactly like other `*SecurityAction` fields in the resource (e.g. `invalid_attestation_action`, `adaptive_frequency_control.action`).

**Rationale**: The SDK `FrequentScanningProtection.Action` is `*SecurityAction`, identical in shape to other action fields. The SDK documentation states its `Name` supports `Deny`, `Monitor`, `JSChallenge` — all covered by `securityActionSchema()`. Reusing the shared schema and the existing `flattenSecurityAction()` / `buildSecurityActionFromMap()` helpers guarantees consistency and minimizes new code.

**Alternatives considered**: Inlining a brand-new action schema was rejected because it would diverge from the established pattern and duplicate maintenance.

### 3. `block_threshold` as `TypeInt`

**Decision**: Define `block_threshold` as `TypeInt` (Optional, Computed), mapping to the SDK's `*int64` field `BlockThreshold`. In expand, use `helper.IntInt64()` (the helper used elsewhere for `*int64` fields); in flatten, read the pointer and set the map value.

**Rationale**: The SDK field is `*int64`. The resource already maps integer fields to `TypeInt` (e.g. `priority`). Using `TypeInt` is the natural fit.

**Alternatives considered**: `TypeString` was rejected — the value is an integer and `TypeInt` gives proper plan-time type checking.

### 4. `enabled`, `count_by`, `counting_period`, `action_duration` as `TypeString`

**Decision**: Define these four fields as `TypeString` (Optional, Computed), mapping to the SDK's `*string` fields. Use `helper.String()` in expand; nil-check before setting in flatten.

**Rationale**: The SDK fields are `*string`. The values include time-with-unit strings (e.g. `5s`, `60s`) and enum strings (`on`/`off`, `http.request.ip`), which are best represented as `TypeString`.

### 5. `id` as `TypeString` (Computed, output-only)

**Decision**: Define `id` as `TypeString`, `Computed` only (not Optional), mapping to the SDK's `*string` field `Id`. Only flatten it on read; never expand it.

**Rationale**: The SDK documentation states `Id` is "仅作为出参返回" (returned only as output). Marking it Computed-only prevents users from setting it and avoids sending it to the cloud API.

### 6. Read flatten: nil-check every nesting level

**Decision**: In the `ManagedRules` flatten block, after the existing `ManagedRuleGroups` handling and before `securityPolicyMap["managed_rules"] = ...`, add: check `respData.ManagedRules.FrequentScanningProtection != nil` before building the map; for each scalar field, check the pointer for nil before setting; for `Action`, reuse `flattenSecurityAction()` and wrap in `[]interface{}` when non-empty.

**Rationale**: Per project rules, the Read operation must check for nil before accessing sub-fields because the cloud API may return nil for optional fields. This matches the existing per-field nil-check style already used in the `ManagedRules` flatten block.

### 7. Expand: build `*FrequentScanningProtection` and assign to `managedRules.FrequentScanningProtection`

**Decision**: In the `ManagedRules` expand block, after the existing `managed_rule_groups` handling and before `securityPolicy.ManagedRules = &managedRules`, read `managedRulesMap["frequent_scanning_protection"]`; if present, construct a `teov20220901.FrequentScanningProtection{}` populated via `helper.String()` / `helper.IntInt64()` / `buildSecurityActionFromMap()`, then assign to `managedRules.FrequentScanningProtection`. Do NOT expand `id` (output-only).

**Rationale**: This follows the established expand style (construct SDK struct from a map, assign to the parent field). The expand happens inside the existing `managed_rules` block so the same `managedRules` value receives `FrequentScanningProtection` alongside already-set fields.

### 8. Unit testing approach

**Decision**: Use the gomonkey mock approach for the new unit tests, mocking the SDK client methods, consistent with the project rule that new resource parameter additions use mock-based unit tests rather than the terraform test suite.

**Rationale**: The change adds a brand-new self-contained parameter block. Mock-based tests keep the suite hermetic and fast and align with the project's testing guidance for new parameters.

## Risks / Trade-offs

- [Risk: nil `FrequentScanningProtection` in API response] → Mitigation: guard the entire flatten with `if respData.ManagedRules.FrequentScanningProtection != nil`, so a nil value simply omits the attribute (state stays empty), consistent with existing optional nested fields.
- [Risk: partial `FrequentScanningProtection` (only some sub-fields set)] → Mitigation: per-field nil checks in flatten ensure only populated fields are written to state, avoiding nil-pointer dereferences and keeping state accurate.
- [Risk: large file further growth] → The resource file is already ~7258 lines. Mitigation: follow existing patterns closely (reuse `securityActionSchema()` and the action helpers) so the net addition is a focused schema block + a flatten block + an expand block, minimizing review surface.
- [Risk: backward compatibility] → Mitigation: `frequent_scanning_protection` is purely additive and Optional; the cloud API preserves existing configuration when `FrequentScanningProtection` is omitted, so pre-existing Terraform configurations and state are unaffected.
