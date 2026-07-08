## Why

The `tencentcloud_cls_machine_group` resource currently does not expose the `OSType` parameter, which is supported by the CLS cloud API across CreateMachineGroup, DescribeMachineGroups (response), and as a filter in DescribeMachineGroups request. Users need to specify the operating system type (Linux or Windows) when creating a machine group to ensure proper LogListener agent configuration.

## What Changes

- Add `OSType` (`TypeInt`, Optional, ForceNew) to `tencentcloud_cls_machine_group` resource schema
- Set `OSType` in `CreateMachineGroup` API request
- Read `OSType` from `MachineGroupInfo` in the Read method
- OSType is ForceNew because `ModifyMachineGroup` API does not support updating this field

## Capabilities

### New Capabilities
- `cls-machine-group-ostype`: Support specifying the operating system type (0: Linux, 1: Windows) when creating a CLS machine group

### Modified Capabilities
<!-- None - this is a new optional parameter, no existing behavior changes -->

## Impact

- **Affected code**: `tencentcloud/services/cls/resource_tc_cls_machine_group.go` (schema, Create, Read)
- **Affected tests**: `tencentcloud/services/cls/resource_tc_cls_machine_group_test.go`
- **Affected docs**: `tencentcloud/services/cls/resource_tc_cls_machine_group.md`
- **API dependencies**: `CreateMachineGroup`, `DescribeMachineGroups` (via `MachineGroupInfo`)
- **Breaking change**: No (new optional parameter with default behavior preserved)