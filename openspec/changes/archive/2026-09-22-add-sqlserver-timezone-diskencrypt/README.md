# Add TimeZone and DiskEncryptFlag to SqlServer Basic Instance

**Change ID**: `add-sqlserver-timezone-diskencrypt`  
**Status**: Proposal  
**Type**: Resource Enhancement  
**Priority**: Medium

---

## Quick Summary

Add two new optional parameters to `tencentcloud_sqlserver_basic_instance`:
- **`time_zone`**: System timezone configuration (e.g., "China Standard Time", "UTC")
- **`disk_encrypt_flag`**: Disk encryption enable/disable (0=disabled, 1=enabled)

Both parameters are:
- ✅ Optional (backward compatible)
- ✅ Computed (populated from API)
- ✅ ForceNew (require instance recreation to change)

---

## Problem

Users cannot:
1. Configure SQL Server instance timezone during creation
2. Enable disk encryption for compliance requirements
3. View current timezone or encryption status in Terraform state

---

## Solution

### New Schema Fields

```hcl
resource "tencentcloud_sqlserver_basic_instance" "example" {
  name              = "my-sqlserver"
  # ... existing fields ...
  
  # NEW: Timezone configuration
  time_zone         = "UTC"
  
  # NEW: Disk encryption
  disk_encrypt_flag = 1  # 0=disabled (default), 1=enabled
}
```

### API Integration

**Create**: Both parameters passed to `CreateBasicDBInstances`

**Read**: 
- `time_zone` from `DescribeDBInstances` (existing call)
- `disk_encrypt_flag` from `DescribeDBInstancesAttribute` (new call)

---

## Key Features

### 1. TimeZone Support
- Control SQL Server system timezone
- Common values: "China Standard Time" (default), "UTC", "Eastern Standard Time"
- Retrieved directly from `DescribeDBInstances` API
- Cannot be changed after creation (ForceNew)

### 2. Disk Encryption Support
- Enable/disable disk encryption (0 or 1)
- Important for compliance and security
- Retrieved from separate `DescribeDBInstancesAttribute` API call
- Cannot be changed after creation (ForceNew)

### 3. Backward Compatibility
- Both parameters are optional
- Existing configurations work unchanged
- Computed values populated on refresh

---

## Usage Examples

### Example 1: Custom Timezone

```hcl
resource "tencentcloud_sqlserver_basic_instance" "utc" {
  name              = "sqlserver-utc"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_PREMIUM"
  cpu               = 2
  memory            = 4
  storage           = 20
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2008R2"
  
  time_zone = "UTC"
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```

### Example 2: Disk Encryption

```hcl
resource "tencentcloud_sqlserver_basic_instance" "encrypted" {
  name              = "sqlserver-encrypted"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_SSD"
  cpu               = 4
  memory            = 8
  storage           = 50
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2016SP1"
  
  disk_encrypt_flag = 1  # Enable encryption
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```

### Example 3: Both Parameters

```hcl
resource "tencentcloud_sqlserver_basic_instance" "full_config" {
  name              = "sqlserver-full"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_SSD"
  cpu               = 4
  memory            = 8
  storage           = 50
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2016SP1"
  
  time_zone         = "UTC"
  disk_encrypt_flag = 1
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```

### Example 4: Import Existing Instance

```bash
terraform import tencentcloud_sqlserver_basic_instance.example mssql-abc123
```

After import, both `time_zone` and `disk_encrypt_flag` will be populated from API.

---

## Implementation Details

### Files Modified

1. **`tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`**
   - Add schema fields
   - Update Create to pass parameters
   - Update Read to fetch both values

2. **`tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`**
   - Add `DescribeSqlserverInstanceAttributeById` method
   - Update `CreateSqlserverBasicInstance` to accept new parameters

3. **`tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go`**
   - Add 3 new acceptance tests

4. **`tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`**
   - Add usage examples
   - Document timezone values
   - Document encryption implications

5. **`website/docs/r/sqlserver_basic_instance.html.markdown`**
   - Auto-generated via `make doc`

---

## Technical Challenges

### Challenge 1: Separate API for Encryption Status
**Issue**: `disk_encrypt_flag` not in `DescribeDBInstances` response  
**Solution**: Call `DescribeDBInstancesAttribute` separately in Read operation  
**Impact**: One additional API call per Read operation

### Challenge 2: Type Conversion
**Issue**: API uses `*int64` but schema uses `int`  
**Solution**: Proper type conversion with nil checks  
**Code**: `int(*attribute.IsDiskEncryptFlag)`

### Challenge 3: ForceNew Behavior
**Issue**: Both parameters cannot be changed after creation  
**Solution**: Add to immutableArgs list + set ForceNew=true in schema  
**Validation**: Terraform enforces recreation on change

---

## Testing Strategy

### Acceptance Tests (3 tests)

1. **`TestAccTencentCloudSqlserverBasicInstance_Timezone`**
   - Create instance with custom timezone
   - Verify state contains correct value
   - Test import functionality

2. **`TestAccTencentCloudSqlserverBasicInstance_DiskEncrypt`**
   - Create instance with encryption enabled
   - Verify disk_encrypt_flag=1
   - Test import functionality

3. **`TestAccTencentCloudSqlserverBasicInstance_TimezoneAndEncrypt`**
   - Create instance with both parameters
   - Verify both values in state
   - Test import functionality

### Manual Tests

- Import existing instances and verify field population
- Test ForceNew behavior with `terraform plan`
- Test various timezone values
- Test invalid disk_encrypt_flag values (validation)

---

## Success Criteria

- ✅ Users can set `time_zone` during creation
- ✅ Users can enable `disk_encrypt_flag` during creation
- ✅ Both values visible in Terraform state
- ✅ Both values visible in `terraform plan`
- ✅ Changing either value triggers recreation (ForceNew)
- ✅ Import populates both fields correctly
- ✅ Omitting parameters uses API defaults
- ✅ Backward compatible (no breaking changes)
- ✅ Tests pass
- ✅ Documentation complete

---

## Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| 1. Schema Definition | 4 tasks | 30 min |
| 2. Service Layer | 5 tasks | 45 min |
| 3. Create Operation | 4 tasks | 15 min |
| 4. Read Operation | 6 tasks | 30 min |
| 5. Testing | 5 tasks | 45 min |
| 6. Documentation | 8 tasks | 30 min |
| 7. Code Quality | 5 tasks | 15 min |
| **Total** | **37 tasks** | **3.5 hours** |

---

## API References

### CreateBasicDBInstances
https://cloud.tencent.com/document/api/238/50262

**Parameters**:
- `TimeZone`: System timezone
- `DiskEncryptFlag`: Disk encryption flag

### DescribeDBInstances
https://cloud.tencent.com/document/api/238/19969

**Response**:
- `DBInstance.TimeZone`: Instance timezone

### DescribeDBInstancesAttribute
https://cloud.tencent.com/document/api/238/73319

**Response**:
- `IsDiskEncryptFlag`: Disk encryption status (0 or 1)

---

## Risk Assessment

### Risk 1: Extra API Call
**Impact**: Low  
**Mitigation**: Only one extra call per Read, acceptable performance impact

### Risk 2: API Field Changes
**Impact**: Medium  
**Mitigation**: Use official SDK, monitor API updates

### Risk 3: Nil Pointer Issues
**Impact**: High (could crash provider)  
**Mitigation**: Proper nil checks for all pointer fields

---

## Next Steps

1. **Review Proposal**: Read `proposal.md` for detailed design
2. **Review Tasks**: Read `tasks.md` for implementation steps
3. **Review Spec**: Read `specs/sqlserver-basic-instance-timezone-encrypt/spec.md` for requirements
4. **Approve**: Approve the proposal if design is acceptable
5. **Implement**: Run `openspec apply add-sqlserver-timezone-diskencrypt` to begin

---

## Documentation

- **Proposal**: `proposal.md` (detailed design, 15+ pages)
- **Tasks**: `tasks.md` (34 implementation tasks)
- **Spec**: `specs/sqlserver-basic-instance-timezone-encrypt/spec.md` (9 requirements, 24 scenarios)
- **This File**: `README.md` (quick reference)

---

## Questions & Answers

### Q: Why ForceNew for both parameters?
**A**: TencentCloud API does not support modifying timezone or encryption after instance creation. These are instance-level settings that are immutable.

### Q: Why separate API call for disk_encrypt_flag?
**A**: The `DescribeDBInstances` API does not return encryption status. TencentCloud requires a separate `DescribeDBInstancesAttribute` call to retrieve this information.

### Q: What happens to existing configurations?
**A**: No changes needed. Both parameters are optional and will be auto-populated from API on next refresh.

### Q: What if the DescribeDBInstancesAttribute API fails?
**A**: A warning is logged, but the Read operation continues. Other fields are still populated normally.

### Q: Can I change timezone after creation?
**A**: No. Changing timezone or encryption requires destroying and recreating the instance. This is enforced by ForceNew.

---

**Status**: ✅ Proposal Complete - Ready for Review

To begin implementation, run:
```bash
openspec apply add-sqlserver-timezone-diskencrypt
```
