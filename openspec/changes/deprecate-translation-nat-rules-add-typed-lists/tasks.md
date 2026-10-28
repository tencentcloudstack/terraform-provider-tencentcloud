# Tasks for `deprecate-translation-nat-rules-add-typed-lists`

## 1. Schema

- [x] 1.1 Add `Deprecated:` annotation to existing `translation_nat_rules`. Change `Required: true` to `Optional: true, Computed: true` (required to coexist with the new fields under ConflictsWith).
- [x] 1.2 Add `local_network_layer_rules` (`TypeList`, Optional, Computed). Inner `Elem`: `translation_ip` (Required), `original_ip` (Required), `description` (Optional). All inner string types; **no ForceNew**.
- [x] 1.3 Add `local_transport_layer_rules` (`TypeList`, Optional, Computed). Inner `Elem`: `translation_ip` (Required), `description` (Optional). Do NOT expose `original_ip`.
- [x] 1.4 Add `peer_network_layer_rules` (`TypeList`, Optional, Computed). Inner `Elem`: `translation_ip` (Required), `original_ip` (Required), `description` (Optional). No ForceNew.
- [x] 1.5 Add `ConflictsWith` between legacy ↔ each new field; do NOT cross-link the three new fields among themselves.

## 2. Create

- [x] 2.1 Aggregate inputs from legacy + the three new fields into a single `[]*vpcv20170312.TranslationNatRuleInput` in canonical bucket order: legacy first (compat), then `local_network` → `local_transport` → `peer_network`.
- [x] 2.2 Reuse the existing batched `MAX_CREATE_RULES_LEN = 20` create flow with retry. No retry-related code duplication.

## 3. Read — order-preserving (the core fix)

- [x] 3.1 Build `apiByKey` map indexed by `ruleKey(item, includeOriginalIp)` for each bucket separately.
- [x] 3.2 Walk `d.Get(field).([]interface{})` (user HCL order) and emit, for each user item whose key is found in `apiByKey`, a flatten of the API rule. Mark each emitted key as `seen`.
- [x] 3.3 Append remaining `apiByKey` entries (those whose key is NOT in `seen`) in API order at the tail of the output list.
- [x] 3.4 Continue to populate `translation_nat_rules` (legacy union) for backward compatibility (full union, no order remapping; the legacy `Set` doesn't have a position-drift problem).
- [x] 3.5 `_ = d.Set("local_network_layer_rules", reorderedList)` etc. for the three new fields.

## 4. Update

- [x] 4.1 Keep the existing `if d.HasChange("translation_nat_rules") { ... }` branch unchanged (`Set.Difference` add/remove).
- [x] 4.2 Add `updateTypedRuleField(field, direction, type, includeOriginalIp)` helper that:
  - Skips when `d.HasChange(field)` is false.
  - `oldI, newI := d.GetChange(field)`; cast each to `[]interface{}`; walk to `[]map[string]interface{}`.
  - Indexes both sides by `ruleKey`. Emits `toCreate`, `toDelete`, `toModify` via the primary-key diff.
  - Routes through `applyRulesDiff(deleteRules, modifyDiffs, createRules)`.
- [x] 4.3 Update body invokes `updateTypedRuleField` three times, once per new field:
  - `local_network_layer_rules` with `direction="LOCAL"`, `type="NETWORK_LAYER"`, `includeOriginalIp=true`
  - `local_transport_layer_rules` with `direction="LOCAL"`, `type="TRANSPORT_LAYER"`, `includeOriginalIp=false`
  - `peer_network_layer_rules` with `direction="PEER"`, `type="NETWORK_LAYER"`, `includeOriginalIp=true`

## 5. Helpers

- [x] 5.1 Bring back the typed-list helpers (some were deleted in the prior TypeSet attempt):
  - `listItemsToMaps(raw interface{}) []map[string]interface{}` for `[]interface{}` inputs
  - `buildTypedRuleInput(m, direction, typ, includeOriginalIp)` for Create
  - `buildTypedRule(m, direction, typ, includeOriginalIp)` for Delete
  - `buildModifyDiff(oldItem, newItem, direction, typ, includeOriginalIp)` for Modify
  - `flattenNetworkLayerRule(item)` and `flattenTransportLayerRule(item)` for Read
  - `ruleKey(m, includeOriginalIp)` for both schema and Update
  - `equalRuleMap(a, b)` for Update content comparison
  - `diffByKey(oldList, newList, includeOriginalIp)` returning `(toCreate, toDelete, toModify)`
  - `applyRulesDiff(ctx, logId, meta, natGatewayId, deleteRules, modifyDiffs, createRules)` issuing calls in DELETE → MODIFY → CREATE order
- [x] 5.2 Keep the legacy helper `buildLegacyRuleInput(m)` for Create's compat path.

## 6. Documentation

- [x] 6.1 Author `.md` with two HCL example blocks: typed-list (recommended) first, then legacy with a `# DEPRECATED` note. Include a paragraph explaining `(direction, type, key_ip)` rule identity and the in-place Modify behavior.
- [x] 6.2 `make doc` regenerates the website file.

## 7. Validation

- [x] 7.1 `go build ./tencentcloud/...` clean.
- [x] 7.2 `go vet ./tencentcloud/services/vpc/...` clean.
- [x] 7.3 `read_lints` shows no new errors/warnings.

## 8. Update — switch typed-list diff from primary-key to positional (index-aligned)

The three new typed-list fields MUST be treated as ordered lists. Diff is computed by index, not by primary key:

- [x] 8.1 Replace `diffByKey` usage in `updateTypedRuleField` with a new `diffByIndex(oldList, newList, includeOriginalIp)` returning `(toCreate, toDelete, toModify)`:
  - For `i` in `[0, min(len(old), len(new)))`: if `equalRuleMap(old[i], new[i], includeOriginalIp)` is false → emit `modifyPair{oldItem: old[i], newItem: new[i]}` (MODIFY).
  - If `len(new) > len(old)`: items at indices `[len(old), len(new))` of `new` are emitted as CREATE in their HCL order.
  - If `len(old) > len(new)`: items at indices `[len(new), len(old))` of `old` are emitted as DELETE in their original order.
- [x] 8.2 Keep call ordering in `applyRulesDiff` unchanged (DELETE → MODIFY → CREATE). MODIFY remains one request per pair (SDK constraint).
- [x] 8.3 Keep `ruleKey` / `equalRuleMap` helpers; `diffByKey` becomes unused — remove it (and its `modifyPair` references that are also produced by `diffByIndex`, which reuses the same `modifyPair` struct).
- [x] 8.4 Keep Read's `alignTypedListToUserOrder` (still needed: the API may return rules in a different order than user HCL after a CREATE batch, and this aligns Read state to user HCL order to avoid spurious diffs at the position-sensitive list level).
- [x] 8.5 Re-run `go build`, `go vet`, `read_lints`.

## 9. Read — drop user-HCL reorder; preserve API authoritative order

The previous `alignTypedListToUserOrder` masked the API's real order with the user's HCL order to avoid spurious diff. That hid drift and would silently break once the API gains user-defined ordering semantics. Switch to honoring the API order and document the append-only contract.

- [x] 9.1 Remove `alignTypedListToUserOrder` from the resource file.
- [x] 9.2 In Read, populate the three typed lists by walking `respData` in API order and appending to each per-bucket list (no map indirection, no reorder). Legacy `translation_nat_rules` Set continues unchanged.
- [x] 9.3 Update `.md` resource doc with an **"Append-only convention"** section: new rules MUST be appended to the end; do NOT insert in the middle, do NOT reorder existing rules; explain that violating the convention triggers connected `~` plan output and per-API uniqueness errors at apply time.
- [x] 9.4 Update OpenSpec `design.md`:
  - Rewrite D5 from "Read reorders typed lists to match the user's HCL order" to "Read preserves the API's authoritative order + append-only user contract".
  - Add a Risks bullet documenting the failure modes that surface when a user violates append-only (mid-insert, reorder, mid-delete) and that the API uniqueness error is the intended user-visible signal.
- [x] 9.5 Update OpenSpec `specs/.../spec.md` Read scenarios accordingly (drop "user-HCL order" wording; add "API order preserved").
- [x] 9.6 `make doc` regenerates the website file.
- [x] 9.7 `go build`, `read_lints` clean; `openspec validate ... --strict` passes.
