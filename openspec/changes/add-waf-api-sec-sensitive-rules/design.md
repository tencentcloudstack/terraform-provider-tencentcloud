## Context

WAF exposes API-security "sensitive data" configuration through a single multiplexed write API, `ModifyApiSecSensitiveRule` (`waf/v20180125`). Its request body carries a top-level `Domain`, `Status`, `RuleName`, plus seven mutually-independent rule structures (one per rule type). The read side is `DescribeApiSecSensitiveRuleList`, which returns the rule lists keyed by a set of `IsQuery*` boolean switches.

We are splitting the seven rule structures into seven independent CRUD Terraform resources. All seven share the same write/read APIs and the same lifecycle pattern; they differ only in their schema (the embedded struct) and which Describe list / `IsQuery*` flag they read from.

The implementation must follow the `tencentcloud_igtm_monitor` code style (var-block CRUD, `resource.Retry` wrappers, nil-safe response handling, `service_tencentcloud_*.go` query helpers).

## Goals / Non-Goals

**Goals:**
- 7 independent CRUD resources, each schema = exactly one sub-struct's fields, strictly validated (no extra fields).
- Composite resource ID `Domain#RuleName`.
- Unify the sub-struct `RuleName` with the top-level `RuleName` (expose a single `rule_name`).
- Expose `status` only as `0`/`1`; use `3` internally for delete.
- All API calls wrapped with `resource.Retry`; all response reads are nil-safe.
- Read uses `DescribeApiSecSensitiveRuleList` with the correct `IsQuery*` flag, then locates the rule by `RuleName`.

**Non-Goals:**
- No batch operations (the `*NameList` / `*RuleName []` batch fields are not exposed).
- No new SDK; no SDK source modification.
- No changes to any existing resource/schema/state.

## Decisions

### D1. Seven resources over one multiplexed resource
Each resource sets only its own sub-struct field on the `ModifyApiSecSensitiveRule` request and leaves the other six nil. Rationale: clean declarative model, independent lifecycle, matches the user's explicit requirement. Alternative (single resource with 7 optional blocks) rejected — confusing semantics and cross-field coupling.

### D2. Composite ID `Domain#RuleName`
Per business rule, the unique identity is `Domain` + `RuleName`. ID is built with `strings.Join([]string{domain, ruleName}, tccommon.FILED_SP)`. Read/Update/Delete split via `tccommon.FILED_SP`; on malformed ID return an error. Import uses the same `Domain#RuleName` form.

### D3. `RuleName` unification
The top-level request `RuleName` is authoritative. Each sub-struct that also has a `RuleName` field is populated from the same `rule_name` schema attribute; we do not expose the sub-struct `RuleName` separately. (`ApiSecCustomSensitiveRule` has no `RuleName` — it relies solely on the top-level one.)

### D4. `status` mapping (0/1 exposed, 3 internal for delete)
- Schema `status` is `TypeInt`, validated to `{0, 1}` (`ValidateAllowedIntValue`).
- Create/Update: `request.Status = uint64(status)`.
- Delete: `request.Status = 3` (hard-coded default), reusing `ModifyApiSecSensitiveRule` as the delete operation.
- Sub-structs that also have an inner `Status` field are kept consistent with the top-level `status` (set from the same attribute) so the rule's own switch matches.

### D5. Read via `DescribeApiSecSensitiveRuleList`
Each resource sets exactly one `IsQuery*` flag to `true` and passes `Domain` (+ `RuleName` where the API supports filtering). Then iterate the corresponding response list (`Data` / `ApiExtractRule` / `ApiSecPrivilegeRule` / `ApiSecSceneRule` / `ApiSecCustomEventRule` / `ApiExcludeRule` / `ApiSecSensitiveWhiteRule`), match by `RuleName`, and flatten. If not found, `d.SetId("")` and return nil (treat as deleted). Read flag mapping:
  - custom_rule → `Data` (`ApiSecSensitiveRule`, read its `CustomRule`); no dedicated `IsQuery*` flag, returned by default.
  - custom_api_extract_rule → `IsQueryApiExtractRule` → `ApiExtractRule`.
  - privilege_rule → `IsQueryApiPrivilegeRule` → `ApiSecPrivilegeRule`.
  - scene_rule → `IsQueryApiSceneRule` → `ApiSecSceneRule`.
  - custom_event_rule → `IsQueryApiCustomEventRule` → `ApiSecCustomEventRule`.
  - custom_api_exclude_rule → `IsQueryApiExcludeRule` → `ApiExcludeRule`.
  - white_rule → `IsQueryApiSensitiveWhiteRule` → `ApiSecSensitiveWhiteRule`.

### D6. Field handling & nested structs
- Read-only/output-only fields (`UpdateTime`, `Timestamp`, `Source`, `Count`, `Label`) are modeled as `Computed` (not user-settable) where they appear in a struct, since they are server-generated. Only meaningful input fields are `Optional/Required`.
- Nested objects are modeled as `TypeList` blocks with `Elem: &schema.Resource{...}`:
  - `ApiSecPrivilegeRule.ApiNameOp` / `ApiSecSensitiveWhiteRule.ApiNameOp` / `ApiSecCustomEventRule.ApiNameOp` → `api_name_op` block, which itself nests `api_name_method`.
  - `ApiSecSceneRule.RuleList`, `ApiSecCustomEventRule.MatchRuleList` / `StatRuleList` → `ApiSecSceneRuleEntry` blocks.
  - `ApiSecSensitiveWhiteRule.WhiteFields` → `ApiSecSensitiveWhiteField` blocks.
- `int`/`bool` scalar inputs are read with `d.GetOkExists`; strings/lists with `d.GetOk`.

### D7. Service-layer helpers
Add one query helper per resource (or a shared `DescribeApiSecSensitiveRuleListByFilter`) in `service_tencentcloud_waf.go`, following the `IgtmService.DescribeIgtmMonitorById` pattern (`resource.Retry` + `ratelimit.Check` + nil-safe). Resources match the target rule by `RuleName` from the returned list.

## Risks / Trade-offs

- [Multiplexed write API used for delete via `Status=3`] → Document clearly; Delete always sends `Status=3` with the sub-struct identified by `RuleName`.
- [`custom_rule` has no `RuleName` in its sub-struct and is returned in the generic `Data` list] → Match the parent `ApiSecSensitiveRule.RuleName` and read its nested `CustomRule`; if the API does not filter by `RuleName`, iterate `Data` and match locally.
- [Server-generated fields drift] → Mark them `Computed` to avoid perpetual diffs.
- [Eventual consistency after write] → All reads/writes wrapped in `resource.Retry`; Read tolerates not-found right after delete.
- [`DescribeApiSecSensitiveRuleList` has no paging params] → Not applicable (returns full lists); no limit/offset to set.

## Migration Plan

Additive only. New resources + docs + provider registration. No state migration. Rollback = revert the additions; existing configs unaffected.

## Open Questions

- Whether `DescribeApiSecSensitiveRuleList` reliably honors the `RuleName` filter for every rule type, or whether local matching by `RuleName` is required for all — implementation will match locally to be safe.
