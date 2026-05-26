## Context

The `tencentcloud_clb_target_group` resource manages CLB target groups. It already supports parameters like `health_check`, `schedule_algorithm`, `tags`, `weight`, `full_listen_switch`, `keepalive_enable`, `session_expire_time`, and `ip_version`. The `SnatEnable` field exists in both `CreateTargetGroupRequest` and `ModifyTargetGroupAttributeRequest` in the SDK but is not yet exposed in the Terraform resource schema.

Key constraint: The `DescribeTargetGroups` API response (`TargetGroupInfo` struct) does NOT include a `SnatEnable` field, meaning the value cannot be read back from the API after creation.

Current service layer method `CreateTargetGroup` accepts 14 parameters. The `ModifyTargetGroupAttribute` method accepts the full request struct directly.

## Goals / Non-Goals

**Goals:**
- Add `snat_enable` as an Optional bool parameter to the resource schema
- Pass `SnatEnable` to `CreateTargetGroup` API during resource creation
- Pass `SnatEnable` to `ModifyTargetGroupAttribute` API during resource update
- Maintain backward compatibility with existing configurations

**Non-Goals:**
- Reading `snat_enable` back from the API (not available in DescribeTargetGroups response)
- Adding any other missing parameters in this change
- Modifying the Delete flow (DeleteTargetGroups does not use this parameter)

## Decisions

### Decision 1: Add `snat_enable` as a new parameter to `CreateTargetGroup` service method

**Choice**: Extend the `CreateTargetGroup` method signature to accept a `snatEnable *bool` parameter.

**Rationale**: This follows the existing pattern used for `fullListenSwitch`, `keepaliveEnable`, and other bool parameters. The method already accepts 14 parameters; adding one more maintains consistency.

**Alternative considered**: Refactoring to use a struct/options pattern. Rejected because it would be a larger refactor outside the scope of this single-parameter addition.

### Decision 2: Handle as write-only attribute (no read-back)

**Choice**: Since `TargetGroupInfo` does not include `SnatEnable`, the Read method will NOT attempt to set this field. Terraform state will retain whatever value the user configured.

**Rationale**: This is the standard Terraform pattern for write-only attributes. The value persists in state from the user's configuration. If the user imports the resource, this field will not be populated (which is acceptable for an Optional field).

**Alternative considered**: Using a separate API call to read the value. Rejected because no such API exists for this specific field.

### Decision 3: Support update via `ModifyTargetGroupAttribute`

**Choice**: Add `snat_enable` change detection in the Update method and pass it to the existing `ModifyTargetGroupAttribute` request.

**Rationale**: The `ModifyTargetGroupAttributeRequest` struct already includes `SnatEnable *bool`. The Update method already constructs this request for other fields. Adding one more `HasChange` check is straightforward.

## Risks / Trade-offs

- **[Risk] State drift undetectable**: Since the API does not return `SnatEnable` in Describe, if someone changes it outside Terraform, the drift won't be detected. → Mitigation: This is acceptable for a write-only attribute; document this limitation.
- **[Risk] Import incomplete**: When importing an existing target group, `snat_enable` will not be populated in state. → Mitigation: The field is Optional, so this is safe. Users can add it to their config after import.
