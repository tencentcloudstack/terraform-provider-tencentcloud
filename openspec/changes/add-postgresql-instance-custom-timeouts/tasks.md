# Implementation Tasks: Add Custom Timeouts for PostgreSQL Instance Resources

## Change ID
`add-postgresql-instance-custom-timeouts`

## Overview
This document provides a detailed task breakdown for adding custom timeout configuration support to PostgreSQL instance resources.

---

## Task Breakdown

### Task 1: Modify `tencentcloud_postgresql_instance` Resource

#### Task 1.1: Add Timeouts Block to Schema
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`  
**Location**: In `ResourceTencentCloudPostgresqlInstance()` function, after `Importer` block

**Action**:
```go
func ResourceTencentCloudPostgresqlInstance() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudPostgresqlInstanceCreate,
        Read:   resourceTencentCloudPostgresqlInstanceRead,
        Update: resourceTencentCloudPostgresqlInstanceUpdate,
        Delete: resourceTencentCLoudPostgresqlInstanceDelete,
        Importer: &schema.ResourceImporter{
            State: helper.ImportWithDefaultValue(map[string]interface{}{
                "delete_protection": false,
            }),
        },
        
        // ADD THIS BLOCK
        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(60 * time.Minute),
            Update: schema.DefaultTimeout(60 * time.Minute),
        },
        
        Schema: map[string]*schema.Schema{
            // ... existing schema ...
        },
    }
}
```

**Checklist**:
- [ ] Add `Timeouts` block after `Importer`
- [ ] Set Create timeout to 60 minutes
- [ ] Set Update timeout to 60 minutes
- [ ] Ensure proper indentation

---

#### Task 1.2: Update Create Function - Initial Instance Creation
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`  
**Function**: `resourceTencentCloudPostgresqlInstanceCreate`  
**Location**: Line ~360-400 (where CreatePostgresqlInstance is called)

**Current Code**:
```go
err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    instanceId, inErr = postgresqlService.CreatePostgresqlInstance(...)
    if inErr != nil {
        return tccommon.RetryError(inErr)
    }
    return nil
})
```

**Updated Code**:
```go
err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    instanceId, inErr = postgresqlService.CreatePostgresqlInstance(...)
    if inErr != nil {
        return tccommon.RetryError(inErr)
    }
    return nil
})
```

**Checklist**:
- [ ] Replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutCreate)`
- [ ] Verify function signature unchanged
- [ ] Verify error handling unchanged

---

#### Task 1.3: Update Create Function - Instance Status Check
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`  
**Function**: `resourceTencentCloudPostgresqlInstanceCreate`  
**Location**: After instance creation, where CheckDBInstanceStatus is called

**Search for patterns**:
```go
resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    // ... status check logic ...
})
```

**Note**: Check if status checking should also use Create timeout. If yes, update to:
```go
resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    // ... status check logic ...
})
```

**Checklist**:
- [ ] Identify all Retry calls in create function
- [ ] Determine if status checks should use Create timeout
- [ ] Update if appropriate

---

#### Task 1.4: Update Update Function - Memory/Storage/CPU Scaling
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`  
**Function**: `resourceTencentCloudPostgresqlInstanceUpdate`  
**Location**: Line ~1410 (memory/storage/cpu upgrade block)

**Current Code**:
```go
// upgrade storage and memory size
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    memory := d.Get("memory").(int)
    storage := d.Get("storage").(int)
    var cpu int
    if v, ok := d.GetOkExists("cpu"); ok {
        cpu = v.(int)
    }

    outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        inErr = postgresqlService.UpgradePostgresqlInstance(ctx, instanceId, memory, storage, cpu, waitSwitch)
        if inErr != nil {
            return tccommon.RetryError(inErr)
        }

        return nil
    })

    if outErr != nil {
        return outErr
    }
}
```

**Updated Code**:
```go
// upgrade storage and memory size
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    memory := d.Get("memory").(int)
    storage := d.Get("storage").(int)
    var cpu int
    if v, ok := d.GetOkExists("cpu"); ok {
        cpu = v.(int)
    }

    outErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
        inErr = postgresqlService.UpgradePostgresqlInstance(ctx, instanceId, memory, storage, cpu, waitSwitch)
        if inErr != nil {
            return tccommon.RetryError(inErr)
        }

        return nil
    })

    if outErr != nil {
        return outErr
    }
}
```

**Checklist**:
- [ ] Locate memory/storage/cpu upgrade block
- [ ] Replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutUpdate)`
- [ ] Verify logic unchanged

---

#### Task 1.5: Update Update Function - Instance Status Check After Scaling
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`  
**Function**: `resourceTencentCloudPostgresqlInstanceUpdate`  
**Location**: After UpgradePostgresqlInstance, where CheckDBInstanceStatus is called

**Search for**:
```go
// After upgrade, check status
checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
```

**Note**: CheckDBInstanceStatus might not use resource.Retry. If it does, consider updating to use Update timeout.

**Checklist**:
- [ ] Check if status verification uses Retry
- [ ] Update if appropriate

---

#### Task 1.6: Format Code
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`

**Command**:
```bash
go fmt tencentcloud/services/postgresql/resource_tc_postgresql_instance.go
```

**Checklist**:
- [ ] Run `go fmt` on the file
- [ ] Verify no formatting errors
- [ ] Save the file

---

### Task 2: Modify `tencentcloud_postgresql_readonly_instance` Resource

#### Task 2.1: Add Timeouts Block to Schema
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`  
**Location**: In `ResourceTencentCloudPostgresqlReadonlyInstance()` function, after `Importer` block

**Action**:
```go
func ResourceTencentCloudPostgresqlReadonlyInstance() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudPostgresqlReadOnlyInstanceCreate,
        Read:   resourceTencentCloudPostgresqlReadOnlyInstanceRead,
        Update: resourceTencentCloudPostgresqlReadOnlyInstanceUpdate,
        Delete: resourceTencentCLoudPostgresqlReadOnlyInstanceDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        
        // ADD THIS BLOCK
        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(60 * time.Minute),
            Update: schema.DefaultTimeout(60 * time.Minute),
        },
        
        Schema: map[string]*schema.Schema{
            // ... existing schema ...
        },
    }
}
```

**Checklist**:
- [ ] Add `Timeouts` block after `Importer`
- [ ] Set Create timeout to 60 minutes
- [ ] Set Update timeout to 60 minutes
- [ ] Ensure proper indentation

---

#### Task 2.2: Update Create Function
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`  
**Function**: `resourceTencentCloudPostgresqlReadOnlyInstanceCreate`

**Search for**:
```go
resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    // ... create readonly instance logic ...
})
```

**Replace with**:
```go
resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    // ... create readonly instance logic ...
})
```

**Checklist**:
- [ ] Locate create Retry call
- [ ] Replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutCreate)`
- [ ] Check for additional Retry calls in create function
- [ ] Update status checking Retry calls if appropriate

---

#### Task 2.3: Update Update Function - Memory/Storage/CPU Scaling
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`  
**Function**: `resourceTencentCloudPostgresqlReadOnlyInstanceUpdate`

**Search for**:
```go
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    // ... scaling logic ...
    resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        // ... upgrade logic ...
    })
}
```

**Replace with**:
```go
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    // ... scaling logic ...
    resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
        // ... upgrade logic ...
    })
}
```

**Checklist**:
- [ ] Locate memory/storage/cpu upgrade block
- [ ] Replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutUpdate)`
- [ ] Verify all scaling-related Retry calls are updated

---

#### Task 2.4: Format Code
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`

**Command**:
```bash
go fmt tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go
```

**Checklist**:
- [ ] Run `go fmt` on the file
- [ ] Verify no formatting errors
- [ ] Save the file

---

### Task 3: Verification and Testing

#### Task 3.1: Code Compilation
**Command**:
```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go build ./tencentcloud/services/postgresql/...
```

**Checklist**:
- [ ] Code compiles without errors
- [ ] No syntax errors
- [ ] No import errors

---

#### Task 3.2: Linter Checks
**Command**:
```bash
# Run linter on modified files
golint tencentcloud/services/postgresql/resource_tc_postgresql_instance.go
golint tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go
```

**Checklist**:
- [ ] No new linter warnings
- [ ] Code follows project conventions
- [ ] No deprecated function usage

---

#### Task 3.3: Code Review Checklist

**Schema Changes**:
- [ ] Timeouts block added correctly
- [ ] Default timeout values set to 60 minutes
- [ ] Positioned correctly in resource definition

**Create Function Changes**:
- [ ] `d.Timeout(schema.TimeoutCreate)` used in create operations
- [ ] Error handling unchanged
- [ ] Logic flow unchanged

**Update Function Changes**:
- [ ] `d.Timeout(schema.TimeoutUpdate)` used in memory/storage/CPU scaling
- [ ] Only scaling operations use Update timeout
- [ ] Other update operations unchanged (if using different timeouts)

**Code Quality**:
- [ ] Code formatted with `go fmt`
- [ ] No breaking changes
- [ ] Backward compatible
- [ ] Follows existing code patterns

---

### Task 4: Documentation Updates (Optional)

#### Task 4.1: Update Resource Documentation
**Files to update**:
- `website/docs/r/postgresql_instance.html.markdown`
- `website/docs/r/postgresql_readonly_instance.html.markdown`

**Add timeout example**:
```markdown
## Timeouts

`tencentcloud_postgresql_instance` provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `60m`) Used for creating instances.
- `update` - (Default `60m`) Used for updating instances, particularly when scaling memory, storage, or CPU.

## Example with Custom Timeouts

```hcl
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example-instance"
  availability_zone = "ap-guangzhou-3"
  # ... other configurations ...
  
  timeouts {
    create = "90m"
    update = "120m"
  }
}
```
```

**Checklist**:
- [ ] Add Timeouts section to documentation
- [ ] Include default values (60m)
- [ ] Add usage example
- [ ] Document which operations use which timeout

---

## Testing Strategy

### Unit Testing
1. **Test Default Timeout Behavior**
   - Verify default 60-minute timeout is applied
   - Test without explicit timeout configuration

2. **Test Custom Timeout Configuration**
   - Set custom timeout values in Terraform config
   - Verify custom values are respected

3. **Test Create Operation**
   - Create new instance
   - Verify timeout is applied correctly
   - Test with both default and custom timeouts

4. **Test Update Operation**
   - Update memory/storage/CPU
   - Verify Update timeout is applied
   - Test with both default and custom timeouts

### Integration Testing
1. **Real Instance Creation**
   - Create actual PostgreSQL instance
   - Monitor operation time
   - Verify timeout behavior

2. **Resource Scaling**
   - Scale memory, storage, CPU
   - Monitor operation time
   - Verify timeout during long-running operations

---

## Implementation Order

### Recommended Sequence
1. ✅ Task 1.1: Add Timeouts to postgresql_instance schema
2. ✅ Task 1.2: Update postgresql_instance create function
3. ✅ Task 1.3: Check postgresql_instance create status logic
4. ✅ Task 1.4: Update postgresql_instance memory/storage/CPU scaling
5. ✅ Task 1.5: Check postgresql_instance update status logic
6. ✅ Task 1.6: Format postgresql_instance file
7. ✅ Task 2.1: Add Timeouts to readonly_instance schema
8. ✅ Task 2.2: Update readonly_instance create function
9. ✅ Task 2.3: Update readonly_instance memory/storage/CPU scaling
10. ✅ Task 2.4: Format readonly_instance file
11. ✅ Task 3.1: Compile and verify code
12. ✅ Task 3.2: Run linter checks
13. ✅ Task 3.3: Code review
14. ⭕ Task 4.1: Update documentation (optional)

---

## Time Estimates

| Task | Estimated Time |
|------|----------------|
| Task 1.1 - Add schema Timeouts | 5 minutes |
| Task 1.2 - Update create function | 10 minutes |
| Task 1.3 - Check create status | 5 minutes |
| Task 1.4 - Update scaling logic | 10 minutes |
| Task 1.5 - Check update status | 5 minutes |
| Task 1.6 - Format code | 2 minutes |
| Task 2.1 - Add schema Timeouts | 5 minutes |
| Task 2.2 - Update create function | 10 minutes |
| Task 2.3 - Update scaling logic | 10 minutes |
| Task 2.4 - Format code | 2 minutes |
| Task 3.1 - Compilation | 5 minutes |
| Task 3.2 - Linter checks | 5 minutes |
| Task 3.3 - Code review | 10 minutes |
| Task 4.1 - Documentation | 15 minutes |
| **Total** | **~99 minutes (~1.5 hours)** |

---

## Notes

1. **Code Placement**: Any new functions or helper code should be added at the end of the file
2. **Formatting**: Always run `go fmt` after making changes to ensure consistent formatting
3. **Testing**: Test both default and custom timeout configurations
4. **Backward Compatibility**: Ensure existing configurations work without modification
5. **Pattern Consistency**: Follow the same pattern used in other TencentCloud resources (CVM, MySQL, CLB)

---

## Success Criteria

- ✅ Both resources have Timeouts block with 60-minute defaults
- ✅ Create operations use Create timeout
- ✅ Memory/Storage/CPU scaling uses Update timeout
- ✅ Code compiles without errors
- ✅ Linter checks pass
- ✅ Code formatted with `go fmt`
- ✅ Backward compatible with existing configurations
- ✅ Follows project conventions

---

## Related Files

### Files to Modify
1. `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`
2. `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`

### Files to Review (Reference)
1. `tencentcloud/services/cvm/resource_tc_instance.go` (timeout pattern)
2. `tencentcloud/services/cdb/resource_tc_mysql_instance.go` (timeout pattern)
3. `tencentcloud/services/clb/resource_tc_clb_instance.go` (timeout pattern)

### Documentation Files (Optional)
1. `website/docs/r/postgresql_instance.html.markdown`
2. `website/docs/r/postgresql_readonly_instance.html.markdown`
