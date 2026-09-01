## 1. Schema Fix

- [x] 1.1 Add `Default: 2` to the `timeout` field in the `health_check` nested schema in `resource_tc_gwlb_target_group.go`
- [x] 1.2 Add `Default: 5` to the `interval_time` field in the `health_check` nested schema in `resource_tc_gwlb_target_group.go`
- [x] 1.3 Add `Default: 3` to the `health_num` field in the `health_check` nested schema in `resource_tc_gwlb_target_group.go`
- [x] 1.4 Add `Default: 3` to the `un_health_num` field in the `health_check` nested schema in `resource_tc_gwlb_target_group.go`

## 2. Test Update

- [x] 2.1 Add a test case in `resource_tc_gwlb_target_group_test.go` that creates a target group with `health_check` block but omits `timeout`, `interval_time`, `health_num`, and `un_health_num`, and verifies the resource is created successfully with the correct defaults
- [x] 2.2 Update the existing test configuration if needed to ensure the default behavior is covered

## 3. Documentation

- [x] 3.1 Update `resource_tc_gwlb_target_group.md` to reflect the default values in the health_check field descriptions