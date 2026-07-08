## 1. Schema Definition

- [x] 1.1 Add `ostype` field (`TypeInt`, Optional, ForceNew) to `resourceTencentCloudClsMachineGroup` schema in `resource_tc_cls_machine_group.go`

## 2. Create Method

- [x] 2.1 Add `OSType` assignment in `resourceTencentCloudClsMachineGroupCreate`: read `ostype` from schema data and set `request.OSType`

## 3. Read Method

- [x] 3.1 Add `OSType` read in `resourceTencentCloudClsMachineGroupRead`: read `machineGroup.OSType` (nil-safe) and set `ostype` in state

## 4. Tests

- [x] 4.1 Add test cases for `OSType` parameter in `resource_tc_cls_machine_group_test.go`

## 5. Documentation

- [x] 5.1 Update `resource_tc_cls_machine_group.md` to include `ostype` usage example

## 6. Verification

- [x] 6.1 Verify the code compiles successfully and all existing tests pass