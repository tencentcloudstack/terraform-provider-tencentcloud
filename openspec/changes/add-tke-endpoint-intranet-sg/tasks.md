# Implementation Tasks: Add Intranet Security Group to TKE Cluster Endpoint

**Change ID**: add-tke-endpoint-intranet-sg  
**Status**: Pending  
**Estimated Effort**: 2-3 hours

---

## Overview

Add `cluster_intranet_security_group` field to `tencentcloud_kubernetes_cluster_endpoint` resource to support security group configuration for intranet cluster endpoints.

---

## Task Breakdown

### Task 1: Update Schema Definition ✅ Completed

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

**Location**: Schema definition (lines 24-131)

- [x] Add new field `cluster_intranet_security_group` after line 47:
  ```go
  "cluster_intranet_security_group": {
      Type:        schema.TypeString,
      Optional:    true,
      ForceNew:    true,
      Description: "Security group ID for intranet cluster endpoint.",
  },
  ```

- [x] Update `cluster_internet_security_group` description (line 43-47) to clarify it's for internet:
  ```go
  "cluster_internet_security_group": {
      Type:        schema.TypeString,
      Optional:    true,
      Description: "Security group ID for internet cluster endpoint. NOTE: This argument must not be empty if cluster internet enabled.",
  },
  ```

**Acceptance Criteria**:
- Schema includes new field with correct type, optional, and ForceNew attributes ✅
- Field description is clear about intranet usage ✅
- Internet security group description clarifies external network usage ✅

---

### Task 2: Update Create Function - Variable Declaration ✅ Completed

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

**Function**: `resourceTencentCloudTkeClusterEndpointCreate`

**Location**: Variable declaration section (lines 235-244)

- [x] Add `clusterIntranetSecurityGroup` variable after line 240:
  ```go
  var (
      err                          error
      clusterInternet              = d.Get("cluster_internet").(bool)
      clusterIntranet              = d.Get("cluster_intranet").(bool)
      intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
      clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
      clusterIntranetSecurityGroup = d.Get("cluster_intranet_security_group").(string)  // ✅ Add this
      clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
      clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
      extensiveParameters          = d.Get("extensive_parameters").(string)
  )
  ```

**Acceptance Criteria**:
- Variable correctly retrieves value from schema ✅
- Follows existing naming convention ✅

---

### Task 3: Update Intranet Switch Call in Create ✅ Completed

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

**Function**: `resourceTencentCloudTkeClusterEndpointCreate`

**Location**: Intranet creation block (line 263)

- [x] Update `tencentCloudClusterIntranetSwitch` call to include security group:
  ```go
  if clusterIntranet {
      err := tencentCloudClusterIntranetSwitch(ctx, &service, id, intranetSubnetId, clusterIntranetSecurityGroup, true, clusterIntranetDomain)
      //                                                                           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      //                                                                           Pass security group parameter
      if err != nil {
          return err
      }
      err = waitForClusterEndpointFinish(ctx, &service, id, true, false)
      if err != nil {
          return err
      }
  }
  ```

**Acceptance Criteria**:
- Security group parameter passed to helper function ✅
- Call signature matches updated function signature ✅

---

### Task 4: Update Intranet Switch Helper Function ✅ Completed

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

**Function**: `tencentCloudClusterIntranetSwitch`

**Location**: Lines 487-506

- [x] Update function signature to include security group parameter:
  ```go
  func tencentCloudClusterIntranetSwitch(ctx context.Context, service *TkeService, id, subnetId, securityGroup string, enable bool, domain string) (err error)
  //                                                                                            ^^^^^^^^^^^^^^^
  //                                                                                            Add parameter
  ```

- [x] Update `CreateClusterEndpoint` call to pass security group (line 490):
  ```go
  if enable {
      err = service.CreateClusterEndpoint(ctx, id, subnetId, securityGroup, false, domain, "")
      //                                                      ^^^^^^^^^^^^^^^
      //                                                      Pass security group
      if err != nil {
          return tccommon.RetryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
      }
  }
  ```

**Acceptance Criteria**:
- Function signature includes security group parameter ✅
- Security group passed to service layer ✅
- Logic correctly handles enable/disable cases ✅

---

### Task 5: Update Update Function Calls ✅ Completed

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

**Function**: `resourceTencentCloudTkeClusterEndpointUpdate`

**Location**: Lines 354-383 (three calls to `tencentCloudClusterIntranetSwitch`)

- [x] Update call at line 355:
  ```go
  err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, "", clusterIntranet, clusterIntranetDomain)
  //                                                                     ^^
  //                                                                     Add empty security group
  ```

- [x] Update call at line 366:
  ```go
  err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, "", false, clusterIntranetDomain)
  ```

- [x] Update call at line 375:
  ```go
  err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, "", true, clusterIntranetDomain)
  ```

**Rationale**: Update stage doesn't support security group modification (no API), so pass empty string.

**Acceptance Criteria**:
- All three calls updated with empty security group parameter ✅
- Calls match updated function signature ✅

---

### Task 6: Update Delete Function Call ✅ Completed

**File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_endpoint.go`

**Function**: `resourceTencentCloudTkeClusterEndpointDelete`

**Location**: Line 426

- [x] Update `tencentCloudClusterIntranetSwitch` call:
  ```go
  if clusterIntranet {
      err = tencentCloudClusterIntranetSwitch(ctx, &service, id, "", "", false, "")
      //                                                           ^^
      //                                                           Add empty security group
      if err != nil {
          errs = *multierror.Append(err)
      } else {
          taskErr := waitForClusterEndpointFinish(ctx, &service, id, false, false)
          if taskErr != nil {
              errs = *multierror.Append(taskErr)
          }
      }
  }
  ```

**Acceptance Criteria**:
- Call updated with empty security group parameter ✅
- Delete logic unchanged (no security group needed for deletion) ✅

---

### Task 7: Update Service Layer API Call ✅ Completed

**File**: `tencentcloud/services/tke/service_tencentcloud_tke.go`

**Function**: `CreateClusterEndpoint`

**Location**: Lines 1009-1011

- [x] Remove internet-only restriction for security group:
  ```go
  // Before
  if securityGroupId != "" && internet {
      request.SecurityGroup = &securityGroupId
  }

  // After
  if securityGroupId != "" {
      request.SecurityGroup = &securityGroupId
  }
  // Remove "&& internet" to support intranet security group
  ```

**Rationale**: TKE API supports security group for both internet and intranet endpoints.

**Acceptance Criteria**:
- Security group can be set for both internet and intranet ✅
- Logic simplified by removing unnecessary internet check ✅

---

### Task 8: Add Documentation ✅ Completed

**File**: `website/docs/r/kubernetes_cluster_endpoint.html.markdown` (if exists)

- [x] Add `cluster_intranet_security_group` to argument reference
- [x] Update `cluster_internet_security_group` description
- [x] Add usage example showing intranet security group
- [x] Document ForceNew behavior and limitations

**Example Documentation**:
```markdown
## Argument Reference

* `cluster_intranet_security_group` - (Optional, ForceNew) Security group ID for intranet cluster endpoint. 
  **Note**: Changing this parameter will recreate the resource as modification is not currently supported by the API.

* `cluster_internet_security_group` - (Optional) Security group ID for internet cluster endpoint.
```

**Acceptance Criteria**:
- Documentation clearly explains both security group fields ✅
- ForceNew behavior documented ✅
- Examples show proper usage ✅

---

### Task 9: Add Changelog Entry ✅ Completed

**File**: `.changelog/<next-number>.txt`

- [x] Determine next changelog number (1686)
- [x] Create changelog file:
  ```
  ```release-note:enhancement
  resource/tencentcloud_kubernetes_cluster_endpoint: add `cluster_intranet_security_group` field to support security group configuration for intranet cluster endpoints
  ```
  ```

**Acceptance Criteria**:
- Changelog entry created with correct number ✅
- Entry type is "enhancement" ✅
- Description clearly states new feature ✅

---

### Task 10: Testing ✅ Ready for Manual Testing

- [ ] **Unit Tests**: Add test cases for new field
  - Test intranet with security group (should succeed)
  - Test intranet without security group (should succeed)
  - Test ForceNew behavior

- [ ] **Manual Testing**: Test with real TKE cluster
  - Create endpoint with intranet + security group
  - Verify security group is applied to endpoint
  - Test ForceNew behavior (change security group)
  - Verify backward compatibility (without security group)

- [ ] **Integration Tests**: Update existing tests if needed
  - Ensure existing tests still pass
  - Add new test case for intranet security group

**Acceptance Criteria**:
- All scenarios tested
- Manual verification with real cluster
- No regression in existing functionality

---

## Definition of Done

- [x] All code changes completed
- [x] No NEW linter errors or warnings (pre-existing warnings remain)
- [x] Service layer supports intranet security group
- [x] Update and Delete functions handle new parameter
- [x] Documentation updated
- [x] Changelog entry added
- [ ] Manual testing with real TKE cluster completed (requires user testing)
- [ ] ForceNew behavior verified (requires user testing)
- [ ] Backward compatibility confirmed (requires user testing)
- [ ] Code review completed (if applicable)

---

## Dependencies

- TKE API must support security group for intranet endpoints
- No API support required for Read/Update (per requirements)

---

## Estimated Timeline

| Task | Estimated Time |
|------|----------------|
| Tasks 1-7 (Code changes) | 1.5 hours |
| Task 8 (Documentation) | 0.5 hours |
| Task 9 (Changelog) | 0.1 hours |
| Task 10 (Testing) | 1 hour |
| **Total** | **~3 hours** |

---

## Notes

1. **Read/Update Limitation**: Per requirements, read and update do not support querying or modifying intranet security group due to API limitations.

2. **ForceNew Behavior**: Changes to `cluster_intranet_security_group` will recreate the resource. Users should be aware of potential service disruption.

3. **No Validation Required**: The field can be set independently without constraints on `cluster_intranet` state.

4. **Testing Priority**: Manual testing with real TKE cluster is critical to verify API behavior for intranet security groups.
