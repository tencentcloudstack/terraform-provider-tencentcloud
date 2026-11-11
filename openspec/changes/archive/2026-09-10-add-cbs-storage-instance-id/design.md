## Context

The `tencentcloud_cbs_storage` resource manages CBS (Cloud Block Storage) disks via the `CreateDisks`, `DescribeDisks`, `ModifyDiskAttributes`, and `DeleteDisks` APIs. The resource currently supports creating disks with various configurations (disk type, size, encryption, prepaid/postpaid billing, tags, etc.) but does not support the auto-mount feature.

The `CreateDisks` API accepts an `AutoMountConfiguration` struct that allows a newly created data disk to be automatically mounted and initialized to a specified CVM instance during creation. The struct contains:
- `InstanceId []*string` - the CVM instance ID(s) to mount to (a list, but typically a single instance for a data disk)
- `MountPoint []*string` - mount point inside the instance
- `FileSystemType *string` - filesystem type (ext4, xfs)

The `DescribeDisks` API returns `Disk.InstanceId` (`*string`) which indicates the CVM instance the disk is currently mounted to.

**Current State:**
- The SDK's `CreateDisksRequest` already contains the `AutoMountConfiguration` field (line 942 in models.go)
- The SDK's `Disk` struct already contains the `InstanceId` field (line 2942 in models.go)
- The `resource_tc_cbs_storage.go` resource does not expose either field in its schema
- The `tags` parameter already exists and is fully handled via the tag service's `ModifyResourceTags` (note: the cbs package has no `ModifyTags` API; tags are managed through the `tag/v20180813` package)

**Constraints:**
- Must maintain backward compatibility (only adding Optional fields)
- Must follow existing patterns in `resource_tc_cbs_storage.go`
- The `AutoMountConfiguration` is only available in `CreateDisks`, NOT in `ModifyDiskAttributes`, so the auto-mount target cannot be changed after creation via this mechanism
- Per the task requirement, only ONE new parameter (`instance_id`) is being added to this resource

## Goals / Non-Goals

**Goals:**
- Add support for the `instance_id` parameter to enable auto-mounting a CBS data disk to a CVM instance during creation
- Populate `instance_id` from the API response in Read so Terraform state reflects the actual mount status
- Follow the existing flat schema pattern used by the resource
- Maintain backward compatibility

**Non-Goals:**
- Not exposing `MountPoint` or `FileSystemType` from `AutoMountConfiguration` (only `instance_id` is in scope per the requirement)
- Not supporting updates to the auto-mount target (no update API exists for `AutoMountConfiguration`); changes require recreating the disk
- Not adding a `tags` parameter (it already exists in the resource schema and is fully implemented)
- Not refactoring existing parameter handling

## Decisions

### Decision 1: Flat Schema Field for instance_id

**Choice:** Add a single flat schema field `instance_id` (Type=String, Optional=true) at the root level.

**Rationale:**
- Consistent with existing flat schema fields in `resource_tc_cbs_storage.go` (e.g., `storage_name`, `snapshot_id`)
- The requirement maps `request.AutoMountConfiguration.InstanceId` → SchemaName `InstanceId`, exposing it as a flat top-level `instance_id` field is the simplest and most user-friendly approach
- Avoids introducing a nested `auto_mount_configuration` block for a single sub-field

**Alternatives Considered:**
- Nested `auto_mount_configuration` block with `instance_id` inside: more API-aligned but over-engineered for a single parameter and inconsistent with the existing flat pattern

### Decision 2: Create Function - Map instance_id to AutoMountConfiguration

**Choice:** In the Create function, when `instance_id` is set, populate `request.AutoMountConfiguration.InstanceId` with the value wrapped in a `[]*string` slice.

**Rationale:**
- The SDK field `AutoMountConfiguration.InstanceId` is `[]*string` (a list), but for a single data disk the typical use case is mounting to one instance
- The Terraform schema field is a single string; converting it to a single-element `[]*string` matches the API expectation
- Only set `AutoMountConfiguration` when `instance_id` is provided, to avoid changing behavior for existing configurations

**Implementation:**
```go
if v, ok := d.GetOk("instance_id"); ok {
    request.AutoMountConfiguration = &cbs.AutoMountConfiguration{
        InstanceId: []*string{helper.String(v.(string))},
    }
}
```

### Decision 3: Read Function - Populate instance_id from Disk.InstanceId

**Choice:** In the Read function, populate `instance_id` from `storage.InstanceId` (the `Disk.InstanceId` field from the `DescribeDisks` response).

**Rationale:**
- `Disk.InstanceId` indicates the CVM instance the disk is currently mounted to
- Populating this in Read enables Terraform to detect drift (e.g., if the disk was detached/attached out-of-band)
- The existing Read function already follows the pattern of setting fields directly from the `storage *cbs.Disk` object
- Must guard against nil before calling `d.Set()` per the project's code requirements

**Implementation:**
```go
if storage.InstanceId != nil {
    _ = d.Set("instance_id", storage.InstanceId)
}
```

### Decision 4: Update Function - No AutoMountConfiguration Update

**Choice:** The Update function does not handle `instance_id` changes via `AutoMountConfiguration` (no update API supports it). The field is not marked ForceNew.

**Rationale:**
- `ModifyDiskAttributes` does not contain `AutoMountConfiguration`, so the auto-mount target cannot be updated
- The `instance_id` field read from `DescribeDisks` reflects the currently mounted instance, which can change if the disk is attached/detached via other means
- Not marking ForceNew: the field is read from API state and represents current mount status; marking ForceNew would force unnecessary recreation when the disk is simply detached and reattached
- If a user wants to change the auto-mount target, they must recreate the disk (destroy + create)

**Alternatives Considered:**
- Mark `instance_id` as ForceNew: rejected because the read-back value represents current mount status which legitimately changes through attach/detach operations; ForceNew would cause unexpected recreation

## Risks / Trade-offs

**Risk:** API validation errors if the specified instance does not exist or is incompatible
- **Mitigation:** Clear documentation; the API will return validation errors during apply

**Risk:** `Disk.InstanceId` may be empty string when the disk is not mounted, causing `d.Set("instance_id", "")` to overwrite a previously set value
- **Mitigation:** Guard with `if storage.InstanceId != nil` check; an empty string is a valid state (disk unmounted) and should be reflected in state

**Trade-off:** The `instance_id` field cannot be updated after creation (no update API for AutoMountConfiguration)
- **Acceptable:** Auto-mount is a creation-time operation; users who need to change the mount target can use the `tencentcloud_cbs_storage_attachment` resource or recreate the disk

**Trade-off:** Only `instance_id` is exposed from `AutoMountConfiguration`, not `MountPoint` or `FileSystemType`
- **Acceptable:** Per the requirement, only one parameter is being added; the API will use default values for the omitted sub-fields
