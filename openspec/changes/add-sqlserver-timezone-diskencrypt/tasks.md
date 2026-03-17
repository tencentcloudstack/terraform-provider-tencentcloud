# Tasks: Add TimeZone and DiskEncryptFlag Support to SqlServer Basic Instance

**Change ID**: `add-sqlserver-timezone-diskencrypt`  
**Total Tasks**: 34  
**Estimated Time**: 3.5 hours

---

## Phase 1: Schema Definition (30 min)

### Task 1.1: Add `time_zone` Schema Field
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Add the schema field after the existing `engine_version` field:

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
- ✅ Field is Optional
- ✅ Field is Computed
- ✅ Field is ForceNew
- ✅ Description is clear and helpful

---

### Task 1.2: Add `disk_encrypt_flag` Schema Field
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Add the schema field after `time_zone`:

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
- ✅ Field is Optional with Default 0
- ✅ Field is Computed
- ✅ Field is ForceNew
- ✅ ValidateFunc restricts to 0-1 range
- ✅ Description explains values

---

### Task 1.3: Update Immutable Args List
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Update the `immutableArgs` list in `resourceTencentCloudSqlserverBasicInstanceUpdate`:

```go
immutableArgs := []string{"collation", "time_zone", "disk_encrypt_flag"}
```

**Validation**:
- ✅ Both new fields added to immutableArgs
- ✅ Change detection will prevent modification

---

### Task 1.4: Verify Schema Compilation
**Command**: `go build ./tencentcloud/services/sqlserver/...`

**Validation**:
- ✅ No compilation errors
- ✅ Schema syntax is correct

---

## Phase 2: Service Layer (45 min)

### Task 2.1: Create `DescribeSqlserverInstanceAttributeById` Method Signature
**File**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

Add method after existing Describe methods (around line 400):

```go
func (me *SqlserverService) DescribeSqlserverInstanceAttributeById(ctx context.Context, instanceId string) (
    attribute *sqlserver.DescribeDBInstancesAttributeResponseParams,
    errRet error,
) {
```

**Validation**:
- ✅ Method signature matches service pattern
- ✅ Returns response params and error

---

### Task 2.2: Implement API Request Construction
**File**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

```go
    logId := tccommon.GetLogId(ctx)
    request := sqlserver.NewDescribeDBInstancesAttributeRequest()
    request.InstanceId = &instanceId
    
    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
                logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()
```

**Validation**:
- ✅ LogId retrieved from context
- ✅ Request initialized correctly
- ✅ InstanceId parameter set
- ✅ Defer error logging added

---

### Task 2.3: Add Rate Limiting and API Call
**File**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

```go
    ratelimit.Check(request.GetAction())
    
    response, err := me.client.UseSqlserverClient().DescribeDBInstancesAttribute(request)
    if err != nil {
        errRet = err
        return
    }
    
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
        logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
```

**Validation**:
- ✅ Rate limiting applied
- ✅ API called with correct client method
- ✅ Error handled and returned
- ✅ Debug logging added

---

### Task 2.4: Extract Response and Return
**File**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

```go
    if response != nil && response.Response != nil {
        attribute = response.Response
    }
    
    return
}
```

**Validation**:
- ✅ Nil check for response
- ✅ Response params extracted
- ✅ Return statement present

---

### Task 2.5: Modify `CreateSqlserverBasicInstance` to Accept New Parameters
**File**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

Find the `CreateSqlserverBasicInstance` method and add parameter handling after existing paramMap reads:

```go
// time_zone
if v, ok := paramMap["time_zone"]; ok {
    request.TimeZone = helper.String(v.(string))
}

// disk_encrypt_flag
if v, ok := paramMap["disk_encrypt_flag"]; ok {
    request.DiskEncryptFlag = helper.IntInt64(v.(int))
}
```

**Validation**:
- ✅ Parameters read from paramMap
- ✅ Proper type conversion (string, int64)
- ✅ Optional parameters handled with ok check
- ✅ Helper functions used correctly

---

## Phase 3: Create Operation (15 min)

### Task 3.1: Read `time_zone` from Schema in Create
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

In `resourceTencentCloudSqlserverBasicInstanceCreate`, add after `collation` parameter:

```go
// time_zone
if v, ok := d.GetOk("time_zone"); ok {
    paramMap["time_zone"] = v.(string)
}
```

**Validation**:
- ✅ Parameter read from schema
- ✅ Added to paramMap
- ✅ Type assertion correct

---

### Task 3.2: Read `disk_encrypt_flag` from Schema in Create
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Add after `time_zone`:

```go
// disk_encrypt_flag
if v, ok := d.GetOk("disk_encrypt_flag"); ok {
    paramMap["disk_encrypt_flag"] = v.(int)
}
```

**Validation**:
- ✅ Parameter read from schema
- ✅ Added to paramMap
- ✅ Type assertion correct (int)

---

### Task 3.3: Verify paramMap Usage
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Verify the paramMap is passed to service layer (should already exist):

```go
instanceId, inErr = sqlserverService.CreateSqlserverBasicInstance(ctx, paramMap, weekSet, voucherIds, securityGroups)
```

**Validation**:
- ✅ paramMap passed to service layer
- ✅ No changes needed (already correct)

---

### Task 3.4: Test Create Compilation
**Command**: `go build ./tencentcloud/services/sqlserver/...`

**Validation**:
- ✅ No compilation errors in Create logic
- ✅ All variables used correctly

---

## Phase 4: Read Operation (30 min)

### Task 4.1: Read `TimeZone` from Existing Response
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

In `resourceTencentCloudSqlserverBasicInstanceRead`, after existing field sets (around line 345):

```go
// time_zone
if instance.TimeZone != nil {
    _ = d.Set("time_zone", instance.TimeZone)
}
```

**Validation**:
- ✅ Nil check for pointer
- ✅ Value set in state
- ✅ Field name matches schema

---

### Task 4.2: Call `DescribeSqlserverInstanceAttributeById`
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

After the existing instance describe call:

```go
// Get disk encryption flag from attributes API
var attribute *sqlserver.DescribeDBInstancesAttributeResponseParams
outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    attribute, inErr = sqlserverService.DescribeSqlserverInstanceAttributeById(ctx, instanceId)
    if inErr != nil {
        return tccommon.RetryError(inErr)
    }
    return nil
})
if outErr != nil {
    log.Printf("[WARN]%s describe sqlserver instance attribute failed, reason: %v", logId, outErr)
    // Don't fail the entire read, just log the warning
}
```

**Validation**:
- ✅ Retry logic applied
- ✅ Error handling with warning log
- ✅ Doesn't fail entire Read on error

---

### Task 4.3: Read `IsDiskEncryptFlag` from Attribute Response
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

After the attribute API call:

```go
// disk_encrypt_flag
if attribute != nil && attribute.IsDiskEncryptFlag != nil {
    _ = d.Set("disk_encrypt_flag", int(*attribute.IsDiskEncryptFlag))
}
```

**Validation**:
- ✅ Double nil check (attribute and field)
- ✅ Type conversion from *int64 to int
- ✅ Value set in state

---

### Task 4.4: Handle Nil Pointers Safely
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Review all pointer dereferences:

```go
// Ensure all pointer fields have nil checks
if instance.TimeZone != nil { ... }
if attribute != nil && attribute.IsDiskEncryptFlag != nil { ... }
```

**Validation**:
- ✅ No bare pointer dereferences
- ✅ All checks in place
- ✅ No potential panics

---

### Task 4.5: Test Read Compilation
**Command**: `go build ./tencentcloud/services/sqlserver/...`

**Validation**:
- ✅ No compilation errors in Read logic
- ✅ All API calls correct

---

### Task 4.6: Verify State Management
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Ensure both fields are properly managed:
- Set during Read
- Not modified during Update (ForceNew)
- Handled during Import

**Validation**:
- ✅ Fields set in Read
- ✅ Fields in immutableArgs
- ✅ Import will call Read automatically

---

## Phase 5: Testing (45 min)

### Task 5.1: Add Test Case - Basic Instance with Custom Timezone
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go`

Add test after existing tests:

```go
func TestAccTencentCloudSqlserverBasicInstance_Timezone(t *testing.T) {
    t.Parallel()
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { tcacctest.AccPreCheck(t) },
        Providers:    tcacctest.AccProviders,
        CheckDestroy: testAccCheckSqlserverBasicInstanceDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccSqlserverBasicInstance_timezone,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckSqlserverBasicInstanceExists("tencentcloud_sqlserver_basic_instance.test"),
                    resource.TestCheckResourceAttr("tencentcloud_sqlserver_basic_instance.test", "time_zone", "UTC"),
                ),
            },
            {
                ResourceName:      "tencentcloud_sqlserver_basic_instance.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

const testAccSqlserverBasicInstance_timezone = `
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name              = "test-sqlserver-timezone"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_PREMIUM"
  cpu               = 2
  memory            = 4
  storage           = 20
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2008R2"
  time_zone         = "UTC"
}
`
```

**Validation**:
- ✅ Test is parallel
- ✅ timezone value is checked
- ✅ Import state is verified

---

### Task 5.2: Add Test Case - Basic Instance with Disk Encryption
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go`

```go
func TestAccTencentCloudSqlserverBasicInstance_DiskEncrypt(t *testing.T) {
    t.Parallel()
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { tcacctest.AccPreCheck(t) },
        Providers:    tcacctest.AccProviders,
        CheckDestroy: testAccCheckSqlserverBasicInstanceDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccSqlserverBasicInstance_diskencrypt,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckSqlserverBasicInstanceExists("tencentcloud_sqlserver_basic_instance.test"),
                    resource.TestCheckResourceAttr("tencentcloud_sqlserver_basic_instance.test", "disk_encrypt_flag", "1"),
                ),
            },
            {
                ResourceName:      "tencentcloud_sqlserver_basic_instance.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

const testAccSqlserverBasicInstance_diskencrypt = `
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name              = "test-sqlserver-encrypt"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_SSD"
  cpu               = 2
  memory            = 4
  storage           = 20
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2008R2"
  disk_encrypt_flag = 1
}
`
```

**Validation**:
- ✅ Test is parallel
- ✅ disk_encrypt_flag=1 is checked
- ✅ Import state is verified

---

### Task 5.3: Add Test Case - Both Parameters Together
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go`

```go
func TestAccTencentCloudSqlserverBasicInstance_TimezoneAndEncrypt(t *testing.T) {
    t.Parallel()
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { tcacctest.AccPreCheck(t) },
        Providers:    tcacctest.AccProviders,
        CheckDestroy: testAccCheckSqlserverBasicInstanceDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccSqlserverBasicInstance_both,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckSqlserverBasicInstanceExists("tencentcloud_sqlserver_basic_instance.test"),
                    resource.TestCheckResourceAttr("tencentcloud_sqlserver_basic_instance.test", "time_zone", "UTC"),
                    resource.TestCheckResourceAttr("tencentcloud_sqlserver_basic_instance.test", "disk_encrypt_flag", "1"),
                ),
            },
            {
                ResourceName:      "tencentcloud_sqlserver_basic_instance.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

const testAccSqlserverBasicInstance_both = `
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name              = "test-sqlserver-both"
  availability_zone = "ap-guangzhou-3"
  machine_type      = "CLOUD_SSD"
  cpu               = 2
  memory            = 4
  storage           = 20
  charge_type       = "POSTPAID_BY_HOUR"
  engine_version    = "2008R2"
  time_zone         = "UTC"
  disk_encrypt_flag = 1
}
`
```

**Validation**:
- ✅ Both parameters tested together
- ✅ Import works with both set

---

### Task 5.4: Compile Tests
**Command**: `go test -c ./tencentcloud/services/sqlserver/... -o /dev/null`

**Validation**:
- ✅ Tests compile successfully
- ✅ No syntax errors

---

### Task 5.5: Run Unit Tests (Optional - Requires Credentials)
**Command**: `TF_ACC=1 go test ./tencentcloud/services/sqlserver/ -v -run="TestAccTencentCloudSqlserverBasicInstance_(Timezone|DiskEncrypt)"`

**Note**: Only run if credentials are available. Otherwise, verify compilation only.

**Validation**:
- ✅ Tests pass (if run)
- ✅ Resources created successfully
- ✅ Import works correctly

---

## Phase 6: Documentation (30 min)

### Task 6.1: Update Resource Documentation File
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`

Add usage examples after existing content:

```markdown
Example with custom timezone:

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
  time_zone         = "UTC"
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```

Example with disk encryption enabled:

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
  disk_encrypt_flag = 1
  time_zone         = "China Standard Time"
  
  vpc_id    = "vpc-xxxxx"
  subnet_id = "subnet-xxxxx"
}
```
```

**Validation**:
- ✅ Examples are syntactically correct
- ✅ Both parameters documented with examples
- ✅ Real-world use cases shown

---

### Task 6.2: Document Common Timezone Values
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`

Add a note about timezone values:

```markdown
## Timezone Values

Common timezone values for `time_zone` parameter:
- `China Standard Time` (default)
- `UTC`
- `Eastern Standard Time`
- `Pacific Standard Time`
- `Central Standard Time`

For a complete list, refer to the SQL Server timezone documentation.
```

**Validation**:
- ✅ Common values listed
- ✅ Reference to official docs

---

### Task 6.3: Document Disk Encryption Implications
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`

Add a note about encryption:

```markdown
## Disk Encryption

The `disk_encrypt_flag` parameter controls whether disk encryption is enabled:
- `0` - Encryption disabled (default)
- `1` - Encryption enabled

**Important Notes**:
- Disk encryption cannot be changed after instance creation (ForceNew)
- Encrypted instances may have slightly lower performance
- Encryption is recommended for sensitive data and compliance requirements
```

**Validation**:
- ✅ Values clearly explained
- ✅ ForceNew behavior documented
- ✅ Performance note included

---

### Task 6.4: Add Note About ForceNew Behavior
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`

Ensure ForceNew is mentioned:

```markdown
## Schema Change Behavior

The following parameters cannot be changed after instance creation and will trigger instance recreation:
- `time_zone`
- `disk_encrypt_flag`
- `engine_version`
- `charge_type`
- `vpc_id`
- `subnet_id`
```

**Validation**:
- ✅ ForceNew parameters listed
- ✅ User expectations set correctly

---

### Task 6.5: Generate Website Documentation
**Command**: `make doc`

**Validation**:
- ✅ Command runs successfully
- ✅ File `website/docs/r/sqlserver_basic_instance.html.markdown` updated
- ✅ New parameters appear in generated docs

---

### Task 6.6: Verify Generated Documentation
**File**: `website/docs/r/sqlserver_basic_instance.html.markdown`

Check that generated docs include:
- `time_zone` in Argument Reference
- `disk_encrypt_flag` in Argument Reference
- Proper formatting

**Validation**:
- ✅ Both parameters in Argument Reference
- ✅ Descriptions match schema
- ✅ Examples included

---

### Task 6.7: Review Documentation Quality
Review all documentation for:
- Spelling and grammar
- Technical accuracy
- Clarity and completeness
- Proper HCL formatting

**Validation**:
- ✅ No spelling errors
- ✅ Technical details correct
- ✅ Examples work

---

### Task 6.8: Add Import Example
**File**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md`

Update import section:

```markdown
Import

SQLServer basic instance can be imported using the id, e.g.

```bash
terraform import tencentcloud_sqlserver_basic_instance.example mssql-abc123
```

After import, both `time_zone` and `disk_encrypt_flag` will be populated from the API.
```

**Validation**:
- ✅ Import command correct
- ✅ Note about field population

---

## Phase 7: Code Quality & Validation (15 min)

### Task 7.1: Format Code with gofmt
**Command**: `gofmt -w tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

**Validation**:
- ✅ Code is properly formatted
- ✅ No formatting changes needed

---

### Task 7.2: Run Linter
**Command**: `make lint`

**Validation**:
- ✅ No linter errors
- ✅ No linter warnings for new code
- ✅ Code follows project style

---

### Task 7.3: Compile Provider
**Command**: `go build -o /tmp/terraform-provider-tencentcloud .`

**Validation**:
- ✅ Provider compiles successfully
- ✅ No compilation errors
- ✅ Binary created

---

### Task 7.4: Check for Read Lints
**Tool**: IDE linter or `read_lints` tool

**Validation**:
- ✅ No new linter errors introduced
- ✅ Existing issues not in modified files

---

### Task 7.5: Final Review Checklist
Review the implementation against requirements:

- ✅ Schema fields added correctly
- ✅ ForceNew behavior implemented
- ✅ Computed values work
- ✅ Create operation passes parameters
- ✅ Read operation fetches both values
- ✅ Service layer method created
- ✅ Nil pointers handled safely
- ✅ Error handling in place
- ✅ Logging added
- ✅ Tests added
- ✅ Documentation complete
- ✅ Code formatted
- ✅ No linter errors
- ✅ Backward compatible

**Validation**:
- ✅ All items checked
- ✅ Implementation complete
- ✅ Ready for PR

---

## Summary

**Total Tasks**: 34  
**Phases**: 7  
**Estimated Time**: 3.5 hours

### Task Breakdown by Phase:
1. Schema Definition: 4 tasks (30 min)
2. Service Layer: 5 tasks (45 min)
3. Create Operation: 4 tasks (15 min)
4. Read Operation: 6 tasks (30 min)
5. Testing: 5 tasks (45 min)
6. Documentation: 8 tasks (30 min)
7. Code Quality: 5 tasks (15 min)

### Key Implementation Points:
1. Both parameters are Optional, Computed, ForceNew
2. `time_zone` from DescribeDBInstances API
3. `disk_encrypt_flag` from DescribeDBInstancesAttribute API (separate call)
4. Proper nil checking for all pointer fields
5. Backward compatible (optional with defaults)

### Success Criteria:
- ✅ Users can set and view both parameters
- ✅ ForceNew behavior works correctly
- ✅ Import populates both fields
- ✅ Tests pass
- ✅ Documentation complete
- ✅ No breaking changes

---

**Ready for Implementation!**

Run `openspec apply add-sqlserver-timezone-diskencrypt` to begin.
