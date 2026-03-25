# Implementation Tasks: Add Import Support to tencentcloud_kubernetes_addon_config

## Task Breakdown

### Phase 1: Code Implementation ✅
**Status**: Completed  
**Estimated Time**: 5 minutes

- [x] **Task 1.1**: Add Importer configuration to resource schema
  - **File**: `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
  - **Location**: Lines 19-24 (after Delete, before Schema)
  - **Change**:
    ```go
    Importer: &schema.ResourceImporter{
        State: schema.ImportStatePassthrough,
    },
    ```
  - **Validation**: ✅ Code compiles without errors

- [x] **Task 1.2**: Run `go fmt` to format the code
  - **Command**: `go fmt ./tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
  - **Validation**: ✅ Formatting completed

---

### Phase 2: Testing ⏳
**Status**: Pending  
**Estimated Time**: 20 minutes

#### 2.1 Acceptance Test ⏳

- [ ] **Task 2.1.1**: Create acceptance test function
  - **File**: `tencentcloud/services/tke/resource_tc_kubernetes_addon_config_test.go`
  - **Function**: `TestAccTencentCloudKubernetesAddonConfig_import`
  - **Test Steps**:
    1. Create addon config via Terraform
    2. Import the same addon config
    3. Verify all fields match (ImportStateVerify: true)

- [ ] **Task 2.1.2**: Run acceptance test
  - **Command**: 
    ```bash
    TF_ACC=1 go test -v ./tencentcloud/services/tke -run TestAccTencentCloudKubernetesAddonConfig_import -timeout 30m
    ```
  - **Success Criteria**: Test passes, no errors

#### 2.2 Manual Testing ⏳

- [ ] **Task 2.2.1**: Test valid import
  - **Setup**: Ensure a TKE cluster with addon exists
  - **Command**: `terraform import tencentcloud_kubernetes_addon_config.test <cluster_id>#<addon_name>`
  - **Validation**: Import succeeds, state populated correctly

- [ ] **Task 2.2.2**: Test plan after import
  - **Command**: `terraform plan`
  - **Validation**: No changes shown (or only expected differences in raw_values formatting)

- [ ] **Task 2.2.3**: Test invalid ID format
  - **Command**: `terraform import tencentcloud_kubernetes_addon_config.test invalid-id`
  - **Validation**: Error message: "id is broken,invalid-id"

- [ ] **Task 2.2.4**: Test non-existent resource
  - **Command**: `terraform import tencentcloud_kubernetes_addon_config.test cls-fake#fake`
  - **Validation**: Error message indicates resource not found

- [ ] **Task 2.2.5**: Test update after import
  - **Action**: Modify `addon_version` or `raw_values` in config
  - **Command**: `terraform apply`
  - **Validation**: Update succeeds, no unexpected errors

---

### Phase 3: Documentation ⏳
**Status**: Pending  
**Estimated Time**: 10 minutes

- [ ] **Task 3.1**: Update resource documentation
  - **File**: `website/docs/r/kubernetes_addon_config.html.markdown`
  - **Section to Add**: "## Import"
  - **Content**:
    ```markdown
    ## Import

    Kubernetes addon configuration can be imported using the id, e.g.

    ```bash
    terraform import tencentcloud_kubernetes_addon_config.example cls-abc123#tcr
    ```

    Where:
    - `cls-abc123` is the cluster ID
    - `tcr` is the addon name
    ```
  - **Validation**: Documentation builds without errors

- [ ] **Task 3.2**: Add import example
  - **File**: `examples/tencentcloud-kubernetes-addon-config/main.tf` (or create if not exists)
  - **Content**: Example showing import usage
  - **Validation**: Example is clear and executable

- [ ] **Task 3.3**: Update CHANGELOG
  - **File**: `CHANGELOG.md`
  - **Entry**:
    ```markdown
    ENHANCEMENTS:
    * **resource/tencentcloud_kubernetes_addon_config**: Add import support
    ```

---

### Phase 4: Code Review & Quality Checks ⏳
**Status**: Pending  
**Estimated Time**: 10 minutes

- [ ] **Task 4.1**: Run linter
  - **Command**: `golangci-lint run ./tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
  - **Validation**: No new linting errors

- [ ] **Task 4.2**: Verify no breaking changes
  - **Check**: Existing tests still pass
  - **Command**: 
    ```bash
    TF_ACC=1 go test -v ./tencentcloud/services/tke -run TestAccTencentCloudKubernetesAddonConfig -timeout 30m
    ```
  - **Validation**: All tests pass

- [ ] **Task 4.3**: Code review checklist
  - [ ] Follows existing code patterns
  - [ ] No hardcoded values
  - [ ] Error handling is appropriate
  - [ ] Documentation is complete
  - [ ] Tests cover main scenarios

---

### Phase 5: Finalization ⏳
**Status**: Pending  
**Estimated Time**: 5 minutes

- [ ] **Task 5.1**: Commit changes
  - **Branch**: `feature/tke-addon-config-import`
  - **Commit Message**: 
    ```
    feat(tke): add import support to tencentcloud_kubernetes_addon_config
    
    - Add Importer configuration using ImportStatePassthrough
    - Add acceptance test for import functionality
    - Update documentation with import examples
    
    Closes #XXXX
    ```

- [ ] **Task 5.2**: Create pull request
  - **Title**: "feat(tke): add import support to tencentcloud_kubernetes_addon_config"
  - **Description**: Link to this OpenSpec proposal
  - **Labels**: `enhancement`, `tke`, `import`

- [ ] **Task 5.3**: Address review feedback
  - **Action**: Respond to reviewer comments
  - **Action**: Make necessary adjustments

---

## Task Dependencies

```
Phase 1 (Code Implementation)
    ↓
Phase 2 (Testing)
    ├─ 2.1 (Acceptance Test)
    └─ 2.2 (Manual Testing)
    ↓
Phase 3 (Documentation)
    ↓
Phase 4 (Code Review)
    ↓
Phase 5 (Finalization)
```

## Estimated Total Time

| Phase | Time |
|-------|------|
| Phase 1: Code Implementation | 5 min |
| Phase 2: Testing | 20 min |
| Phase 3: Documentation | 10 min |
| Phase 4: Code Review | 10 min |
| Phase 5: Finalization | 5 min |
| **Total** | **50 min** |

## Current Status Summary

| Category | Status | Progress |
|----------|--------|----------|
| Code Implementation | ✅ Completed | 2/2 tasks |
| Testing | ⏳ Pending | 0/7 tasks |
| Documentation | ⏳ Pending | 0/3 tasks |
| Code Review | ⏳ Pending | 0/3 tasks |
| Finalization | ⏳ Pending | 0/3 tasks |
| **Overall** | **🔄 In Progress** | **2/18 tasks (11%)** |

---

## Quick Start Checklist

For immediate implementation, follow these minimal steps:

1. ✅ **Proposal Created** - Review and approve this proposal
2. ⏳ **Code Change** - Add 3 lines to resource definition
3. ⏳ **Format Code** - Run `go fmt`
4. ⏳ **Manual Test** - Import one addon config
5. ⏳ **Verify** - Run `terraform plan` to ensure no diff
6. ⏳ **Document** - Add import section to docs
7. ⏳ **Commit** - Push changes for review

---

## Notes

- **Low Risk**: This is a purely additive feature with no breaking changes
- **Standard Pattern**: Uses the same import pattern as other resources in the provider
- **Already Tested**: The underlying Read function is already battle-tested in production
- **Quick Win**: High user value with minimal implementation effort

## Blockers

None identified. All dependencies (Terraform SDK, existing Read function, API access) are already in place.

## Related Work

- [x] Prior change: `fix-tke-addon-config-raw-values-json-diff` - Added JSON diff suppression
  - This ensures imported `raw_values` won't show spurious diffs due to JSON ordering

## Rollback Plan

If issues arise:
1. Revert the 3-line change
2. Import functionality becomes unavailable
3. No impact on existing resources or state
