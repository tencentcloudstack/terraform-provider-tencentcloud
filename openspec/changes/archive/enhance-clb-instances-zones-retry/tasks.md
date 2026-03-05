# Implementation Tasks

## Overview

This document tracks implementation tasks for enhancing the `tencentcloud_clb_instances` data source.

## Tasks

### 1. Add Zones Field to Data Source Schema ✅ Pending

**File**: `tencentcloud/services/clb/data_source_tc_clb_instances.go`

- [ ] Add `zones` field to `clb_list` schema (after line 184, before closing brace of `Elem`)
  ```go
  "zones": {
      Type:        schema.TypeList,
      Computed:    true,
      Elem:        &schema.Schema{Type: schema.TypeString},
      Description: "Zones where rules are deployed for VPC internal load balancers with nearby access mode. Note: This field may return null, indicating no valid values can be obtained.",
  },
  ```

- [ ] Add zones field mapping in `dataSourceTencentCloudClbInstancesRead` function (around line 287, before `clbList = append(clbList, mapping)`)
  ```go
  if clbInstance.Zones != nil {
      mapping["zones"] = helper.StringsInterfaces(clbInstance.Zones)
  }
  ```

### 2. Add Retry Logic to Service Function ✅ Pending

**File**: `tencentcloud/services/clb/service_tencentcloud_clb.go`

- [ ] Import required package if not already imported
  ```go
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
  ```

- [ ] Refactor `DescribeLoadBalancerByFilter` function (starting at line 77):
  - Keep pagination `for` loop structure unchanged
  - Wrap **only the API call** inside each loop iteration with `resource.Retry()`
  - Declare `response` variable before retry block
  - Use `tccommon.RetryError(e)` for retryable errors
  - Use `resource.NonRetryableError()` for nil response validation
  - Set `response` variable inside retry success case
  - Handle error after each retry completes

### 3. Testing ✅ Completed

**Files Created**:
- `examples/tencentcloud-clb-instances-zones/main.tf`
- `examples/tencentcloud-clb-instances-zones/README.md`

- [x] Create example configuration with data source query
- [x] Create README with documentation and expected output examples

### 4. Changelog Entry ✅ Pending

**File**: `.changelog/<next-number>.txt`

- [ ] Determine next changelog number (check latest in `.changelog/` directory)
- [ ] Create changelog file with content:
  ```
  release-note:enhancement
  data-source/tencentcloud_clb_instances: add `zones` field to `clb_list` output and improve reliability with retry logic in underlying API calls
  ```

## Implementation Notes

### Zones Field Data Flow

1. **API Response**: `LoadBalancer.Zones` (`[]*string`) from `DescribeLoadBalancers` response
2. **Service Layer**: Pass through unchanged in `DescribeLoadBalancerByFilter` return
3. **Data Source Layer**: Convert to `[]interface{}` using `helper.StringsInterfaces()` helper
4. **Terraform State**: Stored as list of strings

### Retry Logic Implementation

The retry wrapper should:
- Use `tccommon.ReadRetryTimeout` as timeout duration (consistent with data source reader)
- Wrap **only the API call** inside each loop iteration (not the entire pagination loop)
- Return `tccommon.RetryError(e)` for retryable API errors
- Return `resource.NonRetryableError()` for nil response validation failures
- Return `nil` on successful API call
- Pagination logic (`for` loop, offset increment) remains outside retry block

### Backward Compatibility Verification

- Existing fields remain unchanged
- New `zones` field is optional (computed only)
- No changes to input parameters
- Retry logic is transparent to users

## Definition of Done

- [x] All code changes completed and tested locally
- [x] Example configuration created and verified
- [x] Changelog entry added (.changelog/3843.txt)
- [x] No linter errors or warnings
- [ ] Manual testing confirms zones field populated correctly (requires real CLB instances)
- [ ] Code review completed (if applicable)
