# Add Custom Timeouts for PostgreSQL Instance Resources

## Overview
This OpenSpec proposal adds custom timeout configuration support to PostgreSQL instance resources (`tencentcloud_postgresql_instance` and `tencentcloud_postgresql_readonly_instance`), allowing users to configure operation timeouts for create and update operations.

---

## Quick Summary

### What's Changing?
- Add `Timeouts` configuration block to both PostgreSQL instance resources
- Set default create timeout to **60 minutes**
- Set default update timeout to **60 minutes** (for memory/storage/CPU scaling)
- Users can override defaults with custom timeout values

### Why?
- PostgreSQL instance creation can take significant time
- Resource scaling (memory/storage/CPU) requires instance restart and extended time
- Users need control over operation timeouts to avoid premature failures
- Consistency with other TencentCloud resources (CVM, MySQL, CLB)

---

## Resources Affected

### 1. `tencentcloud_postgresql_instance`
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`

**Changes**:
- ✅ Add Timeouts block to schema (60m create, 60m update)
- ✅ Update create function to use `d.Timeout(schema.TimeoutCreate)`
- ✅ Update memory/storage/CPU scaling to use `d.Timeout(schema.TimeoutUpdate)`

### 2. `tencentcloud_postgresql_readonly_instance`
**File**: `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`

**Changes**:
- ✅ Add Timeouts block to schema (60m create, 60m update)
- ✅ Update create function to use `d.Timeout(schema.TimeoutCreate)`
- ✅ Update memory/storage/CPU scaling to use `d.Timeout(schema.TimeoutUpdate)`

---

## Usage Examples

### Default Behavior (60-minute timeouts)
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
  
  # Default 60-minute timeouts apply automatically
}
```

### Custom Timeout Configuration
```hcl
resource "tencentcloud_postgresql_instance" "large_instance" {
  name              = "large-instance"
  availability_zone = "ap-guangzhou-3"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-xxxxx"
  subnet_id         = "subnet-xxxxx"
  engine_version    = "13.3"
  root_password     = "your-password"
  charset           = "UTF8"
  project_id        = 0
  memory            = 32  # Large instance
  storage           = 2000
  
  # Custom timeouts for large instance
  timeouts {
    create = "120m"  # 2 hours for creation
    update = "90m"   # 1.5 hours for updates
  }
}
```

### Readonly Instance Example
```hcl
resource "tencentcloud_postgresql_readonly_instance" "readonly" {
  master_db_instance_id = "postgres-xxxxx"
  name                  = "readonly-instance"
  zone                  = "ap-guangzhou-3"
  project_id            = 0
  vpc_id                = "vpc-xxxxx"
  subnet_id             = "subnet-xxxxx"
  db_version            = "13.3"
  memory                = 8
  storage               = 500
  
  timeouts {
    create = "90m"
    update = "60m"
  }
}
```

---

## Technical Details

### Timeout Configuration
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(60 * time.Minute),  // 60 minutes default
    Update: schema.DefaultTimeout(60 * time.Minute),  // 60 minutes default
},
```

### When Timeouts Apply

#### Create Timeout
- Used during initial instance creation
- Applies to the entire creation process including:
  - API call to create instance
  - Waiting for instance to become available
  - Initial configuration

#### Update Timeout
- Applies specifically to resource scaling operations when:
  - `memory` field changes
  - `storage` field changes
  - `cpu` field changes
- Used during the upgrade process including:
  - API call to modify specifications
  - Instance restart
  - Waiting for instance to become available after restart

---

## Implementation Details

### Key Code Changes

#### 1. Schema Addition
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
        
        // Added Timeouts configuration
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

#### 2. Create Function Update
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

#### 3. Update Function - Resource Scaling
```go
// Before
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        inErr = postgresqlService.UpgradePostgresqlInstance(...)
        if inErr != nil {
            return tccommon.RetryError(inErr)
        }
        return nil
    })
}

// After
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    outErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
        inErr = postgresqlService.UpgradePostgresqlInstance(...)
        if inErr != nil {
            return tccommon.RetryError(inErr)
        }
        return nil
    })
}
```

---

## Benefits

### For Users
1. **Control**: Users can customize timeouts based on their instance size and requirements
2. **Reliability**: Appropriate timeouts prevent premature failure of long-running operations
3. **Flexibility**: Different timeout values can be set for create vs update operations
4. **Consistency**: Matches timeout behavior of other TencentCloud resources

### For Developers
1. **Maintainability**: Follows established patterns in the codebase
2. **Simplicity**: Minimal code changes, low complexity
3. **Backward Compatibility**: Existing configurations work without modification
4. **Testability**: Easy to test with different timeout values

---

## Backward Compatibility

✅ **Fully Backward Compatible**

- Existing resource configurations work without any changes
- If no timeout block is specified, default 60-minute timeouts apply
- No breaking changes to resource schema or behavior
- No changes to required/optional fields

### Migration Path
**No migration required!** Existing users can continue using their current configurations, and the resources will automatically use the new default 60-minute timeouts.

---

## Testing Checklist

### Unit Tests
- [ ] Test resource with default timeouts (no timeout block specified)
- [ ] Test resource with custom timeout configuration
- [ ] Test create operation respects Create timeout
- [ ] Test update operation (memory/storage/CPU) respects Update timeout
- [ ] Test backward compatibility with existing configurations

### Integration Tests
- [ ] Create actual PostgreSQL instance and verify timeout behavior
- [ ] Scale instance resources and verify Update timeout applies
- [ ] Test with custom timeout values
- [ ] Verify timeout errors are properly reported

---

## Documentation

### Terraform Configuration
Users can configure timeouts in their Terraform files:

```hcl
resource "tencentcloud_postgresql_instance" "example" {
  # ... resource configuration ...
  
  timeouts {
    create = "90m"   # Optional: defaults to 60m
    update = "120m"  # Optional: defaults to 60m
  }
}
```

### Timeout Values
- Specified in string format: `"30m"`, `"1h"`, `"90m"`, `"2h30m"`, etc.
- Common values:
  - `"30m"` = 30 minutes
  - `"60m"` or `"1h"` = 1 hour
  - `"90m"` or `"1h30m"` = 1.5 hours
  - `"120m"` or `"2h"` = 2 hours

---

## Implementation Timeline

**Estimated Total Time**: 1-2 hours

### Breakdown
1. Schema modifications: **15 minutes**
2. Create function updates: **20 minutes**
3. Update function modifications: **30 minutes**
4. Code formatting and linting: **10 minutes**
5. Testing and verification: **30 minutes**
6. Documentation: **15 minutes**

---

## Files Modified

### Primary Changes
1. `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`
   - Add Timeouts block
   - Update create function
   - Update memory/storage/CPU scaling logic

2. `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go`
   - Add Timeouts block
   - Update create function
   - Update memory/storage/CPU scaling logic

### Documentation (Optional)
1. `website/docs/r/postgresql_instance.html.markdown`
2. `website/docs/r/postgresql_readonly_instance.html.markdown`

---

## Related Resources

### Similar Implementations
This change follows the pattern used in other TencentCloud resources:

1. **CVM Instance** (`resource_tc_instance.go`)
   - Create timeout: 15 minutes

2. **MySQL Instance** (`resource_tc_mysql_instance.go`)
   - Create timeout: 20 minutes
   - Delete timeout: 20 minutes

3. **CLB Instance** (`resource_tc_clb_instance.go`)
   - Create timeout: 10 minutes
   - Update timeout: 10 minutes

### Reference Documentation
- [Terraform Schema Timeouts](https://www.terraform.io/plugin/sdkv2/schemas/schema-behaviors#timeouts)
- [TencentCloud PostgreSQL Documentation](https://cloud.tencent.com/document/product/409)
- [Terraform Resource Timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts)

---

## FAQs

### Q: Do I need to update my existing Terraform configurations?
**A**: No! Your existing configurations will work as-is. The resources will automatically use 60-minute default timeouts.

### Q: Can I specify different timeouts for create and update?
**A**: Yes! You can customize each timeout independently:
```hcl
timeouts {
  create = "90m"
  update = "120m"
}
```

### Q: What happens if I don't specify a timeout?
**A**: Default 60-minute timeouts are used for both create and update operations.

### Q: Does the update timeout apply to all updates?
**A**: No, the update timeout specifically applies to memory, storage, and CPU scaling operations. Other updates may use different timeout mechanisms.

### Q: What if my operation exceeds the timeout?
**A**: Terraform will report a timeout error. You can increase the timeout value in your configuration and retry.

### Q: Can I set the timeout to 0 (unlimited)?
**A**: While technically possible, it's not recommended. Always set a reasonable timeout value based on your instance size and requirements.

---

## Support

### Getting Help
- Review the detailed [proposal.md](./proposal.md) for technical specifications
- Check [tasks.md](./tasks.md) for implementation details
- Refer to reference implementations in other TencentCloud resources

### Reporting Issues
If you encounter any issues with timeout configuration:
1. Verify your timeout syntax is correct
2. Check that timeout values are reasonable for your instance size
3. Review Terraform logs for detailed error messages

---

## Status

**Current Status**: ✅ Proposed  
**Proposed Date**: 2026-03-24  
**Target Implementation**: 2026-03-24  
**Priority**: Medium  
**Complexity**: Low  

---

## Approval Checklist

- [ ] Technical specification reviewed
- [ ] Code changes reviewed
- [ ] Testing strategy approved
- [ ] Documentation plan approved
- [ ] Ready for implementation

---

## Next Steps

1. Review and approve this proposal
2. Implement changes following [tasks.md](./tasks.md)
3. Run tests and verification
4. Update documentation
5. Submit for code review
6. Merge to main branch

---

**Proposal Created**: 2026-03-24  
**Last Updated**: 2026-03-24  
**Version**: 1.0
