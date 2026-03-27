# Add decode_values Field and Optimize ID for tencentcloud_kubernetes_addons

## Summary

Enhance the `tencentcloud_kubernetes_addons` data source with two improvements:
1. Add a `decode_values` computed field to provide base64-decoded addon configuration values
2. Use `cluster_id` as the resource ID instead of a hash-based ID

## Quick Overview

**Type:** Enhancement  
**Status:** Ready for Implementation  
**Impact:** Low Risk, Backward Compatible  
**Files Changed:** 1 file (`data_source_tc_kubernetes_addons.go`)

## What's Changing

### 1. New decode_values Field

**Before:**
```hcl
data "tencentcloud_kubernetes_addons" "example" {
  cluster_id = "cls-xxx"
}

# Only base64-encoded values available
output "config" {
  value = data.tencentcloud_kubernetes_addons.example.addons[0].raw_values
  # Output: "eyJrZXkiOiJ2YWx1ZSJ9" (need manual decode)
}
```

**After:**
```hcl
data "tencentcloud_kubernetes_addons" "example" {
  cluster_id = "cls-xxx"
}

# Both encoded and decoded values available
output "encoded_config" {
  value = data.tencentcloud_kubernetes_addons.example.addons[0].raw_values
  # Output: "eyJrZXkiOiJ2YWx1ZSJ9"
}

output "decoded_config" {
  value = data.tencentcloud_kubernetes_addons.example.addons[0].decode_values
  # Output: {"key":"value"} (automatically decoded!)
}
```

### 2. Simplified Resource ID

**Before:** ID is a hash of multiple addon names  
**After:** ID is simply the `cluster_id`

This makes the resource ID more predictable and stable.

## Why This Matters

- 🎯 **Better UX**: No need to manually decode base64 strings
- 🔄 **Backward Compatible**: All existing code continues to work
- 🆔 **Simpler IDs**: More intuitive resource identification
- ✅ **Error Handling**: Gracefully handles invalid base64 data

## Implementation Checklist

- [x] Proposal created (proposal.md)
- [x] Implementation tasks defined (tasks.md)
- [ ] Code changes applied
- [ ] Tests verified
- [ ] Documentation updated

## Files Involved

### Modified
- `tencentcloud/services/tke/data_source_tc_kubernetes_addons.go`
  - Add `decode_values` field to schema
  - Add base64 decoding logic
  - Change resource ID to use `cluster_id`

### Documents
- `proposal.md` - Detailed design and rationale
- `tasks.md` - Step-by-step implementation guide
- `README.md` - This file

## Key Changes Summary

| Change | Description | Lines Affected |
|--------|-------------|----------------|
| Schema Addition | Add `decode_values` field | ~51-55 |
| Import | Add `encoding/base64` | ~4-13 |
| Decoding Logic | Implement base64 decode | ~123-130 |
| ID Update | Use `cluster_id` as ID | ~142 |
| Formatting | Run `go fmt` | Entire file |

## Testing Strategy

1. ✅ Compilation verification
2. ✅ Manual Terraform testing with real cluster
3. ✅ Edge case testing (empty values, invalid base64)
4. ✅ State migration testing (optional)

## Next Steps

Ready to implement! Run:

```bash
/opsx:apply add-kubernetes-addons-decode-values
```

## References

- **Proposal**: See `proposal.md` for detailed design
- **Tasks**: See `tasks.md` for implementation steps
- **Current Code**: `tencentcloud/services/tke/data_source_tc_kubernetes_addons.go`
- **TKE SDK**: `tencentcloud-sdk-go/tencentcloud/tke/v20180525`

## Migration Notes

Users with existing state may see an ID change on next `terraform plan`. Simply run:

```bash
terraform refresh
```

This is a one-time operation with no configuration changes needed.

## Questions?

- Review `proposal.md` for design rationale
- Check `tasks.md` for detailed implementation steps
- Test with your own cluster to verify functionality

---

**Status**: ✅ Ready for implementation  
**Priority**: Medium  
**Estimated Time**: 30-45 minutes
