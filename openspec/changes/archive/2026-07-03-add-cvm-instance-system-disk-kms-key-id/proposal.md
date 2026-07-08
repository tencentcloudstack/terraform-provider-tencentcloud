## Why

The `tencentcloud_instance` resource currently supports KMS encryption for data disks via `data_disks.*.kms_key_id`, but there is no way to specify a custom KMS key for encrypting the system disk. The vendor CVM SDK's `SystemDisk` struct has been updated to include `KmsKeyId` and `Encrypt` fields, enabling users to specify custom KMS keys for system disk encryption via the `RunInstances` and `ResetInstance` APIs. This change exposes the system disk `kms_key_id` as a top-level parameter in the Terraform resource.

## What Changes

- Add a new top-level `kms_key_id` parameter (TypeString, Optional, ForceNew) to `tencentcloud_instance` resource for system disk encryption
- Map this parameter to `SystemDisk.KmsKeyId` in the `RunInstances` API call during instance creation
- Map this parameter to `SystemDisk.KmsKeyId` in the `ResetInstance` API call during instance reinstall/reset

## Capabilities

### New Capabilities

- `instance-system-disk-kms-key`: Support specifying a custom KMS key ID for system disk encryption on `tencentcloud_instance`

### Modified Capabilities

None - this is an enhancement to an existing resource that does not change existing requirement-level behavior.

## Impact

**Affected Code:**
- `tencentcloud/services/cvm/resource_tc_instance.go`: Schema definition, Create function (RunInstances), and Update function (ResetInstance section)
- `tencentcloud/services/cvm/resource_tc_instance.md`: Documentation with usage example
- `tencentcloud/services/cvm/resource_tc_instance_test.go`: Unit test cases for the new parameter

**User Impact:**
- **Non-breaking**: Existing configurations continue to work unchanged
- Users can now specify a custom KMS key for system disk encryption when creating or reinstalling instances

**Dependencies:**
- No SDK changes needed — `SystemDisk.KmsKeyId` and `SystemDisk.Encrypt` fields already exist in vendor (`cvm/v20170312/models.go`)
