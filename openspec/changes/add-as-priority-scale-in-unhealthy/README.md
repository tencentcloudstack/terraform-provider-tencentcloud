# Change: Add PriorityScaleInUnhealthy Parameter to AS Scaling Group

## Status
📋 **PROPOSAL** - Awaiting Review and Approval

## Quick Links
- [Proposal Document](./proposal.md) - Complete change proposal with motivation and impact analysis
- [Implementation Tasks](./tasks.md) - Detailed task breakdown (11 tasks across 7 phases)
- [Design Document](./design.md) - Technical design decisions and architecture
- [Spec Delta](./specs/as-scaling-group-service-settings/spec.md) - Formal specification changes

## Summary
Add support for the `priority_scale_in_unhealthy` boolean parameter in the `tencentcloud_as_scaling_group` resource to control whether unhealthy instances should be prioritized during scale-in operations.

## Change ID
`add-as-priority-scale-in-unhealthy`

## Type
Feature Addition - Low Complexity

## Impact
- ✅ No breaking changes
- ✅ Backward compatible
- ✅ Aligns provider with TencentCloud API capabilities
- ✅ Completes ServiceSettings parameter support

## Implementation Estimate
1-2 days

## Files to be Modified
1. `tencentcloud/services/as/resource_tc_as_scaling_group.go` - Add schema field and CRUD operations
2. `tencentcloud/services/as/resource_tc_as_scaling_group_test.go` - Add acceptance tests
3. `tencentcloud/services/as/resource_tc_as_scaling_group.md` - Update documentation

## Validation Checklist
- [x] Proposal document created
- [x] Tasks breakdown completed
- [x] Design document created
- [x] Spec delta written with scenarios
- [ ] OpenSpec validation passed (requires openspec CLI)
- [ ] Peer review completed
- [ ] Approval received

## Next Steps
1. **Review**: Review all documents in this directory
2. **Validate**: Run `openspec validate add-as-priority-scale-in-unhealthy --strict` if CLI available
3. **Approve**: Get stakeholder approval
4. **Implement**: Follow tasks.md sequentially
5. **Test**: Run acceptance tests with TF_ACC=1
6. **Archive**: Move to archive/ after deployment

## API References
- Create: https://cloud.tencent.com/document/product/377/20440
- Describe: https://cloud.tencent.com/document/product/377/20438
- Modify: https://cloud.tencent.com/document/product/377/20433

## Related Work
This completes the ServiceSettings parameters support that includes:
- replace_monitor_unhealthy ✅
- scaling_mode ✅
- replace_load_balancer_unhealthy ✅
- replace_mode ✅
- desired_capacity_sync_with_max_min_size ✅
- priority_scale_in_unhealthy ⏳ (this change)
