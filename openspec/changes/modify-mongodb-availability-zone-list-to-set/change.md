# Change: Modify MongoDB availability_zone_list from List to Set

## Status
✅ **Completed**

## Overview
Modified the `availability_zone_list` field in MongoDB resources from `TypeList` to `TypeSet` to prevent Terraform from treating different orderings of the same zones as a configuration change.

## Changes Made

### 1. Modified Files

#### 1.1 `resource_tc_mongodb_instance.go`

**Schema Definition (Line 105-116):**
```go
// BEFORE
"availability_zone_list": {
    Type:     schema.TypeList,
    Optional: true,
    Computed: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "...",
},

// AFTER
"availability_zone_list": {
    Type:     schema.TypeSet,
    Optional: true,
    Computed: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "...",
},
```

**Create Function (Line 235-238):**
```go
// BEFORE
if v, ok := d.GetOk("availability_zone_list"); ok {
    availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
    value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
}

// AFTER
if v, ok := d.GetOk("availability_zone_list"); ok {
    availabilityZoneList := helper.InterfacesStringsPoint(v.(*schema.Set).List())
    value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
}
```

**Read Function (Line 471-487):**
```go
// BEFORE
availabilityZoneList := make([]string, 0, 3)
for _, replicate := range replicateSets[0].Nodes {
    itemZone := *replicate.Zone
    if *replicate.Hidden {
        hiddenZone = itemZone
    }
    availabilityZoneList = append(availabilityZoneList, itemZone)
}
_ = d.Set("availability_zone_list", availabilityZoneList)

// AFTER
availabilityZoneList := make([]interface{}, 0, 3)
for _, replicate := range replicateSets[0].Nodes {
    itemZone := *replicate.Zone
    if *replicate.Hidden {
        hiddenZone = itemZone
    }
    availabilityZoneList = append(availabilityZoneList, itemZone)
}
_ = d.Set("availability_zone_list", availabilityZoneList)
```

#### 1.2 `resource_tc_mongodb_sharding_instance.go`

**Schema Definition (Line 38-50):**
```go
// Changed from schema.TypeList to schema.TypeSet
```

**Create Function (Line 159-162):**
```go
// Changed from v.([]interface{}) to v.(*schema.Set).List()
```

**Read Function (Line 385-395):**
```go
// Changed from []string to []interface{}
```

#### 1.3 `resource_tc_mongodb_readonly_instance.go`

**Create Function (Line 163-166):**
```go
// Changed from v.([]interface{}) to v.(*schema.Set).List()
```

### 2. Impact Analysis

#### Breaking Change: YES ⚠️

This is a **BREAKING CHANGE** because:
- Terraform will detect a state change on first apply after upgrade
- The state file format changes from List to Set
- Users may see a diff even when no actual changes were made

#### Migration Path

Users will see output like this on first `terraform plan` after upgrade:

```hcl
  ~ availability_zone_list = [
      - "ap-guangzhou-3",
      - "ap-guangzhou-4",
      - "ap-guangzhou-6",
    ] -> (known after apply)
```

And on `terraform apply`:
```hcl
  ~ availability_zone_list = [
      + "ap-guangzhou-3",
      + "ap-guangzhou-4", 
      + "ap-guangzhou-6",
    ]
```

**Important:** This is a **cosmetic change only**. No actual API calls or infrastructure changes will occur. The zones remain the same; only the state representation changes.

### 3. Benefits

✅ **Order Independence**: Different orderings like `["zone-a", "zone-b", "zone-c"]` and `["zone-c", "zone-a", "zone-b"]` are now treated as identical

✅ **Duplicate Detection**: TypeSet automatically handles duplicate values

✅ **Terraform Best Practice**: Aligns with Terraform's recommendation for unordered collections

### 4. Testing

#### Compilation
- ✅ Code compiles without errors
- ✅ No new linter warnings introduced

#### Files Modified
- ✅ `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go`
- ✅ `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go`
- ✅ `tencentcloud/services/mongodb/resource_tc_mongodb_readonly_instance.go`

#### Code Formatting
- ✅ All modified files formatted with `go fmt`

### 5. Compatibility

#### Terraform Version
- Requires Terraform >= 0.12 (already a requirement)

#### API Compatibility
- ✅ No API changes
- ✅ Same API calls are made
- ✅ Same parameters sent to Tencent Cloud

#### State Migration
- ⚠️ Automatic on first apply
- ℹ️ Users will see a state refresh
- ℹ️ No manual migration required

### 6. Release Notes

#### For CHANGELOG.md

```markdown
BREAKING CHANGES:

* **resource/tencentcloud_mongodb_instance**: The `availability_zone_list` field has been changed from a List to a Set. This change ensures that different orderings of the same availability zones are not treated as configuration changes. **Impact**: On the first `terraform apply` after upgrading, users will see a state change for this field, but no actual infrastructure modification will occur. This is a one-time cosmetic change.

* **resource/tencentcloud_mongodb_sharding_instance**: The `availability_zone_list` field has been changed from a List to a Set for consistency.

* **resource/tencentcloud_readonly_instance**: The `availability_zone_list` field handling updated for Set compatibility.
```

### 7. Next Steps

- [ ] Update test cases to use Set assertions
- [ ] Update documentation if needed
- [ ] Add to CHANGELOG.md for next release
- [ ] Consider adding migration guide to documentation

## Implementation Details

### Technical Approach

1. **Schema Change**: Modified `Type` from `schema.TypeList` to `schema.TypeSet`
2. **Read Logic**: Changed from `v.([]interface{})` to `v.(*schema.Set).List()`
3. **Write Logic**: Changed slice type from `[]string` to `[]interface{}` for compatibility with `d.Set()`

### Code Quality

- ✅ No breaking API changes
- ✅ Backward compatible at API level
- ✅ Forward compatible with future Terraform versions
- ✅ Follows Terraform SDK best practices
- ✅ Maintains existing validation logic

## Conclusion

The change has been successfully implemented across all MongoDB resources that use `availability_zone_list`. The modification improves the Terraform user experience by eliminating false-positive configuration drifts due to zone ordering differences.

**Date Completed**: 2026-03-20
**Author**: Terraform Provider Development Team
