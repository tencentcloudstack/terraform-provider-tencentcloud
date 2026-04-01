# Proposal: Fix JSON Order Diff Issue in TKE Addon Config raw_values Field

## Overview

Add a custom diff suppression function to the `raw_values` field in the `tencentcloud_kubernetes_addon_config` resource to prevent false-positive diffs caused by JSON element ordering differences between user input and API response.

## Motivation

The current `tencentcloud_kubernetes_addon_config` resource's `raw_values` field experiences unnecessary diff detection due to JSON element ordering:

- The `raw_values` field accepts and returns JSON strings
- The TKE API may return JSON with elements in a different order than the user's input
- While the JSON content is semantically identical, Terraform's default string comparison treats different ordering as a change
- This causes Terraform to report spurious diffs on every `plan`, prompting unnecessary resource updates
- Users are forced to run `apply` repeatedly with no actual changes

### Example Scenario

**User Input:**
```json
{"replicas": 2, "image": "nginx:latest", "port": 80}
```

**API Response:**
```json
{"image": "nginx:latest", "port": 80, "replicas": 2}
```

**Current Behavior:** Terraform shows a diff even though the content is identical.

**Desired Behavior:** Terraform should recognize these as equivalent and show no diff.

## Current State

### Existing Schema
```go
"raw_values": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    Description: "Params of addon, base64 encoded json format.",
},
```

### Problem
- No diff suppression logic
- Simple string comparison fails to recognize JSON semantic equality
- Causes poor user experience with constant false-positive diffs

## Proposed Solution

Add a custom `DiffSuppressFunc` that performs semantic JSON comparison instead of string comparison.

### Implementation Approach

1. **Add Required Imports**
   - `encoding/json` - for JSON parsing
   - `reflect` - for deep equality comparison

2. **Create Diff Suppression Function**
   ```go
   func suppressJSONOrderDiff(k, old, new string, d *schema.ResourceData) bool {
       // Handle empty strings
       if old == "" && new == "" {
           return true
       }
       if old == "" || new == "" {
           return false
       }

       // Parse both JSON strings
       var oldJSON, newJSON interface{}
       if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
           // Fallback to string comparison if parse fails
           return old == new
       }
       if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
           return old == new
       }

       // Compare parsed JSON objects (ignoring order)
       return reflect.DeepEqual(oldJSON, newJSON)
   }
   ```

3. **Update Schema**
   ```go
   "raw_values": {
       Type:             schema.TypeString,
       Optional:         true,
       Computed:         true,
       Description:      "Params of addon, base64 encoded json format.",
       DiffSuppressFunc: suppressJSONOrderDiff,
   },
   ```

4. **Function Placement**
   - Place the helper function at the end of the file, after all CRUD functions
   - Follows project convention for utility functions

### Why This Solution Works

- **Semantic Comparison**: `json.Unmarshal` + `reflect.DeepEqual` compares the actual data structure, not string representation
- **Order-Independent**: Go maps and the comparison ignore key ordering
- **Graceful Fallback**: If JSON parsing fails, falls back to string comparison (safe default)
- **Backward Compatible**: Only affects diff detection; doesn't change API calls or data storage
- **Terraform Standard Pattern**: Uses the built-in `DiffSuppressFunc` mechanism

## Impact Analysis

### Files Changed
- `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`

### Backward Compatibility
- ✅ **Fully Backward Compatible**: All existing configurations continue working
- ✅ **No State Migration Required**: Only affects diff logic, not stored state
- ✅ **No Breaking Changes**: Existing behavior preserved for non-JSON strings

### User Impact
- **Positive**: Eliminates false-positive diffs
- **Positive**: Reduces unnecessary `terraform apply` operations
- **Positive**: Improves user experience and confidence in Terraform state
- **No Negative Impact**: No changes to actual resource behavior

### Performance Impact
- **Minimal**: JSON parsing only occurs during plan/diff operations
- **Acceptable Overhead**: Parsing small JSON configs is fast (<1ms typically)

## Testing Strategy

### Unit Test Scenarios
1. **Empty strings**: Both empty → no diff
2. **One empty, one non-empty**: → diff detected
3. **Same JSON, different order**: `{"a":1,"b":2}` vs `{"b":2,"a":1}` → no diff
4. **Different JSON content**: `{"a":1}` vs `{"a":2}` → diff detected
5. **Invalid JSON**: Falls back to string comparison
6. **Nested objects**: Verify deep comparison works
7. **Arrays**: Verify array element order is respected (arrays ARE order-sensitive)

### Integration Test
1. Create addon config with specific `raw_values` JSON
2. Run `terraform plan` immediately after apply
3. Verify no diff is shown (even if API returns different ordering)
4. Update `raw_values` with actual content change
5. Verify diff IS shown for real changes

### Manual Verification
1. Apply a configuration with `raw_values`
2. Check API response ordering
3. Run `terraform plan`
4. Confirm no spurious diff

## Implementation Steps

1. ✅ **Code Changes**
   - ✅ Add imports: `encoding/json`, `reflect`
   - ✅ Add `suppressJSONOrderDiff` function at end of file
   - ✅ Update `raw_values` schema with `DiffSuppressFunc`

2. ✅ **Code Formatting**
   - ✅ Run `go fmt` on modified file

3. ⏳ **Testing**
   - ⏳ Add unit tests for diff suppression function
   - ⏳ Run acceptance tests: `TF_ACC=1 go test -v ./tencentcloud/services/tke -run TestAccTencentCloudKubernetesAddonConfig`
   - ⏳ Manual integration testing

4. ⏳ **Documentation**
   - ⏳ Update resource documentation if needed (likely no doc change required)
   - ⏳ Add changelog entry in `.changelog/`

5. ⏳ **Code Review**
   - ⏳ Submit PR with clear description
   - ⏳ Address review feedback

## Risk Assessment

- **Risk Level**: Low
- **Rationale**:
  - Only changes diff detection logic
  - No API call modifications
  - Graceful fallback for edge cases
  - Standard Terraform pattern
  - Easy to revert if issues arise

## Alternative Approaches Considered

### Alternative 1: Normalize JSON Before Storage
**Rejected** - Would require state migration and could break existing workflows

### Alternative 2: Sort JSON Keys Before Comparison
**Rejected** - More complex; `reflect.DeepEqual` already handles this elegantly

### Alternative 3: Ignore Diffs (Always Return True)
**Rejected** - Would hide legitimate content changes

## Timeline Estimate

- Development: 0.5 days (✅ Complete)
- Testing: 0.5 days
- Review: 0.5 days
- **Total**: 1.5 days

## Conclusion

This change improves user experience by eliminating false-positive diffs in the `raw_values` field while maintaining full backward compatibility and following Terraform best practices. The implementation is straightforward, low-risk, and uses standard Terraform patterns.
