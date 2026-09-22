## Context

CLS extraction keys are positional. For delimiter logs, each key names the
corresponding delimited field. For full regular-expression logs, each key names
the corresponding capture group. The TencentCloud SDK represents
`ExtractRuleInfo.Keys` as an ordered slice, but the existing Terraform schema
uses a set and Create/Update enumerate that set.

Changing a released SDKv2 collection type also changes the resource state
shape. Existing state therefore needs an explicit schema version and upgrader.
The old set state has already lost the user's original HCL order, so migration
cannot reconstruct that intent.

## Goals / Non-Goals

**Goals:**

- Preserve declared key order through Create and Update.
- Preserve CLS API order through Read and Import.
- Make a key-only reorder an in-place Terraform diff.
- Upgrade version 0 JSON and legacy flatmap state without decode errors or
  resource replacement.
- Keep the change local to `tencentcloud_cls_config_extra` and avoid dependency
  changes.

**Non-Goals:**

- Reconstruct order that was already discarded in version 0 set state.
- Suppress the one-time corrective diff when migrated state and configuration
  order differ.
- Change `tencentcloud_cls_config`, CLS data sources, or any other collection
  schema.
- Perform or authorize changes to existing CLS resources.

## Decisions

### Model keys as `schema.TypeList`

`TypeList` matches the positional API contract and makes a reorder observable
to Terraform. Create and Update convert the list directly to the SDK string
pointer slice. Read passes the API slice to `d.Set`, which encodes it in the
same order.

### Use a local version 0 schema copy

The resource declares schema version 1. Its version 0 upgrader type is derived
from a copy of the current schema where only `extract_rule.keys` is restored to
`TypeSet`. The copy clones the `extract_rule` schema, nested resource, nested
schema map, and `keys` leaf before mutation, so the current schema is not
modified through shared pointers.

### Let the upgrader return the decoded state map

SDKv2 represents both sets and lists as arrays in JSON state. The legacy type
decodes version 0 state and the current schema re-encodes the returned map as a
list. The protocol-level regression test exercises this complete path instead
of only invoking the upgrade callback.

For legacy flatmap state, the set contains values keyed by hashes and has no
declaration order to preserve. Migration retains all values. A subsequent Read
uses CLS API order, and Terraform may plan one in-place correction if the HCL
list differs.

## Risks / Trade-offs

- A migrated instance may show a one-time in-place diff because version 0 state
  cannot recover original declaration order. Hiding that diff would preserve a
  potentially incorrect capture mapping.
- Downgrading after version 1 state is written is not a supported rollback.
  Users should restore a pre-upgrade state backup if migration itself fails.
- Cloud acceptance requires TencentCloud credentials and must run only in an
  isolated test account.

## Migration Plan

1. Upgrade version 0 state through the registered legacy schema type.
2. Refresh from CLS so state reflects the API's current key order.
3. If configuration order differs, apply the resulting in-place update once.
4. Verify the following plan is empty.

No resource replacement or API mutation is performed by state upgrade alone.

## Open Questions

None.
