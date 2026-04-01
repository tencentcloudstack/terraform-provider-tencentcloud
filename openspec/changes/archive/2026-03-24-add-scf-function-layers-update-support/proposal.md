# Proposal: Add Layers Update Support to tencentcloud_scf_function

## Metadata

- **Proposal ID**: add-scf-function-layers-update-support
- **Author**: CodeBuddy AI Assistant
- **Date**: 2026-03-24
- **Status**: 📋 Proposed
- **Target Resource**: `tencentcloud_scf_function`
- **Change Type**: Enhancement - Add update support for existing field

---

## Problem Statement

### Current Situation

The `tencentcloud_scf_function` resource has a `layers` field defined in its schema:

```go
"layers": {
    Type:        schema.TypeList,
    Optional:    true,
    Description: "The list of association layers.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "layer_name": {...},
            "layer_version": {...},
        },
    },
},
```

**Issues**:
1. ❌ The `layers` field is read during **Create** operation (line 614-624)
2. ❌ The `layers` field is **NOT** handled in the **Update** function
3. ❌ Users cannot modify layers after the function is created
4. ❌ Any changes to layers require resource recreation (destroy + create)

### User Impact

**Scenario**: A user wants to add, remove, or update a layer for an existing SCF function.

**Current Behavior** ❌:
```bash
# User modifies layers in Terraform config
resource "tencentcloud_scf_function" "example" {
  name    = "my-function"
  runtime = "Python3.6"
  
  layers {
    layer_name    = "my-layer"
    layer_version = 2  # Changed from 1 to 2
  }
}

$ terraform apply
# Terraform detects change but doesn't apply it
# Or shows drift on next plan
```

**Expected Behavior** ✅:
```bash
$ terraform apply
# Terraform calls UpdateFunctionConfiguration API
# Layers are updated without recreating the function
```

### Why This Matters

1. **Operational Efficiency**: Users should be able to update layers without function downtime
2. **API Capability**: Tencent Cloud SCF API **already supports** updating layers via `UpdateFunctionConfiguration`
3. **Consistency**: Other mutable fields (e.g., `environment`, `vpc_id`) already support updates
4. **Best Practice**: Terraform resources should support updates for mutable cloud resource properties

---

## Proposed Solution

### Overview

Add update support for the `layers` field in `resourceTencentCloudScfFunctionUpdate()` function, following the same pattern as other updateable fields like `l5_enable`.

### Implementation Plan

#### 1. Resource File Changes

**File**: `tencentcloud/services/scf/resource_tc_scf_function.go`

**Location**: After line 1279 (after `l5_enable` block)

**Code to Add**:
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

#### 2. Service File Changes

**File**: `tencentcloud/services/scf/service_tencentcloud_scf.go`

**Function**: `ModifyFunctionConfig()`

**Location**: After line 316 (after `l5Enable` block, before `DnsCache`)

**Code to Add**:
```go
if info.layers != nil {
    request.Layers = info.layers
}
```

### Why This Approach?

1. ✅ **Consistent Pattern**: Follows the exact same pattern as existing updateable fields
2. ✅ **Leverages Existing Code**: Reuses layer parsing logic from Create function
3. ✅ **API Compatible**: Uses existing `UpdateFunctionConfiguration` API call
4. ✅ **Minimal Change**: Only adds ~20 lines of code
5. ✅ **Backward Compatible**: Doesn't affect existing resources or behavior

---

## Technical Analysis

### API Verification

The Tencent Cloud SCF `UpdateFunctionConfiguration` API **already supports** the `Layers` parameter:

```go
type UpdateFunctionConfigurationRequest struct {
    FunctionName *string                  `json:"FunctionName"`
    // ... other fields ...
    Layers       []*LayerVersionSimple   `json:"Layers"`  // ✅ Supported
    // ... other fields ...
}
```

**Reference**: `tencentcloud-sdk-go/tencentcloud/scf/v20180416/models.go`

### Data Flow

**Create Flow** (already working):
```
User Config → d.GetOk("layers") → Parse to LayerVersionSimple[] → 
    → functionInfo.layers → CreateFunctionRequest.Layers → API
```

**Update Flow** (to be implemented):
```
User Config Change → d.HasChange("layers") → d.GetOk("layers") → 
    → Parse to LayerVersionSimple[] → functionInfo.layers → 
    → UpdateFunctionConfigurationRequest.Layers → API
```

### Edge Cases Handled

1. **Add Layers**: User adds new layers → Parsed and sent to API ✅
2. **Remove Layers**: User removes layers block → Send empty array to API ✅
3. **Update Layer Version**: User changes version → New values sent to API ✅
4. **Multiple Layers**: User has multiple layers → All parsed correctly ✅
5. **No Change**: User doesn't modify layers → `HasChange` returns false, no API call ✅

---

## Impact Analysis

### Benefits

| Benefit | Impact |
|---------|--------|
| **User Experience** | 🟢 High - Users can now update layers without recreation |
| **Operational** | 🟢 High - Reduces function downtime and deployment complexity |
| **Code Quality** | 🟢 Medium - Follows existing patterns, increases consistency |
| **API Utilization** | 🟢 Medium - Uses API capabilities that were previously unused |

### Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| **Breaking Changes** | 🟢 None | Purely additive feature, doesn't change existing behavior |
| **State Drift** | 🟢 Low | Read function already handles layers correctly |
| **API Errors** | 🟢 Low | Using well-tested API endpoint with existing error handling |
| **Backward Compatibility** | 🟢 None | Existing resources continue to work unchanged |

**Overall Risk**: 🟢 **Very Low**

---

## Testing Strategy

### Manual Testing Checklist

- [ ] **Test 1**: Add layers to existing function
- [ ] **Test 2**: Remove layers from existing function
- [ ] **Test 3**: Update layer version
- [ ] **Test 4**: Change layer name
- [ ] **Test 5**: Multiple layers management
- [ ] **Test 6**: Verify no drift after apply
- [ ] **Test 7**: Test with other concurrent updates (e.g., environment)

### Test Scenarios

#### Scenario 1: Add Layers
```hcl
# Initial state: no layers
resource "tencentcloud_scf_function" "test" {
  name    = "test-function"
  runtime = "Python3.6"
  handler = "index.main"
}

# Add layers
resource "tencentcloud_scf_function" "test" {
  name    = "test-function"
  runtime = "Python3.6"
  handler = "index.main"
  
  layers {
    layer_name    = "my-layer"
    layer_version = 1
  }
}
```

**Expected**: `terraform apply` updates layers without recreation

#### Scenario 2: Update Layer Version
```hcl
# Change version 1 → 2
layers {
  layer_name    = "my-layer"
  layer_version = 2  # Changed from 1
}
```

**Expected**: `terraform apply` updates to new version

#### Scenario 3: Remove Layers
```hcl
# Remove entire layers block
resource "tencentcloud_scf_function" "test" {
  name    = "test-function"
  runtime = "Python3.6"
  handler = "index.main"
  # layers block removed
}
```

**Expected**: `terraform apply` clears all layers

---

## Implementation Checklist

### Code Changes

- [ ] Add `HasChange("layers")` check in Update function
- [ ] Add layer parsing logic in Update function
- [ ] Add `request.Layers` assignment in `ModifyFunctionConfig()`
- [ ] Run `go fmt` on modified files

### Testing

- [ ] Manual testing with real SCF function
- [ ] Verify all edge cases
- [ ] Test backward compatibility
- [ ] Verify no state drift

### Documentation

- [ ] Update resource documentation (if needed)
- [ ] Add example for updating layers (optional)

---

## Timeline

| Phase | Duration | Description |
|-------|----------|-------------|
| **Code Implementation** | 10 minutes | Add code in resource and service files |
| **Code Formatting** | 2 minutes | Run `go fmt` |
| **Manual Testing** | 20 minutes | Test all scenarios |
| **Documentation** | 5 minutes | Update docs if needed |
| **Total** | **~40 minutes** | End-to-end implementation |

---

## Alternatives Considered

### Alternative 1: Do Nothing
**Pros**: No development effort  
**Cons**: Poor user experience, requires resource recreation  
**Decision**: ❌ Rejected - User experience is important

### Alternative 2: Mark as ForceNew
**Pros**: Makes behavior explicit  
**Cons**: Forces recreation, API supports updates  
**Decision**: ❌ Rejected - API already supports updates

### Alternative 3: Current Proposal (Add Update Support)
**Pros**: Best user experience, uses API capabilities, follows existing patterns  
**Cons**: Requires code changes (minimal)  
**Decision**: ✅ **Selected** - Best overall solution

---

## Success Criteria

1. ✅ Users can add layers to existing functions
2. ✅ Users can remove layers from existing functions
3. ✅ Users can update layer versions
4. ✅ No unexpected diffs after apply
5. ✅ Backward compatible with existing configurations
6. ✅ Code follows project conventions
7. ✅ All edge cases handled correctly

---

## References

- **Resource File**: `tencentcloud/services/scf/resource_tc_scf_function.go`
- **Service File**: `tencentcloud/services/scf/service_tencentcloud_scf.go`
- **SDK Models**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416/models.go`
- **Similar Pattern**: See `l5_enable` field update logic (lines 1276-1279)

---

## Approval

**Proposal Status**: 📋 Awaiting Approval

**Next Steps**:
1. Review this proposal
2. Approve for implementation
3. Execute code changes
4. Test and validate
5. Archive as spec

---

## Notes

- This change is purely additive and has no breaking changes
- The API already supports this feature - we're just exposing it
- Implementation follows established patterns in the codebase
- Very low risk, high user value
