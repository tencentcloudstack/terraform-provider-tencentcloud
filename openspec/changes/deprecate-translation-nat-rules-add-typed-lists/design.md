# Design: Three typed `TypeList` fields with order-preserving Read

## Context

Rules under a private NAT gateway are partitioned by `(direction, type)` into three buckets, each with its own primary key (three-layer keys on `original_ip`, four-layer keys on `translation_ip`). Two earlier modeling attempts failed:

1. Original `translation_nat_rules` `TypeSet` with default hash — every field participated in the hash, so any edit became DELETE+CREATE. No Modify path.
2. `TypeSet` with custom hash on the primary key — fragile due to `*string` vs `string` type assertions in Read, and plan output renders elements in unstable hash order.

The decisive product-team feedback: keep set semantics (no business order), but render plan output in **the order the user wrote**. The cleanest way is three explicit `TypeList` fields, with a Read-time reorder to fix the only place where Terraform's TypeList position-sensitive diff causes spurious churn — namely, when the API returns rules in an order that differs from the user's HCL order.

## Decisions

### D1 — Three `TypeList` fields, plus a deprecated legacy field

Schema declares (in source order):

* `nat_gateway_id`
* `translation_nat_rules` (legacy, deprecated, kept for backward compat)
* `local_network_layer_rules`
* `local_transport_layer_rules`
* `peer_network_layer_rules`

The three new fields are `TypeList` with `Optional: true, Computed: true`. Inner shape varies by bucket per D2.

### D2 — Inner schema per bucket

* `local_network_layer_rules`: `translation_ip` (Required), `original_ip` (Required), `description` (Optional)
* `local_transport_layer_rules`: `translation_ip` (Required), `description` (Optional). **No** `original_ip` — the API rejects it on four-layer rules.
* `peer_network_layer_rules`: `translation_ip` (Required), `original_ip` (Required), `description` (Optional)

None of the inner IP fields carry `ForceNew`. Editing them is routed through Update to MODIFY (non-key IP changes) or DELETE+CREATE (key IP changes), per the primary-key diff in D7.

### D3 — Pairwise `ConflictsWith` between legacy and each new field

Each of the four nested-block fields lists the other (only across the legacy↔new boundary; the three new fields are not mutually exclusive among themselves). Mixing legacy and any typed list in one config fails at plan time.

### D4 — Custom Set hash is not used

Because the new fields are `TypeList`, no custom Set function is needed. Position-sensitivity is handled by D5 instead. The legacy `translation_nat_rules` keeps its default Set hash unchanged (we don't fix the legacy field's Modify behavior; users wanting Modify must migrate).

### D5 — Read preserves the API's authoritative order; user follows an append-only contract

The most important decision. Earlier iterations attempted to reorder Read state to match the user's HCL order (so plans never showed spurious diffs even when the API's response order differed). That cosmetic alignment masked the API's real order — which is currently creation-time, but may gain user-defined ordering semantics in the future (e.g. matching priority). Once that happens, a provider-side reorder would silently break: users editing HCL order would see no effect, because Read would overwrite their intent on every refresh.

The chosen design instead keeps the API order authoritative and makes the user's contract explicit:

```
For each typed list field F:
    out = empty list
    for each item in respData (API authoritative order):
        if item belongs to F's (direction, type) bucket:
            out.append(flattenRule(item))
    d.Set(F, out)
```

**User-facing contract (documented in the resource `.md`):**

* New rules MUST be appended to the end of the corresponding typed list.
* Existing rules MUST NOT be inserted into the middle, MUST NOT be reordered, and SHOULD NOT be removed from the middle (remove from the tail when possible).
* Violating the contract triggers a connected `~` plan and may surface a uniqueness error at apply time. The SDK's error message identifies the offending IP, instructing the user to fix HCL.

**Why this works**:

* Initial Apply: Create issues rules in HCL order; the API records them in that order; Read returns them in that order. State == HCL → 0 diff.
* Append: Adding a rule at the tail of HCL produces a single `+` slot; Update emits one CREATE; Read writes the new tail back to state. Stable.
* Edit in place (no order change): Update emits one MODIFY for the touched slot; Read returns the updated content at the same position. Stable.
* External rule additions appear at the tail of state on next Refresh; the user is expected to mirror them at the tail of their HCL.
* Future API ordering semantics: when the backend adds a "move rule" capability, the provider already mirrors the API order, so user reordering in HCL composes correctly with that capability without provider changes.

### D6 — Create aggregates legacy + typed lists in canonical bucket order

Create reads any rules from `translation_nat_rules` first (legacy compat), then concatenates `local_network_layer_rules` → `local_transport_layer_rules` → `peer_network_layer_rules` in that canonical order. The aggregated `[]TranslationNatRuleInput` is sent through the existing `MAX_CREATE_RULES_LEN = 20` batched create with retry. The exclusivity between legacy and new is enforced by D3 at plan time, so at most one path produces inputs in practice.

### D7 — Update routes by field, **positional (index-aligned) diff** per typed list

* Legacy `translation_nat_rules` branch: unchanged from the original implementation (`Set.Difference` add/remove). We do not retrofit Modify on the legacy field; users wanting Modify must migrate.
* Each typed list gets its own `d.HasChange(field)` branch. Within a branch:
  1. `oldI, newI := d.GetChange(field)`; cast each to `[]interface{}`.
  2. Diff by **index position** (not by primary key) — the typed lists are ordered lists, so a rule's identity within the list is its slot index, not its IP value:
     - For `i` in `[0, min(len(old), len(new)))`: if `equalRuleMap(old[i], new[i], includeOriginalIp)` is false → emit a MODIFY pair `(old[i] → new[i])`. This MUST cover the case where `original_ip` (three-layer) or `translation_ip` (four-layer) at slot `i` is changed in place; `buildModifyDiff` populates both `OldOriginalIp/OriginalIp` and `OldTranslationIp/TranslationIp` so the SDK can perform the rename.
     - Items at indices `[len(old), len(new))` of `new` (when new is longer) → CREATE in HCL order.
     - Items at indices `[len(new), len(old))` of `old` (when old is longer) → DELETE in original order.
  3. Issue API calls in this fixed order: **DELETE → MODIFY → CREATE**.

**Rationale.** Length parity is therefore handled exclusively via DELETE/CREATE; equal-length in-place edits are handled exclusively via MODIFY. This matches the user mental model that the typed lists are ordered: "I edited slot N" → MODIFY slot N; "I appended/removed" → CREATE/DELETE the tail; never silent surgery elsewhere.

### D8 — MODIFY is one rule per request

Per the SDK comment "only supports modifying a single translation rule", `ModifyPrivateNatGatewayTranslationNatRule` is invoked once per modify pair. Three-layer requests populate both `OldOriginalIp/OriginalIp` and `OldTranslationIp/TranslationIp`; four-layer requests only set `OldTranslationIp/TranslationIp`. Each call is `resource.Retry(WriteRetryTimeout, ...)`-wrapped.

### D9 — `.md` example shows the typed-list path first

Two HCL examples in source order:

1. **Recommended**: typed lists across all three buckets.
2. **Deprecated**: legacy `translation_nat_rules` block, prefixed with `# DEPRECATED ...`.

The `.md` also documents the `(direction, type, key_ip)` rule identity convention so users understand which edits route to MODIFY vs DELETE+CREATE.

## Risks

* **Append-only contract is enforced at runtime, not at plan time.** A user who inserts a rule in the middle of HCL or reorders existing rules will see Plan render a chain of `~` blocks and may receive a uniqueness error at Apply time when the positional MODIFY tries to rename slot N's IP to a value still held by slot N+1. This is intentional: the SDK error is the authoritative signal, with a precise IP reference, instructing the user to restructure HCL into an append-only edit. The provider does NOT silently DELETE+CREATE behind the user's back, because that would alter audit semantics and break future API ordering capabilities.
* **TypeList position-sensitive diff still applies for inserts/deletes inside the user's HCL list.** Inserting a rule in the middle of HCL appears in plan output as several `~` blocks plus one `+` block. This is part of the append-only enforcement story above — it surfaces the violation visibly so the user can correct course before Apply.
* **Pure swap (e.g. `[A, B] → [B, A]`) MUST surface an API error rather than be silently rewritten by the provider.** The positional diff in D7 turns a swap into two MODIFY calls (`A→B`, `B→A`); the first call collides with the still-existing peer rule and the SDK returns a uniqueness violation. This is intentional: the provider does not invent DELETE+CREATE behind the user's back, because that would alter audit semantics and risk losing rule-attached state. The error message itself instructs the user to break the swap into two applies (e.g. rename one side to a temporary value first, then complete the swap).
* **Legacy `translation_nat_rules` keeps its DELETE+CREATE behavior** — it's deprecated, and users with Modify needs are expected to migrate to the typed lists.
