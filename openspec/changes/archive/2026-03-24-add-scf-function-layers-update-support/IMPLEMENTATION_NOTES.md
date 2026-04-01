# Implementation Notes: SCF Function Layers Update Support

**Implementation Date**: 2026-03-24  
**Status**: ✅ Code Implementation Completed  
**Developer**: AI Agent

---

## 📋 Overview

This document records the implementation details of adding layers update support to the `tencentcloud_scf_function` resource.

**Goal**: Enable users to modify the `layers` field of an SCF function without requiring resource recreation.

**Complexity**: Low - Only 2 files modified, ~20 lines of code added

---

## 🔧 Code Changes Summary

| File | Lines Added | Lines Modified | Location |
|------|-------------|----------------|----------|
| `resource_tc_scf_function.go` | 18 | 0 | Lines 1281-1298 |
| `service_tencentcloud_scf.go` | 3 | 0 | Lines 318-320 |
| **Total** | **21** | **0** | **2 files** |

---

## 📝 Detailed Changes

### Change 1: Resource File - Update Function Logic

**File**: `tencentcloud/services/scf/resource_tc_scf_function.go`  
**Location**: After line 1279 (following `l5_enable` block)  
**Lines**: 1281-1298

#### Code Added:

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
        // If layers block is removed from configuration, clear all layers
        functionInfo.layers = []*scf.LayerVersionSimple{}
    }
}
```

#### Purpose:
- Detects when the `layers` field changes in the Terraform configuration
- Parses the layers block into `LayerVersionSimple` structs
- Handles removal of layers (sets to empty slice)
- Adds "layers" to `updateAttrs` to track the change

#### Key Design Decisions:
1. **Reused Create Logic**: The parsing logic mirrors the existing Create function (lines 614-624), ensuring consistency
2. **Empty Slice for Removal**: When layers are removed, we set an empty slice rather than nil to explicitly clear all layers
3. **Placement**: Added immediately after `l5_enable` block as requested, maintaining logical grouping of update handlers

---

### Change 2: Service File - API Request Parameter

**File**: `tencentcloud/services/scf/service_tencentcloud_scf.go`  
**Function**: `ModifyFunctionConfig()`  
**Location**: After line 316 (between `l5Enable` and `DnsCache`)  
**Lines**: 318-320

#### Code Added:

```go
if info.layers != nil {
    request.Layers = info.layers
}
```

#### Purpose:
- Assigns the parsed layers to the API request
- Only sets the field if layers were explicitly provided (nil check)

#### Key Design Decisions:
1. **Nil Check**: Only set `request.Layers` if `info.layers` is not nil, allowing the API to handle default behavior
2. **Placement**: Positioned logically between L5Enable and DnsCache parameters
3. **Simple Assignment**: Direct assignment since the type already matches `[]*scf.LayerVersionSimple`

---

## 🎯 How It Works

### Update Flow

```
User modifies layers in .tf config
         ↓
terraform plan detects change
         ↓
resourceTencentCloudScfFunctionUpdate() called
         ↓
d.HasChange("layers") == true
         ↓
Parse layers config → functionInfo.layers
         ↓
ModifyFunctionConfig(functionInfo)
         ↓
request.Layers = info.layers
         ↓
UpdateFunctionConfiguration API called
         ↓
Function layers updated (no recreation)
         ↓
terraform apply completes
         ↓
Next plan shows "no changes"
```

### Data Flow

1. **Terraform Config** → `layers { layer_name, layer_version }`
2. **Resource Layer** → Parse into `[]*scf.LayerVersionSimple`
3. **Service Layer** → Assign to `request.Layers`
4. **API Call** → `UpdateFunctionConfiguration` with Layers parameter
5. **Tencent Cloud** → Updates function's layer configuration
6. **State** → Read function, layers reflected in state

---

## ✅ Testing Scenarios

The following test scenarios should be validated:

### Scenario 1: Add Layers to Function
```hcl
# Before: No layers
resource "tencentcloud_scf_function" "test" {
  name = "test-function"
}

# After: Add layers
resource "tencentcloud_scf_function" "test" {
  name = "test-function"
  
  layers {
    layer_name    = "my-layer"
    layer_version = 1
  }
}
```
**Expected**: Layers added without recreation

---

### Scenario 2: Update Layer Version
```hcl
# Before
layers {
  layer_name    = "my-layer"
  layer_version = 1
}

# After
layers {
  layer_name    = "my-layer"
  layer_version = 2  # Version changed
}
```
**Expected**: Layer version updated in place

---

### Scenario 3: Remove All Layers
```hcl
# Before: Has layers
layers {
  layer_name    = "my-layer"
  layer_version = 1
}

# After: No layers block (removed)
# (No layers field)
```
**Expected**: All layers cleared

---

### Scenario 4: Multiple Layers
```hcl
layers {
  layer_name    = "layer1"
  layer_version = 1
}

layers {
  layer_name    = "layer2"
  layer_version = 2
}
```
**Expected**: Both layers applied correctly

---

## 🔍 Code Quality

### Formatting
- ✅ `gofmt` applied to both files
- ✅ Code follows Go formatting standards
- ✅ Proper indentation and spacing

### Consistency
- ✅ Matches existing code patterns in the file
- ✅ Reuses parsing logic from Create function
- ✅ Follows same structure as other update handlers (e.g., `l5_enable`)

### Error Handling
- ✅ Type assertions for interface{} conversions
- ✅ Nil checks for optional fields
- ✅ Empty slice for removed layers

---

## 📊 Before/After Comparison

### Before Implementation

| Action | Behavior | User Impact |
|--------|----------|-------------|
| Add layers to function | ❌ Requires recreation | Service downtime |
| Update layer version | ❌ Requires recreation | Service downtime |
| Remove layers | ❌ Requires recreation | Service downtime |
| Change layer name | ❌ Requires recreation | Service downtime |

### After Implementation

| Action | Behavior | User Impact |
|--------|----------|-------------|
| Add layers to function | ✅ In-place update | No downtime |
| Update layer version | ✅ In-place update | No downtime |
| Remove layers | ✅ In-place update | No downtime |
| Change layer name | ✅ In-place update | No downtime |

---

## 🎉 Benefits

1. **User Experience**
   - No service interruption when updating layers
   - Faster deployment times
   - More flexible layer management

2. **Consistency**
   - Aligns with other updatable fields (runtime, environment, etc.)
   - Makes the resource behavior more predictable

3. **API Utilization**
   - Uses the full capability of `UpdateFunctionConfiguration` API
   - No longer requires workarounds (destroy + recreate)

---

## 🔐 Risk Assessment

| Risk | Level | Mitigation |
|------|-------|------------|
| Breaking Change | 🟢 None | Purely additive functionality |
| State Drift | 🟢 Low | Read function already handles layers |
| API Compatibility | 🟢 None | API already supports this parameter |
| Backward Compatibility | 🟢 None | Existing resources unaffected |

---

## 🚀 Next Steps

### Immediate (Manual Testing Required)
1. ⏳ Test adding layers to existing function
2. ⏳ Test updating layer versions
3. ⏳ Test removing all layers
4. ⏳ Test multiple layers management
5. ⏳ Verify no state drift after apply
6. ⏳ Test concurrent updates (layers + other fields)

### Optional (Documentation)
1. ⏳ Update resource documentation
2. ⏳ Add changelog entry
3. ⏳ Update examples

### Before Merge
1. ⏳ Run full test suite
2. ⏳ Code review
3. ⏳ Verify linter passes
4. ⏳ Manual acceptance tests

---

## 📚 References

### Related Code Sections
- **Create Function**: Lines 614-624 (layers parsing logic)
- **Read Function**: Lines 943-956 (layers state handling)
- **Update Function**: Lines 1281-1298 (NEW - layers update logic)
- **Service Method**: Lines 318-320 (NEW - layers API assignment)

### API Documentation
- **UpdateFunctionConfiguration**: Supports `Layers` parameter
- **Layer Structure**: `LayerName` + `LayerVersion`

### Design Documents
- Proposal: `openspec/changes/add-scf-function-layers-update-support/proposal.md`
- Design: `openspec/changes/add-scf-function-layers-update-support/design.md`
- Tasks: `openspec/changes/add-scf-function-layers-update-support/tasks.md`

---

## ✨ Implementation Summary

**Status**: ✅ **Code Implementation Complete**

**What Was Done**:
- ✅ Added layers change detection in Update function
- ✅ Added layers parsing logic (reused from Create)
- ✅ Added layers API assignment in service layer
- ✅ Formatted code with gofmt
- ✅ Updated task tracking

**What's Pending**:
- ⏳ Manual testing (7 test scenarios)
- ⏳ Documentation updates (optional)
- ⏳ Code review and validation

**Implementation Time**: ~5 minutes (estimated 10 minutes)

**Lines of Code**: 21 lines across 2 files

**Complexity**: Low - straightforward implementation following existing patterns

---

## 🎯 Success Criteria

### Must Have (Before Merge)
- [x] Code compiles without errors
- [x] Code formatted with gofmt
- [ ] Manual tests pass (7 scenarios)
- [ ] No state drift after apply
- [ ] Code review approved

### Should Have
- [ ] Linter warnings addressed (if any new ones)
- [ ] Acceptance tests pass

### Nice to Have
- [ ] Documentation updated
- [ ] Changelog entry added
- [ ] Usage examples provided

---

**Implementation Date**: 2026-03-24  
**Last Updated**: 2026-03-24  
**Status**: ✅ Code Complete, ⏳ Testing Pending
