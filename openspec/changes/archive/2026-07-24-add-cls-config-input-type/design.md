## Context

The `tencentcloud_cls_config` resource manages CLS (Cloud Log Service) collection configurations. The cloud API supports an `InputType` parameter in `CreateConfigRequest`, `ModifyConfigRequest`, and the `ConfigInfo` response, but the Terraform resource schema does not expose it.

The `InputType` parameter specifies the type of log input:
- `file`: File-based log collection (default)
- `windows_event`: Windows event log collection
- `syslog`: System log collection (syslog protocol)

This is a straightforward addition of a single Optional string field, following the existing patterns in the resource.

## Goals / Non-Goals

**Goals:**
- Add `input_type` (Optional, TypeString) to the `tencentcloud_cls_config` resource schema
- Wire it into Create, Read, and Update operations
- Update the resource documentation (`.md` file)

**Non-Goals:**
- No changes to the `tencentcloud_cls_config_extra` or `tencentcloud_cls_config_attachment` resources
- No validation logic beyond what the API provides
- No changes to the `DeleteConfig` operation (does not use this parameter)

## Decisions

1. **Field type: `TypeString`** — The API defines `InputType` as `*string`, so using `TypeString` with no default value is the natural choice. This is consistent with how other string parameters like `log_type` and `name` are handled.

2. **Optional, not Computed** — The field is Optional. It is not marked as Computed because the API does not return a default value when the field is not set; it simply returns `nil`. The Read function will only set the field when the API response includes it.

3. **No ForceNew** — The `input_type` can be modified via `ModifyConfig`, so the field does not need `ForceNew`. Users can update it in-place.

4. **Read function: nil-safe** — Following the existing pattern, the Read function only calls `d.Set("input_type", ...)` when `config.InputType != nil`, preventing unnecessary diffs when the API returns null.

5. **Update function: `d.HasChange` guard** — Following the existing pattern, the Update function only includes `input_type` in the `ModifyConfigRequest` when it has changed, using `d.HasChange("input_type")`.

## Risks / Trade-offs

- **Minimal risk**: This is a single Optional field addition with no breaking changes. Existing configurations continue to work unchanged.
- **API compatibility**: The `InputType` field is already present in the vendor SDK, so no SDK update is needed.