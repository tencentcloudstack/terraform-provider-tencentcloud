# Proposal: Add TimeZone and DiskEncryptFlag Support to SqlServer Basic Instance

## Metadata
- **Change ID**: `add-sqlserver-timezone-diskencrypt`
- **Status**: Proposal
- **Created**: 2026-02-06
- **Author**: AI Assistant
- **Type**: Enhancement
- **Estimated Effort**: 2-3 hours

---

## Problem Statement

The `tencentcloud_sqlserver_basic_instance` resource currently does not expose two important configuration parameters:

1. **TimeZone**: The system timezone setting for the SQL Server instance (e.g., "China Standard Time", "UTC")
2. **DiskEncryptFlag**: Whether disk encryption is enabled for the instance (0=disabled, 1=enabled)

### Current Limitations

- **TimeZone**: 
  - Available in `CreateBasicDBInstances` API request parameter
  - Returned in `DescribeDBInstances` API response (in `DBInstance.TimeZone` field)
  - **Not exposed** in the Terraform resource schema
  - Users cannot set or view the timezone configuration

- **DiskEncryptFlag**:
  - Available in `CreateBasicDBInstances` API request parameter
  - **NOT** returned in `DescribeDBInstances` API response
  - Must be retrieved via separate `DescribeDBInstancesAttribute` API call (field: `IsDiskEncryptFlag`)
  - **Not exposed** in the Terraform resource schema
  - Users cannot set or view disk encryption status

### User Impact

Without these parameters:
- Users cannot control timezone settings during instance creation
- Users cannot enable disk encryption for compliance requirements
- Users cannot view current timezone or encryption status in Terraform state
- Infrastructure-as-code configurations are incomplete

---

## Proposed Solution

Add two new optional parameters to `tencentcloud_sqlserver_basic_instance`:

### 1. `time_zone` (Optional, ForceNew, Computed)
- **Type**: String
- **Default**: "China Standard Time" (API default)
- **ForceNew**: Yes (requires instance recreation to change)
- **Computed**: Yes (read from API if not set)
- **Validation**: None (API will validate)
- **Description**: System timezone for the SQL Server instance
- **Examples**: "China Standard Time", "UTC", "Eastern Standard Time"

### 2. `disk_encrypt_flag` (Optional, ForceNew, Computed)
- **Type**: Int
- **Allowed Values**: 0 (disabled), 1 (enabled)
- **Default**: 0 (API default)
- **ForceNew**: Yes (encryption cannot be changed after creation)
- **Computed**: Yes (read from API if not set)
- **Validation**: ValidateIntegerInRange(0, 1)
- **Description**: Whether disk encryption is enabled (0=disabled, 1=enabled)

---

## Technical Design

### API Integration

#### Create Operation
Both parameters will be passed to `CreateBasicDBInstances`:
```go
request := sqlserver.NewCreateBasicDBInstancesRequest()
// ... other parameters ...

if v, ok := paramMap["time_zone"]; ok {
    request.TimeZone = helper.String(v.(string))
}

if v, ok := paramMap["disk_encrypt_flag"]; ok {
    request.DiskEncryptFlag = helper.IntInt64(v.(int))
}
```

#### Read Operation
Two different API calls needed:

**For TimeZone**:
```go
// Already called in existing Read function
instance, err := DescribeDBInstances(instanceId)
timeZone := instance.TimeZone  // Direct access
```

**For DiskEncryptFlag**:
```go
// NEW: Additional API call required
attribute, err := DescribeDBInstancesAttribute(instanceId)
diskEncryptFlag := attribute.IsDiskEncryptFlag  // Convert *int64 to int
```

### Schema Changes

Add to `resource_tc_sqlserver_basic_instance.go` Schema:

```go
"time_zone": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    ForceNew:    true,
    Description: "System timezone for the SQL Server instance. Default is `China Standard Time`. Common values: `China Standard Time`, `UTC`. This setting cannot be changed after creation.",
},
"disk_encrypt_flag": {
    Type:         schema.TypeInt,
    Optional:     true,
    Computed:     true,
    ForceNew:     true,
    Default:      0,
    ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
    Description:  "Disk encryption flag. `0` - Disabled (default), `1` - Enabled. Disk encryption cannot be changed after instance creation.",
},
```

### Service Layer Changes

**New Method Required** in `service_tencentcloud_sqlserver.go`:

```go
func (me *SqlserverService) DescribeSqlserverInstanceAttributeById(ctx context.Context, instanceId string) (
    attribute *sqlserver.DescribeDBInstancesAttributeResponseParams, 
    errRet error,
) {
    logId := tccommon.GetLogId(ctx)
    request := sqlserver.NewDescribeDBInstancesAttributeRequest()
    request.InstanceId = &instanceId
    
    ratelimit.Check(request.GetAction())
    response, err := me.client.UseSqlserverClient().DescribeDBInstancesAttribute(request)
    if err != nil {
        log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
            logId, request.GetAction(), request.ToJsonString(), err.Error())
        errRet = err
        return
    }
    
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
        logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
    
    if response.Response != nil {
        attribute = response.Response
    }
    
    return
}
```

### Implementation Files

1. **Resource File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`
   - Add schema fields
   - Modify `CreateSqlserverBasicInstance` to pass new parameters
   - Modify `resourceTencentCloudSqlserverBasicInstanceRead` to read both fields

2. **Service File**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`
   - Add `DescribeSqlserverInstanceAttributeById` method
   - Modify `CreateSqlserverBasicInstance` to accept new parameters

3. **Documentation**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`
   - Add usage examples
   - Document both parameters

4. **Test File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go`
   - Add test cases for timezone setting
   - Add test cases for disk encryption

5. **Website Documentation**: Auto-generated via `make doc`

---

## Implementation Plan

### Phase 1: Schema Definition (30 min)
1. Add `time_zone` schema field with ForceNew, Computed, Optional
2. Add `disk_encrypt_flag` schema field with ForceNew, Computed, Optional, Default
3. Add validation for `disk_encrypt_flag` (0-1 range)
4. Update immutableArgs list (both are ForceNew)

### Phase 2: Service Layer (45 min)
1. Create `DescribeSqlserverInstanceAttributeById` method
2. Add error handling and logging
3. Add rate limiting
4. Modify `CreateSqlserverBasicInstance` to accept `time_zone` and `disk_encrypt_flag`
5. Pass both parameters to `CreateBasicDBInstancesRequest`

### Phase 3: Create Operation (15 min)
1. Read `time_zone` from schema in Create function
2. Read `disk_encrypt_flag` from schema in Create function
3. Add both to `paramMap` if present
4. Pass to service layer

### Phase 4: Read Operation (30 min)
1. Read `TimeZone` from existing `DescribeDBInstances` response
2. Call new `DescribeSqlserverInstanceAttributeById` method
3. Read `IsDiskEncryptFlag` from attribute response
4. Convert `*int64` to `int` for disk_encrypt_flag
5. Handle nil pointers safely
6. Set both values in Terraform state

### Phase 5: Testing (45 min)
1. Add test case: Basic instance with custom timezone
2. Add test case: Basic instance with disk encryption enabled
3. Add test case: Basic instance with both parameters
4. Add test case: Import existing instance (verify read)
5. Verify ForceNew behavior (plan shows recreation on change)

### Phase 6: Documentation (30 min)
1. Update resource documentation with usage examples
2. Document timezone values (common timezones)
3. Document disk_encrypt_flag values and implications
4. Add note about ForceNew behavior
5. Generate website documentation via `make doc`

### Phase 7: Code Quality & Validation (15 min)
1. Run `gofmt` and `goimports`
2. Run `make lint`
3. Compile provider: `go build`
4. Run tests: `go test`
5. Final review

---

## Impact Analysis

### Breaking Changes
**None** - Both parameters are optional with sensible defaults.

### Backward Compatibility
✅ **Fully backward compatible**:
- Existing configurations without these parameters will continue to work
- Both fields are `Computed`, so existing state will be populated on next `terraform apply`
- Default values match API defaults

### API Calls
**Additional API Call in Read**:
- Current: 1 call to `DescribeDBInstances`
- After: 2 calls (`DescribeDBInstances` + `DescribeDBInstancesAttribute`)
- Impact: Minimal (only affects Read operation, not Create/Update/Delete)
- Mitigation: Both calls are necessary to provide complete resource information

---

## Alternative Approaches

### Alternative 1: Skip DiskEncryptFlag (Not Recommended)
- **Pro**: Avoids extra API call
- **Con**: Incomplete resource representation, no way to verify encryption status
- **Verdict**: Rejected - encryption status is important for compliance

### Alternative 2: Make Parameters Required (Not Recommended)
- **Pro**: Forces users to make explicit choice
- **Con**: Breaking change for existing users
- **Verdict**: Rejected - backward compatibility is important

### Alternative 3: Use Separate Data Source for Attributes (Not Recommended)
- **Pro**: Separates concerns, no extra call in resource Read
- **Con**: Poor user experience, requires two resources to get complete information
- **Verdict**: Rejected - attributes belong in resource

---

## Success Criteria

### Functional Requirements
- ✅ Users can set `time_zone` during instance creation
- ✅ Users can enable `disk_encrypt_flag` during instance creation
- ✅ Both values are visible in Terraform state after creation
- ✅ Both values are visible in `terraform plan` output
- ✅ Changing either value triggers instance recreation (ForceNew)
- ✅ Import existing instances populates both fields correctly
- ✅ Omitting both parameters uses API defaults

### Non-Functional Requirements
- ✅ Code follows project conventions (file naming, error handling, logging)
- ✅ Tests pass (unit + acceptance)
- ✅ Documentation is complete and clear
- ✅ No linter errors or warnings
- ✅ Backward compatible with existing configurations

---

## Risks and Mitigation

### Risk 1: Extra API Call Performance
- **Impact**: Low - only affects Read operation
- **Mitigation**: Call is necessary for complete information, caching not needed

### Risk 2: API Field Changes
- **Impact**: Medium - if TencentCloud renames fields
- **Mitigation**: Use official SDK types, monitor API documentation

### Risk 3: Nil Pointer Dereference
- **Impact**: High - could crash provider
- **Mitigation**: Proper nil checks for all pointer fields

---

## Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| 1. Schema Definition | 4 tasks | 30 min |
| 2. Service Layer | 5 tasks | 45 min |
| 3. Create Operation | 4 tasks | 15 min |
| 4. Read Operation | 6 tasks | 30 min |
| 5. Testing | 5 tasks | 45 min |
| 6. Documentation | 5 tasks | 30 min |
| 7. Code Quality | 5 tasks | 15 min |
| **Total** | **34 tasks** | **3.5 hours** |

---

## Related Documentation

- **TencentCloud API**:
  - CreateBasicDBInstances: https://cloud.tencent.com/document/api/238/50262
  - DescribeDBInstances: https://cloud.tencent.com/document/api/238/19969
  - DescribeDBInstancesAttribute: https://cloud.tencent.com/document/api/238/73319

- **SDK Source**:
  - Models: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328/models.go`
  - Client: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328/client.go`

- **Existing Implementation**:
  - Resource: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`
  - Service: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

---

## Appendix: Code Examples

### Example 1: Basic Instance with Custom Timezone

```hcl
resource "tencentcloud_sqlserver_basic_instance" "example" {
  name              = "example-sqlserver"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_PREMIUM"
  cpu               = 2
  memory            = 4
  storage           = 20
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2008R2"
  
  # New parameters
  time_zone         = "UTC"
  disk_encrypt_flag = 0
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```

### Example 2: Basic Instance with Disk Encryption

```hcl
resource "tencentcloud_sqlserver_basic_instance" "encrypted" {
  name              = "encrypted-sqlserver"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_SSD"
  cpu               = 4
  memory            = 8
  storage           = 50
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2016SP1"
  
  # Enable disk encryption
  disk_encrypt_flag = 1
  time_zone         = "China Standard Time"
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```

### Example 3: Import Existing Instance

```bash
# Import will populate both time_zone and disk_encrypt_flag from API
terraform import tencentcloud_sqlserver_basic_instance.example mssql-abc123
```

---

## Approval

This proposal is ready for review and implementation.

**Next Steps**:
1. Review and approve this proposal
2. Run `openspec apply add-sqlserver-timezone-diskencrypt` to begin implementation
3. Follow the task list in `tasks.md`
4. Complete all phases sequentially
5. Submit PR after validation

---

**END OF PROPOSAL**
