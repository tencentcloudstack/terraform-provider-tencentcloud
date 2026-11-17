## Why

The `tencentcloud_cbs_storage` resource currently does not support the `instance_id` parameter, which allows users to specify a CVM instance for auto-mounting and auto-initializing a data disk during disk creation. Without this parameter, users must separately attach and initialize disks after creation, which adds operational overhead. The `CreateDisks` API already supports `AutoMountConfiguration.InstanceId` for this purpose, and `DescribeDisks` returns the `InstanceId` field indicating the mounted CVM instance, but the Terraform resource exposes neither.

## What Changes

### Auto-Mount Instance ID Parameter
- Add a new optional schema field `instance_id` (string) to the `tencentcloud_cbs_storage` resource
- In the Create function, when `instance_id` is set, populate `request.AutoMountConfiguration.InstanceId` with the instance ID so the disk is auto-mounted and auto-initialized during creation
- In the Read function, populate `instance_id` from `Disk.InstanceId` (the CVM instance the disk is currently mounted to) so Terraform state reflects the actual mount status
- The `instance_id` field is **read-only after creation** (ForceNew not required, but no update API supports changing the mount via AutoMountConfiguration): since `ModifyDiskAttributes` does not contain `AutoMountConfiguration`, changes to `instance_id` are not supported via update; users must recreate the disk to change the auto-mount target

## Capabilities

### New Capabilities
- `cbs-storage-auto-mount`: Support for specifying the CVM instance ID to auto-mount and auto-initialize a CBS data disk during creation, and reading back the mounted instance ID from the API response.

### Modified Capabilities
<!-- No existing spec-level capability changes -->

## Impact

**Affected Code:**
- `tencentcloud/services/cbs/resource_tc_cbs_storage.go`: Schema definition (add `instance_id` field), Create function (populate `AutoMountConfiguration.InstanceId`), Read function (populate `instance_id` from `Disk.InstanceId`)
- `tencentcloud/services/cbs/resource_tc_cbs_storage.md`: Documentation with usage example demonstrating auto-mount
- `tencentcloud/services/cbs/resource_tc_cbs_storage_test.go`: Test coverage for the new parameter

**User Impact:**
- **Non-breaking**: Existing configurations continue to work unchanged (the new field is Optional)
- Users can now create a CBS data disk and have it auto-mounted and auto-initialized to a specified CVM instance in a single Terraform apply
- The `instance_id` field reflects the actual mounted instance after Read, enabling drift detection

**Dependencies:**
- Requires no changes to the SDK (`AutoMountConfiguration.InstanceId` already exists in `cbs.CreateDisksRequest` and `Disk.InstanceId` exists in the `DescribeDisks` response)
- Works with the existing `CreateDisks` and `DescribeDisks` APIs
