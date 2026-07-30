## 1. Schema Definition

- [x] 1.1 Add `ostype` field (`TypeInt`, Optional) to `resourceTencentCloudClsMachineGroup` schema in `resource_tc_cls_machine_group.go`
- [x] 1.2 Add `meta_tags` field (`TypeMap`, Optional) to `resourceTencentCloudClsMachineGroup` schema in `resource_tc_cls_machine_group.go`

## 2. Create Method

- [x] 2.1 Add `OSType` assignment in `resourceTencentCloudClsMachineGroupCreate`: read `ostype` from schema data and set `request.OSType`
- [x] 2.2 Add `MetaTags` assignment in `resourceTencentCloudClsMachineGroupCreate`: read `meta_tags` map from schema data and set `request.MetaTags`

## 3. Read Method

- [x] 3.1 Add `OSType` read in `resourceTencentCloudClsMachineGroupRead`: read `machineGroup.OSType` (nil-safe) and set `ostype` in state
- [x] 3.2 Add `MetaTags` read in `resourceTencentCloudClsMachineGroupRead`: read `machineGroup.MetaTags` (nil-safe) and set `meta_tags` map in state

## 4. Update Method

- [x] 4.1 Add `immutableArgs` check in `resourceTencentCloudClsMachineGroupUpdate`: check if `ostype` has changed and return error if so
- [x] 4.2 Add `MetaTags` update in `resourceTencentCloudClsMachineGroupUpdate`: when `meta_tags` changes, send updated MetaTags to `ModifyMachineGroup` API

## 5. Tests

- [x] 5.1 Add test cases for `OSType` and `MetaTags` parameters in `resource_tc_cls_machine_group_test.go`

## 6. Documentation

- [x] 6.1 Update `resource_tc_cls_machine_group.md` to include `ostype` and `meta_tags` usage examples

## 7. Verification

- [x] 7.1 Verify the code compiles successfully and all existing tests pass