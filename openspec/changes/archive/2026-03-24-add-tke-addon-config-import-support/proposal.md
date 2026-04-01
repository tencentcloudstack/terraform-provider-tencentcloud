# Proposal: Add Import Support to tencentcloud_kubernetes_addon_config Resource

## Overview

Add Terraform import capability to the `tencentcloud_kubernetes_addon_config` resource, allowing users to import existing TKE cluster addon configurations into their Terraform state.

## Motivation

The `tencentcloud_kubernetes_addon_config` resource currently lacks import functionality, which creates several issues for users:

1. **Cannot Import Existing Resources**: Users with addon configurations created outside Terraform (via console, API, or other tools) cannot bring them under Terraform management without recreating them.

2. **State Migration Challenges**: Teams migrating to Terraform face difficulties when they need to manage existing addon configurations.

3. **Disaster Recovery**: In case of state file loss or corruption, users cannot recover by importing existing resources.

4. **Inconsistent UX**: Most other TKE resources (e.g., `tencentcloud_kubernetes_addon`, `tencentcloud_kubernetes_cluster`, `tencentcloud_kubernetes_node_pool`) support import, creating an inconsistent user experience.

## Current State

### Resource ID Format
The resource uses a composite ID format:
```
<cluster_id>#<addon_name>
```
Example: `cls-abc123#tcr`

### Read Function Status
✅ The `resourceTencentCloudKubernetesAddonConfigRead` function is **fully implemented** and can:
- Parse the composite ID correctly (lines 91-96)
- Call the API to fetch addon configuration (line 101)
- Set all schema fields from API response:
  - `cluster_id` (line 98)
  - `addon_name` (line 99)
  - `addon_version` (line 113)
  - `raw_values` (line 120) - with base64 decoding
  - `phase` (line 124)
  - `reason` (line 128)
- Handle resource not found scenarios (lines 106-110)

### Missing Component
❌ The resource schema lacks an `Importer` configuration, preventing Terraform from using the existing Read function for import operations.

## Proposed Solution

Add import support using the standard Terraform `ImportStatePassthrough` method, which leverages the existing Read function.

### Implementation Approach

**Single Change Required**: Add `Importer` field to resource schema:

```go
func ResourceTencentCloudKubernetesAddonConfig() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudKubernetesAddonConfigCreate,
        Read:   resourceTencentCloudKubernetesAddonConfigRead,
        Update: resourceTencentCloudKubernetesAddonConfigUpdate,
        Delete: resourceTencentCloudKubernetesAddonConfigDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{
            // ... existing schema
        },
    }
}
```

### Why ImportStatePassthrough?

1. **Simple and Standard**: The most common import method in Terraform providers
2. **No Additional Logic Needed**: The existing Read function already handles all necessary operations:
   - ID parsing
   - API calls
   - Field population
   - Error handling
3. **Consistency**: Aligns with similar resources in the same service:
   - `tencentcloud_kubernetes_addon` - uses `ImportStatePassthrough`
   - `tencentcloud_kubernetes_auth_attachment` - uses `ImportStatePassthrough`
   - `tencentcloud_kubernetes_native_node_pool` - uses `ImportStatePassthrough`

### User Experience

After implementation, users can import existing addon configurations:

```bash
# Import format
terraform import tencentcloud_kubernetes_addon_config.example <cluster_id>#<addon_name>

# Real example
terraform import tencentcloud_kubernetes_addon_config.tcr cls-abc123#tcr
```

The import will populate the Terraform state with:
- `cluster_id`: "cls-abc123"
- `addon_name`: "tcr"
- `addon_version`: (from API)
- `raw_values`: (from API, base64 decoded)
- `phase`: (from API)
- `reason`: (from API)

## Impact Analysis

### User Impact
- **Positive**: Enables import functionality for existing resources
- **Backward Compatibility**: ✅ 100% compatible - no breaking changes to existing configurations
- **Migration**: ✅ No migration needed - purely additive feature

### Code Impact
- **Files Modified**: 1 file (`resource_tc_kubernetes_addon_config.go`)
- **Lines Changed**: ~3 lines (add Importer block)
- **Risk Level**: 🟢 **Very Low** - uses standard Terraform pattern with no custom logic

### Testing Impact
- **Existing Tests**: ✅ No impact on existing tests
- **New Tests Required**: Acceptance test for import functionality

## Dependencies

- **Terraform SDK**: Already using `github.com/hashicorp/terraform-plugin-sdk/v2` - no new dependencies
- **API Requirements**: No API changes needed - uses existing `DescribeExtensionAddon` API

## Alternatives Considered

### Alternative 1: Custom Import Function
```go
Importer: &schema.ResourceImporter{
    StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
        // Custom logic
    },
}
```

**Rejected**: Unnecessary complexity - the standard Read function already handles all requirements.

### Alternative 2: Wait for Provider-Wide Import Refactor
**Rejected**: Import is a standard, stable feature. No reason to delay.

## Success Criteria

1. ✅ Users can successfully import existing addon configurations
2. ✅ All schema fields are correctly populated after import
3. ✅ Imported resources show no diff when `terraform plan` is run
4. ✅ Subsequent updates to imported resources work correctly
5. ✅ Documentation clearly explains the import format
6. ✅ Acceptance tests validate import functionality

## Timeline

- **Effort Estimate**: 30 minutes
  - Code change: 5 minutes (3 lines)
  - Testing: 15 minutes
  - Documentation: 10 minutes
- **Complexity**: 🟢 Trivial - standard pattern

## Documentation Updates

1. **Resource Documentation**: Add import section with example
2. **Example Code**: Show import command usage
3. **Migration Guide**: Document for teams adopting Terraform

## References

- [Terraform Provider Development - Import](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/import)
- Similar implementations in this provider:
  - `resource_tc_kubernetes_addon.go:26-28`
  - `resource_tc_kubernetes_auth_attachment.go:22-24`
  - `resource_tc_kubernetes_native_node_pool.go:27-29`
