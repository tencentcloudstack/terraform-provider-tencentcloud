## Why

The `tencentcloud_cls_machine_group` resource currently does not expose the `OSType` and `MetaTags` parameters, which are supported by the CLS cloud API across CreateMachineGroup, ModifyMachineGroup, and DescribeMachineGroups (response). Users need to specify the operating system type (Linux or Windows) when creating a machine group to ensure proper LogListener agent configuration, and need MetaTags to attach metadata key-value pairs to machine groups.

## What Changes

- Add `OSType` (`TypeInt`, Optional) to `tencentcloud_cls_machine_group` resource schema
- Add `MetaTags` (`TypeMap`, Optional) to `tencentcloud_cls_machine_group` resource schema
- Set `OSType` in `CreateMachineGroup` API request
- Set `MetaTags` in `CreateMachineGroup` API request
- Read `OSType` from `MachineGroupInfo` in the Read method
- Read `MetaTags` from `MachineGroupInfo` in the Read method
- Set `MetaTags` in `ModifyMachineGroup` API request when changed
- Use `immutableArgs` check in Update method to prevent `ostype` from being changed (instead of ForceNew), returning an error if the user attempts to modify it

## Capabilities

### New Capabilities
- `cls-machine-group-ostype`: Support specifying the operating system type (0: Linux, 1: Windows) when creating a CLS machine group
- `cls-machine-group-meta-tags`: Support specifying metadata key-value pairs (MetaTags) when creating or updating a CLS machine group

### Modified Capabilities
<!-- None - these are new optional parameters, no existing behavior changes -->

## Impact

- **Affected code**: `tencentcloud/services/cls/resource_tc_cls_machine_group.go` (schema, Create, Read, Update)
- **Affected tests**: `tencentcloud/services/cls/resource_tc_cls_machine_group_test.go`
- **Affected docs**: `tencentcloud/services/cls/resource_tc_cls_machine_group.md`
- **API dependencies**: `CreateMachineGroup`, `ModifyMachineGroup`, `DescribeMachineGroups` (via `MachineGroupInfo`)
- **Breaking change**: No (new optional parameters with default behavior preserved)