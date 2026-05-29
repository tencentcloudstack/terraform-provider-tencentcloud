## Context

The `tencentcloud_clb_target_group` resource already supports multiple parameters (target_group_name, vpc_id, port, type, protocol, health_check, schedule_algorithm, tags, weight, full_listen_switch, keepalive_enable, session_expire_time, ip_version). However, the `SnatEnable` parameter is not yet exposed in the Terraform schema despite being available in the underlying cloud APIs:

- `CreateTargetGroupRequest.SnatEnable` (Create)
- `ModifyTargetGroupAttributeRequest.SnatEnable` (Update)
- `TargetGroupInfo.SnatEnable` (Read via DescribeTargetGroups)

The SDK package is `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317` and already includes `SnatEnable *bool` in all relevant structs. No SDK upgrade is needed.

Current resource file: `tencentcloud/services/clb/resource_tc_clb_target_group.go`
Current service file: `tencentcloud/services/clb/service_tencentcloud_clb.go`

## Goals / Non-Goals

**Goals:**
- Add `snat_enable` (bool, Optional) parameter to the `tencentcloud_clb_target_group` resource schema.
- Pass `SnatEnable` during resource creation via `CreateTargetGroup` API.
- Read `SnatEnable` from `DescribeTargetGroups` API response in the Read function.
- Support updating `SnatEnable` via `ModifyTargetGroupAttribute` API in the Update function.
- Update resource documentation (.md file) with the new parameter.
- Add unit tests using gomonkey mock for the new parameter.

**Non-Goals:**
- Adding any other missing parameters (only `snat_enable` is in scope).
- Modifying the Delete logic (DeleteTargetGroups does not use SnatEnable).
- Changing existing parameter behavior or schema structure.

## Decisions

### Decision 1: Schema field type and attributes

`snat_enable` will be defined as:
```go
"snat_enable": {
    Type:        schema.TypeBool,
    Optional:    true,
    Computed:    true,
    Description: "Whether to enable SNAT (source IP replacement). true: enable, false: disable. Default: false. Note: When SnatEnable is enabled, client source IP will be replaced, and the pass-through client source IP option is disabled, and vice versa.",
}
```

**Rationale**: The field is Optional (not required for creation) and Computed (the API returns a value even if not set by user). This follows the same pattern as other bool fields in this resource (e.g., `keepalive_enable`).

### Decision 2: Service layer method signature change

Add `snatEnable *bool` parameter to the `ClbService.CreateTargetGroup` method signature.

**Rationale**: The existing pattern passes all parameters explicitly to the service layer method. Adding a new `*bool` parameter maintains consistency.

### Decision 3: Update support

`SnatEnable` is supported in `ModifyTargetGroupAttributeRequest`, so it can be updated in-place without ForceNew.

**Rationale**: The `ModifyTargetGroupAttribute` API explicitly includes `SnatEnable` as a modifiable field (confirmed in SDK struct at line 8098 of models.go).

### Decision 4: No ForceNew needed

Unlike `ip_version` and `full_listen_switch` which are in the `immutableFields` list, `snat_enable` is updatable and should NOT be added to `immutableFields`.

**Rationale**: The ModifyTargetGroupAttribute API supports changing SnatEnable after creation.

## Risks / Trade-offs

- **[Risk] Backward compatibility** → Mitigated: The field is Optional+Computed, so existing configurations without `snat_enable` will continue to work. The API default is `false`.
- **[Risk] State drift on import** → Mitigated: The Read function will read `SnatEnable` from the API response, so imported resources will have the correct state.
