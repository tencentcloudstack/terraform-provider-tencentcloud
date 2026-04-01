# Implementation Tasks: Fix raw_values JSON Diff

## Status Overview

| Phase | Tasks | Completed | Status |
|-------|-------|-----------|--------|
| Implementation | 2 | 0/2 | ⏳ Pending |
| Testing | 5 | 0/5 | ⏳ Pending |
| Review | 2 | 0/2 | ⏳ Pending |
| **Total** | **9** | **0/9 (0%)** | **⏳ Not Started** |

**Estimated Total Time**: 30 minutes

---

## Phase 1: Implementation ⏳

**Estimated Time**: 10 minutes  
**Status**: Pending

### Task 1.1: Update raw_values Schema Definition

- [ ] **File**: `tencentcloud/services/tke/resource_tc_kubernetes_addon.go`
- [ ] **Location**: Lines 56-60
- [ ] **Action**: Add `DiffSuppressFunc: helper.DiffSupressJSON,` to the field
- [ ] **Before**:
  ```go
  "raw_values": {
      Type:        schema.TypeString,
      Optional:    true,
      Description: "Params of addon, base64 encoded json format.",
  },
  ```
- [ ] **After**:
  ```go
  "raw_values": {
      Type:             schema.TypeString,
      Optional:         true,
      Description:      "Params of addon, base64 encoded json format.",
      DiffSuppressFunc: helper.DiffSupressJSON,
  },
  ```
- [ ] **Verification**: 
  - [ ] Import statement exists (line 17: already imported ✅)
  - [ ] Field alignment matches schema style
  - [ ] No syntax errors
- **Estimated Time**: 5 minutes
- **Priority**: P0 (Critical)

### Task 1.2: Format and Validate Code

- [ ] **Run go fmt**:
  ```bash
  go fmt tencentcloud/services/tke/resource_tc_kubernetes_addon.go
  ```
- [ ] **Verify compilation**:
  ```bash
  go build ./tencentcloud/services/tke/resource_tc_kubernetes_addon.go
  ```
- [ ] **Check for syntax errors**
- [ ] **Verify no new linter warnings**:
  ```bash
  golangci-lint run tencentcloud/services/tke/resource_tc_kubernetes_addon.go
  ```
- **Estimated Time**: 5 minutes
- **Priority**: P0 (Critical)

---

## Phase 2: Testing ⏳

**Estimated Time**: 15 minutes  
**Status**: Pending

### Task 2.1: Test - Different Key Ordering

- [ ] **Setup**: Create test configuration
  ```hcl
  resource "tencentcloud_kubernetes_addon" "test" {
    cluster_id   = "cls-test123"
    addon_name   = "nginx-ingress"
    raw_values   = jsonencode({
      key1 = "value1"
      key2 = "value2"
      key3 = "value3"
    })
  }
  ```
- [ ] **Action**: 
  1. Apply configuration
  2. Simulate API returning different order (mock if needed)
  3. Run `terraform plan`
- [ ] **Expected Result**: No diff shown
- [ ] **Actual Result**: _____
- **Estimated Time**: 3 minutes
- **Priority**: P0 (Critical)

### Task 2.2: Test - Different Whitespace/Formatting

- [ ] **Test Case**: Compact vs formatted JSON
  ```json
  Input:  {"key":"value","nested":{"a":1}}
  Output: {  "key" : "value" , "nested" : { "a" : 1 } }
  ```
- [ ] **Expected Result**: No diff shown
- [ ] **Actual Result**: _____
- **Estimated Time**: 2 minutes
- **Priority**: P1 (High)

### Task 2.3: Test - Actual Value Change

- [ ] **Test Case**: Modify actual value
  ```hcl
  # Change: key1 = "value1" -> key1 = "new_value"
  ```
- [ ] **Action**: Run `terraform plan`
- [ ] **Expected Result**: Diff IS shown (confirming real changes detected)
- [ ] **Actual Result**: _____
- [ ] **Verification**: Ensure DiffSuppressFunc doesn't hide real changes
- **Estimated Time**: 3 minutes
- **Priority**: P0 (Critical)

### Task 2.4: Test - Complex Nested JSON

- [ ] **Test Case**: Real-world addon configuration
  ```hcl
  raw_values = jsonencode({
    controller = {
      replicas = 2
      resources = {
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
        requests = {
          cpu    = "250m"
          memory = "256Mi"
        }
      }
    }
    service = {
      type = "LoadBalancer"
      annotations = {
        "service.kubernetes.io/loadbalance-id" = "lb-xxx"
      }
    }
  })
  ```
- [ ] **Action**: Apply, then plan (no changes)
- [ ] **Expected Result**: No diff for nested structures
- [ ] **Actual Result**: _____
- **Estimated Time**: 5 minutes
- **Priority**: P1 (High)

### Task 2.5: Test - Edge Cases

- [ ] **Test Case 1**: Empty JSON object `{}`
  - [ ] Expected: No diff
  - [ ] Result: _____
- [ ] **Test Case 2**: JSON array values `{"list": [1,2,3]}`
  - [ ] Expected: No diff (same order), diff (different order)
  - [ ] Result: _____
- [ ] **Test Case 3**: Invalid JSON (if somehow set)
  - [ ] Expected: Falls back to string comparison, no crash
  - [ ] Result: _____
- **Estimated Time**: 2 minutes
- **Priority**: P2 (Medium)

---

## Phase 3: Code Review ⏳

**Estimated Time**: 5 minutes  
**Status**: Pending

### Task 3.1: Self Review Checklist

- [ ] **Code Quality**:
  - [ ] Follows existing code style
  - [ ] Proper indentation and alignment
  - [ ] Consistent with similar fields in codebase
  - [ ] No commented-out code
- [ ] **Best Practices**:
  - [ ] Uses existing helper function (not reinventing)
  - [ ] Follows pattern from similar resources
  - [ ] No new functions added (per requirement)
- [ ] **Completeness**:
  - [ ] Change is minimal and focused
  - [ ] No unrelated modifications
  - [ ] go fmt applied
  - [ ] No debug statements left
- **Estimated Time**: 3 minutes
- **Priority**: P0 (Critical)

### Task 3.2: Final Verification

- [ ] **Compilation Check**:
  ```bash
  go build ./tencentcloud/services/tke/...
  ```
- [ ] **Linter Check**:
  ```bash
  golangci-lint run tencentcloud/services/tke/resource_tc_kubernetes_addon.go
  ```
- [ ] **No New Warnings**: Verify error count hasn't increased
- [ ] **File Size**: Check lines changed (should be ~1-2 lines)
- [ ] **Git Diff Review**:
  ```bash
  git diff tencentcloud/services/tke/resource_tc_kubernetes_addon.go
  ```
- **Estimated Time**: 2 minutes
- **Priority**: P0 (Critical)

---

## Verification Checklist

### Before Implementation
- [x] helper.DiffSupressJSON exists in codebase ✅
- [x] helper package imported in target file ✅
- [x] Similar pattern exists in other resources ✅
- [x] No breaking changes required ✅

### During Implementation
- [ ] Schema field updated
- [ ] Code formatted with go fmt
- [ ] Compilation successful
- [ ] No new linter errors

### After Implementation
- [ ] All test cases pass
- [ ] No regression in existing functionality
- [ ] Code follows established patterns
- [ ] Ready for merge

---

## Risk Mitigation

### Critical Path Items (Must Complete)
1. ✅ Task 1.1 - Schema update
2. ✅ Task 1.2 - Format and validate
3. ✅ Task 2.1 - Test key ordering (core issue)
4. ✅ Task 2.3 - Test real changes still detected
5. ✅ Task 3.2 - Final verification

### Nice-to-Have Items (Can Skip if Time Limited)
- Task 2.2 - Whitespace testing
- Task 2.5 - Edge case testing

### Rollback Plan
If any critical test fails:
1. Revert the one-line change
2. Run go fmt
3. Verify compilation
4. Document issue for further investigation

---

## Dependencies

### Required
- ✅ `helper.DiffSupressJSON` function exists
- ✅ `helper` package imported
- ✅ Go development environment

### Optional
- Test Kubernetes cluster (for integration testing)
- Real addon configurations (for comprehensive testing)

---

## Notes

### Code Placement Requirement
> "在主资源如果新增函数，代码块应该放在整个资源代码下方"

**Status**: ✅ **Not Applicable**
- No new functions are being added
- Reusing existing `helper.DiffSupressJSON`
- Only modifying schema definition

### Go Fmt Requirement
> "每更新完毕一个go文件，只要有代码变动，最后结束前都要执行一下go fmt对代码进行格式化并保存"

**Status**: ✅ **Covered in Task 1.2**
- Will run `go fmt` after code changes
- Included in implementation checklist

---

## Progress Tracking

### Session 1: Implementation (Planned)
- [ ] Start: _____
- [ ] Task 1.1 complete: _____
- [ ] Task 1.2 complete: _____
- [ ] End: _____

### Session 2: Testing (Planned)
- [ ] Start: _____
- [ ] All tests complete: _____
- [ ] End: _____

### Session 3: Review (Planned)
- [ ] Start: _____
- [ ] Ready for merge: _____
- [ ] End: _____

---

## Related Files

### Files to Modify
- `tencentcloud/services/tke/resource_tc_kubernetes_addon.go` (1 field change)

### Files to Reference
- `tencentcloud/internal/helper/helper.go` (helper function definition)
- `tencentcloud/services/tke/resource_tc_kubernetes_cluster.go` (reference pattern)
- `tencentcloud/services/cdn/resource_tc_cdn_domain.go` (reference pattern)
- `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy.go` (reference pattern)

### Files NOT Modified
- No new files created
- No test files modified (manual testing only)
- No documentation changes required

---

**Task List Version**: 1.0  
**Created**: 2026-03-24  
**Last Updated**: 2026-03-24  
**Estimated Completion**: 30 minutes after start
