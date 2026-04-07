# Add Layers Update Support to tencentcloud_scf_function

## Quick Summary

**What**: Enable users to update the `layers` field of existing SCF functions without recreation

**Why**: Currently, users cannot modify layers after function creation, requiring destroy-and-recreate operations

**How**: Add `HasChange("layers")` check in Update function, following the same pattern as `l5_enable`

**Effort**: ~60 minutes total | **Risk**: 🟢 Very Low | **Impact**: 🟢 High User Value

---

## Status

- **Proposal Date**: 2026-03-24
- **Status**: 📋 **Awaiting Approval**
- **Implementation Status**: ⏳ Not Started
- **Target Version**: Next minor release

---

## Problem

### Current Situation ❌

```hcl
resource "tencentcloud_scf_function" "example" {
  name    = "my-function"
  runtime = "Python3.6"
  
  layers {
    layer_name    = "my-layer"
    layer_version = 1  # User wants to update to version 2
  }
}
```

**User changes version 1 → 2 in config**

```bash
$ terraform apply
# ❌ Terraform either:
#   - Ignores the change (drift), or
#   - Forces recreation (downtime)
```

### Expected Behavior ✅

```bash
$ terraform apply
# ✅ Terraform calls UpdateFunctionConfiguration API
# ✅ Layers updated without recreation
# ✅ No downtime
```

---

## Solution

Add update support for `layers` field by implementing:

1. **Resource Update Logic** (~15 lines in `resource_tc_scf_function.go`)
   - Check `d.HasChange("layers")`
   - Parse new layers value
   - Assign to `functionInfo.layers`

2. **Service Update Logic** (~3 lines in `service_tencentcloud_scf.go`)
   - Assign `request.Layers = info.layers`
   - API call already exists, just add layers to request

**Total Code**: ~20 lines  
**API Support**: ✅ Already exists (UpdateFunctionConfiguration supports Layers parameter)

---

## Key Benefits

| Benefit | Description |
|---------|-------------|
| 🎯 **Better UX** | Users can update layers without destroying functions |
| ⚡ **No Downtime** | Updates don't require recreation |
| 🔧 **Operational** | Simpler deployment workflows |
| 📊 **Consistency** | Aligns with other updateable fields (environment, vpc_id, etc.) |

---

## Technical Details

### Code Changes

**File 1**: `tencentcloud/services/scf/resource_tc_scf_function.go`
- **Location**: After line 1279 (after `l5_enable` block)
- **Change**: Add `HasChange("layers")` check and parsing logic
- **Pattern**: Mirrors existing `l5_enable` pattern

**File 2**: `tencentcloud/services/scf/service_tencentcloud_scf.go`
- **Location**: After line 316 (in `ModifyFunctionConfig` function)
- **Change**: Add `request.Layers = info.layers`
- **Pattern**: Follows existing field assignment pattern

### Why This Works

1. ✅ **Schema Already Exists**: `layers` field is already defined in schema
2. ✅ **Create Works**: Layers are already parsed and sent during Create
3. ✅ **Read Works**: Read function already populates layers from API
4. ✅ **API Supports It**: `UpdateFunctionConfiguration` API accepts Layers parameter
5. ✅ **Only Update Missing**: Just need to add update logic

---

## Test Scenarios

- ✅ Add layers to function without layers
- ✅ Remove all layers from function
- ✅ Update layer version (1 → 2)
- ✅ Change layer name
- ✅ Manage multiple layers
- ✅ Verify no state drift after apply
- ✅ Test concurrent updates (layers + environment)

---

## Documents

1. **[proposal.md](./proposal.md)** - Complete proposal with motivation, solution, and impact analysis
2. **[design.md](./design.md)** - Technical design, architecture, and implementation details
3. **[tasks.md](./tasks.md)** - Detailed task breakdown and progress tracking (15 tasks)

---

## Implementation Timeline

| Phase | Duration | Tasks |
|-------|----------|-------|
| Code Implementation | 10 min | 3 tasks |
| Testing | 30 min | 7 tests |
| Documentation | 10 min | 2 tasks |
| Code Review | 10 min | 3 checks |
| **Total** | **~60 min** | **15 tasks** |

---

## Risk Assessment

| Risk Factor | Level | Notes |
|-------------|-------|-------|
| Breaking Changes | 🟢 None | Purely additive feature |
| State Drift | 🟢 Low | Read function already handles layers |
| API Compatibility | 🟢 None | API already supports this |
| Backward Compatibility | 🟢 None | Existing resources unaffected |
| **Overall Risk** | **🟢 Very Low** | Simple, safe change |

---

## Success Criteria

- [x] Proposal documents created
- [ ] Proposal approved
- [ ] Code implemented (~20 lines)
- [ ] Code formatted with `go fmt`
- [ ] All manual tests pass
- [ ] No state drift after updates
- [ ] Documentation updated (if needed)
- [ ] Code review approved
- [ ] Merged to main branch

---

## Example Usage

### Before This Change

```hcl
# User must destroy and recreate to change layers
resource "tencentcloud_scf_function" "example" {
  name    = "my-function"
  runtime = "Python3.6"
  
  layers {
    layer_name    = "my-layer"
    layer_version = 1
  }
  
  # To update to version 2:
  # 1. Remove resource or use taint
  # 2. terraform destroy (downtime!)
  # 3. Update version to 2
  # 4. terraform apply (creates new function)
}
```

### After This Change

```hcl
# User simply updates the config
resource "tencentcloud_scf_function" "example" {
  name    = "my-function"
  runtime = "Python3.6"
  
  layers {
    layer_name    = "my-layer"
    layer_version = 2  # Just change this
  }
}

# terraform apply → Updates in-place ✅
```

---

## Dependencies

### Prerequisites

- None - all required infrastructure already exists

### Related Changes

- None - independent change

### Follows Pattern Of

- `l5_enable` field update logic (lines 1276-1279)
- `environment` field update logic (lines 1229-1232)
- `vpc_id` field update logic (lines 1239-1253)

---

## Approval Process

### Required Reviews

1. **Technical Review**: Verify code follows project conventions
2. **Testing Review**: Ensure all test scenarios covered
3. **Documentation Review**: Check if docs need updates

### Approval Checklist

- [ ] Proposal reviewed and approved
- [ ] Technical design approved
- [ ] Implementation plan approved
- [ ] Ready to proceed with `/opsx: apply`

---

## Quick Start (After Approval)

```bash
# 1. Implement code changes
# (See tasks.md for detailed steps)

# 2. Format code
go fmt tencentcloud/services/scf/resource_tc_scf_function.go
go fmt tencentcloud/services/scf/service_tencentcloud_scf.go

# 3. Test manually
cd examples/dist
# Edit main.tf to test layers updates
terraform plan
terraform apply

# 4. Verify no drift
terraform plan  # Should show "No changes"
```

---

## Questions & Answers

**Q: Why wasn't this implemented initially?**  
A: Likely an oversight. The API supports it, and the Create function already implements the parsing logic.

**Q: Will this break existing configurations?**  
A: No. This is purely additive. Existing functions continue working unchanged.

**Q: What happens if a layer doesn't exist?**  
A: The API returns an error, which Terraform propagates to the user with a clear message.

**Q: Can we update multiple layers at once?**  
A: Yes. The implementation handles arrays of layers.

**Q: Does this require provider version upgrade?**  
A: Users must upgrade to the version that includes this change to use the feature.

---

## Related Issues

- Similar to how `environment`, `vpc_id`, and other fields support updates
- Closes gap between API capabilities and provider functionality

---

## Next Steps

1. **Review this proposal** and the linked documents
2. **Provide feedback** or approve
3. **Execute implementation** with `/opsx: apply`
4. **Test and validate** all scenarios
5. **Archive as spec** with `/opsx: archive`

---

## Contact

**Proposal Author**: CodeBuddy AI Assistant  
**Date**: 2026-03-24  
**Change ID**: add-scf-function-layers-update-support

For questions or feedback about this proposal, please review the detailed documents:
- **[proposal.md](./proposal.md)** - Why and what
- **[design.md](./design.md)** - How and architecture
- **[tasks.md](./tasks.md)** - Implementation tasks
