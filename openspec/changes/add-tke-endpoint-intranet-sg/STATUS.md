# Change Status

**Change ID**: add-tke-endpoint-intranet-sg  
**Status**: Draft  
**Created**: 2026-03-05  
**Last Updated**: 2026-03-05

## Current Status

📝 **Draft** - Proposal created, awaiting review and approval

## Progress Tracking

### Proposal Stage
- [x] Problem statement documented
- [x] Solution designed
- [x] Tasks broken down
- [ ] Proposal reviewed
- [ ] Proposal approved

### Implementation Stage
- [ ] Schema updated with new field
- [ ] Create function updated
- [ ] Helper function signature updated
- [ ] Update function calls updated
- [ ] Delete function call updated
- [ ] Service layer updated
- [ ] Validation logic added
- [ ] Documentation updated
- [ ] Changelog entry added
- [ ] Manual testing completed
- [ ] Code review completed

### Completion Stage
- [ ] All tasks completed
- [ ] Ready for merge
- [ ] Merged to main branch

## Next Steps

1. **Review the proposal** - Ensure the approach aligns with TKE API capabilities
2. **Verify API support** - Confirm TKE API accepts security group for intranet endpoints
3. **Get approval** - Confirm the design and implementation plan
4. **Begin implementation** - Start with schema and validation changes

## Implementation Summary

### Key Changes

1. **New Field**: `cluster_intranet_security_group` (ForceNew)
   - Allows security group configuration for intranet endpoints
   - Required when `cluster_intranet` is true
   - Cannot be modified after creation (ForceNew)

2. **Updated Descriptions**:
   - `cluster_internet_security_group` - Clarified as "internet" (external network)
   - `cluster_intranet_security_group` - New field for "intranet" (internal network)

3. **Service Layer Update**:
   - Remove internet-only restriction in `CreateClusterEndpoint`
   - Support security group for both internet and intranet

4. **Validation Logic**:
   - Enforce security group when intranet enabled
   - Prevent security group when intranet disabled

### Files to Modify

- `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`
- `tencentcloud/services/tke/service_tencentcloud_tke.go`
- `website/docs/r/kubernetes_cluster_endpoint.html.markdown` (if exists)
- `.changelog/<next-number>.txt`

## Risks and Considerations

1. **API Compatibility**: Need to verify TKE API supports security group for intranet endpoints
2. **ForceNew Impact**: Users will need to recreate endpoints if security group changes
3. **Validation Breaking**: New validation may require existing users to add security group field
4. **Testing Coverage**: Requires real TKE cluster for comprehensive testing

## Related Issues

- Addresses feature parity between internet and intranet security group configuration
- Aligns with security best practices for cluster endpoint access control

## Dependencies

- TKE API support for intranet security groups (assumed available, needs verification)
- No breaking changes to existing API calls

## Notes

- Read/Update stages intentionally do not support security group (per API limitations)
- ForceNew ensures correct behavior when security group changes
- Backward compatible for existing configurations (optional field)
