## Why

The `tencentcloud_clb_target_group` resource currently does not support the `SnatEnable` parameter, which controls whether SNAT (source IP replacement) is enabled for the target group. Users who need to enable SNAT must configure it manually through the console or API, reducing Infrastructure as Code coverage.

## What Changes

- Add `snat_enable` (bool, Optional) parameter to the `tencentcloud_clb_target_group` resource schema
- Pass `SnatEnable` to `CreateTargetGroup` API during resource creation
- Read `SnatEnable` from `DescribeTargetGroups` API response in the Read function
- Pass `SnatEnable` to `ModifyTargetGroupAttribute` API during resource update

## Capabilities

### New Capabilities
- `target-group-snat-enable`: Support for enabling/disabling SNAT (source IP replacement) on CLB target groups via the `snat_enable` parameter

### Modified Capabilities

## Impact

- `tencentcloud/services/clb/resource_tc_clb_target_group.go` - Add schema field and CRUD logic for `snat_enable`
- `tencentcloud/services/clb/resource_tc_clb_target_group_test.go` - Add unit tests for the new parameter
- `tencentcloud/services/clb/resource_tc_clb_target_group.md` - Update documentation with new parameter example
