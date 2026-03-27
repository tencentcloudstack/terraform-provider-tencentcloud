# Proposal: Add decode_values Field and Optimize ID for tencentcloud_kubernetes_addons

## Overview

This proposal aims to enhance the `tencentcloud_kubernetes_addons` data source by:
1. Adding a new computed field `decode_values` to provide the base64-decoded version of `raw_values`
2. Optimizing the resource ID to use `cluster_id` as the unique identifier since it's a required parameter

## Motivation

### Problem 1: raw_values Field is Base64 Encoded

Currently, the `raw_values` field in the `tencentcloud_kubernetes_addons` data source returns base64-encoded JSON strings. Users need to manually decode this field to read the actual addon configuration values, which is inconvenient.

**Current Behavior:**
```hcl
data "tencentcloud_kubernetes_addons" "example" {
  cluster_id = "cls-xxx"
}

# Users must decode raw_values manually
# raw_values = "eyJrZXkiOiJ2YWx1ZSJ9"  (base64 encoded)
```

**Expected Behavior:**
Users should have access to both the encoded `raw_values` (for compatibility) and a decoded `decode_values` field for direct consumption.

### Problem 2: Resource ID Not Using cluster_id

The current implementation generates a hash-based ID from multiple addon names, but since `cluster_id` is already a required parameter and uniquely identifies the query scope, it would be more intuitive and stable to use `cluster_id` directly as the resource ID.

**Current Implementation:**
```go
d.SetId(helper.DataResourceIdsHash(ids))  // Hash of cluster_id+addon_name pairs
```

**Proposed Implementation:**
```go
d.SetId(clusterId)  // Direct use of cluster_id
```

## Why These Changes Matter

1. **Improved User Experience**: Users can directly access decoded addon configurations without additional processing
2. **Better Compatibility**: Keeps `raw_values` for backward compatibility while adding convenience
3. **Simpler Resource Identification**: Using `cluster_id` as ID is more predictable and aligns with Terraform best practices
4. **Reduced Complexity**: Eliminates unnecessary hash computation for ID generation

## Proposed Solution

### Change 1: Add decode_values Field

**Schema Changes:**

Add a new computed field in the `addons` list:

```go
"decode_values": {
    Type:        schema.TypeString,
    Computed:    true,
    Description: "Add-on parameters in JSON format (decoded from raw_values). Note: This field may return empty string if raw_values is not set.",
},
```

**Data Processing Logic:**

In the read function, after setting `raw_values`, decode it and set `decode_values`:

```go
if addons.RawValues != nil {
    addonsMap["raw_values"] = addons.RawValues
    
    // Decode base64 to get decode_values
    if decodedBytes, err := base64.StdEncoding.DecodeString(*addons.RawValues); err == nil {
        addonsMap["decode_values"] = string(decodedBytes)
    } else {
        // If decode fails, set empty string
        addonsMap["decode_values"] = ""
    }
}
```

**Import Required:**
```go
"encoding/base64"
```

### Change 2: Use cluster_id as Resource ID

**Current Code (Line 142):**
```go
d.SetId(helper.DataResourceIdsHash(ids))
```

**New Code:**
```go
d.SetId(clusterId)
```

**Rationale:**
- `cluster_id` is required and uniquely identifies the data source query
- More predictable and stable than hash-based IDs
- Aligns with Terraform conventions for data sources
- Simplifies the code by removing unnecessary ID construction

## Technical Details

### File to Modify
- `tencentcloud/services/tke/data_source_tc_kubernetes_addons.go`

### Changes Summary

1. **Import Statement** (add if not present):
   ```go
   "encoding/base64"
   ```

2. **Schema Definition** (around line 47-51):
   - Add `decode_values` field after `raw_values`

3. **Data Processing** (around line 123-125):
   - Add base64 decoding logic after setting `raw_values`

4. **Resource ID** (line 142):
   - Change from `helper.DataResourceIdsHash(ids)` to `clusterId`

5. **Code Formatting**:
   - Run `go fmt` on the file after all changes

### Error Handling

The base64 decoding includes error handling:
- If decoding succeeds: set the decoded JSON string
- If decoding fails: set empty string (graceful degradation)

This ensures the data source continues to work even if `raw_values` contains invalid base64 data.

## Backward Compatibility

### ✅ Fully Backward Compatible

1. **decode_values is a new field**: Existing code won't break
2. **raw_values remains unchanged**: All existing references continue to work
3. **ID change is internal**: Data source queries will still work, though state may need refresh

### Migration Notes

Users who have this data source in their state may see a plan diff on the next `terraform plan` due to the ID change. They should run:
```bash
terraform refresh
```

This is a one-time operation and doesn't require any configuration changes.

## Impact Assessment

### 🟢 Low Risk Changes

1. **Schema Addition**: Adding a computed field has no breaking impact
2. **ID Simplification**: More stable and predictable
3. **Decoding Logic**: Safe with error handling

### Benefits

- ✅ Easier addon configuration inspection
- ✅ No manual base64 decoding needed
- ✅ Cleaner resource identification
- ✅ Maintains full backward compatibility

### Potential Issues

- ⚠️ State refresh may be needed after upgrade (one-time, minor)
- ℹ️ If `raw_values` contains non-JSON data, `decode_values` will show the raw string

## Testing Recommendations

1. **Unit Tests**: Verify base64 decoding logic
2. **Integration Tests**: Test with real cluster addons
3. **Edge Cases**:
   - Empty `raw_values`
   - Invalid base64 data
   - Special characters in decoded JSON

## Example Usage

After implementation:

```hcl
data "tencentcloud_kubernetes_addons" "example" {
  cluster_id = "cls-xxx"
  addon_name = "nginx-ingress"
}

output "raw_config" {
  value = data.tencentcloud_kubernetes_addons.example.addons[0].raw_values
}

output "decoded_config" {
  # New field - no manual decoding needed!
  value = data.tencentcloud_kubernetes_addons.example.addons[0].decode_values
}

output "data_source_id" {
  # Now simply the cluster_id
  value = data.tencentcloud_kubernetes_addons.example.id
}
```

## References

- TKE SDK: `tencentcloud-sdk-go/tencentcloud/tke/v20180525`
- Current Implementation: `tencentcloud/services/tke/data_source_tc_kubernetes_addons.go`
- Terraform Schema Documentation: https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas

## Conclusion

This proposal introduces two small but valuable improvements to the `tencentcloud_kubernetes_addons` data source:
1. A user-friendly `decode_values` field for easier addon configuration access
2. A simplified resource ID using `cluster_id`

Both changes are backward compatible and follow Terraform best practices.
