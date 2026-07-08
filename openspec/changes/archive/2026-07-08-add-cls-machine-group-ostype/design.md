## Context

The `tencentcloud_cls_machine_group` resource manages CLS (Cloud Log Service) machine groups. The cloud API (`CreateMachineGroup`, `DescribeMachineGroups` response) supports an `OSType` parameter to specify the operating system type (0: Linux, 1: Windows) for the machine group. The `ModifyMachineGroup` API does NOT support updating `OSType`, so once set, this field cannot be changed without recreating the resource.

The existing resource code is at `tencentcloud/services/cls/resource_tc_cls_machine_group.go` and uses the `cls/v20201016` SDK package.

## Goals / Non-Goals

**Goals:**
- Add `OSType` (`TypeInt`, Optional, ForceNew) to the resource schema
- Support setting `OSType` in `CreateMachineGroup` API request
- Read `OSType` from `MachineGroupInfo` response in the Read method
- Keep full backward compatibility (existing configurations continue to work)

**Non-Goals:**
- Update `OSType` in-place (API does not support it)
- Add `OSType` filtering to any datasource
- Add any other parameters beyond `OSType`

## Decisions

1. **ForceNew for OSType**: Since `ModifyMachineGroup` API does not include `OSType`, the field must be `ForceNew: true`. Any change to `OSType` will trigger resource recreation.

2. **TypeInt for OSType**: The cloud API uses `*uint64`, so Terraform schema uses `TypeInt`. The values are 0 (Linux) and 1 (Windows).

3. **Optional with no default in schema**: The field is `Optional` only. The default behavior (Linux) is handled by the cloud API, not by Terraform. This avoids state drift issues.

4. **No update logic needed**: Because `ForceNew` is set, the Update method does not need to handle `OSType`. Terraform will automatically destroy and recreate the resource when `OSType` changes.

5. **Read method handles nil**: Since `OSType` is optional in the cloud API response, the Read method must check if `machineGroup.OSType` is nil before setting it in state.

## Risks / Trade-offs

- **[Risk] Nil pointer panic in Read**: If the cloud API response does not include `OSType` for older machine groups, `*machineGroup.OSType` could panic. → Mitigation: Check for nil before dereferencing, use `d.Set("ostype", int(*machineGroup.OSType))` only when `machineGroup.OSType != nil`.
- **[Risk] ForceNew causes resource recreation**: Users who add `OSType` to an existing resource will see a destroy + create. → Mitigation: This is documented in the description and is the expected behavior for a field that the API does not support updating.