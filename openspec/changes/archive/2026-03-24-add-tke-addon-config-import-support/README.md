# Add Import Support to tencentcloud_kubernetes_addon_config

## Quick Summary

**What**: Add Terraform import capability to `tencentcloud_kubernetes_addon_config` resource

**Why**: Enable users to import existing TKE addon configurations into Terraform state

**How**: Add `Importer` configuration with `ImportStatePassthrough` (3 lines of code)

**Effort**: ~50 minutes total | **Risk**: 🟢 Very Low | **Impact**: 🟢 High User Value

---

## Status

- **Proposal Date**: 2026-03-24
- **Status**: 🔄 **Implementation In Progress**
- **Implementation Status**: ✅ Code Complete, ⏳ Testing Pending
- **Target Version**: Next minor release

---

## Documents

1. **[proposal.md](./proposal.md)** - Complete proposal with motivation, solution, and impact analysis
2. **[design.md](./design.md)** - Technical design, architecture diagrams, and implementation details
3. **[tasks.md](./tasks.md)** - Detailed task breakdown and progress tracking

---

## Key Points

### Why This Matters

Users with existing addon configurations (created via console/API) **cannot** currently bring them under Terraform management without:
- Deleting and recreating them (risky)
- Manually crafting state files (error-prone)

This creates a significant adoption barrier for teams migrating to Infrastructure as Code.

### The Solution

Add standard Terraform import support:

```bash
terraform import tencentcloud_kubernetes_addon_config.example cls-abc123#tcr
```

### Implementation Simplicity

**Only 3 lines of code needed**:
```go
Importer: &schema.ResourceImporter{
    State: schema.ImportStatePassthrough,
},
```

The existing `Read` function already does everything needed for import.

---

## Why It's Low Risk

1. ✅ **Purely Additive**: Zero impact on existing functionality
2. ✅ **Standard Pattern**: Used by 15+ other resources in this provider
3. ✅ **No New Logic**: Leverages existing, battle-tested Read function
4. ✅ **Opt-In**: Users must explicitly run import command
5. ✅ **Easy Rollback**: Can be reverted without affecting existing resources

---

## Success Criteria

- [x] Proposal documents created
- [x] Proposal approved (via `/opsx: apply`)
- [x] Code implemented (3 lines)
- [ ] Acceptance test passes
- [ ] Manual testing validates import works
- [ ] Documentation updated
- [ ] Code review approved
- [ ] Merged to main branch

---

## Example Usage

### Current Pain Point ❌
```bash
$ terraform import tencentcloud_kubernetes_addon_config.tcr cls-abc123#tcr

Error: resource tencentcloud_kubernetes_addon_config doesn't support import
```

### After Implementation ✅
```bash
$ terraform import tencentcloud_kubernetes_addon_config.tcr cls-abc123#tcr

tencentcloud_kubernetes_addon_config.tcr: Importing from ID "cls-abc123#tcr"...
tencentcloud_kubernetes_addon_config.tcr: Import prepared!
tencentcloud_kubernetes_addon_config.tcr: Refreshing state... [id=cls-abc123#tcr]

Import successful!

$ terraform plan
No changes. Your infrastructure matches the configuration.
```

---

## Timeline

| Phase | Duration | Status |
|-------|----------|--------|
| Proposal Review | 1 day | ⏳ Pending |
| Code Implementation | 5 min | ⏳ Waiting for approval |
| Testing | 20 min | ⏳ Waiting for approval |
| Documentation | 10 min | ⏳ Waiting for approval |
| Code Review | 1-2 days | ⏳ Waiting for approval |
| **Total** | **~3 days** | |

---

## Related Resources

### Similar Resources with Import Support
- `tencentcloud_kubernetes_addon` - Uses same pattern
- `tencentcloud_kubernetes_cluster` - Uses custom importer
- `tencentcloud_kubernetes_node_pool` - Uses custom importer
- `tencentcloud_kubernetes_auth_attachment` - Uses same pattern

### Prior Related Work
- **fix-tke-addon-config-raw-values-json-diff** (2026-03-24)
  - Added JSON diff suppression for `raw_values` field
  - Ensures imported resources won't show spurious diffs

---

## Questions & Answers

### Q: Why wasn't import supported from the beginning?
**A**: Likely an oversight during initial implementation. The Read function is already complete.

### Q: Will this affect existing resources?
**A**: No. Import is purely additive and opt-in.

### Q: What happens if import fails?
**A**: Standard Terraform error handling - no state changes if import fails.

### Q: Can imported resources be updated?
**A**: Yes, exactly the same as resources created via Terraform.

### Q: What about raw_values field differences?
**A**: Already handled by the `suppressJSONOrderDiff` function added in the prior change.

---

## Next Steps

1. **Review this proposal** - Check for any concerns or questions
2. **Approve or request changes** - Provide feedback
3. **Implement** - If approved, proceed with the 3-line code change
4. **Test** - Validate import functionality works as expected
5. **Document** - Update resource documentation
6. **Merge** - Include in next release

---

## Contact

- **Proposer**: CodeBuddy AI Assistant
- **Reviewer**: TBD
- **Approver**: TBD

---

## References

- [Terraform Provider Development - Import](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/import)
- [Terraform Plugin SDK - ImportStatePassthrough](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ImportStatePassthrough)
- [TKE DescribeExtensionAddon API](https://cloud.tencent.com/document/product/457)
