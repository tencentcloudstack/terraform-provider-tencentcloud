# Quick Reference: PostgreSQL Instance Custom Timeouts

## TL;DR (Too Long; Didn't Read)

**Goal**: Add 60-minute default timeouts for PostgreSQL instance create and update operations

**Files to modify**:
1. `resource_tc_postgresql_instance.go` - Add Timeouts, update 2 locations
2. `resource_tc_postgresql_readonly_instance.go` - Add Timeouts, update 2 locations

**Time Required**: ~1-2 hours

---

## Quick Implementation Guide

### Step 1: Add Timeouts to Schema (Both Resources)

**Location**: In `ResourceTencentCloudPostgresqlInstance()` and `ResourceTencentCloudPostgresqlReadonlyInstance()`

**Add after `Importer` block**:
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(60 * time.Minute),
    Update: schema.DefaultTimeout(60 * time.Minute),
},
```

### Step 2: Update Create Function (Both Resources)

**Find**:
```go
resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
```

**Replace with**:
```go
resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
```

### Step 3: Update Scaling Logic (Both Resources)

**Find**:
```go
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    // ...
    resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
```

**Replace with**:
```go
if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
    // ...
    resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
```

### Step 4: Format and Verify

```bash
# Format both files
go fmt tencentcloud/services/postgresql/resource_tc_postgresql_instance.go
go fmt tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go

# Compile
go build ./tencentcloud/services/postgresql/...
```

---

## Code Snippets

### Complete Timeouts Block
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(60 * time.Minute),
    Update: schema.DefaultTimeout(60 * time.Minute),
},
```

### Usage in Create
```go
err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    instanceId, inErr = postgresqlService.CreatePostgresqlInstance(...)
    if inErr != nil {
        return tccommon.RetryError(inErr)
    }
    return nil
})
```

### Usage in Update (Scaling)
```go
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

---

## Checklist

### For `resource_tc_postgresql_instance.go`
- [ ] Add `Timeouts` block to schema (after `Importer`)
- [ ] Update create function: replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutCreate)`
- [ ] Update scaling logic: replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutUpdate)` in memory/storage/CPU block
- [ ] Run `go fmt` on file
- [ ] Verify compilation

### For `resource_tc_postgresql_readonly_instance.go`
- [ ] Add `Timeouts` block to schema (after `Importer`)
- [ ] Update create function: replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutCreate)`
- [ ] Update scaling logic: replace `tccommon.WriteRetryTimeout` with `d.Timeout(schema.TimeoutUpdate)` in memory/storage/CPU block
- [ ] Run `go fmt` on file
- [ ] Verify compilation

---

## User Example

### Basic Usage (Default 60m)
```hcl
resource "tencentcloud_postgresql_instance" "example" {
  name              = "my-postgres"
  availability_zone = "ap-guangzhou-3"
  memory            = 4
  storage           = 100
  # ... other config ...
}
```

### Custom Timeout
```hcl
resource "tencentcloud_postgresql_instance" "large" {
  name              = "large-postgres"
  availability_zone = "ap-guangzhou-3"
  memory            = 32
  storage           = 2000
  # ... other config ...
  
  timeouts {
    create = "90m"
    update = "120m"
  }
}
```

---

## Important Notes

1. **Default timeout**: 60 minutes for both create and update
2. **Update timeout applies to**: Memory, storage, CPU changes only
3. **Backward compatible**: Existing configs work without changes
4. **Code placement**: Add Timeouts block AFTER Importer, BEFORE Schema
5. **Format**: Always run `go fmt` after changes

---

## Testing Commands

```bash
# Compile
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go build ./tencentcloud/services/postgresql/...

# Run linter
golint tencentcloud/services/postgresql/resource_tc_postgresql_instance.go
golint tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go

# Format
go fmt tencentcloud/services/postgresql/resource_tc_postgresql_instance.go
go fmt tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance.go
```

---

## Reference Examples

### CVM Instance (15m create)
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(15 * time.Minute),
},
```

### MySQL Instance (20m create/delete)
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(20 * time.Minute),
    Delete: schema.DefaultTimeout(20 * time.Minute),
},
```

### CLB Instance (10m create/update)
```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(10 * time.Minute),
    Update: schema.DefaultTimeout(10 * time.Minute),
},
```

---

## Common Mistakes to Avoid

❌ **Don't**:
- Forget to add `Timeouts` block
- Put `Timeouts` in the wrong location (must be after `Importer`)
- Miss any `resource.Retry` calls that should use the new timeout
- Forget to run `go fmt`

✅ **Do**:
- Add `Timeouts` block after `Importer`, before `Schema`
- Update ALL relevant `resource.Retry` calls
- Run `go fmt` on modified files
- Test compilation after changes

---

## Time Estimates

| Task | Time |
|------|------|
| Add schema Timeouts (2 files) | 10 min |
| Update create functions (2 files) | 20 min |
| Update scaling logic (2 files) | 20 min |
| Format and verify | 10 min |
| **Total** | **~60 min** |

---

## Success Criteria

✅ Code compiles without errors  
✅ Both resources have Timeouts block  
✅ Create operations use Create timeout  
✅ Scaling operations use Update timeout  
✅ Files formatted with `go fmt`  
✅ No new linter warnings  
✅ Backward compatible  

---

## Questions?

- See [proposal.md](./proposal.md) for detailed technical spec
- See [tasks.md](./tasks.md) for step-by-step instructions
- See [README.md](./README.md) for overview and examples
