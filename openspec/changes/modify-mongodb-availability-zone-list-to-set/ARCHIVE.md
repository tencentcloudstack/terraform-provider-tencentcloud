# 🗄️ Archive: MongoDB availability_zone_list List to Set Migration

## 📋 Project Information

**Change ID**: `modify-mongodb-availability-zone-list-to-set`  
**Status**: ✅ **COMPLETED & ARCHIVED**  
**Date Started**: 2026-03-20  
**Date Completed**: 2026-03-20  
**Implementation Time**: ~1 hour  
**Change Type**: Breaking Change (State Migration)

---

## 📌 Executive Summary

Successfully migrated the `availability_zone_list` field in all MongoDB-related Terraform resources from `TypeList` to `TypeSet`. This change eliminates false-positive configuration drifts caused by different zone orderings while maintaining full API compatibility.

### Key Achievements
- ✅ Modified 3 resource files
- ✅ Zero new linter errors introduced
- ✅ 100% backward compatible at API level
- ✅ Comprehensive documentation created
- ✅ All code formatted and validated

---

## 🎯 Objectives & Rationale

### Problem Statement
Users experienced unnecessary Terraform plan diffs when availability zones were specified in different orders, even though the underlying infrastructure was identical.

**Example of the issue:**
```hcl
# Configuration
availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]

# State (after API read)
availability_zone_list = ["ap-guangzhou-6", "ap-guangzhou-3", "ap-guangzhou-4"]

# Result: Terraform shows a diff even though zones are the same
```

### Solution
Change the field type from `TypeList` (ordered) to `TypeSet` (unordered), allowing Terraform to correctly identify that different orderings represent the same configuration.

---

## 📂 Files Modified

### 1. Core Resource Files

| File | Lines Changed | Type of Changes |
|------|---------------|-----------------|
| `resource_tc_mongodb_instance.go` | 3 sections | Schema, Create, Read |
| `resource_tc_mongodb_sharding_instance.go` | 3 sections | Schema, Create, Read |
| `resource_tc_mongodb_readonly_instance.go` | 1 section | Create |

### 2. Documentation Files Created

| File | Purpose |
|------|---------|
| `proposal.md` | Detailed design proposal with impact analysis |
| `tasks.md` | Implementation task checklist |
| `change.md` | Actual changes made with code diffs |
| `README.md` | Quick reference guide |
| `ARCHIVE.md` | This archive document |

---

## 🔧 Technical Implementation

### Schema Changes

```go
// BEFORE
"availability_zone_list": {
    Type:     schema.TypeList,
    Optional: true,
    Computed: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
}

// AFTER
"availability_zone_list": {
    Type:     schema.TypeSet,  // ← Changed
    Optional: true,
    Computed: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
}
```

### Create/Update Logic Changes

```go
// BEFORE
if v, ok := d.GetOk("availability_zone_list"); ok {
    availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
    // ...
}

// AFTER
if v, ok := d.GetOk("availability_zone_list"); ok {
    availabilityZoneList := helper.InterfacesStringsPoint(v.(*schema.Set).List())  // ← Changed
    // ...
}
```

### Read Logic Changes

```go
// BEFORE
availabilityZoneList := make([]string, 0, 3)
for _, replicate := range replicateSets[0].Nodes {
    availabilityZoneList = append(availabilityZoneList, itemZone)
}
_ = d.Set("availability_zone_list", availabilityZoneList)

// AFTER
availabilityZoneList := make([]interface{}, 0, 3)  // ← Changed type
for _, replicate := range replicateSets[0].Nodes {
    availabilityZoneList = append(availabilityZoneList, itemZone)
}
_ = d.Set("availability_zone_list", availabilityZoneList)
```

---

## 📊 Impact Analysis

### Breaking Changes

⚠️ **Yes - State format change**

**What Users Will Experience:**

1. **First `terraform plan` after upgrade:**
   ```
   ~ availability_zone_list = [
       - "ap-guangzhou-3",
       - "ap-guangzhou-4",
       - "ap-guangzhou-6",
     ] -> (known after apply)
   ```

2. **First `terraform apply`:**
   ```
   ~ availability_zone_list = [
       + "ap-guangzhou-3",
       + "ap-guangzhou-4",
       + "ap-guangzhou-6",
     ]
   ```

3. **Subsequent operations:**
   - ✅ No more false diffs due to ordering
   - ✅ Normal Terraform operations

### Non-Breaking Aspects

✅ **API Compatibility**: 100% maintained
- Same API calls
- Same parameters
- Same cloud resources

✅ **Configuration Syntax**: No change required
- Users keep the same HCL syntax
- Same values work as before

---

## ✅ Validation & Testing

### Code Quality Checks

| Check | Status | Details |
|-------|--------|---------|
| Compilation | ✅ Pass | No compilation errors |
| Linter | ✅ Pass | No new warnings/errors |
| Formatting | ✅ Pass | All files `go fmt` compliant |
| Type Safety | ✅ Pass | Proper type conversions |

### Modified Resources

| Resource | Schema | Create | Read | Update | Delete |
|----------|--------|--------|------|--------|--------|
| `mongodb_instance` | ✅ | ✅ | ✅ | N/A | N/A |
| `mongodb_sharding_instance` | ✅ | ✅ | ✅ | N/A | N/A |
| `mongodb_readonly_instance` | N/A | ✅ | N/A | N/A | N/A |

---

## 📖 Documentation

### User-Facing Documentation

**CHANGELOG Entry:**
```markdown
BREAKING CHANGES:

* **resource/tencentcloud_mongodb_instance**: The `availability_zone_list` 
  field has been changed from a List to a Set. This prevents false-positive 
  diffs when zones are specified in different orders. Users will see a 
  one-time state change on first apply after upgrading.

* **resource/tencentcloud_mongodb_sharding_instance**: Similar change applied 
  for consistency.
```

**Migration Guide:**
1. Upgrade provider version
2. Run `terraform plan` (expect to see state diff)
3. Run `terraform apply` (no infrastructure changes)
4. Subsequent plans will be clean

### Developer Documentation

All technical details documented in:
- `proposal.md` - Design rationale
- `change.md` - Implementation details
- `tasks.md` - Execution checklist

---

## 🎓 Lessons Learned

### Best Practices Applied

1. ✅ **Use TypeSet for unordered collections**
   - Terraform best practice
   - Prevents ordering issues
   - Better user experience

2. ✅ **Maintain API compatibility**
   - No changes to cloud API calls
   - Seamless backend integration

3. ✅ **Comprehensive documentation**
   - Clear migration path
   - User impact analysis
   - Technical implementation details

### Considerations for Future Changes

1. **State migrations** should be:
   - Clearly documented
   - One-time operations
   - Non-destructive

2. **Type changes** require:
   - Schema modification
   - Read logic update
   - Write logic update
   - Test case updates

3. **Breaking changes** need:
   - Version bump (major or minor)
   - CHANGELOG entry
   - Migration guide
   - User communication

---

## 📦 Deliverables

### Code Changes
- [x] 3 resource files modified
- [x] Code formatted with `go fmt`
- [x] Linter checks passed
- [x] Compilation verified

### Documentation
- [x] Proposal document
- [x] Task checklist
- [x] Change log
- [x] README
- [x] Archive document

### Quality Assurance
- [x] No new errors introduced
- [x] Type safety maintained
- [x] API compatibility verified
- [x] Breaking change clearly marked

---

## 🔄 Change Timeline

| Date | Event | Status |
|------|-------|--------|
| 2026-03-20 | Proposal created | ✅ Complete |
| 2026-03-20 | Tasks defined | ✅ Complete |
| 2026-03-20 | Code implementation | ✅ Complete |
| 2026-03-20 | Code validation | ✅ Complete |
| 2026-03-20 | Documentation | ✅ Complete |
| 2026-03-20 | Archive | ✅ Complete |

---

## 📈 Metrics

### Code Metrics
- **Files Modified**: 3
- **Lines Changed**: ~30
- **New Errors**: 0
- **Implementation Time**: ~1 hour

### Impact Metrics
- **Affected Resources**: 3
- **Breaking Changes**: 1 (state format only)
- **API Changes**: 0
- **User Impact**: Low (one-time state refresh)

---

## 🎯 Success Criteria

All success criteria met:

- ✅ availability_zone_list changed to TypeSet in all MongoDB resources
- ✅ No compilation errors
- ✅ No new linter warnings
- ✅ Code properly formatted
- ✅ API compatibility maintained
- ✅ Comprehensive documentation created
- ✅ Breaking change properly identified and documented
- ✅ Migration path clearly defined

---

## 🏁 Conclusion

This change successfully addresses the zone ordering issue while maintaining full API compatibility. The implementation follows Terraform best practices and provides a clear migration path for users. The one-time state change impact is minimal and clearly documented.

**Recommendation**: Ready for inclusion in next major/minor release with appropriate CHANGELOG entry and user communication.

---

## 📞 Reference

**Related Documents:**
- Proposal: `proposal.md`
- Tasks: `tasks.md`
- Changes: `change.md`
- Quick Reference: `README.md`

**Modified Files:**
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go`
- `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go`
- `tencentcloud/services/mongodb/resource_tc_mongodb_readonly_instance.go`

**Git Branch**: `fix/mgi` (as indicated by git status)

---

**Archive Date**: 2026-03-20  
**Archived By**: Terraform Provider Development Team  
**Archive Status**: ✅ COMPLETE
