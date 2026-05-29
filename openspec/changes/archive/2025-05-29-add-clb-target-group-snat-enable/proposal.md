## Why

The `tencentcloud_clb_target_group` resource currently does not support the `SnatEnable` parameter, which controls whether SNAT (Source Network Address Translation / source IP replacement) is enabled for a target group. Users who need to enable SNAT on their CLB target groups must configure this manually through the console or API, reducing Infrastructure as Code coverage.

## What Changes

- Add a new optional `snat_enable` (bool) parameter to the `tencentcloud_clb_target_group` resource schema.
- Pass `SnatEnable` to the `CreateTargetGroup` API during resource creation.
- Read `SnatEnable` from the `DescribeTargetGroups` API response (`TargetGroupInfo.SnatEnable`) during resource read.
- Pass `SnatEnable` to the `ModifyTargetGroupAttribute` API during resource update.
- Update resource documentation with the new parameter and usage example.

## Capabilities

### New Capabilities
- `target-group-snat`: Support for enabling/disabling SNAT (source IP replacement) on CLB target groups via the `snat_enable` parameter.

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/clb/resource_tc_clb_target_group.go` — schema definition, Create, Read, Update functions.
- **Code**: `tencentcloud/services/clb/service_tencentcloud_clb.go` — service layer method for CreateTargetGroup (if parameters are passed through service layer).
- **Documentation**: `tencentcloud/services/clb/resource_tc_clb_target_group.md` — add `snat_enable` parameter description and example.
- **Tests**: `tencentcloud/services/clb/resource_tc_clb_target_group_test.go` — add unit test for the new parameter.
- **APIs used**: `CreateTargetGroup`, `DescribeTargetGroups`, `ModifyTargetGroupAttribute` (all in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317`).
- **No breaking changes**: The new parameter is optional with a default of `false`, maintaining backward compatibility.
