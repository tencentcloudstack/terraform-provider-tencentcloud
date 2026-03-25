# Tasks: Fix JSON Order Diff in TKE Addon Config

## Implementation Tasks

### 1. Code Implementation ✅

#### 1.1 Update Imports ✅
- **File:** `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
- **Action:** Add required imports
- **Status:** ✅ Complete
- **Details:**
  ```go
  import (
      // ... existing imports ...
      "encoding/json"  // Added for JSON parsing
      "reflect"        // Added for deep comparison
  )
  ```

#### 1.2 Add Diff Suppression Function ✅
- **File:** `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
- **Action:** Add `suppressJSONOrderDiff` function at end of file
- **Status:** ✅ Complete
- **Location:** After `resourceTencentCloudKubernetesAddonConfigDelete` function
- **Function:**
  ```go
  func suppressJSONOrderDiff(k, old, new string, d *schema.ResourceData) bool
  ```

#### 1.3 Update Schema Definition ✅
- **File:** `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
- **Action:** Add `DiffSuppressFunc` to `raw_values` field
- **Status:** ✅ Complete
- **Change:**
  ```go
  "raw_values": {
      Type:             schema.TypeString,
      Optional:         true,
      Computed:         true,
      Description:      "Params of addon, base64 encoded json format.",
      DiffSuppressFunc: suppressJSONOrderDiff,  // Added this line
  },
  ```

#### 1.4 Code Formatting ✅
- **Action:** Run `go fmt` on modified file
- **Status:** ✅ Complete
- **Command:** `gofmt -w tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`

### 2. Testing Tasks ⏳

#### 2.1 Unit Tests ⏳
- **File:** `tencentcloud/services/tke/resource_tc_kubernetes_addon_config_test.go`
- **Action:** Add unit tests for `suppressJSONOrderDiff` function
- **Status:** ⏳ To Do
- **Test Cases:**
  - [ ] Empty strings (both empty → true)
  - [ ] One empty string → false
  - [ ] Same JSON, different order → true
  - [ ] Different JSON content → false
  - [ ] Nested objects with different order → true
  - [ ] Arrays with different order → false (order matters for arrays)
  - [ ] Invalid JSON fallback to string comparison
  - [ ] Mixed valid/invalid JSON

**Example test structure:**
```go
func TestSuppressJSONOrderDiff(t *testing.T) {
    // Test cases as documented in design.md
}
```

#### 2.2 Acceptance Tests ⏳
- **File:** `tencentcloud/services/tke/resource_tc_kubernetes_addon_config_test.go`
- **Action:** Add or update acceptance test
- **Status:** ⏳ To Do
- **Requirements:**
  - [ ] Create addon config with specific `raw_values`
  - [ ] Verify no diff on subsequent plan (even if API reorders)
  - [ ] Verify diff IS shown when content actually changes

**Test command:**
```bash
TF_ACC=1 go test -v ./tencentcloud/services/tke -run TestAccTencentCloudKubernetesAddonConfig
```

#### 2.3 Manual Testing ⏳
- **Status:** ⏳ To Do
- **Steps:**
  1. [ ] Create a test TKE cluster
  2. [ ] Apply addon config with specific `raw_values` JSON
  3. [ ] Verify addon is created successfully
  4. [ ] Run `terraform plan` immediately after apply
  5. [ ] Confirm no diff is shown (test JSON order independence)
  6. [ ] Modify `raw_values` with actual content change
  7. [ ] Verify diff IS shown
  8. [ ] Clean up test resources

### 3. Code Quality Tasks ⏳

#### 3.1 Linter Checks ⏳
- **Action:** Run linter on modified file
- **Status:** ⏳ To Do (Note: Currently shows 1 pre-existing deprecated warning, not related to this change)
- **Command:**
  ```bash
  golangci-lint run ./tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go
  ```

#### 3.2 Build Verification ⏳
- **Action:** Verify provider builds successfully
- **Status:** ⏳ To Do
- **Command:**
  ```bash
  go build -o terraform-provider-tencentcloud
  ```

#### 3.3 Vet Check ⏳
- **Action:** Run `go vet`
- **Status:** ⏳ To Do
- **Command:**
  ```bash
  go vet ./tencentcloud/services/tke/
  ```

### 4. Documentation Tasks ⏳

#### 4.1 Changelog Entry ⏳
- **File:** `.changelog/<PR_NUMBER>.txt` (create after PR number assigned)
- **Status:** ⏳ To Do
- **Content:**
  ```
  ```enhancement:tke```: Fix raw_values field showing false-positive diffs due to JSON element ordering
  ```

#### 4.2 Resource Documentation Review ⏳
- **File:** `website/docs/r/kubernetes_addon_config.html.markdown`
- **Action:** Review if documentation needs update
- **Status:** ⏳ To Do
- **Expected:** Likely no change needed (behavior improvement, not feature change)

#### 4.3 Update OpenSpec Status ⏳
- **Action:** Mark this change as implemented
- **Status:** ⏳ To Do
- **Update:** OpenSpec metadata when merged

### 5. Review and Submission Tasks ⏳

#### 5.1 Self Review ⏳
- **Status:** ⏳ To Do
- **Checklist:**
  - [ ] Code follows project conventions
  - [ ] Function placement correct (at end of file)
  - [ ] Imports properly organized
  - [ ] Error handling is robust
  - [ ] Comments are clear and accurate
  - [ ] No debugging code left behind

#### 5.2 Create Pull Request ⏳
- **Status:** ⏳ To Do
- **PR Title:** `fix(tke): Suppress false-positive diffs for kubernetes_addon_config raw_values field`
- **PR Description Template:**
  ```markdown
  ## Description
  Fixes false-positive diffs in `tencentcloud_kubernetes_addon_config` resource when the TKE API returns JSON with elements in a different order than user input.
  
  ## Changes
  - Added `suppressJSONOrderDiff` function to perform semantic JSON comparison
  - Updated `raw_values` schema to use the diff suppression function
  - Added unit tests for the diff suppression logic
  
  ## Motivation
  The TKE API may return JSON responses with keys in a different order than the user's input. While semantically identical, Terraform's default string comparison treats this as a change, causing unnecessary diffs and update prompts.
  
  ## Testing
  - [x] Unit tests added
  - [x] Acceptance tests pass
  - [x] Manual testing completed
  - [x] Linter passes
  
  ## Backward Compatibility
  Fully backward compatible - only affects diff detection logic, not actual resource behavior.
  
  ## Related Issues
  Resolves: #<issue_number> (if applicable)
  
  ## OpenSpec
  See: `openspec/changes/fix-tke-addon-config-raw-values-json-diff/`
  ```

#### 5.3 Code Review ⏳
- **Status:** ⏳ To Do
- **Action:** Address review comments
- **Reviewers:** TBD

#### 5.4 Merge ⏳
- **Status:** ⏳ To Do
- **Prerequisites:**
  - All tests passing
  - Approvals received
  - No conflicts

## Validation Checklist

### Pre-Merge Validation
- [ ] All unit tests pass
- [ ] All acceptance tests pass
- [ ] Linter checks pass (no new issues)
- [ ] Build succeeds
- [ ] Manual testing completed successfully
- [ ] Documentation reviewed
- [ ] Changelog entry added
- [ ] Code review approved

### Post-Merge Validation
- [ ] CI/CD pipeline passes
- [ ] No regression reports
- [ ] OpenSpec status updated

## Timeline

| Phase                  | Estimated Time | Status      |
|------------------------|----------------|-------------|
| Code Implementation    | 0.5 days       | ✅ Complete |
| Unit Testing           | 0.5 days       | ⏳ Pending  |
| Integration Testing    | 0.5 days       | ⏳ Pending  |
| Documentation          | 0.25 days      | ⏳ Pending  |
| Code Review            | 0.5 days       | ⏳ Pending  |
| **Total**              | **2.25 days**  | **In Progress** |

## Dependencies

- No external dependencies
- No SDK updates required
- No breaking changes

## Rollback Plan

If critical issues are discovered:

1. Revert PR
2. Remove `DiffSuppressFunc` from schema
3. Remove `suppressJSONOrderDiff` function
4. Remove added imports if not used elsewhere
5. Issue hotfix release

**Risk Level:** Low (isolated change, easy to revert)

## Notes

- The code implementation is complete and formatted
- Pre-existing linter warning (deprecated `resource.Retry`) is unrelated to this change
- Focus now shifts to comprehensive testing
- Consider adding integration test that explicitly verifies JSON ordering independence
