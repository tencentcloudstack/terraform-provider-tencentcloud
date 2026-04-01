# Technical Design: Add Import Support to tencentcloud_kubernetes_addon_config

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│  User executes: terraform import                            │
│  tencentcloud_kubernetes_addon_config.example cls-123#tcr  │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Terraform Core                                             │
│  - Parses import command                                    │
│  - Extracts ID: "cls-123#tcr"                              │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  ImportStatePassthrough (Terraform SDK)                     │
│  - Calls: d.SetId("cls-123#tcr")                           │
│  - No custom logic needed                                   │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  resourceTencentCloudKubernetesAddonConfigRead              │
│  (EXISTING FUNCTION - NO CHANGES NEEDED)                    │
│                                                             │
│  1. Parse ID                                                │
│     idSplit := strings.Split(d.Id(), "#")                  │
│     clusterId := idSplit[0]  // "cls-123"                  │
│     addonName := idSplit[1]  // "tcr"                      │
│                                                             │
│  2. Call TKE API                                            │
│     respData := service.DescribeKubernetesAddonById(...)   │
│                                                             │
│  3. Set State Fields                                        │
│     d.Set("cluster_id", clusterId)                         │
│     d.Set("addon_name", addonName)                         │
│     d.Set("addon_version", respData.AddonVersion)          │
│     d.Set("raw_values", decodedValues)                     │
│     d.Set("phase", respData.Phase)                         │
│     d.Set("reason", respData.Reason)                       │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Terraform State File Updated                               │
│  {                                                          │
│    "cluster_id": "cls-123",                                │
│    "addon_name": "tcr",                                    │
│    "addon_version": "v1.0.0",                              │
│    "raw_values": "{...}",                                  │
│    "phase": "Running",                                     │
│    "reason": ""                                            │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

## Code Changes

### File: `resource_tc_kubernetes_addon_config.go`

**Location**: Lines 19-68 (Resource definition)

**Change Type**: Addition (3 lines)

**Before**:
```go
func ResourceTencentCloudKubernetesAddonConfig() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudKubernetesAddonConfigCreate,
        Read:   resourceTencentCloudKubernetesAddonConfigRead,
        Update: resourceTencentCloudKubernetesAddonConfigUpdate,
        Delete: resourceTencentCloudKubernetesAddonConfigDelete,
        Schema: map[string]*schema.Schema{
            // ... schema definition
        },
    }
}
```

**After**:
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
            // ... schema definition
        },
    }
}
```

**Diff**:
```diff
func ResourceTencentCloudKubernetesAddonConfig() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudKubernetesAddonConfigCreate,
        Read:   resourceTencentCloudKubernetesAddonConfigRead,
        Update: resourceTencentCloudKubernetesAddonConfigUpdate,
        Delete: resourceTencentCloudKubernetesAddonConfigDelete,
+       Importer: &schema.ResourceImporter{
+           State: schema.ImportStatePassthrough,
+       },
        Schema: map[string]*schema.Schema{
```

## Import Flow Validation

### Step-by-Step Execution

#### 1. User Initiates Import
```bash
$ terraform import tencentcloud_kubernetes_addon_config.tcr cls-abc123#tcr
```

#### 2. Terraform Core Processing
- Parses command line arguments
- Identifies resource type: `tencentcloud_kubernetes_addon_config`
- Extracts import ID: `cls-abc123#tcr`
- Looks up resource's Importer configuration

#### 3. ImportStatePassthrough Execution
```go
// Pseudo-code of what ImportStatePassthrough does
func ImportStatePassthrough(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    // Simply sets the ID - that's it!
    d.SetId("cls-abc123#tcr")
    return []*schema.ResourceData{d}, nil
}
```

#### 4. Read Function Invocation
```go
func resourceTencentCloudKubernetesAddonConfigRead(d *schema.ResourceData, meta interface{}) error {
    // Get ID set by importer
    id := d.Id() // "cls-abc123#tcr"
    
    // Parse ID (EXISTING LOGIC - lines 91-96)
    idSplit := strings.Split(id, tccommon.FILED_SP) // FILED_SP = "#"
    if len(idSplit) != 2 {
        return fmt.Errorf("id is broken,%s", id)
    }
    clusterId := idSplit[0] // "cls-abc123"
    addonName := idSplit[1] // "tcr"
    
    // Set identifier fields (lines 98-99)
    d.Set("cluster_id", clusterId)
    d.Set("addon_name", addonName)
    
    // Fetch from API (line 101)
    respData, err := service.DescribeKubernetesAddonById(ctx, clusterId, addonName)
    if err != nil {
        return err
    }
    
    // Handle not found (lines 106-110)
    if respData == nil {
        d.SetId("")
        return nil
    }
    
    // Set all fields from API (lines 112-129)
    d.Set("addon_version", respData.AddonVersion)
    d.Set("raw_values", decodedBase64JSON)
    d.Set("phase", respData.Phase)
    d.Set("reason", respData.Reason)
    
    return nil
}
```

#### 5. State File Updated
```hcl
# Resulting state
resource "tencentcloud_kubernetes_addon_config" "tcr" {
  cluster_id     = "cls-abc123"
  addon_name     = "tcr"
  addon_version  = "v1.2.3"  # from API
  raw_values     = "{\"key\":\"value\"}"  # from API, decoded
  phase          = "Running"  # from API
  reason         = ""  # from API
}
```

#### 6. Import Success Message
```bash
tencentcloud_kubernetes_addon_config.tcr: Importing from ID "cls-abc123#tcr"...
tencentcloud_kubernetes_addon_config.tcr: Import prepared!
  Prepared tencentcloud_kubernetes_addon_config for import
tencentcloud_kubernetes_addon_config.tcr: Refreshing state... [id=cls-abc123#tcr]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

## Edge Cases and Error Handling

### Case 1: Invalid ID Format
**Input**: `terraform import tencentcloud_kubernetes_addon_config.example invalid-id`

**Behavior**:
```go
// In Read function (lines 91-94)
idSplit := strings.Split("invalid-id", "#")
if len(idSplit) != 2 {
    return fmt.Errorf("id is broken,%s", "invalid-id")
}
```

**Output**:
```
Error: id is broken,invalid-id
```

**Status**: ✅ Already handled by existing Read function

---

### Case 2: Addon Does Not Exist
**Input**: `terraform import tencentcloud_kubernetes_addon_config.example cls-123#nonexistent`

**Behavior**:
```go
// In Read function (lines 106-110)
if respData == nil {
    d.SetId("")
    log.Printf("[WARN] resource `kubernetes_addon_config` [%s] not found", d.Id())
    return nil
}
```

**Output**:
```
Error: Cannot import non-existent remote object
```

**Status**: ✅ Already handled - Terraform detects empty ID and fails import

---

### Case 3: API Permission Error
**Input**: `terraform import tencentcloud_kubernetes_addon_config.example cls-123#tcr`
(User lacks API permissions)

**Behavior**:
```go
// In Read function (lines 101-104)
respData, err := service.DescribeKubernetesAddonById(ctx, clusterId, addonName)
if err != nil {
    return err
}
```

**Output**:
```
Error: [TencentCloudSDKError] Code=AuthFailure.UnauthorizedOperation, 
Message=You are not authorized to perform this operation
```

**Status**: ✅ Already handled - API errors propagate to user

---

### Case 4: Cluster Does Not Exist
**Input**: `terraform import tencentcloud_kubernetes_addon_config.example cls-nonexistent#tcr`

**Behavior**: API returns nil or error

**Output**: Same as Case 2 or Case 3 depending on API behavior

**Status**: ✅ Already handled

---

### Case 5: Empty ID Components
**Input**: `terraform import tencentcloud_kubernetes_addon_config.example cls-123#`

**Behavior**:
```go
idSplit := strings.Split("cls-123#", "#")
// idSplit = ["cls-123", ""]
addonName := idSplit[1] // ""
```

**API Call**: `DescribeKubernetesAddonById(ctx, "cls-123", "")`

**Output**: API error (invalid addon name)

**Status**: ✅ Handled - API validates and returns error

## Compatibility Matrix

| Scenario | Supported | Behavior |
|----------|-----------|----------|
| Import with valid ID | ✅ Yes | Success |
| Import with invalid ID format | ✅ Yes | Error: "id is broken" |
| Import non-existent addon | ✅ Yes | Error: "Cannot import non-existent remote object" |
| Import without permissions | ✅ Yes | Error: API permission error |
| Subsequent `terraform plan` | ✅ Yes | Shows no diff |
| Subsequent `terraform apply` | ✅ Yes | Updates work normally |
| Import then modify | ✅ Yes | Changes detected correctly |
| Re-import same resource | ✅ Yes | Overwrites state |

## Testing Strategy

### Unit Tests
Not applicable - ImportStatePassthrough has no custom logic to unit test.

### Acceptance Tests

**Test File**: `resource_tc_kubernetes_addon_config_test.go`

**New Test Function**:
```go
func TestAccTencentCloudKubernetesAddonConfig_import(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { tcacctest.AccPreCheck(t) },
        Providers:    tcacctest.AccProviders,
        CheckDestroy: testAccCheckKubernetesAddonConfigDestroy,
        Steps: []resource.TestStep{
            {
                // Step 1: Create addon config
                Config: testAccKubernetesAddonConfig_basic(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_config.test", 
                        "addon_name", "tcr"),
                    resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_config.test", 
                        "addon_version"),
                ),
            },
            {
                // Step 2: Import and verify
                ResourceName:      "tencentcloud_kubernetes_addon_config.test",
                ImportState:       true,
                ImportStateVerify: true,
                // raw_values may have formatting differences, but suppressJSONOrderDiff handles it
            },
        },
    })
}
```

**Test Execution**:
```bash
go test -v ./tencentcloud/services/tke -run TestAccTencentCloudKubernetesAddonConfig_import
```

### Manual Testing Checklist

- [ ] Import existing addon config
- [ ] Run `terraform plan` after import (should show no diff)
- [ ] Modify imported resource
- [ ] Run `terraform apply` to verify update works
- [ ] Try invalid import IDs
- [ ] Try importing non-existent resource

## Performance Considerations

- **Import Time**: Same as Read operation (~1-2 seconds per resource)
- **API Calls**: 1 API call per import (`DescribeExtensionAddon`)
- **State Size**: No change - same fields as Create operation
- **Memory Usage**: Negligible - standard Terraform operation

## Security Considerations

- **Authentication**: Uses existing provider authentication
- **Authorization**: Requires same permissions as Read operation
  - Action: `tke:DescribeExtensionAddon`
  - Resource: `qcs::tke:${region}:uin/${uin}:cluster/${cluster-id}`
- **Data Exposure**: No additional data exposed beyond Read operation
- **Audit Trail**: Import operations logged via CloudAudit (if enabled)

## Rollback Plan

If issues are discovered after deployment:

1. **Immediate**: No rollback needed - feature is opt-in
   - Users must explicitly run `terraform import`
   - No impact on existing resources or workflows

2. **If Critical Bug Found**:
   ```go
   // Temporarily disable import
   Importer: nil,
   ```

3. **No Data Loss Risk**: Import only reads data, never writes to cloud

## Documentation Updates

### Resource Documentation Page

**Add Import Section**:
````markdown
## Import

Kubernetes addon configuration can be imported using the id, e.g.

```bash
terraform import tencentcloud_kubernetes_addon_config.example cls-abc123#tcr
```

Where:
- `cls-abc123` is the cluster ID
- `tcr` is the addon name

**Note**: The ID format is `<cluster_id>#<addon_name>`, using `#` as the separator.
````

### Example Code

**Add to examples**:
```hcl
# Import existing TCR addon configuration
# terraform import tencentcloud_kubernetes_addon_config.tcr cls-abc123#tcr

resource "tencentcloud_kubernetes_addon_config" "tcr" {
  cluster_id    = "cls-abc123"
  addon_name    = "tcr"
  addon_version = "v1.2.3"
  raw_values    = jsonencode({
    # configuration...
  })
}
```

## Success Metrics

1. ✅ Code compiles without errors
2. ✅ Acceptance test passes
3. ✅ Manual import test successful
4. ✅ No diff after import + plan
5. ✅ Documentation complete
6. ✅ Code review approved

## References

- [Terraform Plugin SDK - Import](https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/schema/resource_importer.go)
- [ImportStatePassthrough Implementation](https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/schema/resource_importer.go#L59-L71)
- Provider Examples:
  - `resource_tc_kubernetes_addon.go` (similar resource with import)
  - `resource_tc_kubernetes_auth_attachment.go` (uses same pattern)
