# Status: Fix GWLB Target Group Port Field Drift

**Change ID**: fix-gwlb-target-group-port-drift  
**Status**: Archived ✅  
**Created**: 2026-03-05  
**Completed**: 2026-03-05  
**Archived**: 2026-03-05  
**Author**: AI Assistant

## Overview
Fix configuration drift issue caused by the `port` field in `tencentcloud_gwlb_target_group` resource by adding the `Computed` attribute.

## Archive Notes
This change has been successfully completed and archived. All implementation tasks were finished, tested, and the fix is ready for deployment.

**Archive Location**: `openspec/changes/archive/2026-03-05-fix-gwlb-target-group-port-drift/`

## Current State
- ✅ Proposal created and approved
- ✅ Implementation completed
- ✅ All tasks finished
- ✅ Documentation updated
- ✅ Changelog entry added
- ✅ Test cases created

## Implementation Summary

### Code Changes
**File**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group.go`
- Added `Computed: true` to the `port` field schema (line 41)

### Documentation Updates
**File**: `.changelog/3841.txt`
- Created changelog entry for the bug fix

### Test Cases
**Directory**: `examples/tencentcloud-gwlb-target-group/`
- Created test configuration with two scenarios
- Added comprehensive test documentation in README.md

## Verification

### Test Results
✅ **Test Case 1** (without port): Creates successfully, no drift detected  
✅ **Test Case 2** (with port=6081): Uses specified value, no drift detected  
✅ **Linter Check**: No new errors introduced (only pre-existing deprecated warnings)

### Scenarios Validated
1. ✅ Port not specified → API default accepted, no drift
2. ✅ Port explicitly specified → User value respected, no drift
3. ✅ Existing resources → Compatible with fix, no breaking changes

## Next Steps
Ready for archive. Use:
```bash
openspec archive fix-gwlb-target-group-port-drift --yes
```
