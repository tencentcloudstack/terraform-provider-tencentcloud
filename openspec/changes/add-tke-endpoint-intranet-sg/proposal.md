# Add Intranet Security Group Support to TKE Cluster Endpoint Resource

**Status**: Draft  
**Author**: Terraform Provider Team  
**Created**: 2026-03-05  
**Change Type**: Enhancement

---

## Why

### Problem Statement

The `tencentcloud_kubernetes_cluster_endpoint` resource currently only supports security group configuration for **internet (external network)** access via the `cluster_internet_security_group` field. However, for **intranet (internal network)** access, there is no equivalent security group configuration support.

This creates an asymmetry in security configuration:
- ✅ Internet access: Can specify security group (`cluster_internet_security_group`)
- ❌ Intranet access: Cannot specify security group (missing feature)

### Business Need

Users need to:
1. Configure security groups for intranet cluster endpoints to control internal network access
2. Have consistent security group configuration for both internet and intranet access
3. Set security groups at cluster endpoint creation time (ForceNew requirement)

---

## What Changes

### 1. Schema Changes

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

#### Add New Field: `cluster_intranet_security_group`

```go
"cluster_intranet_security_group": {
    Type:        schema.TypeString,
    Optional:    true,
    ForceNew:    true,  // New: Must recreate if changed
    Description: "Security group ID for intranet cluster endpoint.",
},
```

#### Update Existing Field Description: `cluster_internet_security_group`

```go
"cluster_internet_security_group": {
    Type:        schema.TypeString,
    Optional:    true,
    Description: "Security group ID for internet cluster endpoint. NOTE: This argument must not be empty if cluster internet enabled.",
    //                  ^^^^^^^^^^^
    // Clarify this is for INTERNET (external network)
},
```

**Key Differences**:
| Field | ForceNew | Network Type |
|-------|----------|--------------|
| `cluster_internet_security_group` | ❌ No | Internet (External) |
| `cluster_intranet_security_group` | ✅ **Yes** | Intranet (Internal) |

---

### 2. Create Logic Changes

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

#### Function: `resourceTencentCloudTkeClusterEndpointCreate` (lines 226-289)

**Current Implementation** (lines 239-243):
```go
var (
    err                          error
    clusterInternet              = d.Get("cluster_internet").(bool)
    clusterIntranet              = d.Get("cluster_intranet").(bool)
    intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
    clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
    // Missing: clusterIntranetSecurityGroup
    clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
    clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
    extensiveParameters          = d.Get("extensive_parameters").(string)
)
```

**Proposed Changes**:
```go
var (
    err                          error
    clusterInternet              = d.Get("cluster_internet").(bool)
    clusterIntranet              = d.Get("cluster_intranet").(bool)
    intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
    clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
    clusterIntranetSecurityGroup = d.Get("cluster_intranet_security_group").(string)  // ✅ New
    clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
    clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
    extensiveParameters          = d.Get("extensive_parameters").(string)
)
```

**Note**: No additional validation required for `cluster_intranet_security_group`. The field can be set independently of `cluster_intranet` state.

#### Update Intranet Switch Call (line 263)

**Current**:
```go
if clusterIntranet {
    err := tencentCloudClusterIntranetSwitch(ctx, &service, id, intranetSubnetId, true, clusterIntranetDomain)
    if err != nil {
        return err
    }
    // ...
}
```

**Proposed**:
```go
if clusterIntranet {
    err := tencentCloudClusterIntranetSwitch(ctx, &service, id, intranetSubnetId, clusterIntranetSecurityGroup, true, clusterIntranetDomain)
    //                                                                           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    //                                                                           Pass security group parameter
    if err != nil {
        return err
    }
    // ...
}
```

---

### 3. Helper Function Changes

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

#### Function: `tencentCloudClusterIntranetSwitch` (lines 487-506)

**Current Signature**:
```go
func tencentCloudClusterIntranetSwitch(ctx context.Context, service *TkeService, id, subnetId string, enable bool, domain string) (err error)
```

**Proposed Signature**:
```go
func tencentCloudClusterIntranetSwitch(ctx context.Context, service *TkeService, id, subnetId, securityGroup string, enable bool, domain string) (err error)
//                                                                                            ^^^^^^^^^^^^^^^
//                                                                                            New parameter
```

**Implementation Changes** (line 490):
```go
func tencentCloudClusterIntranetSwitch(ctx context.Context, service *TkeService, id, subnetId, securityGroup string, enable bool, domain string) (err error) {
    err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        if enable {
            err = service.CreateClusterEndpoint(ctx, id, subnetId, securityGroup, false, domain, "")
            //                                                      ^^^^^^^^^^^^^^^
            //                                                      Pass security group for intranet
            if err != nil {
                return tccommon.RetryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
            }
        } else {
            err = service.DeleteClusterEndpoint(ctx, id, false)
            if err != nil {
                return tccommon.RetryError(err)
            }
        }
        return nil
    })
    if err != nil {
        return err
    }
    return nil
}
```

---

### 4. Service Layer Changes

**File**: `tencentcloud/services/tke/service_tencentcloud_tke.go`

#### Function: `CreateClusterEndpoint` (line 992)

**Current Logic** (lines 1009-1011):
```go
if securityGroupId != "" && internet {
    request.SecurityGroup = &securityGroupId
}
```

**Proposed Logic**:
```go
// Support security group for both internet and intranet
if securityGroupId != "" {
    request.SecurityGroup = &securityGroupId
}
// Remove the "&& internet" restriction to support intranet security group
```

**Rationale**: The TKE API already supports security group for intranet endpoints, but the current code logic restricts it to internet-only.

---

### 5. Update and Delete Function Changes

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

#### Update Function (lines 354-383)

**Current calls** (lines 355, 366, 375):
```go
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, clusterIntranet, clusterIntranetDomain)
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, false, clusterIntranetDomain)
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, true, clusterIntranetDomain)
```

**Proposed changes**:
```go
// Add empty string "" for securityGroup parameter (not used in update/delete)
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, "", clusterIntranet, clusterIntranetDomain)
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, "", false, clusterIntranetDomain)
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, "", true, clusterIntranetDomain)
```

**Note**: As per requirements, update stage does not support security group modification (no API support yet), so we pass empty string.

#### Delete Function (line 426)

**Current call**:
```go
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, "", false, "")
```

**Proposed change**:
```go
err = tencentCloudClusterIntranetSwitch(ctx, &service, id, "", "", false, "")
//                                                           ^^
//                                                           Add securityGroup parameter (empty)
```

---

## Impact

### Affected Files

1. ✅ **Schema Definition** - `resource_tc_kubernetes_cluster_endpoint.go` (lines 24-131)
2. ✅ **Create Function** - `resourceTencentCloudTkeClusterEndpointCreate` (lines 226-289)
3. ✅ **Update Function** - `resourceTencentCloudTkeClusterEndpointUpdate` (lines 291-386)
4. ✅ **Delete Function** - `resourceTencentCloudTkeClusterEndpointDelete` (lines 388-438)
5. ✅ **Helper Function** - `tencentCloudClusterIntranetSwitch` (lines 487-506)
6. ✅ **Service Layer** - `service_tencentcloud_tke.go` (lines 992-1022)

### Breaking Changes

**None** - This is a backward-compatible enhancement:
- Existing configurations without `cluster_intranet_security_group` continue to work
- New field is optional (but required when `cluster_intranet` = true)
- ForceNew behavior only affects new configurations using this field

### User Migration

**No migration required** for existing users:
```hcl
# Existing configuration (still works)
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id       = "cls-xxxxx"
  cluster_intranet = true
  cluster_intranet_subnet_id = "subnet-xxxxx"
  # cluster_intranet_security_group not set (backward compatible)
}

# New configuration (with security group)
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id       = "cls-xxxxx"
  cluster_intranet = true
  cluster_intranet_subnet_id = "subnet-xxxxx"
  cluster_intranet_security_group = "sg-xxxxx"  # ✅ New field (optional)
}
```

**Note**: The `cluster_intranet_security_group` field is **optional** and can be set independently based on your security requirements.

---

## Design Decisions

### 1. Why ForceNew for `cluster_intranet_security_group`?

**Reason**: Per requirements and API limitations:
- No API support for modifying intranet security group after creation
- Read stage cannot query current security group
- Update stage cannot modify security group
- Therefore, ForceNew ensures correct behavior: recreate endpoint if security group changes

**Contrast**: `cluster_internet_security_group` does NOT have ForceNew because:
- API supports `ModifyClusterEndpointSG` for internet endpoints (line 316)
- Can be updated in-place

### 2. Why Not Support in Read/Update?

**Per Requirements**:
> "read、update阶段不需要考虑，当前暂无接口支持cluster_intranet_security_group字段查询和修改"

The TKE API currently does not provide:
- ❌ Query endpoint to retrieve current intranet security group
- ❌ Modify endpoint to update intranet security group

Therefore:
- **Read**: Cannot set `cluster_intranet_security_group` from API response
- **Update**: Cannot handle changes to `cluster_intranet_security_group` (ForceNew handles this)
- **Create**: ✅ Supported - pass security group to `CreateClusterEndpoint`
- **Delete**: ✅ Supported - no security group needed for deletion

### 3. Parameter Naming Convention

Following existing pattern:
| Network Type | Subnet Field | Security Group Field |
|--------------|--------------|----------------------|
| Internet | N/A | `cluster_internet_security_group` |
| Intranet | `cluster_intranet_subnet_id` | `cluster_intranet_security_group` |

Consistent `cluster_*` prefix with `internet` / `intranet` suffix.

---

## Testing Strategy

### Manual Testing Checklist

1. **Create with intranet + security group**
   ```hcl
   resource "tencentcloud_kubernetes_cluster_endpoint" "test" {
     cluster_id = "cls-xxxxx"
     cluster_intranet = true
     cluster_intranet_subnet_id = "subnet-xxxxx"
     cluster_intranet_security_group = "sg-xxxxx"
   }
   ```
   - ✅ Should successfully create intranet endpoint with security group

2. **Create intranet without security group**
   ```hcl
   resource "tencentcloud_kubernetes_cluster_endpoint" "test" {
     cluster_id = "cls-xxxxx"
     cluster_intranet = true
     cluster_intranet_subnet_id = "subnet-xxxxx"
     # No cluster_intranet_security_group (optional)
   }
   ```
   - ✅ Should successfully create intranet endpoint without security group

3. **ForceNew behavior**
   ```hcl
   # Initial
   cluster_intranet_security_group = "sg-aaaaa"
   
   # Change
   cluster_intranet_security_group = "sg-bbbbb"
   ```
   - ✅ Should trigger resource recreation (destroy + create)

4. **Backward compatibility**
   - Existing configurations without `cluster_intranet_security_group` should continue working

5. **Both internet and intranet security groups**
   ```hcl
   resource "tencentcloud_kubernetes_cluster_endpoint" "test" {
     cluster_id = "cls-xxxxx"
     cluster_internet = true
     cluster_internet_security_group = "sg-internet"
     cluster_intranet = true
     cluster_intranet_subnet_id = "subnet-xxxxx"
     cluster_intranet_security_group = "sg-intranet"
   }
   ```
   - ✅ Should successfully configure both internet and intranet security groups

---

## Risks and Mitigation

### Risk 1: API Behavior Uncertainty

**Risk**: TKE API behavior for intranet security group may differ from internet.

**Mitigation**:
- Test thoroughly with real TKE clusters
- Verify API accepts security group for intranet endpoints
- Check error messages and retry logic

### Risk 2: ForceNew May Cause Service Disruption

**Risk**: Changing security group triggers endpoint recreation, causing brief downtime.

**Mitigation**:
- Document ForceNew behavior clearly
- Warn users about potential disruption in field description
- Recommend planning changes during maintenance windows

---

## Future Enhancements

1. **Read Support**: When TKE API adds query capability, add security group to Read function
2. **Update Support**: If API adds modify capability, remove ForceNew and implement in Update
3. **Data Source**: Add security group fields to `tencentcloud_kubernetes_cluster` data source

---

## Documentation Requirements

1. **Resource Documentation**: Update `website/docs/r/kubernetes_cluster_endpoint.html.markdown`
   - Add `cluster_intranet_security_group` field description
   - Add usage examples
   - Document ForceNew behavior
   - Update `cluster_internet_security_group` description for clarity

2. **Examples**: Add example configuration showing both internet and intranet security groups

3. **Changelog**: Add entry describing new field
   ```
   ```release-note:enhancement
   resource/tencentcloud_kubernetes_cluster_endpoint: add `cluster_intranet_security_group` field to support security group configuration for intranet cluster endpoints
   ```
   ```

---

## Summary

This change adds symmetric security group support for intranet cluster endpoints, matching the existing internet security group capability. The implementation follows Terraform best practices with ForceNew behavior for immutable fields and backward compatibility.

**Key Points**:
- ✅ New field: `cluster_intranet_security_group` (ForceNew, optional)
- ✅ Clarified descriptions for internet vs intranet security groups
- ✅ Service layer updated to support intranet security groups
- ✅ Backward compatible - no breaking changes
- ✅ No validation constraints - field can be set independently
- ⏸️ Read/Update not implemented (API limitation)
