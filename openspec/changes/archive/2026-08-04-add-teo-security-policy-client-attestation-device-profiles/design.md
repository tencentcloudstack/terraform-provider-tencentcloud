## Context

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policies via the cloud API's `DescribeSecurityPolicy` (read) and `ModifySecurityPolicy` (create/update) interfaces. Both interfaces share the same `SecurityPolicy` struct in the vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).

Within `SecurityPolicy.BotManagement.ClientAttestationRules.Rules` (element type `ClientAttestationRule`), the resource already exposes `id`, `name`, `enabled`, `priority`, `condition`, `attester_id`, and `invalid_attestation_action`. The SDK's `ClientAttestationRule` struct additionally has a `DeviceProfiles []*DeviceProfile` field that is currently NOT exposed in the Terraform schema, flatten (read), or expand (create/update) logic.

The `DeviceProfile` struct in the SDK contains:
- `ClientType *string` — device type: `iOS`, `Android`, `WebView`, `WeChatMiniProgram`
- `HighRiskMinScore *uint64` — minimum score (1-99) for high-risk; default 50
- `HighRiskRequestAction *SecurityAction` — action for high-risk requests (Name: Deny/Monitor/Redirect/Challenge; default Monitor)
- `MediumRiskMinScore *uint64` — minimum score (1-99) for medium-risk; default 15
- `MediumRiskRequestAction *SecurityAction` — action for medium-risk requests (default Monitor)

The existing resource file is large (~7113 lines) and follows consistent patterns: a shared `securityActionSchema()` helper is reused for any `*SecurityAction` field (e.g. `invalid_attestation_action`, `basic_bot_settings.*.base_action`). The flatten and expand helpers `flattenSecurityAction()` / `buildSecurityActionFromMap()` are the corresponding read/write helpers.

## Goals / Non-Goals

**Goals:**
- Add `device_profiles` as an Optional TypeList parameter under each rule of the `client_attestation_rules` schema block (inside `bot_management`)
- Expose the 5 sub-fields of each `DeviceProfile`: `client_type`, `high_risk_min_score`, `high_risk_request_action`, `medium_risk_min_score`, `medium_risk_request_action`
- Implement Read (flatten) logic to populate `device_profiles` from `DescribeSecurityPolicy` response, with nil checks at every nesting level (`ClientAttestationRules`, `Rules`, each `rule.DeviceProfiles`, each `DeviceProfile` and its nested `SecurityAction` fields)
- Implement Create/Update (expand) logic to expand `device_profiles` into the `ModifySecurityPolicy` request's `SecurityPolicy.BotManagement.ClientAttestationRules.Rules[].DeviceProfiles`
- Add unit tests using the gomonkey mock approach
- Update the resource `.md` documentation

**Non-Goals:**
- Modifying any existing schema fields or behavior of `client_attestation_rules` or other blocks
- Adding `block_ip_action_parameters` / `return_custom_page_action_parameters` to the shared `securityActionSchema()` (these are out of scope; `HighRiskRequestAction` / `MediumRiskRequestAction` only support Deny/Monitor/Redirect/Challenge per the SDK documentation, so the existing `securityActionSchema()` subset is sufficient)
- Changing the resource's ID format, import logic, or retry semantics
- Touching the `security_config` (write-only) top-level field

## Decisions

### 1. Schema structure for `device_profiles`

**Decision**: Use `TypeList` (Optional) for `device_profiles`, with each element being a `schema.Resource` exposing the 5 sub-fields. Do NOT impose `MaxItems` (the cloud API allows multiple device profiles, one per `ClientType`).

**Rationale**: The SDK `DeviceProfiles` is `[]*DeviceProfile` (a list). The cloud API allows one profile per client type. Using a plain `TypeList` (no MaxItems) mirrors the SDK cardinality and matches the existing list-based rules patterns in the resource (e.g. `client_attestation_rules` itself is a `TypeList` with no MaxItems). This keeps the schema faithful to the API.

**Alternatives considered**: `MaxItems: 4` (one per supported client type) was considered but rejected because the supported client-type set is API-defined and could evolve; imposing a hard cap risks rejecting valid future configurations.

### 2. Reuse `securityActionSchema()` for request-action sub-fields

**Decision**: Define `high_risk_request_action` and `medium_risk_request_action` as `TypeList` with `MaxItems: 1` and `Elem: securityActionSchema()`, exactly like the existing `invalid_attestation_action` field.

**Rationale**: The SDK types `HighRiskRequestAction` and `MediumRiskRequestAction` are both `*SecurityAction`, identical in shape to `InvalidAttestationAction`. The SDK documentation states their `Name` only supports `Deny`, `Monitor`, `Redirect`, `Challenge` — the exact subset that `securityActionSchema()` already covers (deny/redirect/allow/challenge action parameters; Monitor has no extra parameters). Reusing the shared schema and the existing `flattenSecurityAction()` / `buildSecurityActionFromMap()` helpers guarantees consistency with the sibling `invalid_attestation_action` field and minimizes new code.

**Alternatives considered**: Inlining a brand-new action schema was rejected because it would diverge from the established pattern and duplicate maintenance.

### 3. Score fields as `TypeInt` with `Computed`

**Decision**: Define `high_risk_min_score` and `medium_risk_min_score` as `TypeInt`, `Optional` + `Computed`, mapping to the SDK's `*uint64` fields (`HighRiskMinScore`, `MediumRiskMinScore`). In expand, use `helper.IntUint64()` (the same helper already used for `priority`).

**Rationale**: The SDK fields are `*uint64`. The resource already maps `priority` (`*uint64`) to a `TypeInt` field using `helper.IntUint64()`, so this is the established pattern. Marking them `Computed` lets the provider populate the API defaults (50 and 15) on read when the user omits them, avoiding drift.

**Alternatives considered**: `TypeString` was rejected — the values are integers, and the existing `priority` field establishes the `TypeInt` + `helper.IntUint64()` precedent.

### 4. `client_type` as Required `TypeString`

**Decision**: Define `client_type` as `TypeString`, `Required`.

**Rationale**: The SDK `ClientType` is the identifying discriminator of a `DeviceProfile` (iOS / Android / WebView / WeChatMiniProgram) and has no default. Marking it Required reflects that a device profile is meaningless without a client type, and lets Terraform's plan/validate surface misconfiguration early.

**Alternatives considered**: Optional+ValidateString was considered but rejected; Required is clearer and the field is genuinely mandatory per the API.

### 5. Read flatten: nil-check every nesting level

**Decision**: In the flatten block for `client_attestation_rules`, after the existing `rule.InvalidAttestationAction` handling, add: check `rule.DeviceProfiles != nil` before iterating; for each `DeviceProfile`, check each scalar/action field for nil before setting it into the map; for action fields, reuse `flattenSecurityAction()`.

**Rationale**: Per project rules, the Read operation must check for nil before accessing sub-fields because the cloud API may return nil for optional fields. This matches the existing per-field nil-check style already used in the `client_attestation_rules` flatten loop.

### 6. Expand: build `[]*DeviceProfile` and assign to `clientAttestationRule.DeviceProfiles`

**Decision**: In the expand block for `client_attestation_rules`, after the existing `invalid_attestation_action` handling, read `ruleMap["device_profiles"]`; for each entry, construct a `teov20220901.DeviceProfile{}` populated via `helper.String()` / `helper.IntUint64()` / `buildSecurityActionFromMap()`, then assign the accumulated slice to `clientAttestationRule.DeviceProfiles`.

**Rationale**: This follows the established expand style (construct SDK struct from a map, assign to the rule field). The expand happens inside the existing `client_attestation_rules` loop so the same `clientAttestationRule` value receives `DeviceProfiles` alongside the already-set fields, then the rule is appended to `ClientAttestationRules.Rules`.

### 7. Unit testing approach

**Decision**: Use the gomonkey mock approach for the new unit tests, mocking the SDK client methods, consistent with the project rule that new resource parameter additions use mock-based unit tests rather than the terraform test suite.

**Rationale**: The resource is an existing RESOURCE_KIND_GENERAL resource, but the change adds a brand-new self-contained parameter block. Mock-based tests keep the suite hermetic and fast and align with the project's testing guidance for new parameters.

## Risks / Trade-offs

- [Risk: nil `DeviceProfiles` in API response] → Mitigation: guard the entire `device_profiles` flatten with `if rule.DeviceProfiles != nil`, so a nil list simply omits the attribute (state stays empty), consistent with existing optional list fields.
- [Risk: partial `DeviceProfile` (only some sub-fields set)] → Mitigation: per-field nil checks in flatten ensure only populated fields are written to state, avoiding nil-pointer dereferences and keeping state accurate.
- [Risk: large file further growth] → The resource file is already ~7113 lines. Mitigation: follow existing patterns closely (reuse `securityActionSchema()` and the action helpers) so the net addition is a focused schema block + a flatten loop + an expand loop, minimizing review surface.
- [Risk: drift on API defaults for scores] → Mitigation: mark `high_risk_min_score` / `medium_risk_min_score` as `Computed` so the API defaults (50 / 15) flow into state on read when the user omits them.
- [Risk: backward compatibility] → Mitigation: `device_profiles` is purely additive and Optional; the cloud API preserves existing device configuration when `DeviceProfiles` is omitted, so pre-existing Terraform configurations and state are unaffected.
