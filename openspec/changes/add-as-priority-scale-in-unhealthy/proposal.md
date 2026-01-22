# Proposal: Add PriorityScaleInUnhealthy Parameter to AS Scaling Group

## Change ID
`add-as-priority-scale-in-unhealthy`

## Summary
Add support for the `PriorityScaleInUnhealthy` parameter in the `tencentcloud_as_scaling_group` resource. This parameter is part of the `ServiceSettings` structure in TencentCloud Auto Scaling API and controls whether unhealthy instances should be prioritized during scale-in operations.

## Motivation
The TencentCloud Auto Scaling API supports the `PriorityScaleInUnhealthy` parameter in `ServiceSettings`, which allows users to configure whether instances marked as unhealthy should be removed first during scale-in operations. Currently, the Terraform provider does not expose this parameter, limiting users' ability to fully configure their auto scaling groups' behavior.

## Background
The `ServiceSettings` structure in Auto Scaling Group contains several configuration options:
- `ReplaceMonitorUnhealthy` - Already supported
- `ScalingMode` - Already supported  
- `ReplaceLoadBalancerUnhealthy` - Already supported
- `ReplaceMode` - Already supported
- `DesiredCapacitySyncWithMaxMinSize` - Already supported
- `PriorityScaleInUnhealthy` - **NOT YET SUPPORTED** (this proposal)

### API References
- **Create Scaling Group**: https://cloud.tencent.com/document/product/377/20440
- **Describe Scaling Groups**: https://cloud.tencent.com/document/product/377/20438
- **Modify Scaling Group**: https://cloud.tencent.com/document/product/377/20433

## Proposed Changes

### Resource Schema Addition
Add a new optional boolean field `priority_scale_in_unhealthy` to the `tencentcloud_as_scaling_group` resource schema.

**Field Specification:**
- **Name**: `priority_scale_in_unhealthy`
- **Type**: `schema.TypeBool`
- **Required**: No (Optional)
- **Default**: Not set (API default applies)
- **Description**: "Whether to enable priority for unhealthy instances during scale-in operations. If set to `true`, unhealthy instances will be removed first when scaling in."

### Implementation Scope
1. **Schema Definition**: Add the new field to the resource schema
2. **Create Operation**: Include the parameter in `CreateAutoScalingGroup` API call within `ServiceSettings`
3. **Read Operation**: Read and set the parameter value from `DescribeAutoScalingGroups` API response
4. **Update Operation**: Support updating the parameter via `ModifyAutoScalingGroup` API call
5. **Documentation**: Update the resource documentation with usage examples

## User Impact

### Benefits
- Users gain full control over scale-in behavior for unhealthy instances
- Aligns Terraform provider capabilities with TencentCloud API features
- No breaking changes - this is a pure addition

### Migration
- Existing configurations continue to work without changes
- Users can optionally add the new parameter to their configurations
- No state migration required

## Implementation Complexity
**Low** - This change follows the existing pattern for other `ServiceSettings` parameters already implemented in the resource.

## Alternatives Considered
None. This is a straightforward feature addition to expose existing API functionality.

## Success Criteria
- [ ] Schema field added and properly validated
- [ ] Create operation includes the parameter
- [ ] Read operation correctly retrieves the parameter value
- [ ] Update operation can modify the parameter
- [ ] Unit tests pass
- [ ] Acceptance tests demonstrate the functionality
- [ ] Documentation includes usage examples
- [ ] Code passes linting and formatting checks

## Timeline
Estimated 1-2 days for implementation and testing.

## Dependencies
- TencentCloud SDK Go v20180419 (already in vendor)
- No new external dependencies required

## Related Changes
None. This is an independent feature addition.
