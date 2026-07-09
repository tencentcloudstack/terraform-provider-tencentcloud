## Context

The `tencentcloud_cls_machine_group` resource manages CLS (Cloud Log Service) machine groups. The cloud API (`CreateMachineGroup`, `ModifyMachineGroup`, `DescribeMachineGroups` response) supports an `OSType` parameter to specify the operating system type (0: Linux, 1: Windows) for the machine group, and a `MetaTags` parameter to attach metadata key-value pairs to the machine group.

The `ModifyMachineGroup` API does NOT support updating `OSType`, so once set, this field cannot be changed. However, `MetaTags` can be updated via `ModifyMachineGroup`.

The existing resource code is at `tencentcloud/services/cls/resource_tc_cls_machine_group.go` and uses the `cls/v20201016` SDK package.

## Goals / Non-Goals

**Goals:**
- Add `OSType` (`TypeInt`, Optional) to the resource schema
- Add `MetaTags` (`TypeMap`, Optional) to the resource schema
- Support setting `OSType` in `CreateMachineGroup` API request
- Support setting `MetaTags` in `CreateMachineGroup` API request
- Read `OSType` from `MachineGroupInfo` response in the Read method
- Read `MetaTags` from `MachineGroupInfo` response in the Read method
- Support updating `MetaTags` in `ModifyMachineGroup` API request
- Use `immutableArgs` check in Update method to prevent `OSType` from being changed, returning an error instead of using ForceNew
- Keep full backward compatibility (existing configurations continue to work)

**Non-Goals:**
- Update `OSType` in-place (API does not support it)
- Add `OSType` filtering to any datasource
- Add any other parameters beyond `OSType` and `MetaTags`

## Decisions

1. **immutableArgs for OSType**: Since `ModifyMachineGroup` API does not include `OSType`, the field uses an `immutableArgs` check in the Update method. Any change to `OSType` will return an error message indicating the argument cannot be changed. This approach is preferred over `ForceNew` to avoid accidental resource recreation.

2. **TypeInt for OSType**: The cloud API uses `*uint64`, so Terraform schema uses `TypeInt`. The values are 0 (Linux) and 1 (Windows).

3. **TypeMap for MetaTags**: The cloud API uses `[]*MetaTagInfo` (Key/Value pairs), so Terraform schema uses `TypeMap` for a simple key-value map representation. This provides a clean user experience.

4. **Optional with no default in schema**: Both fields are `Optional` only. The default behavior is handled by the cloud API, not by Terraform. This avoids state drift issues.

5. **Update logic for MetaTags**: Because `MetaTags` is supported by `ModifyMachineGroup`, the Update method handles `MetaTags` changes by sending the updated map to the API.

6. **Read method handles nil**: Since `OSType` and `MetaTags` are optional in the cloud API response, the Read method must check for nil before setting them in state.

## Risks / Trade-offs

- **[Risk] Nil pointer panic in Read**: If the cloud API response does not include `OSType` or `MetaTags` for older machine groups, dereferencing could panic. → Mitigation: Check for nil before dereferencing.
- **[Risk] immutableArgs error message**: Users who try to change `OSType` will get an error instead of a plan showing destroy + create. → Mitigation: The error message clearly states which argument cannot be changed, and this is documented in the description.