# OpenSpec Proposal: Add Custom Timeouts for PostgreSQL Instance Resources

## Change ID
`add-postgresql-instance-custom-timeouts`

## Status
**Proposed** - 2026-03-24

## Summary
Add custom timeout configuration support to `tencentcloud_postgresql_instance` and `tencentcloud_postgresql_readonly_instance` resources to allow users to configure operation timeouts for create and update operations, especially for resource scaling scenarios.

---

## Background

### Current Situation
Currently, both PostgreSQL instance resources (`tencentcloud_postgresql_instance` and `tencentcloud_postgresql_readonly_instance`) do not support custom timeout configuration. They rely on the default Terraform timeouts which may not be sufficient for:

1. **Instance creation**: Creating PostgreSQL instances can take a significant amount of time depending on instance specifications and region load
2. **Resource scaling**: Upgrading memory, storage, or CPU specifications requires instance restart and can take extended time

### Problem Statement
Users have no control over operation timeouts, leading to:
- Premature timeout errors during long-running operations
- Inconsistent user experience compared to other cloud resources
- Need for manual retry mechanisms

---

## Objectives

### Primary Goals
1. Add `Timeouts` configuration block to both PostgreSQL instance resources
2. Set appropriate default timeout values for create and update operations
3. Maintain backward compatibility with existing configurations
4. Follow the timeout pattern established in other TencentCloud resources

### Success Criteria
- Users can customize timeout values for create and update operations
- Default timeouts are set to 60 minutes as requested
- Update timeout applies specifically to resource scaling operations (memory, storage, CPU)
- No breaking changes to existing resource configurations
- Code follows project conventions and passes all linters

---

## Technical Specification

### Resources to Modify

#### 1. `tencentcloud_postgresql_instance`
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`

**Changes Required**:
- Add `Timeouts` block to resource schema
- Set default create timeout: 60 minutes
- Set default update timeout: 60 minutes (for memory/storage/CPU changes)

#### 2. `tencentcloud_postgresql_readonly_instance`
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`

**Changes Required**:
- Add `Timeouts` block to resource schema
- Set default create timeout: 60 minutes
- Set default update timeout: 60 minutes (for memory/storage/CPU changes)

### Timeout Configuration Specification

```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(60 * time.Minute),
    Update: schema.DefaultTimeout(60 * time.Minute),
},
```

### Usage in Code

#### For Create Operations
Replace:
```go
err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    // create logic
})
```

With:
```go
err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    // create logic
})
```

#### For Update Operations (Memory/Storage/CPU changes)
Replace:
```go
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        // upgrade logic
    })
}
```

With:
```go
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    outErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
        // upgrade logic
    })
}
```

---

## Implementation Details

### 1. Schema Changes

Both resources need to add the `Timeouts` field in their schema definition:

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

### 2. Create Function Updates

**File**: `resource_tc_postgresql_instance.go`

**Location**: Line ~360 (resourceTencentCloudPostgresqlInstanceCreate function)

**Changes**:
```go
// Before
err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    instanceId, inErr = postgresqlService.CreatePostgresqlInstance(...)
    if inErr != nil {
        return tccommon.RetryError(inErr)
    }
    return nil
})

// After
err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    instanceId, inErr = postgresqlService.CreatePostgresqlInstance(...)
    if inErr != nil {
        return tccommon.RetryError(inErr)
    }
    return nil
})
```

### 3. Update Function - Resource Scaling

**File**: `resource_tc_postgresql_instance.go`

**Location**: Line ~1410 (memory/storage/cpu upgrade logic)

**Changes**:
```go
// Before
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

// After
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

### 4. Readonly Instance - Similar Changes

Apply the same pattern to `resource_tc_postgresql_readonly_instance.go`:
- Add Timeouts block to schema
- Update create function to use `d.Timeout(schema.TimeoutCreate)`
- Update memory/storage/CPU upgrade logic to use `d.Timeout(schema.TimeoutUpdate)`

---

## Reference Implementations

### Example 1: CVM Instance (resource_tc_instance.go)
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(15 * time.Minute),
},
```

### Example 2: MySQL Instance (resource_tc_mysql_instance.go)
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(20 * time.Minute),
    Delete: schema.DefaultTimeout(20 * time.Minute),
},
```

### Example 3: CLB Instance (resource_tc_clb_instance.go)
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(10 * time.Minute),
    Update: schema.DefaultTimeout(10 * time.Minute),
},
```

---

## User Documentation

### Usage Example

#### Basic Usage with Default Timeouts
```hcl
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example-instance"
  availability_zone = "ap-guangzhou-3"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-xxxxx"
  subnet_id         = "subnet-xxxxx"
  engine_version    = "13.3"
  root_password     = "your-password"
  charset           = "UTF8"
  project_id        = 0
  memory            = 4
  storage           = 100
  
  # Default timeouts (60 minutes) will be used
}
```

#### Custom Timeout Configuration
```hcl
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example-instance"
  availability_zone = "ap-guangzhou-3"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-xxxxx"
  subnet_id         = "subnet-xxxxx"
  engine_version    = "13.3"
  root_password     = "your-password"
  charset           = "UTF8"
  project_id        = 0
  memory            = 4
  storage           = 100
  
  # Custom timeout configuration
  timeouts {
    create = "90m"  # 90 minutes for creation
    update = "120m" # 120 minutes for updates
  }
}
```

#### Readonly Instance Example
```hcl
resource "tencentcloud_postgresql_readonly_instance" "example" {
  master_db_instance_id = "postgres-xxxxx"
  name                  = "readonly-instance"
  zone                  = "ap-guangzhou-3"
  project_id            = 0
  vpc_id                = "vpc-xxxxx"
  subnet_id             = "subnet-xxxxx"
  db_version            = "13.3"
  memory                = 4
  storage               = 100
  
  # Custom timeout configuration
  timeouts {
    create = "90m"
    update = "120m"
  }
}
```

---

## Impact Analysis

### Backward Compatibility
✅ **Fully backward compatible**
- Existing configurations without timeout block will use default values
- No breaking changes to resource schema
- No changes to required/optional fields

### Performance Impact
- Minimal: Only adds timeout configuration parsing
- No impact on API call patterns
- Better user experience with appropriate timeout values

### Testing Requirements
1. Test default timeout behavior (60 minutes)
2. Test custom timeout configuration
3. Test create operation with timeout
4. Test update operation (memory/storage/CPU) with timeout
5. Verify backward compatibility with existing configurations

---

## Implementation Checklist

### Phase 1: Code Implementation
- [ ] Add Timeouts block to `tencentcloud_postgresql_instance` schema
- [ ] Update create function in `tencentcloud_postgresql_instance`
- [ ] Update memory/storage/CPU upgrade logic in `tencentcloud_postgresql_instance`
- [ ] Add Timeouts block to `tencentcloud_postgresql_readonly_instance` schema
- [ ] Update create function in `tencentcloud_postgresql_readonly_instance`
- [ ] Update memory/storage/CPU upgrade logic in `tencentcloud_postgresql_readonly_instance`
- [ ] Run `go fmt` on all modified files

### Phase 2: Testing
- [ ] Verify code compiles without errors
- [ ] Run linter checks
- [ ] Test create operation with default timeout
- [ ] Test update operation with default timeout
- [ ] Test custom timeout configuration

### Phase 3: Documentation
- [ ] Update resource documentation
- [ ] Add timeout configuration examples
- [ ] Document default values

---

## Timeline

**Estimated Implementation Time**: 1-2 hours

### Breakdown
- Schema modifications: 15 minutes
- Create function updates: 20 minutes
- Update function modifications: 30 minutes
- Code formatting and linting: 10 minutes
- Testing and verification: 30 minutes
- Documentation: 15 minutes

---

## Dependencies

### Code Dependencies
- `github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema`
- `time` package (standard library)

### No Breaking Changes
- No changes to required fields
- No changes to API interactions
- No changes to resource IDs or naming

---

## Risks and Mitigations

### Risk 1: Timeout Too Short
**Risk**: Default 60-minute timeout may still be insufficient for very large instances

**Mitigation**: 
- Users can override with custom timeout values
- 60 minutes is reasonable for most scenarios based on typical PostgreSQL instance operations

### Risk 2: Backward Compatibility
**Risk**: Existing users might expect different timeout behavior

**Mitigation**:
- Default behavior unchanged (uses appropriate timeout)
- Only enhancement, no breaking changes
- Fully backward compatible

---

## Future Enhancements

1. **Delete Timeout**: Consider adding delete timeout in future if needed
2. **Read Timeout**: Generally not needed for read operations
3. **Documentation**: Add troubleshooting guide for timeout scenarios

---

## References

### Internal References
- CVM Instance: `tencentcloud/services/cvm/resource_tc_instance.go`
- MySQL Instance: `tencentcloud/services/cdb/resource_tc_mysql_instance.go`
- CLB Instance: `tencentcloud/services/clb/resource_tc_clb_instance.go`

### External References
- [Terraform Schema Timeouts](https://www.terraform.io/plugin/sdkv2/schemas/schema-behaviors#timeouts)
- [TencentCloud PostgreSQL API Documentation](https://cloud.tencent.com/document/product/409)

---

## Approval

**Proposed by**: Development Team  
**Date**: 2026-03-24  
**Priority**: Medium  
**Complexity**: Low  

**Required Approvals**:
- [ ] Technical Lead
- [ ] Code Review
- [ ] QA Testing

---

## Notes

1. All new functions or code blocks should be placed at the end of the resource file
2. Every modified Go file must be formatted with `go fmt` before committing
3. Follow existing code patterns and conventions in the codebase
4. Timeout configuration is optional for users - if not specified, default values are used
