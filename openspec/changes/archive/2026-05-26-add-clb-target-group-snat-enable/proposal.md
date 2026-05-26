## Why

The `tencentcloud_clb_target_group` resource currently does not support the `snat_enable` parameter, which controls whether SNAT (Source Network Address Translation / source IP replacement) is enabled. Users who need to enable SNAT on their target groups must configure it manually through the console or API, reducing Infrastructure as Code coverage.

## What Changes

- Add `snat_enable` parameter (bool, Optional) to the `tencentcloud_clb_target_group` resource schema.
- Support setting `SnatEnable` during target group creation via `CreateTargetGroup` API.
- Support updating `SnatEnable` after creation via `ModifyTargetGroupAttribute` API.
- Note: The `DescribeTargetGroups` response (`TargetGroupInfo`) does not return this field, so the Read method cannot read it back from the API. The parameter will be managed as a write-only attribute in Terraform state.

## Capabilities

### New Capabilities
- `snat-enable`: Add SNAT (source IP replacement) toggle support to the CLB target group resource, allowing users to enable or disable client source IP replacement.

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/clb/resource_tc_clb_target_group.go` - Add schema field, Create logic, Update logic
- **Code**: `tencentcloud/services/clb/service_tencentcloud_clb.go` - Extend `CreateTargetGroup` method signature to accept `snatEnable` parameter
- **Documentation**: `tencentcloud/services/clb/resource_tc_clb_target_group.md` - Add example usage with `snat_enable`
- **Tests**: `tencentcloud/services/clb/resource_tc_clb_target_group_test.go` - Add unit tests for the new parameter
- **Dependencies**: No new SDK dependencies required; `SnatEnable` field already exists in `CreateTargetGroupRequest` and `ModifyTargetGroupAttributeRequest`
