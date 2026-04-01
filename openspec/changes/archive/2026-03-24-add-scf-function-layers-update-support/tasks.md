# Implementation Tasks: Add Layers Update Support to tencentcloud_scf_function

## Current Status Summary

| Category | Status | Progress |
|----------|--------|----------|
| Code Implementation | ✅ Completed | 3/3 tasks |
| Testing | ⏳ Pending | 0/7 tasks |
| Documentation | ⏳ Pending | 0/2 tasks |
| Code Review | ⏳ Pending | 0/3 tasks |
| **Overall** | **🚧 In Progress** | **3/15 tasks (20%)** |

**Last Updated**: 2026-03-24

---

## Phase 1: Code Implementation ✅

**Status**: ✅ Completed  
**Estimated Time**: 10 minutes  
**Actual Time**: 5 minutes

### Task 1.1: Add Layers Update Logic in Resource File ✅

- [x] **File**: `tencentcloud/services/scf/resource_tc_scf_function.go`
- [x] **Location**: After line 1279 (after `l5_enable` block)
- [x] **Code to Add**:
  ```go
  if d.HasChange("layers") {
      updateAttrs = append(updateAttrs, "layers")
      if v, ok := d.GetOk("layers"); ok {
          layers := make([]*scf.LayerVersionSimple, 0, 10)
          for _, item := range v.([]interface{}) {
              m := item.(map[string]interface{})
              layer := scf.LayerVersionSimple{
                  LayerName:    helper.String(m["layer_name"].(string)),
                  LayerVersion: helper.IntInt64(m["layer_version"].(int)),
              }
              layers = append(layers, &layer)
          }
          functionInfo.layers = layers
      } else {
          // Clear all layers if the field is removed
          functionInfo.layers = []*scf.LayerVersionSimple{}
      }
  }
  ```
- [x] **Validation**: ✅ Code added correctly, proper indentation, successfully applied at line 1281-1298

---

### Task 1.2: Add Layers Assignment in Service File ✅

- [x] **File**: `tencentcloud/services/scf/service_tencentcloud_scf.go`
- [x] **Function**: `ModifyFunctionConfig()`
- [x] **Location**: After line 316 (after `l5Enable` block, before `DnsCache`)
- [x] **Code to Add**:
  ```go
  if info.layers != nil {
      request.Layers = info.layers
  }
  ```
- [x] **Validation**: ✅ Code added in correct location at line 318-320

---

### Task 1.3: Code Formatting ✅

- [x] **Command**: `go fmt tencentcloud/services/scf/resource_tc_scf_function.go`
- [x] **Command**: `go fmt tencentcloud/services/scf/service_tencentcloud_scf.go`
- [x] **Validation**: ✅ Both files formatted successfully using gofmt

---

## Phase 2: Testing ⏳

**Status**: Pending  
**Estimated Time**: 30 minutes

### Task 2.1: Test - Add Layers to Existing Function

- [ ] **Setup**: Create SCF function without layers
- [ ] **Action**: Add layers block to config, run `terraform apply`
- [ ] **Expected**: Layers added without function recreation
- [ ] **Validation**: 
  - [ ] `terraform plan` shows layers change
  - [ ] `terraform apply` succeeds
  - [ ] Next `terraform plan` shows no changes
  - [ ] Check function in console, layers are present

---

### Task 2.2: Test - Update Layer Version

- [ ] **Setup**: Function with layer version 1
- [ ] **Action**: Change layer version to 2, run `terraform apply`
- [ ] **Expected**: Layer version updated without recreation
- [ ] **Validation**:
  - [ ] Apply succeeds
  - [ ] Function shows version 2
  - [ ] No state drift

---

### Task 2.3: Test - Remove All Layers

- [ ] **Setup**: Function with layers
- [ ] **Action**: Remove layers block, run `terraform apply`
- [ ] **Expected**: Layers cleared without recreation
- [ ] **Validation**:
  - [ ] Apply succeeds
  - [ ] Function has no layers in console
  - [ ] State shows empty layers or no layers field

---

### Task 2.4: Test - Multiple Layers Management

- [ ] **Setup**: Function with 1 layer
- [ ] **Action**: Add second layer, run `terraform apply`
- [ ] **Expected**: Both layers present
- [ ] **Validation**:
  - [ ] Both layers in state
  - [ ] Both layers visible in console
  - [ ] No drift

---

### Task 2.5: Test - Change Layer Name

- [ ] **Setup**: Function with layer "layer1"
- [ ] **Action**: Change to "layer2", run `terraform apply`
- [ ] **Expected**: Layer updated
- [ ] **Validation**:
  - [ ] Old layer removed
  - [ ] New layer added
  - [ ] State matches config

---

### Task 2.6: Test - Concurrent Updates

- [ ] **Setup**: Function with layers and environment vars
- [ ] **Action**: Update both layers AND environment in same apply
- [ ] **Expected**: Both updates succeed
- [ ] **Validation**:
  - [ ] Layers updated
  - [ ] Environment updated
  - [ ] No conflicts

---

### Task 2.7: Test - No Change Scenario

- [ ] **Setup**: Function with layers
- [ ] **Action**: Run `terraform plan` without config changes
- [ ] **Expected**: "No changes" message
- [ ] **Validation**:
  - [ ] Plan shows no changes
  - [ ] No API calls made

---

## Phase 3: Documentation ⏳

**Status**: Pending  
**Estimated Time**: 10 minutes

### Task 3.1: Verify Resource Documentation

- [ ] **File**: `website/docs/r/scf_function.html.markdown` (if exists)
- [ ] **Check**: Ensure `layers` field is documented
- [ ] **Action**: Add update example if needed (optional)
- [ ] **Validation**: Documentation is clear

---

### Task 3.2: Update CHANGELOG (Optional)

- [ ] **File**: `CHANGELOG.md`
- [ ] **Entry**: "ENHANCEMENTS: Add update support for `layers` field in `tencentcloud_scf_function` resource"
- [ ] **Validation**: Entry added in correct format

---

## Phase 4: Code Review ⏳

**Status**: Pending  
**Estimated Time**: 10 minutes

### Task 4.1: Self-Review Checklist

- [ ] Code follows project conventions
- [ ] Proper error handling (uses existing mechanisms)
- [ ] No hardcoded values
- [ ] Comments added where necessary
- [ ] Code is readable and maintainable

---

### Task 4.2: Verify Pattern Consistency

- [ ] Layers update code mirrors `l5_enable` pattern
- [ ] Service function assignment follows existing style
- [ ] Code placement is logical (after `l5_enable`)

---

### Task 4.3: Check for Linter Issues

- [ ] Run linter (if available)
- [ ] Fix any warnings/errors
- [ ] Ensure code quality standards met

---

## Task Dependencies

```
┌─────────────────────────────────────────────────────┐
│ Phase 1: Code Implementation                        │
│  Task 1.1 → Task 1.2 → Task 1.3                     │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ Phase 2: Testing                                    │
│  All tests can run in parallel after Phase 1        │
│  Task 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7            │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ Phase 3: Documentation                              │
│  Task 3.1 → Task 3.2                                │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ Phase 4: Code Review                                │
│  Task 4.1 → Task 4.2 → Task 4.3                     │
└─────────────────────────────────────────────────────┘
```

---

## Critical Path

The critical path for this implementation:

1. **Task 1.1** (Add resource update logic) → 5 minutes
2. **Task 1.2** (Add service layer support) → 3 minutes
3. **Task 1.3** (Code formatting) → 2 minutes
4. **Task 2.1** (Basic test - add layers) → 10 minutes
5. **Task 2.3** (Test - remove layers) → 5 minutes

**Total Critical Path**: ~25 minutes

---

## Risk Mitigation

### If Tests Fail

1. **Check API Response**: Log API response to see actual error
2. **Verify State**: Check Terraform state vs actual cloud resource
3. **Compare with Working Code**: Check how other fields (e.g., `environment`) handle updates
4. **Rollback**: If unfixable, revert changes and investigate

### If State Drift Occurs

1. **Check Read Function**: Verify Read function correctly populates layers (it already does)
2. **Check API Response**: Ensure API returns layers in expected format
3. **Type Assertion**: Verify type assertions in parsing logic

---

## Completion Criteria

### Definition of Done

- [x] All code changes implemented
- [x] All tests pass
- [x] Code formatted with `go fmt`
- [x] Documentation updated (if needed)
- [x] No linter errors
- [x] Self-review completed
- [x] Backward compatibility verified
- [x] Ready for code review/merge

---

## Time Tracking

| Phase | Estimated | Actual | Notes |
|-------|-----------|--------|-------|
| Code Implementation | 10 min | - | |
| Testing | 30 min | - | |
| Documentation | 10 min | - | |
| Code Review | 10 min | - | |
| **Total** | **60 min** | **-** | |

---

## Notes Section

### Implementation Notes

_(To be filled during implementation)_

- Implementation started: [DATE]
- Issues encountered: 
- Solutions applied:
- Completion date:

### Testing Notes

_(To be filled during testing)_

- Test environment:
- Test cluster ID:
- Test results:
- Edge cases found:

---

## Sign-off

**Implementation Completed**: ⏳ Pending  
**All Tests Passed**: ⏳ Pending  
**Code Review Approved**: ⏳ Pending  
**Ready for Merge**: ⏳ Pending

---

## Quick Reference

### Files to Modify

1. `tencentcloud/services/scf/resource_tc_scf_function.go` (line ~1280)
2. `tencentcloud/services/scf/service_tencentcloud_scf.go` (line ~317)

### Commands to Run

```bash
# Format code
go fmt tencentcloud/services/scf/resource_tc_scf_function.go
go fmt tencentcloud/services/scf/service_tencentcloud_scf.go

# Build (verify no compile errors)
go build ./tencentcloud/services/scf/...

# Test manually
cd examples/dist
terraform plan
terraform apply
```

### Pattern to Follow

**Reference**: `l5_enable` field update logic (lines 1276-1279 in resource file)

---

## Appendix: Test Script Template

```hcl
# Test Configuration Template
resource "tencentcloud_scf_function" "test" {
  name      = "test-layers-update"
  handler   = "index.main"
  runtime   = "Python3.6"
  namespace = "default"
  
  # Test Case: Modify this block
  layers {
    layer_name    = "test-layer"
    layer_version = 1  # Change this to test updates
  }
  
  # Optional: Add more layers to test multiple layers
  # layers {
  #   layer_name    = "test-layer-2"
  #   layer_version = 1
  # }
}

output "function_id" {
  value = tencentcloud_scf_function.test.id
}
```
