# Spec: SqlServer Basic Instance - TimeZone and DiskEncryptFlag Support

**Capability ID**: `sqlserver-basic-instance-timezone-encrypt`  
**Change ID**: `add-sqlserver-timezone-diskencrypt`  
**Type**: Resource Enhancement  
**Priority**: Medium

---

## Overview

This specification defines the requirements for adding `time_zone` and `disk_encrypt_flag` parameters to the `tencentcloud_sqlserver_basic_instance` resource, enabling users to configure timezone settings and disk encryption for SQL Server instances.

---

## ADDED Requirements

### Requirement 1: TimeZone Parameter Support

**ID**: REQ-SQLSERVER-TIMEZONE-001  
**Priority**: High  
**Status**: New

The `tencentcloud_sqlserver_basic_instance` resource MUST support a `time_zone` parameter that allows users to specify the system timezone for SQL Server instances.

#### Scenario 1.1: Create Instance with Custom Timezone

**Given** a user wants to create a SQL Server basic instance with a specific timezone  
**When** the user specifies `time_zone = "UTC"` in the resource configuration  
**Then** the instance MUST be created with UTC timezone  
**And** the timezone MUST be retrievable via `DescribeDBInstances` API  
**And** the Terraform state MUST show `time_zone = "UTC"`

**Test Data**:
```hcl
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name         = "test-timezone"
  time_zone    = "UTC"
  # ... other required fields
}
```

**Expected API Request**:
```json
{
  "TimeZone": "UTC",
  // ... other parameters
}
```

**Validation**:
- ✅ TimeZone parameter sent in CreateBasicDBInstances request
- ✅ Instance created successfully
- ✅ State contains time_zone = "UTC"

---

#### Scenario 1.2: Create Instance with Default Timezone

**Given** a user creates a SQL Server instance without specifying `time_zone`  
**When** the resource is created  
**Then** the API default timezone MUST be used ("China Standard Time")  
**And** the Terraform state MUST be populated with the actual timezone from API  
**And** no error MUST occur

**Test Data**:
```hcl
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name = "test-default-tz"
  # time_zone not specified
}
```

**Expected Behavior**:
- API receives no TimeZone parameter (uses default)
- Read operation populates state with "China Standard Time"
- User can see the actual timezone value

**Validation**:
- ✅ Instance created with API default
- ✅ State shows actual timezone from API
- ✅ Computed value works correctly

---

#### Scenario 1.3: Read Existing Instance Timezone

**Given** an existing SQL Server instance with timezone set  
**When** Terraform reads the instance state  
**Then** the `time_zone` MUST be retrieved from `DescribeDBInstances` response  
**And** the value MUST match the instance configuration  
**And** the state MUST be updated correctly

**API Response**:
```json
{
  "DBInstance": {
    "InstanceId": "mssql-abc123",
    "TimeZone": "UTC",
    // ... other fields
  }
}
```

**Expected State**:
```
time_zone = "UTC"
```

**Validation**:
- ✅ DBInstance.TimeZone field accessed
- ✅ Nil pointer check performed
- ✅ State updated with correct value

---

#### Scenario 1.4: Timezone Change Triggers Recreation

**Given** an existing SQL Server instance with `time_zone = "China Standard Time"`  
**When** the user changes `time_zone` to "UTC" in configuration  
**Then** Terraform plan MUST show the instance will be destroyed and recreated  
**And** the plan MUST indicate `time_zone` change forces replacement  
**And** applying the plan MUST recreate the instance with new timezone

**Terraform Plan Output**:
```
-/+ tencentcloud_sqlserver_basic_instance.test (forces new resource)
      time_zone: "China Standard Time" -> "UTC"
```

**Validation**:
- ✅ ForceNew behavior triggered
- ✅ Instance destroyed
- ✅ New instance created with new timezone
- ✅ State updated correctly

---

#### Scenario 1.5: Import Populates Timezone

**Given** an existing SQL Server instance not managed by Terraform  
**When** the user imports the instance with `terraform import`  
**Then** the Read function MUST be called  
**And** the `time_zone` MUST be populated from API  
**And** the state MUST contain the correct timezone value

**Import Command**:
```bash
terraform import tencentcloud_sqlserver_basic_instance.example mssql-abc123
```

**Expected Behavior**:
- Read function queries DescribeDBInstances
- TimeZone field extracted and stored
- State file contains time_zone value

**Validation**:
- ✅ Import succeeds
- ✅ time_zone populated
- ✅ Subsequent plan shows no changes

---

#### Scenario 1.6: Timezone Validation

**Given** a user specifies a timezone value  
**When** the resource is created  
**Then** the API MUST validate the timezone string  
**And** invalid timezones MUST result in API error  
**And** error MUST be surfaced to the user

**Invalid Example**:
```hcl
time_zone = "Invalid/Timezone"
```

**Expected Behavior**:
- Terraform sends request to API
- API returns validation error
- Error propagated to user

**Validation**:
- ✅ No client-side validation (API handles it)
- ✅ API errors properly surfaced
- ✅ User sees clear error message

---

### Requirement 2: Disk Encryption Flag Support

**ID**: REQ-SQLSERVER-DISKENCRYPT-001  
**Priority**: High  
**Status**: New

The `tencentcloud_sqlserver_basic_instance` resource MUST support a `disk_encrypt_flag` parameter that allows users to enable disk encryption for SQL Server instances.

#### Scenario 2.1: Create Instance with Encryption Enabled

**Given** a user wants to create a SQL Server instance with disk encryption  
**When** the user specifies `disk_encrypt_flag = 1` in the configuration  
**Then** the instance MUST be created with encryption enabled  
**And** the encryption status MUST be retrievable via `DescribeDBInstancesAttribute` API  
**And** the Terraform state MUST show `disk_encrypt_flag = 1`

**Test Data**:
```hcl
resource "tencentcloud_sqlserver_basic_instance" "encrypted" {
  name              = "test-encrypted"
  disk_encrypt_flag = 1
  # ... other required fields
}
```

**Expected API Request**:
```json
{
  "DiskEncryptFlag": 1,
  // ... other parameters
}
```

**Validation**:
- ✅ DiskEncryptFlag=1 sent in CreateBasicDBInstances
- ✅ Instance created successfully
- ✅ State contains disk_encrypt_flag = 1

---

#### Scenario 2.2: Create Instance with Encryption Disabled (Default)

**Given** a user creates a SQL Server instance without specifying `disk_encrypt_flag`  
**When** the resource is created  
**Then** encryption MUST be disabled (default value 0)  
**And** the Terraform state MUST show `disk_encrypt_flag = 0`  
**And** no error MUST occur

**Test Data**:
```hcl
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name = "test-no-encrypt"
  # disk_encrypt_flag not specified (uses default 0)
}
```

**Expected API Request**:
```json
{
  "DiskEncryptFlag": 0,
  // ... other parameters
}
```

**Validation**:
- ✅ Default value 0 used
- ✅ Instance created without encryption
- ✅ State shows disk_encrypt_flag = 0

---

#### Scenario 2.3: Read Existing Instance Encryption Status

**Given** an existing SQL Server instance with encryption enabled  
**When** Terraform reads the instance state  
**Then** the `disk_encrypt_flag` MUST be retrieved via `DescribeDBInstancesAttribute` API  
**And** the `IsDiskEncryptFlag` field MUST be accessed  
**And** the value MUST be converted from *int64 to int  
**And** the state MUST be updated correctly

**API Call Sequence**:
1. `DescribeDBInstances` (for basic info)
2. `DescribeDBInstancesAttribute` (for encryption status)

**API Response**:
```json
{
  "InstanceId": "mssql-abc123",
  "IsDiskEncryptFlag": 1,
  // ... other fields
}
```

**Expected State**:
```
disk_encrypt_flag = 1
```

**Validation**:
- ✅ Separate API call made
- ✅ IsDiskEncryptFlag field accessed
- ✅ Type conversion performed (*int64 -> int)
- ✅ Nil pointer check performed
- ✅ State updated correctly

---

#### Scenario 2.4: Read Operation Handles API Failure Gracefully

**Given** an existing SQL Server instance  
**When** the `DescribeDBInstancesAttribute` API call fails  
**Then** a warning MUST be logged  
**And** the Read operation MUST NOT fail entirely  
**And** other fields MUST still be populated  
**And** `disk_encrypt_flag` MAY be empty in state

**Expected Behavior**:
- API call fails (network issue, permissions, etc.)
- Warning logged: "describe sqlserver instance attribute failed"
- Read operation continues
- Other fields populated normally

**Validation**:
- ✅ Error caught and logged
- ✅ Read doesn't fail
- ✅ User sees warning in logs
- ✅ Other fields work

---

#### Scenario 2.5: Encryption Change Triggers Recreation

**Given** an existing SQL Server instance with `disk_encrypt_flag = 0`  
**When** the user changes `disk_encrypt_flag` to 1 in configuration  
**Then** Terraform plan MUST show the instance will be destroyed and recreated  
**And** the plan MUST indicate `disk_encrypt_flag` change forces replacement  
**And** applying the plan MUST recreate the instance with encryption enabled

**Terraform Plan Output**:
```
-/+ tencentcloud_sqlserver_basic_instance.test (forces new resource)
      disk_encrypt_flag: 0 -> 1
```

**Validation**:
- ✅ ForceNew behavior triggered
- ✅ Instance destroyed
- ✅ New instance created with encryption
- ✅ State updated correctly

---

#### Scenario 2.6: Import Populates Encryption Status

**Given** an existing SQL Server instance with encryption enabled  
**When** the user imports the instance  
**Then** the Read function MUST call both APIs  
**And** the `disk_encrypt_flag` MUST be populated from `DescribeDBInstancesAttribute`  
**And** the state MUST contain the correct encryption status

**Import Command**:
```bash
terraform import tencentcloud_sqlserver_basic_instance.example mssql-abc123
```

**Expected Behavior**:
- DescribeDBInstances called
- DescribeDBInstancesAttribute called
- IsDiskEncryptFlag extracted
- State contains disk_encrypt_flag = 1

**Validation**:
- ✅ Import succeeds
- ✅ disk_encrypt_flag populated correctly
- ✅ Subsequent plan shows no changes

---

#### Scenario 2.7: Encryption Flag Validation

**Given** a user specifies an invalid `disk_encrypt_flag` value  
**When** Terraform validates the configuration  
**Then** a validation error MUST occur  
**And** the error MUST indicate valid values are 0 or 1  
**And** the resource creation MUST NOT proceed

**Invalid Example**:
```hcl
disk_encrypt_flag = 2  # Invalid
```

**Expected Error**:
```
Error: expected disk_encrypt_flag to be in the range (0 - 1), got 2
```

**Validation**:
- ✅ ValidateFunc triggered
- ✅ Clear error message
- ✅ No API call made

---

### Requirement 3: Combined Parameter Support

**ID**: REQ-SQLSERVER-COMBINED-001  
**Priority**: Medium  
**Status**: New

The resource MUST support using both `time_zone` and `disk_encrypt_flag` parameters simultaneously without conflicts.

#### Scenario 3.1: Create Instance with Both Parameters

**Given** a user wants to specify both timezone and encryption  
**When** both parameters are provided in configuration  
**Then** the instance MUST be created with both settings  
**And** both values MUST be retrievable  
**And** the state MUST contain both values correctly

**Test Data**:
```hcl
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name              = "test-both"
  time_zone         = "UTC"
  disk_encrypt_flag = 1
  # ... other fields
}
```

**Expected API Request**:
```json
{
  "TimeZone": "UTC",
  "DiskEncryptFlag": 1,
  // ... other parameters
}
```

**Validation**:
- ✅ Both parameters sent in create request
- ✅ Instance created successfully
- ✅ Both values in state
- ✅ No conflicts or errors

---

#### Scenario 3.2: Change Either Parameter Triggers Recreation

**Given** an instance with both parameters set  
**When** either parameter is changed  
**Then** the instance MUST be recreated  
**And** the plan MUST show which parameter changed  
**And** both parameters MUST be set correctly on new instance

**Test Scenarios**:
1. Change time_zone only -> recreation
2. Change disk_encrypt_flag only -> recreation
3. Change both -> recreation

**Validation**:
- ✅ ForceNew works for both
- ✅ Independent of each other
- ✅ New instance gets correct values

---

### Requirement 4: Schema Design

**ID**: REQ-SQLSERVER-SCHEMA-001  
**Priority**: High  
**Status**: New

The schema definition MUST follow Terraform and project conventions.

#### Scenario 4.1: Schema Properties for time_zone

**Given** the schema definition for `time_zone`  
**Then** it MUST have the following properties:
- Type: TypeString
- Optional: true
- Computed: true
- ForceNew: true
- Description: Clear and helpful

**Schema Definition**:
```go
"time_zone": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    ForceNew:    true,
    Description: "System timezone for the SQL Server instance. Default is `China Standard Time`. Common values: `China Standard Time`, `UTC`, `Eastern Standard Time`. This setting cannot be changed after creation.",
},
```

**Validation**:
- ✅ All properties set correctly
- ✅ Description mentions common values
- ✅ ForceNew behavior documented

---

#### Scenario 4.2: Schema Properties for disk_encrypt_flag

**Given** the schema definition for `disk_encrypt_flag`  
**Then** it MUST have the following properties:
- Type: TypeInt
- Optional: true
- Computed: true
- ForceNew: true
- Default: 0
- ValidateFunc: IntegerInRange(0, 1)
- Description: Clear with value meanings

**Schema Definition**:
```go
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

**Validation**:
- ✅ All properties set correctly
- ✅ Default value defined
- ✅ Validation enforces 0-1 range
- ✅ Description explains values

---

#### Scenario 4.3: Immutable Args Configuration

**Given** both parameters are ForceNew  
**When** a user attempts to modify them  
**Then** the immutableArgs list MUST include both parameters  
**And** Terraform MUST prevent in-place updates  
**And** recreation MUST be required

**Implementation**:
```go
immutableArgs := []string{"collation", "time_zone", "disk_encrypt_flag"}
```

**Validation**:
- ✅ Both in immutableArgs list
- ✅ Update function validates
- ✅ Changes blocked

---

### Requirement 5: API Integration

**ID**: REQ-SQLSERVER-API-001  
**Priority**: High  
**Status**: New

The implementation MUST correctly integrate with TencentCloud SQL Server APIs.

#### Scenario 5.1: Create API Parameter Mapping

**Given** the `CreateBasicDBInstances` API  
**When** creating an instance  
**Then** the following mappings MUST be applied:
- `time_zone` (schema) -> `TimeZone` (API)
- `disk_encrypt_flag` (schema) -> `DiskEncryptFlag` (API)

**Service Layer Code**:
```go
if v, ok := paramMap["time_zone"]; ok {
    request.TimeZone = helper.String(v.(string))
}

if v, ok := paramMap["disk_encrypt_flag"]; ok {
    request.DiskEncryptFlag = helper.IntInt64(v.(int))
}
```

**Validation**:
- ✅ Correct API field names
- ✅ Type conversions applied
- ✅ Optional parameters handled

---

#### Scenario 5.2: Read API - TimeZone from DescribeDBInstances

**Given** the `DescribeDBInstances` API  
**When** reading instance state  
**Then** the `TimeZone` field MUST be accessed from `DBInstance` struct  
**And** nil pointer check MUST be performed  
**And** value MUST be set in state

**Read Code**:
```go
if instance.TimeZone != nil {
    _ = d.Set("time_zone", instance.TimeZone)
}
```

**Validation**:
- ✅ Correct field accessed
- ✅ Nil check performed
- ✅ State updated

---

#### Scenario 5.3: Read API - DiskEncryptFlag from DescribeDBInstancesAttribute

**Given** the `DescribeDBInstancesAttribute` API  
**When** reading instance state  
**Then** a separate API call MUST be made  
**And** the `IsDiskEncryptFlag` field MUST be accessed  
**And** type conversion from *int64 to int MUST be performed  
**And** nil checks MUST be performed  
**And** value MUST be set in state

**Service Layer Method**:
```go
func (me *SqlserverService) DescribeSqlserverInstanceAttributeById(ctx context.Context, instanceId string) (
    attribute *sqlserver.DescribeDBInstancesAttributeResponseParams,
    errRet error,
) {
    // Implementation
}
```

**Read Code**:
```go
attribute, inErr = sqlserverService.DescribeSqlserverInstanceAttributeById(ctx, instanceId)
if attribute != nil && attribute.IsDiskEncryptFlag != nil {
    _ = d.Set("disk_encrypt_flag", int(*attribute.IsDiskEncryptFlag))
}
```

**Validation**:
- ✅ Separate method created
- ✅ API call correct
- ✅ Type conversion performed
- ✅ Double nil check
- ✅ State updated

---

### Requirement 6: Error Handling

**ID**: REQ-SQLSERVER-ERROR-001  
**Priority**: Medium  
**Status**: New

The implementation MUST handle errors gracefully and provide useful error messages.

#### Scenario 6.1: DescribeDBInstancesAttribute API Failure

**Given** the `DescribeDBInstancesAttribute` API call fails  
**When** reading instance state  
**Then** a warning MUST be logged  
**And** the error MUST NOT fail the entire Read operation  
**And** other fields MUST still be populated

**Error Handling**:
```go
if outErr != nil {
    log.Printf("[WARN]%s describe sqlserver instance attribute failed, reason: %v", logId, outErr)
    // Don't fail - just log warning
}
```

**Validation**:
- ✅ Error caught
- ✅ Warning logged
- ✅ Read continues
- ✅ Graceful degradation

---

#### Scenario 6.2: Nil Pointer Safety

**Given** any API response with pointer fields  
**When** accessing fields  
**Then** nil checks MUST be performed  
**And** no nil pointer dereferences MUST occur  
**And** missing fields MUST be handled gracefully

**Safe Access Patterns**:
```go
// Single pointer check
if instance.TimeZone != nil {
    // use value
}

// Double pointer check
if attribute != nil && attribute.IsDiskEncryptFlag != nil {
    // use value
}
```

**Validation**:
- ✅ All pointers checked
- ✅ No panics possible
- ✅ Defensive programming

---

### Requirement 7: Testing

**ID**: REQ-SQLSERVER-TEST-001  
**Priority**: High  
**Status**: New

Comprehensive tests MUST be provided to verify all functionality.

#### Scenario 7.1: Acceptance Test - Timezone

**Given** an acceptance test for timezone  
**When** the test runs  
**Then** it MUST:
- Create instance with custom timezone
- Verify timezone in state
- Test import functionality
- Verify ForceNew behavior

**Test Implementation**:
```go
func TestAccTencentCloudSqlserverBasicInstance_Timezone(t *testing.T) {
    // Test implementation
}
```

**Validation**:
- ✅ Test is parallel
- ✅ Resource created
- ✅ Timezone verified
- ✅ Import tested

---

#### Scenario 7.2: Acceptance Test - Disk Encryption

**Given** an acceptance test for disk encryption  
**When** the test runs  
**Then** it MUST:
- Create instance with encryption enabled
- Verify disk_encrypt_flag=1 in state
- Test import functionality
- Verify ForceNew behavior

**Test Implementation**:
```go
func TestAccTencentCloudSqlserverBasicInstance_DiskEncrypt(t *testing.T) {
    // Test implementation
}
```

**Validation**:
- ✅ Test is parallel
- ✅ Encryption enabled
- ✅ Flag verified
- ✅ Import tested

---

#### Scenario 7.3: Acceptance Test - Both Parameters

**Given** an acceptance test with both parameters  
**When** the test runs  
**Then** it MUST:
- Create instance with both timezone and encryption
- Verify both values in state
- Test import
- Verify no conflicts

**Test Implementation**:
```go
func TestAccTencentCloudSqlserverBasicInstance_TimezoneAndEncrypt(t *testing.T) {
    // Test implementation
}
```

**Validation**:
- ✅ Both parameters set
- ✅ Both verified
- ✅ No conflicts
- ✅ Import works

---

### Requirement 8: Documentation

**ID**: REQ-SQLSERVER-DOC-001  
**Priority**: High  
**Status**: New

Complete and accurate documentation MUST be provided.

#### Scenario 8.1: Resource Documentation File

**Given** the resource documentation file  
**Then** it MUST include:
- Usage examples with timezone
- Usage examples with encryption
- Common timezone values
- Encryption implications
- ForceNew behavior notes
- Import examples

**Documentation Sections**:
1. Example Usage (multiple scenarios)
2. Timezone Values (common list)
3. Disk Encryption (values and implications)
4. Schema Change Behavior (ForceNew fields)
5. Import (with field population note)

**Validation**:
- ✅ All sections present
- ✅ Examples are correct
- ✅ Clear and helpful

---

#### Scenario 8.2: Website Documentation Generation

**Given** the `make doc` command  
**When** documentation is generated  
**Then** the website docs MUST include:
- Both parameters in Argument Reference
- Correct descriptions
- Proper formatting

**Generated File**: `website/docs/r/sqlserver_basic_instance.html.markdown`

**Validation**:
- ✅ Parameters listed
- ✅ Descriptions match schema
- ✅ Formatting correct

---

### Requirement 9: Backward Compatibility

**ID**: REQ-SQLSERVER-COMPAT-001  
**Priority**: Critical  
**Status**: New

The changes MUST be fully backward compatible with existing configurations.

#### Scenario 9.1: Existing Configurations Work Unchanged

**Given** an existing Terraform configuration without new parameters  
**When** the configuration is applied with the updated provider  
**Then** no changes MUST be required  
**And** the instance MUST continue to work  
**And** the state MUST be populated with computed values

**Existing Config**:
```hcl
resource "tencentcloud_sqlserver_basic_instance" "existing" {
  name    = "existing-instance"
  # No time_zone or disk_encrypt_flag
  # ... other fields
}
```

**Expected Behavior**:
- No plan changes
- Computed values populated on next read
- No errors

**Validation**:
- ✅ No breaking changes
- ✅ Existing configs work
- ✅ Smooth upgrade path

---

#### Scenario 9.2: State Migration Not Required

**Given** existing Terraform state files  
**When** the provider is updated  
**Then** no state migration MUST be required  
**And** the state format MUST remain compatible  
**And** Terraform refresh MUST populate new fields

**Validation**:
- ✅ State format unchanged
- ✅ No migration needed
- ✅ Fields added on refresh

---

## Non-Functional Requirements

### NFR-1: Performance

**Requirement**: The additional API call for `DescribeDBInstancesAttribute` MUST NOT significantly impact performance.

**Acceptance Criteria**:
- Read operation completes in reasonable time (<5 seconds under normal conditions)
- API calls are sequential (not problematic for single resource reads)
- Retry logic prevents transient failures

---

### NFR-2: Code Quality

**Requirement**: Code MUST follow project conventions and pass all quality checks.

**Acceptance Criteria**:
- ✅ gofmt formatting applied
- ✅ golangci-lint passes
- ✅ No deprecated functions used
- ✅ Consistent error handling
- ✅ Proper logging
- ✅ File naming conventions followed
- ✅ Function naming conventions followed

---

### NFR-3: Maintainability

**Requirement**: Code MUST be maintainable and follow established patterns.

**Acceptance Criteria**:
- Service layer separation maintained
- Schema patterns consistent with other resources
- Error handling consistent with project style
- Logging consistent with project style
- Comments where necessary

---

## API Mappings

### Create Operation: CreateBasicDBInstances

| Schema Field | API Field | Type Conversion | Required | Default |
|--------------|-----------|-----------------|----------|---------|
| `time_zone` | `TimeZone` | string -> *string | No | "China Standard Time" (API) |
| `disk_encrypt_flag` | `DiskEncryptFlag` | int -> *int64 | No | 0 |

---

### Read Operation: DescribeDBInstances

| API Field | Schema Field | Type Conversion | Always Present |
|-----------|--------------|-----------------|----------------|
| `DBInstance.TimeZone` | `time_zone` | *string -> string | Yes |

---

### Read Operation: DescribeDBInstancesAttribute

| API Field | Schema Field | Type Conversion | Always Present |
|-----------|--------------|-----------------|----------------|
| `IsDiskEncryptFlag` | `disk_encrypt_flag` | *int64 -> int | Yes |

---

## Test Strategy

### Unit Tests
- Schema validation
- Type conversions
- Nil pointer handling

### Acceptance Tests
1. **TestAccTencentCloudSqlserverBasicInstance_Timezone**
   - Create with custom timezone
   - Verify state
   - Test import

2. **TestAccTencentCloudSqlserverBasicInstance_DiskEncrypt**
   - Create with encryption enabled
   - Verify state
   - Test import

3. **TestAccTencentCloudSqlserverBasicInstance_TimezoneAndEncrypt**
   - Create with both parameters
   - Verify both values
   - Test import

### Manual Tests
- Import existing instances
- Verify ForceNew behavior with `terraform plan`
- Test with various timezone values
- Test error scenarios

---

## Success Metrics

1. **Functional Completeness**: All 9 requirements implemented
2. **Test Coverage**: 3+ acceptance tests passing
3. **Documentation**: Complete with examples
4. **Backward Compatibility**: No breaking changes
5. **Code Quality**: All linters pass
6. **Performance**: Read operation < 5 seconds

---

## Related Specifications

None (this is a new capability)

---

## References

- **TencentCloud API Documentation**:
  - CreateBasicDBInstances: https://cloud.tencent.com/document/api/238/50262
  - DescribeDBInstances: https://cloud.tencent.com/document/api/238/19969
  - DescribeDBInstancesAttribute: https://cloud.tencent.com/document/api/238/73319

- **Terraform Provider SDK**:
  - Schema: https://www.terraform.io/plugin/sdkv2/schemas
  - ForceNew: https://www.terraform.io/plugin/sdkv2/schemas/schema-behaviors#forcenew

- **Project Files**:
  - Resource: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`
  - Service: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`
  - SDK Models: `vendor/.../sqlserver/v20180328/models.go`

---

**END OF SPECIFICATION**
