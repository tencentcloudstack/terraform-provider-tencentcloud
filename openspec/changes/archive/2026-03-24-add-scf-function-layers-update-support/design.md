# Technical Design: Add Layers Update Support to tencentcloud_scf_function

## Overview

This document provides detailed technical design for adding update support for the `layers` field in the `tencentcloud_scf_function` resource.

---

## Architecture

### Current Architecture (Create Only)

```
┌─────────────────────────────────────────────────────────────────┐
│ Terraform User Config                                           │
│                                                                  │
│ resource "tencentcloud_scf_function" "example" {                │
│   layers {                                                       │
│     layer_name    = "my-layer"                                   │
│     layer_version = 1                                            │
│   }                                                              │
│ }                                                                │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│ resourceTencentCloudScfFunctionCreate()                         │
│   Lines 614-624: Parse layers from d.GetOk("layers")           │
│   → functionInfo.layers = []*LayerVersionSimple{...}           │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│ ScfService.CreateFunction()                                     │
│   Line 103: request.Layers = info.layers                       │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│ Tencent Cloud API: CreateFunction                               │
│   Layers field is sent and stored                               │
└─────────────────────────────────────────────────────────────────┘
```

### Target Architecture (Create + Update)

```
┌─────────────────────────────────────────────────────────────────┐
│ Terraform User Updates Config                                   │
│                                                                  │
│ resource "tencentcloud_scf_function" "example" {                │
│   layers {                                                       │
│     layer_name    = "my-layer"                                   │
│     layer_version = 2  # Updated from 1 to 2                    │
│   }                                                              │
│ }                                                                │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│ resourceTencentCloudScfFunctionUpdate()                         │
│   ✨ NEW: d.HasChange("layers") check                           │
│   ✨ NEW: Parse layers from d.GetOk("layers")                   │
│   ✨ NEW: functionInfo.layers = []*LayerVersionSimple{...}      │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│ ScfService.ModifyFunctionConfig()                               │
│   ✨ NEW: request.Layers = info.layers                          │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│ Tencent Cloud API: UpdateFunctionConfiguration                  │
│   Layers field is updated (already supported by API)            │
└─────────────────────────────────────────────────────────────────┘
```

---

## Detailed Code Changes

### Change 1: Resource Update Function

**File**: `tencentcloud/services/scf/resource_tc_scf_function.go`

**Location**: After line 1279 (immediately after `l5_enable` block)

**Current Code** (lines 1276-1279):
```go
if d.HasChange("l5_enable") {
    updateAttrs = append(updateAttrs, "l5_enable")
    functionInfo.l5Enable = helper.Bool(d.Get("l5_enable").(bool))
}
```

**New Code to Insert**:
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
        // Clear all layers if the block is removed from config
        functionInfo.layers = []*scf.LayerVersionSimple{}
    }
}
```

**Explanation**:
1. `d.HasChange("layers")` - Detects if user changed the layers configuration
2. `updateAttrs = append(...)` - Marks this attribute for update (used for logging)
3. `d.GetOk("layers")` - Gets the new layers value
4. Loop through layers and parse each one into `LayerVersionSimple` struct
5. If layers block is removed entirely, set to empty array to clear layers

**Pattern Match**: This exactly mirrors the parsing logic in Create function (lines 614-624)

---

### Change 2: Service Update Function

**File**: `tencentcloud/services/scf/service_tencentcloud_scf.go`

**Function**: `ModifyFunctionConfig()`

**Location**: After line 316 (after `L5Enable` block, before `DnsCache`)

**Current Code** (lines 311-319):
```go
if info.l5Enable != nil {
    request.L5Enable = helper.String("FALSE")
    if *info.l5Enable {
        request.L5Enable = helper.String("TRUE")
    }
}

request.DnsCache = info.dnsCache
request.IntranetConfig = info.intranetConfig
```

**New Code to Insert** (between L5Enable and DnsCache):
```go
if info.layers != nil {
    request.Layers = info.layers
}
```

**Explanation**:
1. Check if `layers` field is set in `functionInfo`
2. If set, assign to request (includes empty array case for clearing)
3. API will update layers accordingly

**Why Simple?**: The `layers` field in `functionInfo` is already of type `[]*scf.LayerVersionSimple`, which matches the API request type exactly. No conversion needed.

---

## Data Structures

### LayerVersionSimple Struct

From SDK: `tencentcloud-sdk-go/tencentcloud/scf/v20180416/models.go`

```go
type LayerVersionSimple struct {
    LayerName    *string `json:"LayerName,omitempty" name:"LayerName"`
    LayerVersion *int64  `json:"LayerVersion,omitempty" name:"LayerVersion"`
}
```

### scfFunctionInfo Struct

From: `tencentcloud/services/scf/service_tencentcloud_scf.go` (line 36)

```go
type scfFunctionInfo struct {
    // ... other fields ...
    layers []*scf.LayerVersionSimple  // ✅ Already exists
    // ... other fields ...
}
```

**Note**: The `layers` field already exists in `scfFunctionInfo`, so no struct changes needed.

---

## State Management

### Terraform State Flow

**Create**:
```
User Config → Create Function → API → Response → Read Function → State
```

**Update (After This Change)**:
```
State → Detect Change → Update Function → API → Response → Read Function → Updated State
```

**Read** (Already Correct):
```go
// Lines 928-937 in resource_tc_scf_function.go
if len(function.Layers) > 0 {
    layers := make([]map[string]interface{}, 0, len(function.Layers))
    for _, v := range function.Layers {
        layers = append(layers, map[string]interface{}{
            "layer_name":    v.LayerName,
            "layer_version": v.LayerVersion,
        })
    }
    _ = d.Set("layers", layers)
}
```

The Read function already correctly populates `layers` from API response, so state will always be in sync after update.

---

## Edge Cases and Handling

### Case 1: Add Layers to Function Without Layers

**Before**:
```hcl
resource "tencentcloud_scf_function" "test" {
  name = "test"
  # No layers block
}
```

**After**:
```hcl
resource "tencentcloud_scf_function" "test" {
  name = "test"
  layers {
    layer_name    = "layer1"
    layer_version = 1
  }
}
```

**Handling**:
- `d.HasChange("layers")` → `true`
- `d.GetOk("layers")` → Returns the new layer
- Parsed and sent to API ✅

---

### Case 2: Remove All Layers

**Before**:
```hcl
resource "tencentcloud_scf_function" "test" {
  name = "test"
  layers {
    layer_name    = "layer1"
    layer_version = 1
  }
}
```

**After**:
```hcl
resource "tencentcloud_scf_function" "test" {
  name = "test"
  # layers block removed
}
```

**Handling**:
- `d.HasChange("layers")` → `true`
- `d.GetOk("layers")` → `ok = false`
- Set `functionInfo.layers = []*scf.LayerVersionSimple{}` (empty array)
- API receives empty array and clears all layers ✅

---

### Case 3: Update Layer Version

**Before**:
```hcl
layers {
  layer_name    = "layer1"
  layer_version = 1
}
```

**After**:
```hcl
layers {
  layer_name    = "layer1"
  layer_version = 2  # Changed
}
```

**Handling**:
- `d.HasChange("layers")` → `true`
- Parse new version (2)
- API updates to new version ✅

---

### Case 4: Multiple Layers

**Config**:
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

**Handling**:
- Loop processes all layers
- All added to `layers` slice
- API receives array with both layers ✅

---

### Case 5: No Change to Layers

**Config**: Layers unchanged

**Handling**:
- `d.HasChange("layers")` → `false`
- Block skipped entirely
- No API call for layers ✅

---

## API Compatibility

### UpdateFunctionConfiguration API

**Request Fields** (relevant excerpt):
```go
type UpdateFunctionConfigurationRequest struct {
    FunctionName *string `json:"FunctionName,omitempty" name:"FunctionName"`
    
    // ✅ Layers is supported
    Layers []*LayerVersionSimple `json:"Layers,omitempty" name:"Layers"`
    
    // Other fields...
    MemorySize   *int64  `json:"MemorySize,omitempty" name:"MemorySize"`
    Timeout      *int64  `json:"Timeout,omitempty" name:"Timeout"`
    Environment  *Environment `json:"Environment,omitempty" name:"Environment"`
    // ...
}
```

**Verification**: ✅ The `Layers` field is present in the API request struct.

**Behavior**:
- If `Layers` is provided → Updates function layers
- If `Layers` is empty array → Clears all layers
- If `Layers` is nil → Doesn't change layers

---

## Error Handling

### Existing Error Handling (Leveraged)

```go
// In ModifyFunctionConfig() - lines 321-332
if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    ratelimit.Check(request.GetAction())
    
    if _, err := client.UpdateFunctionConfiguration(request); err != nil {
        return tccommon.RetryError(errors.WithStack(err), tccommon.InternalError)
    }
    return nil
}); err != nil {
    return err
}

return waitScfFunctionReady(ctx, info.name, *info.namespace, client)
```

**Error Scenarios**:
1. **Invalid Layer**: API returns error → Retry logic handles it → Error propagated to user ✅
2. **Layer Not Found**: API returns error → User sees clear error message ✅
3. **Network Error**: Retry logic handles transient failures ✅
4. **Function Not Ready**: `waitScfFunctionReady` polls until ready or timeout ✅

**No New Error Handling Needed**: Existing mechanisms cover all cases.

---

## Testing Plan

### Unit Testing (Manual)

Since this project doesn't appear to have extensive unit tests for resources, we'll rely on manual integration testing.

### Integration Testing Scenarios

#### Test 1: Basic Layer Addition
```bash
# Step 1: Create function without layers
# Step 2: Add layer to config
# Step 3: terraform apply
# Expected: Layer added, no recreation
```

#### Test 2: Layer Version Update
```bash
# Step 1: Create function with layer v1
# Step 2: Update to layer v2
# Step 3: terraform apply
# Expected: Layer version updated, no recreation
```

#### Test 3: Layer Removal
```bash
# Step 1: Create function with layers
# Step 2: Remove layers from config
# Step 3: terraform apply
# Expected: Layers cleared, no recreation
```

#### Test 4: Multiple Layers
```bash
# Step 1: Create function with 1 layer
# Step 2: Add second layer
# Step 3: terraform apply
# Expected: Both layers present
```

#### Test 5: No Drift After Apply
```bash
# Step 1: Apply config with layers
# Step 2: Run terraform plan
# Expected: "No changes" message
```

#### Test 6: Concurrent Updates
```bash
# Step 1: Update both layers AND environment
# Step 2: terraform apply
# Expected: Both updated successfully
```

---

## Performance Considerations

### API Calls

**Before**: Update operations → 1-2 API calls (depending on what changes)

**After**: Same - no additional API calls

**Reason**: Layers are included in the same `UpdateFunctionConfiguration` call used for other config updates.

### Parsing Performance

**Layer Parsing**:
- Time Complexity: O(n) where n = number of layers
- Typical n: 1-5 layers
- Impact: Negligible

---

## Backward Compatibility

### Existing Resources

**Scenario**: User has existing function with layers, upgrades provider

**Behavior**:
1. Existing state contains layers (from Read function)
2. Config matches state → No changes detected
3. Everything continues working ✅

**Scenario**: User has existing function, never specified layers in config

**Behavior**:
1. State may or may not have layers (depending on how function was created)
2. Read function populates current layers from API
3. If config doesn't specify layers, Terraform uses current state
4. No unexpected changes ✅

### Schema Compatibility

**No Schema Changes**: The `layers` field already exists in schema with correct definition. This change only adds update logic.

**Result**: ✅ 100% backward compatible

---

## Code Quality

### Following Project Conventions

1. **Pattern Consistency**: ✅ Mirrors `l5_enable` update pattern
2. **Code Location**: ✅ Placed in logical order after `l5_enable`
3. **Error Handling**: ✅ Uses existing error handling mechanisms
4. **Helper Functions**: ✅ Uses existing helpers (`helper.String`, `helper.IntInt64`)
5. **Logging**: ✅ Integrates with existing `updateAttrs` tracking

### Code Formatting

**Required**: Run `go fmt` after changes

```bash
go fmt tencentcloud/services/scf/resource_tc_scf_function.go
go fmt tencentcloud/services/scf/service_tencentcloud_scf.go
```

---

## Security Considerations

### Input Validation

**Layer Name**:
- Validated by API
- No special handling needed in provider

**Layer Version**:
- Type: `int` in schema → Terraform validates type
- API validates version exists

**No Security Concerns**: This change doesn't introduce new attack vectors.

---

## Rollback Plan

If issues arise:

```bash
# 1. Identify the commit
git log --oneline tencentcloud/services/scf/resource_tc_scf_function.go

# 2. Revert the changes
git revert <commit-hash>

# 3. Rebuild
go build

# 4. Test
go test ./tencentcloud/services/scf/...
```

**Impact of Rollback**: Users lose ability to update layers, must recreate functions to change layers (previous behavior).

---

## Success Metrics

1. ✅ Code compiles without errors
2. ✅ `go fmt` runs successfully
3. ✅ Manual tests pass
4. ✅ No state drift after updates
5. ✅ Backward compatibility verified
6. ✅ API calls succeed

---

## Future Enhancements

### Potential Improvements (Out of Scope)

1. **Acceptance Tests**: Add automated acceptance tests for layer updates
2. **Validation**: Add client-side validation for layer names/versions
3. **Documentation**: Add examples to resource docs showing layer updates
4. **Drift Detection**: Enhance drift detection for layers (already works via Read)

---

## Summary

This design adds layers update support by:
1. Adding `HasChange("layers")` check in Update function (~15 lines)
2. Adding `request.Layers` assignment in service function (~3 lines)
3. Leveraging existing API support
4. Following established patterns

**Total New Code**: ~20 lines  
**Risk Level**: Very Low  
**User Value**: High  
**Complexity**: Low

This is a straightforward enhancement that fills a gap in the resource's capabilities.
