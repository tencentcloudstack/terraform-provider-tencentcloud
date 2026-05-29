## Context

The `tencentcloud_clb_target_group` resource manages CLB target groups. It already supports parameters like `health_check`, `schedule_algorithm`, `tags`, `weight`, `full_listen_switch`, `keepalive_enable`, `session_expire_time`, and `ip_version`. However, the `snat_enable` parameter (SNAT / source IP replacement) is not yet exposed.

The SDK already supports `SnatEnable *bool` in:
- `CreateTargetGroupRequest` - for setting at creation time
- `ModifyTargetGroupAttributeRequest` - for updating after creation
- `TargetGroupInfo` (response of `DescribeTargetGroups`) - for reading current state

## Goals / Non-Goals

**Goals:**
- Add `snat_enable` as an Optional bool parameter to the `tencentcloud_clb_target_group` resource
- Support setting `SnatEnable` during Create via `CreateTargetGroup` API
- Support reading `SnatEnable` from `DescribeTargetGroups` API response
- Support updating `SnatEnable` via `ModifyTargetGroupAttribute` API
- Add unit tests using gomonkey mock approach
- Update resource documentation

**Non-Goals:**
- Adding any other missing parameters
- Changing existing parameter behavior
- Modifying the service layer method signatures (the Create service method already accepts all needed params; we just need to pass the new field)

## Decisions

1. **Parameter type: Optional bool without ForceNew**
   - Rationale: `SnatEnable` is supported in both `CreateTargetGroupRequest` and `ModifyTargetGroupAttributeRequest`, so it can be set at creation and updated in-place. No need for ForceNew.

2. **Pass SnatEnable through existing CreateTargetGroup service method**
   - The current `ClbService.CreateTargetGroup` method builds the request directly. We need to add `SnatEnable` to the request construction in the service layer.
   - Alternative: Add a new parameter to the service method signature. Rejected because it would require changing all callers. Instead, we set it directly on the request in the resource Create function or extend the service method.

3. **Use `d.GetOkExists` for bool parameter extraction**
   - Rationale: For bool parameters, `GetOkExists` correctly distinguishes between "not set" and "set to false". This matches the pattern used for `full_listen_switch` and `keepalive_enable` in the existing code.

## Risks / Trade-offs

- [Risk: SNAT and transparent source IP are mutually exclusive] → The API documentation states that enabling SNAT disables transparent client source IP forwarding and vice versa. Mitigation: Document this constraint clearly in the schema description. Rely on API-side validation.
- [Risk: Backward compatibility] → Adding an Optional field with no default does not affect existing configurations. Existing resources without `snat_enable` will simply not send the field, preserving current behavior.
