# Change: Preserve CLS ConfigExtra Extract Rule Key Order

## Why

`tencentcloud_cls_config_extra.extract_rule.keys` maps extracted fields to
delimiter columns or regular-expression capture groups by position. The
released SDKv2 resource models this attribute as `schema.TypeSet`, so the
provider discards configuration order before Create or Update sends the keys
to CLS. Reordering keys also cannot be represented reliably as a Terraform
change.

## What Changes

- Change `extract_rule.keys` from `schema.TypeSet` to `schema.TypeList`.
- Preserve list order in Create, Read, and Update.
- Increment the resource schema version from 0 to 1 and register a version 0
  state upgrader.
- Cover ordered Create/Update/Import behavior and the real SDKv2 state upgrade
  protocol for JSON and legacy flatmap state.
- Document `keys` as an ordered list.

## Capabilities

### New Capabilities

- `cls-config-extra-ordered-keys`: Defines ordered key handling and backward
  compatible state migration for `tencentcloud_cls_config_extra`.

### Modified Capabilities

None.

## Impact

- Resource: `tencentcloud_cls_config_extra` only.
- Schema: `extract_rule.keys` changes from `set(string)` to `list(string)`.
- State: existing schema version 0 instances are upgraded to version 1 without
  resource replacement.
- API: Create and Update send keys in configuration order; Read stores the CLS
  API order.
- Dependencies: no `go.mod`, `go.sum`, or vendor changes.
- Out of scope: other CLS resources with similar collection schemas.
